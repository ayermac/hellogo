package main

import (
    "database/sql"
    "fmt"
    "github.com/astaxie/beego/orm"
    _ "github.com/go-sql-driver/mysql"
)

type User struct {
    Id int
    Name string `orm:"size(100)"`
}

func init() {
    // 注册驱动
    orm.RegisterDriver("mysql", orm.DR_MySQL)
    // 设置默认数据库
    orm.RegisterDataBase("default", "mysql", "dev_web:dev_web@tcp(127.0.0.1:3307)/dev_web?charset=utf8", 30)
    // 注册定义的 model
    orm.RegisterModel(new(User))

    // 设置数据库的最大空闲连接
    orm.SetMaxIdleConns("default", 30)

    // 设置数据库的最大数据库连接 (go>= 1.2)
    orm.SetMaxOpenConns("default", 30)

    orm.Debug = true

    // 创建 table
    orm.RunSyncdb("default", false, true)
}

func main()  {
    o := orm.NewOrm

    user := User{Name: "Jason"}

    // 插入表
    id, err := o.Insert(&user)
    fmt.Printf("ID: %d, ERR: %v\n", id, err)

    // 更新表
    user.Name = "JasonUpdate"
    num, err = o.Update(&user)
    fmt.Printf("ID: %d, ERR: %v\n", num, err)

    // 读取 one
    u := User{Id: user.Id}
    err = o.Read(&u)
    fmt.Printf("ERR: %v\n", err)

    // 删除表
    num, err = o.Delete(&u)
    fmt.Printf("NUM: %d, ERR: %v\n", num, err)
}