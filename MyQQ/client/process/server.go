package process

import (
	"fmt"
	"os"
	"net"
	"gocode/MyQQ/client/utils"

	"gocode/MyQQ/Common/Message"
	"gocode/MyQQ/Common/MyFmt"
	"encoding/json"

)

type Server struct{

}

//写一个展示登录成功的函数
func (this * Server)ShowMenu(){
	fmt.Println()
	fmt.Println("-----------------------------欢迎登录QQ青春版------------------------------------")
	fmt.Println("                             1.显示在线用户列表                                    ")
	fmt.Println("                             2.发送群消息                                    ")
	fmt.Println("                             3.发送单独消息                                   ")
	fmt.Println("                             4.退出系统                                    ")

	var content string
	//后面总会用到SmsProcess实例，因此将SmsProcess定义在外部
	smsProcess:=&SmsProcess{}

	var key int
	fmt.Println("请选择功能：")
	fmt.Scanf("%d\n",&key)
	switch key{
		case 1 :
			//fmt.Println("你选择了显示在线用户列表功能")
			outputOnlineUser()
		case 2:
			fmt.Println("你选择了发送信息功能")
			fmt.Printf("你想对大家说什么：")
			fmt.Scanf("%s\n",&content)
			smsProcess.SendGroupMes(content)
		case 3:
			var userCode string
			var content_p2p string
			fmt.Println("你选择了发送单独消息功能")
			//outputOnlineUser_outSelf()
			//process := &Process{}
			GetAllUser()

			fmt.Printf("你想对谁说：\n")
			fmt.Scanf("%s\n",&userCode)
			fmt.Println("你想说些什么：")

			
			MyFmt.Scanf(&content_p2p)
			//fmt.Scanf("%s\n",&content_p2p)
			//发送函数（）
			smsProcess.SendPppMes(content_p2p,userCode)
			fmt.Println("发送成功啦")


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
		
		//如果读到了消息，则进行下一步
		switch mes.Type{
		case Message.NotifyUserStatusMesType:
			//有人上线了,或者下线了
			//取出NotifyUserStatusMes

			var notifyUserStatusMes Message.NotifyUserStatusMes
			
			json.Unmarshal([]byte(mes.Data),&notifyUserStatusMes)
			//把这个用户的信息传到客户端维护的map[string] User中
			fmt.Println("开始更改信息 notifyUserStatusMes",notifyUserStatusMes)
			updateUserStatus(&notifyUserStatusMes)
		case Message.SmsMesType://有人群发消息
			outputGroupMes(&mes)
		case Message.PppMesType://有人发来点对点消息
			outputPppMes(&mes)
		default:
			fmt.Println("服务器端发来了未知的信息")
		}




		//如果读取到了，那么展示这个信息
		//fmt.Printf("协程serverProcessMes读到了来自服务器的信息，mes=%v",mes)

	}

}