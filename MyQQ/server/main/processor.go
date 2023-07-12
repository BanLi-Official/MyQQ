package main

import(
	"fmt"
	"net"
	"gocode/MyQQ/Common/Message"
	"gocode/MyQQ/server/process"
	"gocode/MyQQ/server/utils"
	"io"
)

//写一个结构体将这些函数串起来
type Processor struct{
	//将连接conn整合到结构体中
	Conn net.Conn
	
}



//编写一个函数ServerProcessMes函数，用来分配各种不同的消息的处理方法
func (this *Processor) ServerProcessMes(mes *Message.Message)(err error){
	switch mes.Type{
	case Message.LoginMesType:
		//登录处理
		//实例化一个userProcess，为了调用其中的ServerProcessLogin函数
		up :=&processes.UserProcess{
			Conn:this.Conn,
		}
		err = up.ServerProcessLogin(mes)
	case Message.RegisterMesType:
		//注册处理
	default :
		fmt.Println("该类信息暂时没有录入信息库，所以也不晓得怎么办........")
	}

	return
}


func (this *Processor) Process2()(err error){
	//读取客户端发送的信息
	for{
		tf:=&utils.Transfer{
			Conn:this.Conn,
		}	
		mes,err := tf.ReadPkg()
		if err !=nil{
			if err==io.EOF{
				fmt.Printf("客户端关闭了连接，err=%v，服务器也正常关闭\n",err)
				return err
			}
			fmt.Printf("客户端数据读取（readPkg()）失败，err=%v\n",err)
			return err
		}

		//调用一个函数专门处理本次接收到的mes
		err= this.ServerProcessMes(&mes)
		if err!=nil{
			fmt.Printf("ServerProcessMes失败！，err=%v\n",err)
			return err
		}
	
	}


}
