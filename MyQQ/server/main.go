package main

import(
	"fmt"
	"net"
)


func process(conn net.Conn){
//延时关闭conn
	defer conn.Close()

	//读取客户端发送的信息
	for{
		buf :=make([]byte,4096)
		fmt.Println("读取客户端发送的数据")
		n,err:=conn.Read(buf[:4])
		if err !=nil||n!=4{
			fmt.Printf("客户端数据长度接收失败，err=%v\n",err)
			return
		}
		fmt.Println("读到的buf=",buf[:4])
	}

}



func main(){
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

		go process(conn)
	}

}

