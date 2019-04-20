package filevault

import (
    "crypto/sha256"
    "database/sql"
    "errors"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "regexp"
    "strconv"
    "strings"
    "time"
)

type FileInfo struct {
    FileID    int
    Path      string
    Name      string
    Size      int64
    Timestamp time.Time
    Hash      string
}

type FileVault struct {
    db         *sql.DB
    root       string
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
    fv := FileVault{db: db, root: root, QueryLimit: 0}
    return &fv
}

func (fv *FileVault) Check() (results []string, err error) {
    var rows *sql.Rows
    rows, err = fv.db.Query("select (select format(count(*), 0) from files) as names, format(count(*), 0) as files, format(sum(size), 0) as size from hashes")
    if err != nil {
        return
    }
    var names, files, size_str string
    if rows.Next() {
        rows.Scan(&names, &files, &size_str)
    }
    rows.Close()
    rows, err = fv.db.Query("select hash_id, size from hashes order by hash_id")
    if err != nil {
        return
    }
    defer rows.Close()
    var hash_id int
    var size int64
    errors := 0
    for rows.Next() {
        rows.Scan(&hash_id, &size)
        hash_str := fmt.Sprintf("%010d", hash_id)
        path := fv.root + hash_str[:1] + "/" + hash_str[1:4] + "/" + hash_str[4:7] + "/"
        fi, e := os.Stat(path + hash_str)
        if e != nil {
            results = append(results, e.Error())
            errors++
            r, e := fv.db.Query("select concat(path, name) as filename from files where hash_id=? order by 1", hash_id)
            if e == nil {
                var filename string
                for r.Next() {
                    r.Scan(&filename)
                    results = append(results, "  "+filename)
                }
            } else {
                fmt.Printf("%s\n", e.Error())
            }
        } else if fi.Size() != size {
            results = append(results, path+hash_str+": Size Mismatch")
            errors++
        }
    }
    results = append(results, "Names: "+names)
    results = append(results, "Files: "+files)
    results = append(results, "Size: "+size_str)
    results = append(results, "Errors: "+strconv.Itoa(errors))
    return
}

func (fv *FileVault) Extract(file_id int, filename string) (dest_filename string, err error) {
    var rows *sql.Rows
    rows, err = fv.db.Query("select hash_id, path, name from files where file_id=?", file_id)
    if err != nil {
        return
    }
    defer rows.Close()
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
        dest_filename = fname + fext
    } else {
        dest_filename = strings.Replace(filename, "{.path}", fpath, -1)
        dest_filename = strings.Replace(dest_filename, "{.name}", fname, -1)
        dest_filename = strings.Replace(dest_filename, "{.ext}", fext, -1)
    }
    hash_str := fmt.Sprintf("%010d", hash_id)
    path := fv.root + hash_str[:1] + "/" + hash_str[1:4] + "/" + hash_str[4:7] + "/"
    err = CopyFile(path+hash_str, dest_filename)
    return
}

func (fv *FileVault) FileId(filename string, hash string) (file_id int, err error) {
    fpath, fname := filepath.Split(filename)
    query := "select file_id from files inner join hashes using(hash_id) where hash=? and path=? and name=? order by timestamp desc, file_id desc limit 1"
    var rows *sql.Rows
    rows, err = fv.db.Query(query, hash, fpath, fname)
    if err == nil {
        defer rows.Close()
        if rows.Next() {
            rows.Scan(&file_id)
        }
    }
    return
}

func (fv *FileVault) Hash(filename string) (hash string, err error) {
    var f *os.File
    if f, err = os.Open(filename); err == nil {
        defer f.Close()
        h := sha256.New()
        if _, err = io.Copy(h, f); err == nil {
            hash = fmt.Sprintf("%x", h.Sum(nil))
        }
    }
    return
}

func (fv *FileVault) HashId(hash string, size int64) (hash_id int, err error) {
    hash_id = 0
    var rows *sql.Rows
    rows, err = fv.db.Query("select hash_id from hashes where hash=? and size=?", hash, size)
    if err == nil {
        defer rows.Close()
        if rows.Next() {
            rows.Scan(&hash_id)
        }
    }
    return
}

func (fv *FileVault) Info(file_id int) (fi FileInfo, err error) {
    var rows *sql.Rows
    rows, err = fv.db.Query("select file_id, path, name, timestamp, size, hash from files inner join hashes using(hash_id) where file_id=?", file_id)
    if err == nil {
        defer rows.Close()
        if rows.Next() {
            var timestamp string
            err = rows.Scan(&fi.FileID, &fi.Path, &fi.Name, &timestamp, &fi.Size, &fi.Hash)
            fi.Timestamp, err = time.Parse("2006-01-02 15:04:05", timestamp)
        } else {
            err = errors.New("Invalid file_id")
        }
    }
    return
}

func (fv *FileVault) Init() (err error) {
    err = errors.New("Not Supported, Yet.")
    return
}

func (fv *FileVault) Import(filename string, path string, timestamp time.Time) (file_id int, err error) {
    var reg *regexp.Regexp
    reg, err = regexp.Compile("[^ -~]+")
    if err != nil {
        return
    }
    path = reg.ReplaceAllString(path, "_")
    var hash string
    if hash, err = fv.Hash(filename); err != nil {
        return
    }
    var fi os.FileInfo
    fi, err = os.Stat(filename)
    if err != nil {
        return
    }
    var hash_id int
    if hash_id, err = fv.HashId(hash, fi.Size()); err != nil {
        return
    }
    if hash_id == 0 {
        r, e := fv.db.Query("select r.hash_id from rehash r left join hashes h using(hash_id) where h.hash_id is null order by 1 limit 1")
        if e == nil {
            defer r.Close()
            if r.Next() {
                r.Scan(&hash_id)
            }
        }
        if hash_id == 0 {
            if _, err = fv.db.Exec("insert into hashes(hash, size) values(?, ?)", hash, fi.Size()); err != nil {
                return
            }
        } else {
            if _, err = fv.db.Exec("insert into hashes(hash_id, hash, size) values(?, ?, ?)", hash_id, hash, fi.Size()); err != nil {
                return
            }
        }
        if hash_id, err = fv.HashId(hash, fi.Size()); err != nil {
            return
        }
    }
    if err = fv.StoreFile(filename, hash_id); err != nil {
        return
    }
    fpath, fname := filepath.Split(path)
    fv.db.Exec("insert into files(hash_id, path, name, timestamp) values(?,?,?,?)", hash_id, fpath, fname, timestamp)
    if file_id, err = fv.FileId(path, hash); err != nil {
        return
    }
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

func (fv *FileVault) ListPath(path string) (file_ids []int, names []string, err error) {
    query := "select file_id, name from files where path=? order by name, file_id desc"
    var rows *sql.Rows
    rows, err = fv.db.Query(query, path)
    if err == nil {
        defer rows.Close()
        var file_id int
        var name string
        for rows.Next() {
            rows.Scan(&file_id, &name)
            file_ids = append(file_ids, file_id)
            names = append(names, name)
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
    query += " order by file_id desc"
    if fv.QueryLimit != 0 {
      query += " limit " + strconv.Itoa(fv.QueryLimit)
    }

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

    if fv.QueryLimit != 0 {
      if len(file_ids) == fv.QueryLimit {
          err = errors.New("Query results truncated at " + strconv.Itoa(fv.QueryLimit) + ".")
      }
    }
    return
}

func (fv *FileVault) QueryFilename(filename string) (file_ids []int, err error) {
    query := "select file_id from files where path=? and name=?"
    var rows *sql.Rows
    fpath, fname := filepath.Split(filename)
    rows, err = fv.db.Query(query, fpath, fname)
    if err == nil {
        defer rows.Close()
        var file_id int
        for rows.Next() {
            rows.Scan(&file_id)
            file_ids = append(file_ids, file_id)
        }
    }
    return
}

func (fv *FileVault) SetQueryLimit(query_limit int) {
  fv.QueryLimit = query_limit
}

func (fv *FileVault) StoreFile(filename string, hash_id int) (err error) {
    hash_str := fmt.Sprintf("%010d", hash_id)
    path := fv.root + hash_str[:1] + "/" + hash_str[1:4] + "/" + hash_str[4:7] + "/"
    _, err = os.Stat(path + hash_str)
    if err == nil {
        return
    }
    if err = os.MkdirAll(path, 0755); err != nil {
        return
    }
    err = CopyFile(filename, path+hash_str)
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
