package main

import (
	"chatroom/common/message"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
)

func serverProcessLogin(conn net.Conn,msg *message.Message)  (err error){
	//msg 中取出data msg.Data 发序列化层LoginMsg
	var loginMsg message.LoginMsg
	err = json.Unmarshal([]byte(msg.Data), &loginMsg)
	if err != nil{
		fmt.Println("json.Unmarshal fail err=",err)
		return
	}

	//
	var resMsg message.Message
	resMsg.Type = message.LoginResultMsgType

	var loginResutlMsg message.LoginResultMsg
	 //如果用户id是100 密码是123456
	 if loginMsg.UserId == 100 && loginMsg.UserPwd == "123456"{
	 	//合法
	 	loginResutlMsg.Code=200
	 }else {
	 	//不合法
	 	loginResutlMsg.Code =500
	 	loginResutlMsg.Error = "用户不存在"
	 }
	 data ,err  := json.Marshal(loginResutlMsg)
	 if err != nil{
	 	fmt.Println("json.Marshal faild err:",err)
	 	return
	 }
	 resMsg.Data = string(data)

	 data,err  = json.Marshal(resMsg)
		if err != nil{
			fmt.Println("json.Marshal faild err:",err)
			return
		}
	//发送到客户端
	err = writePkg(conn,data)

	return err
}
//编写一个server ProcessMsg
//根据客户端发送消息的内容不同，决定调用哪个函数
func serverProcessMsg(conn net.Conn,msg *message.Message) (err error){
	switch msg.Type {
	case message.LoginMsgType:
		//处理登录逻辑
		err = serverProcessLogin(conn,msg)

	case message.RegisterMsgType:
		//
	default:
		fmt.Println("消息类型不存在，无法处理...")
	}
	return
}

func process(conn net.Conn){
	//延时关闭
	defer  conn.Close()

	//读取客户端发送的信息
	//  循环读取
	for{
		msg,err := readPkg(conn)
		if err !=nil{
			if err == io.EOF {
				fmt.Println("客户端退出，服务器端也退出")
				return
			}else{
				fmt.Println("readPkg err=",err)
				return
			}
		}

		err = serverProcessMsg(conn,&msg)
		if err != nil{
			return
		}
		fmt.Println("msg:",msg)
	}
}

func writePkg(conn net.Conn,data []byte) (err error) {
	//先发送一个长度给对方
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

	n,err = conn.Write(data)
	if n != int(pkgLen) || err != nil{
		fmt.Println("conn.Write 内容 发送错误",err)
		return
	}

	return err
}

func readPkg(conn net.Conn) (msg message.Message,err error) {
	buf := make([]byte,8096)
	fmt.Println("读取客户端发送的数据...")
	n, err := conn.Read(buf[:4])
	if n != 4 || err != nil {
		err = errors.New("read pkg header err")
		return
	}
	fmt.Println("读到buf=", buf[:4])
	//根据buf[:4] 转成一个unit32
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(buf[:4])
	fmt.Println(pkgLen)
	//根据pkgLen读取消息内容
	n,err =conn.Read(buf[:pkgLen])
	if n!= int(pkgLen)|| err !=nil{
		return
	}
	//pkgLen 发序列化 成message
	err = json.Unmarshal(buf[:pkgLen],&msg)
	if err != nil{
		fmt.Println("json.Unmarshal err",err)
		return
	}
	return
}

func main(){
	//提示信息
	 fmt.Println("服务器在8889 端口监听...")
	 listern ,err := net.Listen("tcp", "0.0.0.0:8889")
	 if err != nil{
	 	fmt.Println("net listern err=",err)
		 return
	 }
	defer  listern.Close()

	 for {
	 	fmt.Println("等待客户端连接服务器...")
	 	conn,err := listern.Accept()
	 	if err !=nil{
	 		fmt.Println("listern accpet err :",err)
		}
		//连接成功，启动一个协程 和客户端保持通讯
		go process(conn)

	 }

}