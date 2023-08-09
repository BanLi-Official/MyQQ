package utils

import(
	"fmt"
	"net"
	"gocode/MyQQ/Common/Message"
	"encoding/json"
	"encoding/binary"

)

//写一个结构体将这些函数串起来
type Transfer struct{
	//将连接conn整合到结构体中
	Conn net.Conn
	//将缓冲区整合到结构体中
	Buf [4096]byte
}



//整一个函数用来接收客户端发来的信息
func (this *Transfer) ReadPkg()(mes Message.Message ,err error){
	
	//创建一个缓冲区用来读取从客户端发来的东西
	//buf :=make([]byte,4096)
	//fmt.Printf("2222222222222222222222222222222222222\n")
	fmt.Println("读取客户端发送的数据中")
		//首先读取等下要发送的信息的长度
	n,err:=this.Conn.Read(this.Buf[:4])
	if err !=nil||n!=4{
		fmt.Printf("客户端数据长度接收失败，err=%v\n",err)
		return
	}
	//fmt.Println("读到的buf=",this.Buf[:4])

	//fmt.Printf("2lllllllllllllllllllllllllllllllllllllllllhhhhhh2\n")
		//接下来读取信息本体
	//先将长度转化为uint32,获取到长度
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(this.Buf[0:4])//将切片中的内容从切片中提出来，存入uint32
	//正式开始读取数据
	n,err=this.Conn.Read(this.Buf[:pkgLen])
	if err !=nil||n!=int(pkgLen){
		fmt.Printf("客户端数据接收失败，err=%v\n",err)
		return
	}

	//将读到的东西反序列化

	err = json.Unmarshal(this.Buf[:pkgLen],&mes)//这个mes在返回值那里就已经声明了，所以不用另外声明
	if err !=nil{
		
		fmt.Printf("客户端数据反序列化失败，err=%v\n",err)
		return
	}


	return

}



//整一个函数用来接收客户端发来的信息
func (this *Transfer) ReadPkgLong()(mes Message.Message ,err error){
	
	//创建一个缓冲区用来读取从客户端发来的东西
	//buf :=make([]byte,4096)
	//fmt.Printf("2222222222222222222222222222222222222\n")
	//fmt.Println("读取客户端发送的数据中")
		//首先读取等下要发送的信息的长度
	n,err:=this.Conn.Read(this.Buf[:8])
	if err !=nil||n!=4{
		fmt.Printf("客户端数据长度接收失败，err=%v\n",err)
		return
	}
	//fmt.Println("读到的buf=",this.Buf[:4])

	//fmt.Printf("2lllllllllllllllllllllllllllllllllllllllllhhhhhh2\n")
		//接下来读取信息本体
	//先将长度转化为uint32,获取到长度
	var pkgLen uint64
	pkgLen = binary.BigEndian.Uint64(this.Buf[0:8])//将切片中的内容从切片中提出来，存入uint32
	//正式开始读取数据
	n,err=this.Conn.Read(this.Buf[:pkgLen])
	if err !=nil||n!=int(pkgLen){
		fmt.Printf("客户端数据接收失败，err=%v\n",err)
		return
	}

	//将读到的东西反序列化

	err = json.Unmarshal(this.Buf[:pkgLen],&mes)//这个mes在返回值那里就已经声明了，所以不用另外声明
	if err !=nil{
		
		fmt.Printf("客户端数据反序列化失败，err=%v\n",err)
		return
	}


	return

}


//编写一个writePkg函数用来专门发送信息给对面
func (this *Transfer) WritePkg(data []byte)(err error){
		//获取到data的长度，同时转成一个表示长度的tpye切片
		var pkgLen uint32
		pkgLen=uint32(len(data))
		//var buf [4]byte
		binary.BigEndian.PutUint32(this.Buf[0:4],pkgLen)//把这个长度数据转化为切片才能进行传输
			//要把data长度发送给对面
		n,err:=this.Conn.Write(this.Buf[0:4])
		if err !=nil||n!=4{
			fmt.Printf("长度发送失败，err=%v\n",err)
			return
		}
		fmt.Println("服务器发送长度成功！len=",len(data))
	
			//把data本体发送给对面
		n,err =this.Conn.Write(data)
		if err !=nil{
			fmt.Printf("信息发送失败，err=%v\n",err)
			return
		}

		return
	
}

