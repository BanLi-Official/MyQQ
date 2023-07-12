package process



import(
	"fmt"
	"net"
	"encoding/json"
	"encoding/binary"
	"gocode/MyQQ/Common/Message"
	"gocode/MyQQ/client/utils"
	_"time"

)

//写一个结构体将这些函数串起来
type Process struct{
	
}



func (this *Process) Login(userId string,userPSW string)(err error){
	//fmt.Printf("账户登录成功，userId=%v  userPSW=%v",userId,userPSW)
	//连接到服务器
	conn,err := net.Dial("tcp","127.0.0.1:8889")
	if err !=nil{
		fmt.Printf("客户端连接到服务器失败，err=%v\n",err)
		return
	}
	defer conn.Close()

	//先声明一个用于传送数据的变量
	tf:=&utils.Transfer{
		Conn:conn,
	}
	
	//准备通过conn发送消息给服务器
	var mes Message.Message
	mes.Type=Message.LoginMesType

	//创建一个LoginMessage的结构体
	var loginMes Message.LoginMes
	loginMes.UserId=userId
	loginMes.UserPSW=userPSW


	//4.将loginMessage序列化，并赋值给mes.data
	data,err :=json.Marshal(loginMes)
	if err !=nil{
		fmt.Printf("将loginMessage序列化失败，err=%v\n",err)
		return
	}
	mes.Data=string(data)

	//5.将mes序列化
	data,err = json.Marshal(mes)
	if err !=nil{
		fmt.Printf("将mes序列化失败，err=%v\n",err)
		return
	}

	//6.data 此时就是需要发送的消息
		//获取到data的长度，同时转成一个表示长度的tpye切片
	var pkgLen uint32
	pkgLen=uint32(len(data))
	var bytes [4]byte
	binary.BigEndian.PutUint32(bytes[0:4],pkgLen)
		//要把data长度发送给服务器
	n,err:=conn.Write(bytes[0:4])
	if err !=nil||n!=4{
		fmt.Printf("长度发送失败，err=%v\n",err)
		return
	}
	fmt.Println("客户端发送长度成功！len=",len(data))

		//把data本体发送给服务器
	n,err =conn.Write(data)
	if err !=nil{
		fmt.Printf("信息发送失败，err=%v\n",err)
		return
	}


		//在这里处理服务器返回的信息
	mes,err =tf.ReadPkg()
	if err !=nil{
		fmt.Printf("readPkg(conn)失败，err=%v\n",err)
		return
	}
	//将mes的data部分反序列化为LoginResMes
	var loginResMes Message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data),&loginResMes)
	if err !=nil{
		fmt.Printf("将mes的data部分反序列化为LoginResMes失败，err=%v\n",err)
		return
	}
	if loginResMes.Code==200{
		
		//在此处还需要客户端启动一个协程，这个协程持续关注服务器端发来信息的情况，如果有信息从服务器端发送给了客户端，则在客户端部分展示这个信息
		server :=&Server{}
		go server.serverProcessMes(conn)
		

		//登录成功，展示登录成功的菜单
		for{
			
			server.ShowMenu()
		}
	}else if loginResMes.Code==500{
		fmt.Println(loginResMes.Error)
	}


	






	return 
}

