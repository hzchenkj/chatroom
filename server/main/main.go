package main

import (
	"fmt"
	"net"
)

func process(conn net.Conn) {
	//延时关闭
	defer conn.Close()

	processor := &Processor{
		Conn: conn,
	}
	err := processor.process2()
	if err != nil {
		fmt.Println("客户端和服务器通讯协程错误 err:", err)
		return
	}

}

func main() {
	//提示信息
	fmt.Println("服务器在8889 端口监听...")
	listen, err := net.Listen("tcp", "0.0.0.0:8889")
	if err != nil {
		fmt.Println("net listern err=", err)
		return
	}
	defer listen.Close()

	for {
		fmt.Println("等待客户端连接服务器...")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listern accpet err :", err)
		}
		//连接成功，启动一个协程 和客户端保持通讯
		go process(conn)

	}

}
