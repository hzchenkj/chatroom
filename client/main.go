package main

import (
	"fmt"
	"os"
)

var userId int
var userPwd string

func main() {
	//用户输入
	var key int
	//是否显示菜单
	var loop = true
	fmt.Println(loop)
	for {

		fmt.Println(">>>>欢迎登录多人聊天系统<<<<")
		fmt.Println("\t1 登陆聊天室")
		fmt.Println("\t2 注册用户")
		fmt.Println("\t3 退出系统")
		fmt.Println("请选择 (1-3):")

		fmt.Scanf("%d\n", &key)
		switch key {
		case 1:
			fmt.Println("登陆聊天室")
			loop = false
		case 2:
			fmt.Println("注册用户")
			loop = false
		case 3:
			fmt.Println("退出系统")
			loop = false
			os.Exit(0)

		default:
			fmt.Println("你的输入有误，请重新输入")
		}

		if !loop{
			break
		}

	}

	//根据用户输入，显示新的提示信息

	if key == 1 {
		//登陆聊天室
		fmt.Println("请输入用户id")
		fmt.Scanf("%d\n", &userId)
		fmt.Println("请输入用户密码")
		fmt.Scanf("%s\n", &userPwd)
		login(userId, userPwd)

	} else if key == 2 {
		fmt.Println("用户注册")
	}
}
