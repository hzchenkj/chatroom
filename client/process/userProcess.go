package process

import (
	"chatroom/common/message"
	"chatroom/server/util"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type UserProcess struct {
}

func (this *UserProcess) Register(userId int, userPwd string, userName string) (err error) {
	fmt.Printf("userId = %d userPwd = %s \n", userId, userPwd)
	// 连接到服务器
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("client dail net err=", err)
		return
	}
	defer conn.Close()

	//准备通过conn 发送消息给服务器
	var msg message.Message
	msg.Type = message.RegisterMsgType
	// 创建Register msg
	var registerMsg message.RegisterMsg
	registerMsg.User.UserId = userId
	registerMsg.User.UserPwd = userPwd
	registerMsg.User.UserName = userName
	//4  学历恶化
	data, err := json.Marshal(registerMsg)
	if err != nil {
		fmt.Println("json Marshal err: ", err)
		return
	}

	msg.Data = string(data)
	//msg 序列化
	data, err = json.Marshal(msg)
	if err != nil {
		fmt.Println("json Marshal err:", err)
		return
	}

	fmt.Println("客户端发送的内容:", string(data))

	tf := &util.Transfer{
		Conn: conn,
	}

	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("注册发送信息错误 err=", err)
		return
	}

	msg, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("readPkg err =", err)
		return
	}

	var registerResultMsg message.RegisterResultMsg
	err = json.Unmarshal([]byte(msg.Data), &registerResultMsg)
	if registerResultMsg.Code == 200 {
		fmt.Println("注册成功,你重新登录一把")
		os.Exit(0)
	} else {
		fmt.Println(registerResultMsg.Error)
		os.Exit(0)
	}
	return
}

//登陆校验

func (this *UserProcess) Login(userId int, userPwd string) (err error) {
	fmt.Printf("userId = %d userPwd = %s \n", userId, userPwd)
	// 连接到服务器
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("client dail net err=", err)
		return
	}
	defer conn.Close()
	//准备通过conn 发送消息给服务器
	var msg message.Message
	msg.Type = message.LoginMsgType
	// 创建login msg
	var loginMsg message.LoginMsg
	loginMsg.UserId = userId
	loginMsg.UserPwd = userPwd

	//loginMsg序列化
	data, err := json.Marshal(loginMsg)
	if err != nil {
		fmt.Println("json Marshal err: ", err)
		return
	}

	msg.Data = string(data)

	//msg 序列化
	data, err = json.Marshal(msg)
	if err != nil {
		fmt.Println("json Marshal err:", err)
		return
	}

	fmt.Println("客户端发送的内容:", string(data))
	//data就是需要发送的消息
	//data长度先发送到服务器 转成表示长度的切片
	var pkgLen uint32
	pkgLen = uint32(len(data))
	//fmt.Println("pkgLen:",pkgLen)
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[0:4], pkgLen)
	//fmt.Println("begin to write buf...",buf)
	n, err := conn.Write(buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write 长度不为四 发送错误", err)
		return
	}
	fmt.Println("客户端，发送消息的长度=%d", len(data))

	//发送消息本身
	fmt.Println("发送消息本身开始...")
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("conn.Write 消息本身 发送错误", err)
		return
	}
	tf := &util.Transfer{
		Conn: conn,
	}
	msg, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("readPkg err =", err)
		return
	}
	var loginResultMsg message.LoginResultMsg
	err = json.Unmarshal([]byte(msg.Data), &loginResultMsg)
	if loginResultMsg.Code == 200 {
		fmt.Println("登录成功")
		//启动一个协程 保持和服务器端通讯
		go serverProcessMsg(conn)

		//显示登陆菜单
		for {
			ShowMenu()
		}
	} else {
		fmt.Println(loginResultMsg.Error)
	}

	return
}

func serverProcessMsg(conn net.Conn) {
	tf := &util.Transfer{
		Conn: conn,
	}
	//不停的读取服务器发送的消息
	for {
		fmt.Println("客服端%s正在等待读取服务器发送的消息")
		msg, err := tf.ReadPkg()
		if err != nil {
			fmt.Println("tf.readPkg err=", err)
			return
		}
		fmt.Println(msg)

	}
}
