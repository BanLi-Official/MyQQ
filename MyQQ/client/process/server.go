package process

import (
	"fmt"
	"os"
	"net"
	"gocode/MyQQ/client/utils"
)

type Server struct{

}

//写一个展示登录成功的函数
func (this * Server)ShowMenu(){
	fmt.Println()
	fmt.Println("-----------------------------欢迎登录QQ青春版------------------------------------")
	fmt.Println("                             1.显示在线用户列表                                    ")
	fmt.Println("                             2.发送消息                                    ")
	fmt.Println("                             3.信息列表                                    ")
	fmt.Println("                             4.退出系统                                    ")

	var key int
	fmt.Println("请选择功能：")
	fmt.Scanf("%d\n",&key)
	switch key{
		case 1 :
			fmt.Println("你选择了显示在线用户列表功能")
		case 2:
			fmt.Println("你选择了发送信息功能")
		case 3:
			fmt.Println("你选择了信息列表功能")
		case 4:
			fmt.Println("你选择了退出系统")
			os.Exit(0)
		default:
			fmt.Println("没有这个功能")
	}
}


//这个协程持续关注服务器端发来信息的情况，如果有信息从服务器端发送给了客户端，则在客户端部分展示这个信息
func (this * Server)serverProcessMes(conn net.Conn){
	tf:=&utils.Transfer{
		Conn :conn,
	}
	fmt.Println("客户端启动了一个协程，这个协程持续关注服务器端发来信息的情况")
	for{
		fmt.Println("客户端通过协程serverProcessMes在不断地等待服务器发来的消息")
		mes,err:=tf.ReadPkg()
		if err !=nil{
			fmt.Printf("协程serverProcessMes出现了错误，err=%v",err)
		}
		//如果读取到了，那么展示这个信息
		fmt.Printf("协程serverProcessMes读到了来自服务器的信息，mes=%v",mes)

	}

}