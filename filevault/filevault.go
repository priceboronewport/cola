package filevault
  
import (
  "database/sql"
  "time"
  "errors"
  "path/filepath"
  "os"
  "crypto/sha256"
  "fmt"
  "io"
  "regexp"
  "strings"
  "strconv"
)

type FileVault struct {
  db *sql.DB
  root string
  QueryLimit int
}

func CopyFile(src, dst string) (err error) {
  sfi, err := os.Stat(src)
  if err != nil {
    return
  }
  if !sfi.Mode().IsRegular() {
    return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
  }
  dfi, err := os.Stat(dst)
  if err != nil {
    if !os.IsNotExist(err) {
      return
    }
  } else {
    if !(dfi.Mode().IsRegular()) {
      return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
    }
    if os.SameFile(sfi, dfi) {
      return
    }
  }
  if err = os.Link(src, dst); err == nil {
    return
  }
  err = copyFileContents(src, dst)
  return
}

func copyFileContents(src, dst string) (err error) {
  in, err := os.Open(src)
  if err != nil {
    return
  }
  defer in.Close()
  out, err := os.Create(dst)
  if err != nil {
    return
  }
  defer func() {
    cerr := out.Close()
    if err == nil {
      err = cerr
    }
  }()
  if _, err = io.Copy(out, in); err != nil {
    return
  }
  err = out.Sync()
  return
}

func New(db *sql.DB, root string) *FileVault {
  fv := FileVault{db: db, root: root, QueryLimit: 500}
  return &fv
}

func (fv *FileVault) Extract(file_id int, filename string) (dest_filename string, err error) {
  rows, err := fv.db.Query("select hash_id, path, name from files where file_id=?", file_id)
  if err != nil {
    return
  }
  var hash_id int
  var fpath, fname, fext string
  if rows.Next() {
    rows.Scan(&hash_id, &fpath, &fname)
    fext = filepath.Ext(fname)
    if fext != "" {
      fname = fname[:len(fname)-len(fext)]
    }
  }
  if filename == "" {
    dest_filename = fname+fext
  } else {
    dest_filename = strings.Replace(filename, "{.path}", fpath, -1)
    dest_filename = strings.Replace(dest_filename, "{.name}", fname, -1)
    dest_filename = strings.Replace(dest_filename, "{.ext}", fext, -1)
  }
  hash_str := fmt.Sprintf("%010d", hash_id)
  path := fv.root + hash_str[:1] + "/" + hash_str[1:4] + "/" + hash_str[4:7] + "/"
  err = CopyFile(path + hash_str, dest_filename)  
  return
}

func (fv *FileVault) FileId(filename string) (file_id int, err error) {
  fpath, fname := filepath.Split(filename)
  if rows, err := fv.db.Query("select file_id from files where path=? and name=? order by timestamp desc, file_id desc limit 1", fpath, fname); err == nil {
    defer rows.Close()
    if rows.Next() {
      rows.Scan(&file_id)
    }
  }
  return
}

func (fv *FileVault) Hash(filename string) (hash string, err error) {
  if f, err := os.Open(filename); err == nil {
    defer f.Close()
    h := sha256.New()
    if _, err = io.Copy(h, f); err == nil {
      hash = fmt.Sprintf("%x", h.Sum(nil))
    }
  }
  return
}

func (fv *FileVault) HashId(hash string) (hash_id int, err error) {
  hash_id = 0
  var rows *sql.Rows
  rows, err = fv.db.Query("select hash_id from hashes where hash=?", hash)
  if err == nil {
    defer rows.Close()
    if rows.Next() {
      rows.Scan(&hash_id)
    }
  }
  return
}

func (fv *FileVault) Init() (err error) {
  err = errors.New("Not Supported, Yet.")
  return
}

func (fv *FileVault) Import(filename string, path string, timestamp time.Time) (file_id int, err error) {
  var hash string
  if hash, err = fv.Hash(path); err != nil {
    return
  }
  var hash_id int
  if hash_id, err = fv.HashId(hash); err != nil {
    return
  }
  if hash_id == 0 {
    var fi os.FileInfo
    if fi, err = os.Stat(filename); err != nil {
      return
    }
    if _, err = fv.db.Exec("insert into hashes(hash, size) values(?, ?)", hash, fi.Size()); err != nil {
      return
    }
    if hash_id, err = fv.HashId(hash); err != nil {
      return
    }
    if err = fv.StoreFile(filename, hash_id); err != nil {
      return
    }
  } 
  fpath, fname := filepath.Split(path)
  fv.db.Exec("insert into files(hash_id, path, name, timestamp) values(?,?,?,?)", hash_id, fpath, fname, timestamp)
  if file_id, err = fv.FileId(path); err != nil {
    return
  }
  var reg *regexp.Regexp
  reg, err = regexp.Compile("[^a-zA-Z0-9]+")
  if err != nil {
    return
  }
  word_list := strings.Split(reg.ReplaceAllString(path, " "), " ")
  word_ids := make(map[string]int)
  for _, w := range word_list {
    if w != "" {
      word_ids[w] = 0
    }
  }
  for word, word_id := range word_ids {
    word_id = fv.WordId(word)
    if word_id == 0 {
      fv.db.Exec("insert into words(word) values(?)", word)
      word_id = fv.WordId(word)
    }
    if word_id != 0 {
      fv.db.Exec("insert into file_words(file_id, word_id) values(?,?)", file_id, word_id)
    }
  }
  return
}

func (fv *FileVault) Query(terms string) (file_ids []int, filenames []string, err error) {

  // Parse terms into list of words
  var reg *regexp.Regexp
  reg, err = regexp.Compile("[^a-zA-Z0-9]+")
  if err != nil {
    return
  }
  word_list := strings.Split(reg.ReplaceAllString(terms, " "), " ")

  // Eliminate duplicate words and lookup word_ids
  words := make(map[string]int)
  for _, w := range word_list {
    if w != "" {
      words[w] = fv.WordId(w)
      if words[w] == 0 {
        err = errors.New("No files contain: '" + w + "'")
        return
      }
    }
  }

  // Construct sql query
  var query string
  for _, v := range words {
    if query == "" {
      query = "select f.file_id, concat(path, name) as filename from file_words inner join files f using(file_id) where word_id=" + strconv.Itoa(v)
    } else {
      query += " and file_id in (select file_id from file_words where word_id=" + strconv.Itoa(v)
    }
  }
  for i := 1; i < len(words); i++ {
    query += ")"
  }
  query += " order by file_id desc limit " + strconv.Itoa(fv.QueryLimit)

  // Execute query & fetch results
  var rows *sql.Rows
  rows, err = fv.db.Query(query)
  if err == nil {
    defer rows.Close()
    var file_id int
    var filename string
    for rows.Next() {
      rows.Scan(&file_id, &filename)
      file_ids = append(file_ids, file_id)
      filenames = append(filenames, filename) 
    }
  }

  if len(file_ids) == fv.QueryLimit {
    err = errors.New("Query results truncated at " + strconv.Itoa(fv.QueryLimit) + ".")
  }
  return
}

func (fv *FileVault) StoreFile(filename string, hash_id int) (err error) {
  hash_str := fmt.Sprintf("%010d", hash_id)
  path := fv.root + hash_str[:1] + "/" + hash_str[1:4] + "/" + hash_str[4:7] + "/"
  if err = os.MkdirAll(path, 0755); err != nil {
    return
  }
  err = CopyFile(filename, path + hash_str)  
  return
}

func (fv *FileVault) WordId(word string) (word_id int) {
  rows, err := fv.db.Query("select word_id from words where word=?", word)
  if err == nil {
    defer rows.Close()
    if rows.Next() {
      rows.Scan(&word_id)
    }
  }
  return
}
