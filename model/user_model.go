package model

import (
	"github.com/gofuncchan/ginger/dao/mysql/user"
	"github.com/gofuncchan/ginger/util/e"
	"time"
)

// 邮箱注册创建用户
func CreateUserByEmail(name, email, passwd, salt string) int64 {
	var data []map[string]interface{}
	data = append(data, map[string]interface{}{
		"name":      name,
		"email":     email,
		"password":  passwd,
		"salt":      salt,
		"create_at": time.Now(),
		"update_at": time.Now(),
	})

	id, err := user.Insert(data)
	if !e.Em(err) {
		return -1
	}
	return id
}

// 手机注册创建用户
func CreateUserByPhone(name, phone, passwd, salt string) int64 {
	var data []map[string]interface{}
	data = append(data, map[string]interface{}{
		"name":   name,
		"phone":  phone,
		"passwd": passwd,
		"salt":   salt,
	})

	id, err := user.Insert(data)
	if !e.Em(err) {
		return -1
	}
	return id
}

// 更新用户信息
func UpdateUserInfo(id int, data map[string]interface{}) bool {
	where := map[string]interface{}{
		"id": uint(id),
	}

	_, err := user.Update(where, data)

	return e.Em(err)
}

// 根据邮箱和密码验证用户登录
func GetUserInfoByEmail(email string) *user.User {
	where := map[string]interface{}{
		"email": email,
	}

	userInfo, err := user.GetOne(where)
	if !e.Em(err) {
		return nil
	}
	return userInfo
}

// 根据手机和密码验证用户登录
func GetUserInfoByPhone(phone string) *user.User {
	where := map[string]interface{}{
		"phone": phone,
	}

	userInfo, err := user.GetOne(where)
	if !e.Em(err) {
		return nil
	}
	return userInfo
}

// 根据邮箱验证用户是否已存在
func IsUserExistByEmail(email string) bool {
	where := map[string]interface{}{
		"email": email,
	}

	count, err := user.GetCount(where)
	if !e.Em(err) {
		return false
	}

	return count == 1
}

// 根据手机验证用户是否已存在
func IsUserExistByPhone(phone string) bool {
	where := map[string]interface{}{
		"phone": phone,
	}

	count, err := user.GetCount(where)
	if !e.Em(err) {
		return false
	}
	return count == 1
}

// 根据用户昵称验证是否已存在
func IsUserExistByName(name string) bool {
	where := map[string]interface{}{
		"name": name,
	}

	count, err := user.GetCount(where)
	if !e.Em(err) {
		return false
	}
	return count == 1
}