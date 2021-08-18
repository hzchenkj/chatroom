package model

import (
	"chatroom/common/message"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

type UserDao struct {
	pool *redis.Pool
}

var (
	MyUserDao *UserDao
)

//获取userdao 实例

func NewUserDao(pool *redis.Pool) (userDao *UserDao) {
	userDao = &UserDao{
		pool: pool,
	}
	return
}

//根据用户id 返回一个user 或error
func (this *UserDao) getUserById(conn redis.Conn, userId int) (user User, err error) {
	res, err := redis.String(conn.Do("HGet", "users", userId))
	if err != nil {
		if err == redis.ErrNil {
			err = ERROR_USER_NOTEXISTS
			return
		}
	}
	//res 反序列化成User 实例

	err = json.Unmarshal([]byte(res), &user)
	if err != nil {
		fmt.Println("json.Unmarshal faild err :", err)
		return
	}
	return
}

func (this *UserDao) Register(user *message.User) (err error) {
	fmt.Println("UserDao>>Register>>")
	conn := this.pool.Get()
	defer conn.Close()
	fmt.Println("UserDao>>Register>>getUserById")
	_, err = this.getUserById(conn, user.UserId)
	fmt.Println("UserDao>>Register>>getUserById done")
	if err == nil {
		err = ERROR_USER_EXISTS
		return
	}
	fmt.Println(user.UserId, user.UserName, "用户不存在")
	//说明用户redis中没有。
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	//入库redis
	_, err = conn.Do("HSet", "users", user.UserId, string(data))
	if err != nil {
		fmt.Println("保存注册用户错误，err", err)
		return
	}
	return
}

//登录校验
func (this *UserDao) Login(userId int, userPwd string) (user User, err error) {
	conn := this.pool.Get()
	defer conn.Close()
	user, err = this.getUserById(conn, userId)
	if err != nil {
		return
	}

	//校验密码
	if user.UserPwd != userPwd {
		err = ERROR_USER_PWD
		return
	}
	return
}
