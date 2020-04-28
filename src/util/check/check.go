package check

import(
	"github.com/gin-gonic/gin"
	myjwt "VPN-Manage/src/util/jwt"
	e "VPN-Manage/src/util/myerror"
	"github.com/garyburd/redigo/redis"
	db "VPN-Manage/src/database"
)

func CheckToken(c *gin.Context) (claims *myjwt.CustomClaims,err error){
	//后台校验token合法性
	token := c.Request.Header.Get("token")
	if token == "" {
		err = &e.Myerror{
			Code : e.ERROR_TOKEN_NOTFOUND_CODE,
			Message : e.ERROR_TOKEN_NOTFOUND_MESSAGE,
		}
		return 
	}
	j := myjwt.NewJWT()
    // parseToken 解析token包含的信息
	claims, err = j.ParseToken(token)
	if err != nil {
		if err == myjwt.TokenExpired {
			err = &e.Myerror{
				Code : e.ERROR_TOKEN_TIMEOUT_CODE,
				Message : e.ERROR_TOKEN_TIMEOUT_MESSAGE,
			}
			return	// 返回过期错误
		}
		err = &e.Myerror{
			Code : e.ERROR_TOKEN_ILLEGEL_CODE,
			Message : err.Error(),
		}
		return  //返回其他token验证的错误
	}
	//redis校验token是否存在
	redistoken,err := redis.String(db.Conn.Do("GET", claims.ID))
	if err != nil || redistoken != token{
		err = &e.Myerror{
			Code : e.ERROR_TOKEN_TIMEOUT_CODE,
			Message : e.ERROR_TOKEN_TIMEOUT_MESSAGE,
		}
		return
	}
	return  
}