package main

import (
    "crypto/rand"
    "encoding/base64"
    "fmt"
    "html/template"
    "log"
    "net/http"
    "net/url"
    "sync"
)

type Manager struct {
    cookieName  string     // private cookiename
    lock        sync.Mutex // protects session
    provider    Provider
    maxLifeTime int64
}

type Provider interface {
    SessionInit(sid string) (Session, error)
    SessionRead(sid string) (Session, error)
    SessionDestroy(sid string) error
    SessionGC(maxLifeTime int64)
}

type Session interface {
    Set(key, value interface{}) error // set session value
    Get(key interface{}) interface{}  // get session value
    Delete(key interface{}) error     // delete session value
    SessionID() string                // back current sessionID
}

var globalSessions *Manager
var provides = make(map[string]Provider)

func Register(name string, provider Provider) {
    if provider == nil {
        panic("session: Register provider is nil")
    }

    if _, dup := provides[name]; dup {
        panic("session: Register called twice for provider " + name)
    }
    provides[name] = provider
}

func NewManager(provideName, cookieName string, maxLifeTime int64) (*Manager, error) {
    provider, ok := provides[provideName]
    if !ok {
        return nil, fmt.Errorf("session: unknown provide %q (forgotten import?)", provideName)
    }
    return &Manager{provider: provider, cookieName: cookieName, maxLifeTime: maxLifeTime}, nil
}

func (manager *Manager) sessionId() string {
    b := make([]byte, 32)
    if _, err := rand.Read(b); err != nil {
        return ""
    }
    return base64.URLEncoding.EncodeToString(b)
}

func (manager *Manager) SessionStart(w http.ResponseWriter, r *http.Request) (session Session) {
    manager.lock.Lock()
    defer manager.lock.Unlock()
    cookie, err := r.Cookie(manager.cookieName)
    if err != nil || cookie.Value == "" {
        sid := manager.sessionId()
        session, _ = manager.provider.SessionInit(sid)
        cookie := http.Cookie{Name: manager.cookieName, Value: url.QueryEscape(sid), Path: "/", HttpOnly: true, MaxAge: int(manager.maxLifeTime)}
        http.SetCookie(w, &cookie)
    } else {
        sid, _ := url.QueryUnescape(cookie.Value)
        session, _ = manager.provider.SessionRead(sid)
    }
    return
}

func init() {
    globalSessions, _ = NewManager("memory", "gosessionid", 3600)
}

func login(w http.ResponseWriter, r *http.Request) {
    sess := globalSessions.SessionStart(w, r)
    r.ParseForm()
    if r.Method == "GET" {
        t, _ := template.ParseFiles("login.gtpl")
        w.Header().Set("Content-Type", "text/html")
        t.Execute(w, sess.Get("username"))
    } else {
        sess.Set("username", r.Form["username"])
        http.Redirect(w, r, "/", 302)
    }
}

func main() {
    http.HandleFunc("/login", login)
    // 设置监听的端口
    err := http.ListenAndServe(":9092", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
