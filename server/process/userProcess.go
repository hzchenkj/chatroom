package process2

import (
	"chatroom/common/message"
	"chatroom/server/model"
	"chatroom/server/util"
	"encoding/json"
	"fmt"
	"net"
)

type UserProcess struct {
	Conn net.Conn
}

func (this *UserProcess) ServerProcessRegister(msg *message.Message) (err error) {
	fmt.Println(">>>userProcess.go ServerProcessRegister")
	var registerMsg message.RegisterMsg
	err = json.Unmarshal([]byte(msg.Data), &registerMsg)
	if err != nil {
		fmt.Println("json.Unmarshal fail err=", err)
		return
	}

	//
	var resMsg message.Message
	resMsg.Type = message.RegisterResultMsgType

	var registerResutlMsg message.RegisterResultMsg

	err = model.MyUserDao.Register(&registerMsg.User)

	if err != nil {
		if err == model.ERROR_USER_EXISTS {
			registerResutlMsg.Code = 505
			registerResutlMsg.Error = model.ERROR_USER_EXISTS.Error()
		} else {
			registerResutlMsg.Code = 506
			registerResutlMsg.Error = "注册发送未知错误。。。"
		}
	} else {
		registerResutlMsg.Code = 200
	}

	data, err := json.Marshal(registerResutlMsg)
	if err != nil {
		fmt.Println("json.Marshal faild err:", err)
		return
	}

	resMsg.Data = string(data)

	data, err = json.Marshal(resMsg)
	if err != nil {
		fmt.Println("json.Marshal faild err:", err)
		return
	}
	//发送到客户端
	//创建transer实例

	tf := &util.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)

	return
}

func (this *UserProcess) ServerProcessLogin(msg *message.Message) (err error) {
	//msg 中取出data msg.Data 发序列化层LoginMsg
	var loginMsg message.LoginMsg
	err = json.Unmarshal([]byte(msg.Data), &loginMsg)
	if err != nil {
		fmt.Println("json.Unmarshal fail err=", err)
		return
	}

	//
	var resMsg message.Message
	resMsg.Type = message.LoginResultMsgType

	var loginResutlMsg message.LoginResultMsg

	user, err := model.MyUserDao.Login(loginMsg.UserId, loginMsg.UserPwd)

	if err != nil {

		if err == model.ERROR_USER_NOTEXISTS {
			loginResutlMsg.Code = 500
			loginResutlMsg.Error = err.Error()
		} else if err == model.ERROR_USER_PWD {
			loginResutlMsg.Code = 403
			loginResutlMsg.Error = err.Error()
		} else {
			loginResutlMsg.Code = 505
			loginResutlMsg.Error = "服务器内部错误"
		}
	} else {
		loginResutlMsg.Code = 200
		fmt.Println(user, "登录成功")
	}

	////如果用户id是100 密码是123456
	//if loginMsg.UserId == 100 && loginMsg.UserPwd == "123456" {
	//	//合法
	//	loginResutlMsg.Code = 200
	//} else {
	//	//不合法
	//	loginResutlMsg.Code = 500
	//	loginResutlMsg.Error = "用户不存在"
	//}
	data, err := json.Marshal(loginResutlMsg)
	if err != nil {
		fmt.Println("json.Marshal faild err:", err)
		return
	}
	resMsg.Data = string(data)

	data, err = json.Marshal(resMsg)
	if err != nil {
		fmt.Println("json.Marshal faild err:", err)
		return
	}
	//发送到客户端
	//创建transer实例

	tf := &util.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)

	return err
}
