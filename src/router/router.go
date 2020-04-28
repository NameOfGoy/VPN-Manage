package router
import (
    . "VPN-Manage/src/controllers"
    "VPN-Manage/src/util/cors"
    "VPN-Manage/src/util/jwt"
    "github.com/gin-gonic/gin"
)
func InitRouter() *gin.Engine {
    router := gin.Default()
    
    // 使用跨域中间件，解决复杂跨域问题。浏览器会发两次请求，第一次是option
    router.Use(cors.Cors())

    router.GET("/", Index)
    router.POST("/wxregister",UserRegisterWithWx)
    //router.POST("/login",UserLogin)
    router.POST("/wxlogin",UserWxLogin)
    //router.GET("/exitLogin",ExitLogin)
    router.GET("/getuser",GetUser)
    //router.POST("/modifypassword",ModifyPassword)
   
    taR := router.Group("/data")
    // 路径是/data时，会把jwt.JWTAuth()调用
    taR.Use(jwt.JWTAuth()) // router.Use是使用中间件的意思

    {

    }

    return router
}