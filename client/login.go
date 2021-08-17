package main

import (
	"chatroom/common/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

//登陆校验

func login(userId int, userPwd string) (err error) {
	fmt.Printf("userId = %d userPwd = %s \n", userId, userPwd)
	// 连接到服务器
	conn ,err := net.Dial("tcp","localhost:8889")
	if err!= nil{
		fmt.Println("client dail net err=",err)
		return
	}
	defer  conn.Close()
	//准备通过conn 发送消息给服务器
	var msg  message.Message
	msg.Type = message.LoginMsgType
	// 创建login msg
	var loginMsg message.LoginMsg
	loginMsg.UserId = userId
	loginMsg.UserPwd = userPwd

	//loginMsg序列化
	data, err := json.Marshal(loginMsg)
	if err != nil{
		fmt.Println("json Marshal err: ",err)
		return
	}

	msg.Data = string(data)

	//msg 序列化
	data, err = json.Marshal(msg)
	if err != nil{
		fmt.Println("json Marshal err:",err)
		return
	}

	fmt.Println("客户端发送的内容:",string(data))
	//data就是需要发送的消息
	//data长度先发送到服务器 转成表示长度的切片
	var pkgLen uint32
	pkgLen = uint32(len(data))
	//fmt.Println("pkgLen:",pkgLen)
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[0:4],pkgLen)
	//fmt.Println("begin to write buf...",buf)
	n,err := conn.Write(buf[:4])
	if n !=4 || err != nil{
		fmt.Println("conn.Write 长度不为四 发送错误",err)
		return
	}
	fmt.Println("客户端，发送消息的长度=%d",len(data))


	//发送消息本身
	fmt.Println("发送消息本身开始...")
	_,err = conn.Write(data)
	if err != nil{
		fmt.Println("conn.Write 消息本身 发送错误",err)
		return
	}

	msg,err = readPkg(conn)
	if err != nil{
		fmt.Println("readPkg err =",err)
		return
	}
	var loginResultMsg message.LoginResultMsg
	err = json.Unmarshal([]byte(msg.Data),&loginResultMsg)
	if loginResultMsg.Code == 200 {
		fmt.Println("登录成功")
	}else if loginResultMsg.Code== 500{
		fmt.Println(loginResultMsg.Error)
	}

	return
}
