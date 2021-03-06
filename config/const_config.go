package config

import "github.com/gin-gonic/gin"

/*
常量级参数配置项，编译前设置值
*/
const (

	// 基本设置
	AppName    = "ginger"
	AppVersion = "v1.0.0"
	DefaultEnv = gin.DebugMode // 环境变量未设置时的默认值，默认使用开发环境配置

	// jwt编码时的私钥字符串
	TokenPrivateKey = "ginger"

	// oauth2 Key and Secret
	// 微信
	WechatSignSwitch = true
	WechatOAuth2AppKey = ""
	WechatOAuth2AppSecret = ""

	// qq
	QQSignSwitch = true
	QQOAuth2AppKey = ""
	QQOAuth2AppSecret = ""

	// 微博
	WeiboSignSwitch = true
	WeiboOAuth2AppKey = ""
	WeiboOAuth2AppSecret = ""
)
