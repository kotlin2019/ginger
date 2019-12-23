package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gofuncchan/ginger/cache"
	"github.com/gofuncchan/ginger/common"
	"github.com/gofuncchan/ginger/model/userModel"
	"github.com/gofuncchan/ginger/util/e"
	"github.com/gofuncchan/ginger/util/jwt"
	"github.com/gofuncchan/ginger/util/logger"
	"net/http"
	"strconv"
)

// TODO 邮箱认证码：注册时发送一封邮件到用户邮箱，缓存验证码，用户从邮箱点击链接时校验验证码，设置邮箱验证状态

/*
注册处理,默认邮箱注册
*/
type SignUpForm struct {
	Name       string `json:"name" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,alphanum"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

func SignUp(c *gin.Context) {
	// 请求注册页面
	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "register.tmpl", gin.H{
			"title": "User Register Page",
		})

	} else if c.Request.Method == "POST" {
		// 表单参数校验
		form := new(SignUpForm)
		// _ = c.Bind(form) //MustBindWith()
		if err := c.ShouldBind(form); err != nil {
			common.ResponseInvalidParam(c, err)
			return
		}

		// 检查邮箱账号是否已经存在
		if userModel.IsUserExistByEmail(form.Email) {
			err := errors.New("user email exist,please try again")
			logger.WarmLog(err.Error())
			common.ResponseClientError(c, err)
			return
		}

		// 密码哈希计算
		passHash, salt := common.GenPassHash(form.Password)
		// 创建用户
		id := userModel.CreateUserByEmail(form.Name, form.Email, passHash, salt)
		if id == -1 {
			err := errors.New("register error,please try again")
			logger.WarmLog(err.Error())
			common.ResponseServerError(c, err)
			return
		}

		// 返回方式一：返回注册结果
		common.ResponseOk(c,gin.H{"result": true})
		return

		// 返回方式二：创建登录token并返回
		// userClaim := jwt.TokenUserClaim{
		// 	Id:     id,
		// 	Name:   form.Name,
		// 	Avatar: "",
		// }
		// tkStr, err := jwt.JwtService.Encode(userClaim)
		// if !e.Eh(err) {
		// 	common.ResponseServerError(c, errors.New("jwt token string encode error"))
		// 	return
		// }
		// // Redis存储token保存登录状态
		// userKey := cache.UserTokenCacheKeyPrefix + strconv.Itoa(int(id))
		// cache.SetToken(userKey, tkStr)
		//
		// common.ResponseOk(c,
		// 	gin.H{
		// 		"token": tkStr,
		// 		"user_info": gin.H{
		// 			"id": id,
		// 			"name":    form.Name,
		// 			"avatar":  "",
		// 		}})
		// return

	} else {
		common.ResponseMethodNotAllowed(c, errors.New("only allow GET Or POST method"))
		return
	}
}

/*
登录处理
*/
type SignInForm struct {
	Email    string `json:"email" binding:"required,email"`
	PassWord string `json:"password" binding:"required,alphanum"`
}

func SignIn(c *gin.Context) {
	if c.Request.Method == "GET" {
		// 请求登录页面
		c.HTML(http.StatusOK, "login.tmpl", gin.H{
			"title": "User Login Page",
		})
	} else if c.Request.Method == "POST" {
		// 请求登录参数验证
		form := new(SignInForm)
		if err := c.ShouldBind(form); err != nil {
			common.ResponseInvalidParam(c, err)
			return
		}

		// 与数据库账号密码鉴权
		// 1.根据邮箱获取哈希密码与盐值
		userInfo := userModel.GetUserInfoByEmail(form.Email)
		if userInfo == nil {
			common.CommonResponse(c, common.ResponseCodeUnAuthorized, http.StatusForbidden, nil, errors.New("email Or password error,please try again"))
			return
		}
		// 2.将用户密码与盐值哈希计算并与数据库密码进行比较
		b := common.IsValidPasswd(form.PassWord, userInfo.Salt, userInfo.Password)
		if !b {
			common.CommonResponse(c, common.ResponseCodeUnAuthorized, http.StatusForbidden, nil, errors.New("email Or password error,please try again"))
			return
		}

		// 校验成功则生成token并返回
		// 创建登录token并返回
		userClaim := jwt.TokenUserClaim{
			Id:     int64(userInfo.ID),
			Name:   userInfo.Name,
			Avatar: userInfo.Avatar,
		}

		tkStr, err := jwt.JwtService.Encode(userClaim)
		if !e.Eh(err) {
			common.ResponseServerError(c, errors.New("jwt token string encode error"))
			return
		}
		// Redis存储token保存登录状态
		userKey := cache.UserTokenCacheKeyPrefix + strconv.Itoa(int(userInfo.ID))
		cache.SetToken(userKey, tkStr)

		common.ResponseOk(c,
			gin.H{
				"token": tkStr,
				"user_info": gin.H{
					"id": int64(userInfo.ID),
					"name":    userInfo.Name,
					"avatar":  userInfo.Avatar,
				}})
		return

	} else {
		common.ResponseMethodNotAllowed(c, errors.New("only allow GET Or POST method"))
		return
	}
}

/*
退出登录处理
*/
func SignOut(c *gin.Context) {
	// 从header获取token字段，在redis删除键
	tkStr := c.GetHeader("Authorization")
	// fmt.Println(tkStr)

	// 解码获取id
	claim, err := jwt.JwtService.Decode(tkStr)
	if err != nil {
		common.ResponseUnAuthorized(c, err)
		return
	}

	key := cache.UserTokenCacheKeyPrefix +  strconv.Itoa(int(claim.TokenUserClaim.Id))
	delCount := cache.DeleteToken(key)

	if delCount > 0 {
		common.ResponseOk(c, gin.H{"message": "Sign Out Successful!"})
	} else {
		common.ResponseServerError(c, errors.New("sign out error"))
	}
}
