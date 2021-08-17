package main

import (
	"chatroom/common/message"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"
)

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