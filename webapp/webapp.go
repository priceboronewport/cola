package webapp

import (
  "../filestore"
  "encoding/base64"
  "crypto/sha256"
  "bytes"
  "database/sql"
  "errors"
  "fmt"
  "io"
  _ "github.com/go-sql-driver/mysql"
  "html/template"
  "io/ioutil"
  "log"
  "math/rand"
  "net/http"
  "net/url"
  "os"
  "path/filepath"
  "strings"
  "time"
)

type DataStore interface {
  Read(...string) string
  Write(string, string) error
}

type Record map[string]string

type HandlerParams struct {
  Session  string
  Instance string
  Username string
}

type HandlerFunc func(http.ResponseWriter, *http.Request, HandlerParams)

var DB *sql.DB
var sessions DataStore
var session_values DataStore
var handlers map[string]HandlerFunc
var require_auths map[string]bool
var Config DataStore
var data_path string

const NullUser = "_"

func ConfigFilename(config_path string) string {
  _, exe_filename := filepath.Split(os.Args[0])
  extension := filepath.Ext(os.Args[0])
  result := config_path + exe_filename + ".conf"
  if extension != "" {
    result = config_path + exe_filename[0:len(exe_filename)-len(extension)]
  }
  fmt.Printf("Configuration: %s\n", result)
  return result
}

func ContentType(filename string) string {
  mime_types := filestore.New(data_path + "mime_types.fs")
  ext := strings.ToLower(filepath.Ext(filename))
  mime_type := ""
  if ext != "" {
    mime_type = mime_types.Read(ext[1:])
  }
  if mime_type == "" {
    f, e := os.Open(filename)
    if e != nil {
      return ""
    }
    defer f.Close()
    buffer := make([]byte, 512)
    n, e := f.Read(buffer)
    if e != nil {
      return ""
    }
    mime_type = http.DetectContentType(buffer[:n])
  }
  return mime_type
}

func DocumentRoot(host string) string {
  result := Config.Read("document_root:" + host)
  if result == "" {
    result = Config.Read("document_root")
  }
  if result == "" {
    result = "."
  }
  if result[len(result)-1:] != "/" {
    return result + "/"
  }
  return result
}

func GetSession(w http.ResponseWriter, r *http.Request) string {
  session := ""
  sysid := Config.Read("sysid")
  c, _ := r.Cookie(sysid + "_session")
  if c != nil {
    session = c.Value
  }
  if len(session) > 0 {
    if sessions.Read(session) == "" {
      session = ""
    }
  }
  if len(session) == 0 {
    for {
      session = fmt.Sprintf("%d", rand.Uint32())
      if sessions.Read(session) == "" {
        break
      }
    }
    sessions.Write(session, NullUser)
    http.SetCookie(w, &http.Cookie{Name: sysid + "_session", Value: session, Path: "/", Expires: time.Now().AddDate(1, 0, 0)})
  }
  return session
}

func ListenAndServe(config_path string) {
  if len(require_auths) > 0 {
    if !HandlerExists("", "/login") {
      Register("", "/login", LoginHandler, false)
    }
    if !HandlerExists("", "/logout") {
      Register("", "/logout", LogoutHandler, false)
    }
  }
  rand.Seed(time.Now().UTC().UnixNano())
  Config = filestore.New(ConfigFilename(config_path))
  data_path = Config.Read("data_path", "./data/")
  sessions = filestore.New(data_path + "sessions.fs")
  session_values = filestore.New(data_path + "session_values.fs")
  if _, err := os.Stat(data_path + "sessions.fs"); os.IsNotExist(err) {
    panic(errors.New(data_path + "sessions.fs missing"))
  }
  db_type := Config.Read("db_type")
  if db_type != "" {
    fmt.Printf("Connecting to database...\n")
    var err error
    DB, err = sql.Open(db_type, Config.Read("db_connect"))
    if err != nil {
      panic(err)
    }
    err = DB.Ping()
    if err != nil {
      panic(err)
    }
    defer DB.Close()
  }
  fmt.Printf("Listening...\n")
  http.HandleFunc("/", Handler)
  tls_address := Config.Read("tls_address")
  if tls_address != "" {
    tls_cert := Config.Read("tls_cert")
    tls_key := Config.Read("tls_key")
    go func() {
      fmt.Printf("Listening TLS...\n")
      err := http.ListenAndServeTLS(tls_address, tls_cert, tls_key, nil)
      log.Fatal(err)
    }()
  }
  err := http.ListenAndServe(Config.Read("address"), nil)
  log.Fatal(err)
}

func Handler(w http.ResponseWriter, r *http.Request) {
  if r.TLS == nil && Config.Read("require_tls:"+r.Host, Config.Read("require_tls")) != "" {
    Redirect(w, r, "https://"+r.Host+r.URL.String())
    return
  }
  server_name := Config.Read(r.Host, r.Host)
  timestamp := time.Now().Format("2006-01-02 15:04:05")
  instance := fmt.Sprintf("%s %p", timestamp, &timestamp)
  session := GetSession(w, r)
  username := sessions.Read(session)
  if username == NullUser {
    username = ""
  }
  if handlers != nil {
    rurl := r.URL.String()
    var candidate_f HandlerFunc
    candidate_f = nil
    candidate_l := 0
    for domain_branch, f := range handlers {
      domain := ""
      branch := ""
      pos := strings.Index(domain_branch, ":")
      if pos > -1 {
        domain = domain_branch[:pos]
        branch = domain_branch[pos+1:]
      } else {
        branch = domain_branch
      }
      if (len(rurl) >= len(branch)) && (rurl[:len(branch)] == branch) && ((server_name == domain) || (domain == "")) {
        if len(branch) > candidate_l {
          candidate_f = f
          candidate_l = len(branch)
        }
      }
    }
    log.Printf("[%s] [%s] [%s] [%s]\n", username, r.RemoteAddr, server_name, rurl)
    if candidate_f != nil {
      if require_auths[fmt.Sprintf("%p", candidate_f)] && username == "" {
        Redirect(w, r, "/login?return="+url.QueryEscape(rurl))
      }
      candidate_f(w, r, HandlerParams{Session: session, Instance: instance, Username: username})
      return
    } else {
      w.WriteHeader(http.StatusNotFound)
      fmt.Fprint(w, "Not Found")
      log.Printf("webapp.Handler: 404 - Not Found: [%s]\n", r.URL.String())
      return
    }
  }
}

func HandlerExists(query_domain string, query_branch string) bool {
  for domain_branch, _ := range handlers {
    domain := ""
    branch := ""
    pos := strings.Index(domain_branch, ":")
    if pos > -1 {
      domain = domain_branch[:pos]
      branch = domain_branch[pos+1:]
    } else {
      branch = domain_branch
    }
    if domain == query_domain && branch == query_branch {
      return true
    }
  }
  return false
}

func IfEmpty(str string, defstr string) string {
  if str == "" {
    return defstr
  }
  return str
}

type LoginParams struct {
  ReturnURL    string
  ErrorMessage string
}

func LoginHandler(w http.ResponseWriter, r *http.Request, p HandlerParams) {
  r.ParseForm()
  action := r.Form.Get("action")
  if action == "Login" {
    username := strings.ToLower(r.Form.Get("username"))
    password := r.Form.Get("password")
    return_url := IfEmpty(r.Form.Get("return_url"), "/")
    passwords := filestore.New(data_path + "/passwords.fs")
    password_rec := strings.Split(passwords.Read(username), ",")
    var hash string
    if len(password_rec) == 1 && password == password_rec[0] {
      salt := fmt.Sprintf("%d", rand.Uint32())
      h := sha256.New()
      io.WriteString(h, salt+password)
      hash = base64.StdEncoding.EncodeToString(h.Sum(nil))
      passwords.Write(username, hash+","+salt)
      password_rec[0] = hash
    } else {
      if len(password_rec) >= 2 {
        salt := password_rec[1]
        h := sha256.New()
        io.WriteString(h, salt+password)
        hash = base64.StdEncoding.EncodeToString(h.Sum(nil))
      }
    }
    if password_rec[0] == hash && hash != "" {
      sessions.Write(p.Session, username)
      Redirect(w, r, return_url)
      log.Printf("[%s] Login\n", username)
      return
    }
    log.Printf("[%s] Login Failed\n", username)
    Render(w, "login.html", LoginParams{ReturnURL: return_url, ErrorMessage: "Login Failed"})
  } else {
    return_url := IfEmpty(r.URL.Query().Get("return"), "/")
    Render(w, "login.html", LoginParams{ReturnURL: return_url, ErrorMessage: ""})
  }
}

func LogoutHandler(w http.ResponseWriter, r *http.Request, p HandlerParams) {
  sessions.Write(p.Session, NullUser)
  return_url := r.URL.Query().Get("return")
  if return_url == "" {
    return_url = "/"
  }
  Redirect(w, r, return_url)
}

func Trunc(s string, d string) string {
  if idx := strings.Index(s, d); idx != -1 {
    return s[:idx]
  }
  return s
}

func Redirect(w http.ResponseWriter, r *http.Request, url string) {
  http.Redirect(w, r, url, http.StatusFound)
}

func Register(domain string, branch string, f HandlerFunc, require_auth bool) {
  if handlers == nil {
    handlers = make(map[string]HandlerFunc)
    require_auths = make(map[string]bool)
  }
  if domain != "" {
    handlers[domain+":"+branch] = f
  } else {
    handlers[branch] = f
  }
  require_auths[fmt.Sprintf("%p", f)] = require_auth
}

func Render(w http.ResponseWriter, template_filename string, render_params interface{}) {
  t, err := template.ParseFiles(Config.Read("template_path", "./templates/") + template_filename)
  if t != nil {
    err = t.Execute(w, render_params)
    if err != nil {
      log.Printf("webapp.Render.template.Execute: %s\n", err.Error())
    }
  } else {
    log.Printf("webapp.Render.template.ParseFiles: %s\n", err.Error())
  }
}

func Script(url string) string {
  return "<script type='text/javascript' src='" + url + "'></script>"
}

func SessionValuesRead(p HandlerParams, key string) string {
  return session_values.Read(p.Session + "_" + key)
}

func SessionValuesWrite(p HandlerParams, key string, value string) (err error) {
  session_values.Write(p.Session+"_"+key, value)
  return nil
}

func StaticHandler(w http.ResponseWriter, r *http.Request, p HandlerParams) {
  document_root := DocumentRoot(r.Host)
  filename := Trunc(r.URL.String()[1:], "?")
  if filename == "" {
    filename = document_root + "index.html"
  } else {
    filename = document_root + filename
  }
  mime_type := ContentType(filename)
  f, err := ioutil.ReadFile(filename)
  if (err != nil) || (mime_type == "") {
    w.WriteHeader(http.StatusNotFound)
    fmt.Fprint(w, "Not Found")
    log.Printf("webapp.StaticHandler: 404 - Not Found: [%s]\n", filename)
    return
  }
  b := bytes.NewBuffer(f)
  w.Header().Set("Content-type", mime_type)
  b.WriteTo(w)
}

func Stylesheet(url string) string {
  return "<link href='" + url + "' type='text/css' rel='stylesheet'/>"
}

func UrlPath(r *http.Request, index int) string {
  url_path := strings.Split(strings.Split(r.URL.String()[1:], "?")[0], "/")
  if len(url_path) > index {
    result, _ := url.QueryUnescape(url_path[index])
    return result
  }
  return ""
}

func User(username string) Record {
  users := filestore.New(data_path + "/users.fs")
  result := make(Record)
  user_rec := strings.Split(users.Read(username), ",")
  for i, v := range user_rec {
    switch i {
      case 0:
        result["first_name"] = v
      case 1:
        result["last_name"] = v
      case 2:
        result["email"] = v
    } 
  }
  return result
}
