
package main

import (
    //这里讲db作为go/databases的一个别名，表示数据库连接池
    db "VPN-Manage/src/database"
    . "VPN-Manage/src/router"
)
func main() {
    
    // 初始化mysql和redis
    db.InitMysql()
    defer db.QuitMysql()

    db.InitRedis()
    defer db.QuitRedis()

    router := InitRouter()
    router.Run(":8088")
}   
