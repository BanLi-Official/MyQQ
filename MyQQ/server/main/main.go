package main

import(
	"fmt"
	"net"
	"time"
	"gocode/MyQQ/server/model"
	"gocode/MyQQ/server/process"
)






func process(conn net.Conn){
//延时关闭conn
	defer conn.Close()
	processor :=&Processor{
		Conn: conn,
	}
	err:=processor.Process2()
	if err!=nil{
		fmt.Printf("客户端与服务器端协程通讯错误,客户端已经离线\n")
		up :=&processes.UserProcess{
			Conn:conn,
		}
		err = up.ServerProcessDelete(conn)
	}


	return 

}



//这里写一个函数，完成对userDAO的初始化任务
func InitUserDao(){
	model.MyUserDao=model.NewUserDao(pool)//这个pool是在redis中定义的一个全局变量
	//这里需要注意两个参数的初始化顺序,先initpool,再inituserDao
}



func main(){ 

	//当服务器启动时初始化数据库连接池
	InitPool("127.0.0.1:6379",16,0,300*time.Second)
	InitUserDao()
	//提示信息
	fmt.Println("服务器启动，并开始监听端口8889")
	//写listen
	listen,err:=net.Listen("tcp","127.0.0.1:8889")
	defer listen.Close()
	if err!=nil{
		fmt.Printf("建立监听失败！，err=%v\n",err)
		return
	}
	//开始等待Accept（）
	for{
		fmt.Println("服务器开始等待客户端的连接")
		conn,err := listen.Accept()
		if err!=nil{
			fmt.Printf("连接建立失败！，err=%v\n",err)
		}
		//建立一个协程，专门来处理本次连接
		go process(conn)
	}

}

