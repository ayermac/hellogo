package main

import (
    "fmt"
    "log"
    "net/http"
    "time"
)

func index(w http.ResponseWriter, r *http.Request) {
    expiration := time.Now()
    expiration = expiration.AddDate(1, 0, 0)
    cookie := http.Cookie{Name: "username", Value: "Jason", Expires: expiration}
    http.SetCookie(w, &cookie)

    username, _ := r.Cookie("username")
    fmt.Fprint(w, username)

    for _, cookies := range r.Cookies() {
        fmt.Fprint(w, cookies.Name)
    }
}

func main() {
    // 设置访问的路由
    http.HandleFunc("/", index)
    // 设置监听的端口
    err := http.ListenAndServe(":9091", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
