package main

import(
	"fmt"
	"net"
	"gocode/MyQQ/Common/Message"
	"encoding/json"
	"encoding/binary"
	"io"
)

//整一个函数用来接收客户端发来的信息
func readPkg(conn net.Conn)(mes Message.Message ,err error){
	
	//创建一个缓冲区用来读取从客户端发来的东西
	buf :=make([]byte,4096)
	fmt.Println("读取客户端发送的数据中")
		//首先读取等下要发送的信息的长度
	n,err:=conn.Read(buf[:4])
	if err !=nil||n!=4{
		fmt.Printf("客户端数据长度接收失败，err=%v\n",err)
		return
	}
	fmt.Println("读到的buf=",buf[:4])


		//接下来读取信息本体
	//先将长度转化为uint32,获取到长度
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(buf[0:4])
	//正式开始读取数据
	n,err=conn.Read(buf[:pkgLen])
	if err !=nil||n!=int(pkgLen){
		fmt.Printf("客户端数据接收失败，err=%v\n",err)
		return
	}

	//将读到的东西反序列化

	err = json.Unmarshal(buf[:pkgLen],&mes)//这个mes在返回值那里就已经声明了，所以不用另外声明
	if err !=nil{
		
		fmt.Printf("客户端数据反序列化失败，err=%v\n",err)
		return
	}

	return

}


//编写一个writePkg函数用来专门发送信息给对面
func writePkg(conn net.Conn,data []byte)(err error){
		//获取到data的长度，同时转成一个表示长度的tpye切片
		var pkgLen uint32
		pkgLen=uint32(len(data))
		var buf [4]byte
		binary.BigEndian.PutUint32(buf[0:4],pkgLen)
			//要把data长度发送给对面
		n,err:=conn.Write(buf[0:4])
		if err !=nil||n!=4{
			fmt.Printf("长度发送失败，err=%v\n",err)
			return
		}
		fmt.Println("服务器发送长度成功！len=",len(data))
	
			//把data本体发送给对面
		n,err =conn.Write(data)
		if err !=nil{
			fmt.Printf("信息发送失败，err=%v\n",err)
			return
		}

		return
	
}


//编写一个ServerProcessLogin函数用来专门处理登录请求
func ServerProcessLogin(conn net.Conn,mes *Message.Message)(err error){
	//核心代码
	//先从mes中取出data，并反序列化为LoginMes
	var loginMes Message.LoginMes
	err= json.Unmarshal([]byte(mes.Data),&loginMes)
	if err!=nil{
		fmt.Printf("mes中的data反序列化失败！，err=%v\n",err)
	}

	//先声明一个M.M，用于存储发过去的信息
	var resMes Message.Message
	resMes.Type=Message.LoginResMesType

	//返回消息给客户端的话，必须要声明一个LoginResMes来存消息
	var loginResMes Message.LoginResMes

	//假设id=100, psw=123
	if loginMes.UserId=="100" && loginMes.UserPSW=="123"{
		//合法
		loginResMes.Code=200
	}else{
		//不合法，用户不存在
		loginResMes.Code=500
		loginResMes.Error="此账号不存在，请注册后使用"
	}

	//将loginMes序列化
	data, err :=json.Marshal(loginResMes)
	if err!=nil{
		fmt.Printf("loginResMes序列化失败！，err=%v\n",err)
	}
	//将序列化后的信息存入resMes
	resMes.Data=string(data)

	//将resMes序列化准备发送
	data, err =json.Marshal(resMes)
	if err!=nil{
		fmt.Printf("resMes序列化失败！，err=%v\n",err)
	}

	//将发送这个功能整合到writePkg函数中

	err=writePkg(conn,data)
	return 

}




//编写一个函数ServerProcessMes函数，用来分配各种不同的消息的处理方法
func ServerProcessMes(conn net.Conn,mes *Message.Message)(err error){
	switch mes.Type{
	case Message.LoginMesType:
		//登录处理
		err = ServerProcessLogin(conn,mes)
	case Message.RegisterMesType:
		//注册处理
	default :
		fmt.Println("该类信息暂时没有录入信息库，所以也不晓得怎么办........")
	}

	return
}


func process(conn net.Conn){
//延时关闭conn
	defer conn.Close()
	
	//读取客户端发送的信息
	for{
		mes,err := readPkg(conn)
		if err !=nil{
			if err==io.EOF{
				fmt.Printf("客户端关闭了连接，err=%v，服务器也正常关闭\n",err)
				return
			}
			fmt.Printf("客户端数据读取（readPkg()）失败，err=%v\n",err)
			return
		}

		//fmt.Printf("Message=%v",mes)
		err= ServerProcessMes(conn,&mes)
		if err!=nil{
			fmt.Printf("ServerProcessMes失败！，err=%v\n",err)
			return
		}
	
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

