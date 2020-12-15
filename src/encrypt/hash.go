package main

import (
    "crypto/md5"
    "crypto/sha1"
    "crypto/sha256"
    "encoding/base64"
    "fmt"
    "golang.org/x/crypto/scrypt"
    "io"
)

func main()  {
    salt := []byte("xxx")

    // import "crypto/sha256"
    h := sha256.New()
    io.WriteString(h, "His money is twice tainted: 'taint yours and 'taint mine.")
    fmt.Printf("% x", h.Sum(nil))
    fmt.Println("")

    // import "crypto/sha1"
    s := sha1.New()
    io.WriteString(s, "His money is twice tainted: 'taint yours and 'taint mine.")
    fmt.Printf("% x", s.Sum(nil))
    fmt.Println("")

    // import "crypto/md5"
    m := md5.New()
    io.WriteString(m, "需要加密的密码")
    fmt.Printf("%x", m.Sum(nil))
    fmt.Println("")

    dk, err := scrypt.Key([]byte("some password"), salt, 16384, 8, 1, 32)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(base64.StdEncoding.EncodeToString(dk))
}
