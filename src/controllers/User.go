package controllers

import (
	db "VPN-Manage/src/database"
	"VPN-Manage/src/models"
	"VPN-Manage/src/util/check"
	myjwt "VPN-Manage/src/util/jwt"
	e "VPN-Manage/src/util/myerror"
	"fmt"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

func Index(c *gin.Context) {
	c.String(http.StatusOK, " 兄 dei ，域名地址或者端口是不是写错啦？ 这不是你想要的接口哦。别老link过来")
}

func UserRegisterWithWx(c *gin.Context) {
	registinfo := struct {
		WXuuid   string `json:"wxuuid"`
		Loginid  string `json:"loginid"`
		Password string `json:"password"`
		Username string `json:"username"`
	}{}
	err := c.ShouldBind(&registinfo) // 绑定json
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"response": &e.Myerror{
				Code:    e.ERROR_JSONBIND_FAILED_CODE,
				Message: e.ERROR_JSONBIND_FAILED_MESSAGE,
			},
		})
		return
	}

	user := models.User{
		Uuid:     strconv.FormatInt(time.Now().Unix(), 10), // 暂时用时间戳表示uuid
		Loginid:  registinfo.Loginid,
		Password: registinfo.Password,
		Username: registinfo.Username,
		Isbindwx: "1",
	}

	// 先查询账号是否存在, nil则表示GetUser方法执行没出错，找得到用户
	_, err = user.GetUser()
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"response": &e.Myerror{
				Code:    e.ERROR_USER_HAVEEXISTED_CODE,
				Message: e.ERROR_USER_HAVEEXISTED_MESSAGE,
			},
		})

		return
	}
	log.Println(err)

	/************** 此处插入多表，暂不做原子处理 *****************/
	err = user.AddUser()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"response": &e.Myerror{
				Code:    e.ERROR_INSERT_TABLE_CODE,
				Message: e.ERROR_INSERT_TABLE_MESSAGE,
			},
		})
		return
	}

	uwm := models.UserWxMapping{
		Loginid: registinfo.Loginid,
		Wxuuid:  registinfo.WXuuid,
	}
	err = uwm.Add()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"response": &e.Myerror{
				Code:    e.ERROR_INSERT_TABLE_CODE,
				Message: e.ERROR_INSERT_TABLE_MESSAGE,
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": &models.ResponseSet{
			Code:    1,
			Message: "注册成功",
		},
	})
	return
}

func UserLogin(c *gin.Context) {
	logininfo := struct {
		Loginid  string `json:"loginid"`
		Password string `json:"password"`
	}{}

	// 绑定json
	err := c.ShouldBind(&logininfo)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"response": &e.Myerror{
				Code:    e.ERROR_JSONBIND_FAILED_CODE,
				Message: e.ERROR_JSONBIND_FAILED_MESSAGE,
			},
		})
		return
	}

	user := models.User{
		Loginid:  logininfo.Loginid,
		Password: logininfo.Password,
	}
	fmt.Println(user)
	// 查找数据库
	user2, err := user.GetUser()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"response": &err,
		})
		return
	}
	// 对比密码
	if user.Password != user2.Password {
		c.JSON(http.StatusOK, gin.H{
			"response": &e.Myerror{
				Code:    e.ERROR_USER_WRONGPASSWORD_CODE,
				Message: e.ERROR_USER_WRONGPASSWORD_MESSAGE,
			},
		})
		return
	}

	//密码正确，生成token
	token, err := generateToken(c, user.Loginid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"response": &err,
		})
		return
	}

	//存token进redis
	_, err = db.Conn.Do("SET", user.Loginid, token, "EX", "3600") // key:登录名  value : token
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"response": &e.Myerror{
				Code:    e.ERROR_INSERT_REDIS_KEY_CODE,
				Message: e.ERROR_INSERT_REDIS_KEY_MESSAGE,
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": &models.ResponseSet{
			Code:    1,
			Message: "登录成功",
			Data:    token,
		},
	})
	return
}

func UserWxLogin(c *gin.Context) {
	logininfo := struct {
		Wxuuid string `json:"wxuuid"`
	}{}

	// 绑定json
	err := c.ShouldBind(&logininfo)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"response": &e.Myerror{
				Code:    e.ERROR_JSONBIND_FAILED_CODE,
				Message: e.ERROR_JSONBIND_FAILED_MESSAGE,
			},
		})
		return
	}

	// 根据wxuuid查找映射表，如果有对应的用户数据，则登录成功，没有则提示没注册或者绑定。前端跳转到用户登录界面
	userwxmapping := models.UserWxMapping{
		Wxuuid: logininfo.Wxuuid,
	}
	// 查找数据库
	uwm, err := userwxmapping.Get()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"response": &e.Myerror{
				Code:    e.ERROR_USER_NOTFOUND_CODE,
				Message: e.ERROR_USER_NOTFOUND_MESSAGE,
			},
		})
		return
	}

	//密码正确，生成token
	token, err := generateToken(c, uwm.Loginid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"response": &err,
		})
		return
	}

	//存token进redis
	_, err = db.Conn.Do("SET", uwm.Loginid, token, "EX", "3600") // key:登录名  value : token
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"response": &e.Myerror{
				Code:    e.ERROR_INSERT_REDIS_KEY_CODE,
				Message: e.ERROR_INSERT_REDIS_KEY_MESSAGE,
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": &models.ResponseSet{
			Code:    1,
			Message: "登录成功",
			Data:    token,
		},
	})
	return
}

func ExitLogin(c *gin.Context) {
	//检查token是否合法
	claims, err := check.CheckToken(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"response": &err,
		})
		return
	}
	//删除redis的token缓存
	_, err = db.Conn.Do("DEL", claims.ID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"response": &e.Myerror{
				Code:    e.ERROR_DELETE_REDIS_KEY_CODE,
				Message: e.ERROR_DELETE_REDIS_KEY_MESSAGE,
			},
		})
		return
	}
	//返回
	c.JSON(http.StatusOK, gin.H{
		"response": &models.ResponseSet{
			Code:    1,
			Message: "已退出登录",
		},
	})
}

//生成令牌
func generateToken(c *gin.Context, loginid string) (token string, err error) {
	j := &myjwt.JWT{
		SigningKey: []byte("newtrekWang"),
	}
	StandardClaims := myjwt.CustomClaims{
		ID: loginid,
		StandardClaims: jwtgo.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 100),  // 签名生效时间
			ExpiresAt: int64(time.Now().Unix() + 3600), // 过期时间 一小时
			Issuer:    "goy",                           //签名的发行者
			Audience:  loginid,
		},
	}

	token, err = j.CreateToken(StandardClaims) // 取得token
	if err != nil {
		err = &e.Myerror{
			Code:    e.ERROR_TOKEN_CREATE_CODE,
			Message: e.ERROR_TOKEN_CREATE_MESSAGE,
		}
		return
	}
	return
}

func ModifyPassword(c *gin.Context) {
	//检查token是否合法
	claims, err := check.CheckToken(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"response": &err,
		})
		return
	}

	oldnewpassword := struct {
		Oldpassword string `json:"oldpassword"`
		Newpassword string `json:"newpassword"`
	}{}
	err = c.ShouldBind(&oldnewpassword)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"response": &e.Myerror{
				Code:    e.ERROR_JSONBIND_FAILED_CODE,
				Message: e.ERROR_JSONBIND_FAILED_MESSAGE,
			},
		})
		return
	}
	// 从数据库中取出用户密码，并且对比，正确就更新旧密码
	user := models.User{
		Loginid: claims.ID,
	}
	user, err = user.GetUser()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"response": &err,
		})
		return
	}
	log.Println(user)
	if user.Password != oldnewpassword.Oldpassword {
		c.JSON(http.StatusOK, gin.H{
			"response": &e.Myerror{
				Code:    e.ERROR_USER_WRONGPASSWORD_CODE,
				Message: e.ERROR_USER_WRONGPASSWORD_MESSAGE,
			},
		})
		return
	}
	user.Password = oldnewpassword.Newpassword
	err = user.UpdateUser() //更新
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"response": &err,
		})
		return
	}
	// 成功
	c.JSON(http.StatusOK, gin.H{
		"response": &models.ResponseSet{
			Code:    1,
			Message: "密码修改成功",
		},
	})
}

//
func GetUser(c *gin.Context) {
	user := models.User{
		Loginid: c.Query("loginid"), // 取得url后面附带的参数loginid
	}
	u, err := user.GetUser() // 调用查询方法
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"response": &err,
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"response": &e.Myerror{
				Code:    e.ERROR_STRUCT_TO_STRING_CODE,
				Message: e.ERROR_STRUCT_TO_STRING_MESSAGE,
			},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"response": &models.ResponseSet{
			Code:    1,
			Message: "",
			Data:    u,
		},
	})
}


