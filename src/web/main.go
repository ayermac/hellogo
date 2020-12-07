package main

import (
    "crypto/md5"
    "fmt"
    "html/template"
    "io"
    "net/http"
    "os"
    "strconv"
    "strings"
    "log"
    "time"
)

func sayhelloName(w http.ResponseWriter, r *http.Request) {
    // 解析参数，默认是不会解析的
    r.ParseForm()
    // 这些信息是输出到服务器端的打印信息
    fmt.Println(r.Form)
    fmt.Println("path", r.URL.Path)
    fmt.Println("scheme", r.URL.Scheme)
    fmt.Println(r.Form["url_long"])
    for k, v := range r.Form {
        fmt.Println("key:", k)
        fmt.Println("val:", strings.Join(v, ""))
    }
    // 这个写入到 w 的是输出到客户端的
    fmt.Fprintf(w, "Hello Jason!")
}

func login(w http.ResponseWriter, r *http.Request) {
    // 获取请求的方法
    fmt.Println("method:", r.Method)
    if r.Method == "GET" {
        crutime := time.Now().Unix()
        h := md5.New()
        io.WriteString(h, strconv.FormatInt(crutime, 10))
        token := fmt.Sprintf("%x", h.Sum(nil))

        t, _ := template.ParseFiles("view/login.gtpl")
        log.Println(t.Execute(w, token))
    } else {
        // 解析 url 传递的参数，对于 POST 则解析响应包的主体（request body）
        r.ParseForm()
        token := r.Form.Get("token")
        if token != "" {
            // 验证 token 合法性
        } else {
            // 不存在 token 报错
        }
        // 请求的是登录数据，那么执行登录的逻辑判断
        fmt.Println("username length:", len(r.Form["username"][0]))
        // 输出到服务器端
        fmt.Println("username:", template.HTMLEscapeString(r.Form.Get("username")))
        fmt.Println("password:", template.HTMLEscapeString(r.Form.Get("password")))
        // 输出到客户端
        template.HTMLEscape(w, []byte(r.Form.Get("username")))
    }
}

func upload (w http.ResponseWriter, r *http.Request) {
    fmt.Println("method:", r.Method)
    if r.Method == "GET" {
        crutime := time.Now().Unix()
        h := md5.New()
        io.WriteString(h, strconv.FormatInt(crutime, 10))
        token := fmt.Sprintf("%x", h.Sum(nil))

        t, _ := template.ParseFiles("view/upload.gtpl")
        log.Println(t.Execute(w, token))
    } else {
        r.ParseMultipartForm(32 << 20)
        file, handler, err := r.FormFile("uploadfile")
        if err != nil {
            fmt.Println(err)
            return
        }

        defer file.Close()
        fmt.Fprintf(w, "%v", handler.Header)
        f, err := os.OpenFile("upload/" + handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
        if err != nil {
            fmt.Println(err)
            return
        }

        defer f.Close()
        io.Copy(f, file)
        return
    }
}

func main() {
    // 设置访问的路由
    http.HandleFunc("/", sayhelloName)
    // 设置访问的路由
    http.HandleFunc("/login", login)
    // 设置访问的路由
    http.HandleFunc("/upload", upload)
    // 设置监听的端口
    err := http.ListenAndServe(":9090", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}