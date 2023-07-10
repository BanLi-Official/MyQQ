package main

import(
	"fmt"
	"net"
	"encoding/json"
	"encoding/binary"
	"gocode/MyQQ/Common/Message"

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

