package process

import (
	"fmt"
	"os"
)

//显示登陆成功后的界面

//显示菜单
func ShowMenu() {
	fmt.Println("-------- 恭喜****登陆成功")
	fmt.Println("\t1 显示在线用户列表")
	fmt.Println("\t2 发送消息")
	fmt.Println("\t3 信息列表")
	fmt.Println("\t4 退出系统")
	fmt.Println("请选择(1-4)")
	var key int
	fmt.Scanf("%d\n", &key)
	switch key {
	case 1:
		fmt.Println("显示在线用户列表")
	case 2:
		fmt.Println("发送消息")
	case 3:
		fmt.Println("信息列表")
	case 4:
		fmt.Println("你退出了系统...")
		os.Exit(0)
	default:
		fmt.Println("你的输入有误，请重新输入")
	}

}
