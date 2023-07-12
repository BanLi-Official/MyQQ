package main

import(
	"fmt"
	"os"
	"gocode/MyQQ/client/process"
	
)

//声明两个变量，存储账户与密码
var (
	userId string
	userPSW string
)

func main(){
	var flag bool=true //判定是否循环
	var function int //存储客户选择的功能

	for flag {
		fmt.Println("-------------------------欢迎使用QQ青春版--------------------------")
		fmt.Println("                           1.登录账号                             ")
		fmt.Println("                           2.注册账号                             ")
		fmt.Println("                           3.退出软件                             ")
		fmt.Println()

		fmt.Println("请选择功能：")
		fmt.Scanf("%d\n",&function)


		switch function{
		case 1:
			fmt.Println("账户登录：")
			fmt.Print("请输入用户账号：")
			fmt.Scanf("%s\n",&userId)
			fmt.Print("请输入用户密码：")
			fmt.Scanf("%s\n",&userPSW)
			process := &process.Process{}
			err:=process.Login(userId,userPSW)
			if err!=nil{
				fmt.Printf("登录时出现错误，err=%v",err)
			}
			//flag=false
			break
		case 2:
			fmt.Println("账户注册：")
			//flag=false
			break
		case 3:
			fmt.Println("退出软件：")
			//flag=false
			os.Exit(0)
		default:
			fmt.Println("没有这个功能，再选一次")
		}

	}


	
}