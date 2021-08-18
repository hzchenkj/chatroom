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
		loginResutlMsg.Code = 500
		loginResutlMsg.Error = "该用户不存在"

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
