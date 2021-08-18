package main

import (
	"chatroom/common/message"
	"chatroom/server/process"
	"chatroom/server/util"
	"fmt"
	"io"
	"net"
)

type Processor struct {
	Conn net.Conn
}

//编写一个server ProcessMsg
//根据客户端发送消息的内容不同，决定调用哪个函数
func (this *Processor) serverProcessMsg(msg *message.Message) (err error) {
	fmt.Println("server processor .serverProcessMsg>>>")
	switch msg.Type {
	case message.LoginMsgType:
		//处理登录逻辑
		up := &process2.UserProcess{
			Conn: this.Conn,
		}
		err = up.ServerProcessLogin(msg)
	case message.RegisterMsgType:
		//处理注册
		fmt.Println("注册处理>>>>")
		up := &process2.UserProcess{
			Conn: this.Conn,
		}
		err = up.ServerProcessRegister(msg)
	default:
		fmt.Println("消息类型不存在，无法处理...")
	}
	return
}

func (this *Processor) process2() (err error) {
	//读取客户端发送的信息
	//  循环读取
	for {
		tf := &util.Transfer{
			Conn: this.Conn,
		}
		msg, err := tf.ReadPkg()
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端退出，服务器端也退出")
				return err
			} else {
				fmt.Println("readPkg err=", err)
				return err
			}
		}

		err = this.serverProcessMsg(&msg)
		if err != nil {
			return err
		}
		fmt.Println("msg:", msg)
	}
}
