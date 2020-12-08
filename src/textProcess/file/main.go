package main

import (
    "fmt"
    "os"
)

// 创建目录
func mkdir()  {
    os.Mkdir("text", 0777)
    os.MkdirAll("text/1/2", 0777)

    err := os.Remove("text")
    if err != nil {
        fmt.Println(err)
    }

    os.RemoveAll("text")
}

// 创建文件
func mkFile()  {
    userFile := "text.txt"
    fout, err := os.Create(userFile)
    if err != nil {
        fmt.Println(userFile, err)
        return
    }

    defer fout.Close()
    for i := 0; i < 10; i++ {
        fout.WriteString("Just a test!\r\n")
        fout.Write([]byte("Just a test!\n"))
    }
}

// 读取文件
func readFile()  {
    userFile := "text.txt"
    fl, err := os.Open(userFile)
    if err != nil {
        fmt.Println(userFile, err)
        return
    }

    defer fl.Close()
    buf := make([]byte, 1024)
    for {
        n, _ := fl.Read(buf)
        if 0 == n {
            break
        }
        os.Stdout.Write(buf[:n])
    }
}

func main()  {
    mkdir()

    mkFile()

    readFile()

    os.Remove("text.txt")
}