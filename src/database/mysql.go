package database
import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"log"
)
//因为我们需要在其他地方使用Engine这个变量，所以需要大写代表public
var Engine *xorm.Engine
//初始化方法，每个package 中的init()方法会被自动加载
func InitMysql() {
    var err error
    if err != nil {
    	log.Fatal("读取配置文件失败")
	}
	Engine, err = xorm.NewEngine("mysql", "root:liang.0804@tcp(127.0.0.1:3306)/VPN?charset=utf8")
	if err != nil{
		log.Fatal("数据库连接失败！！",err)
	}
    //连接检测
    err = Engine.Ping()
    if err != nil {
        log.Fatal(err.Error())
    }
}

func QuitMysql(){
	Engine.Close()
}