package processes

import(
	"fmt"
	"net"
	"gocode/MyQQ/Common/Message"
	"gocode/MyQQ/server/utils"
	"encoding/json"


)

//写一个结构体将这些函数串起来
type UserProcess struct{
	//将连接conn整合到结构体中
	Conn net.Conn
	
}

//编写一个ServerProcessLogin函数用来专门处理登录请求
func (this *UserProcess) ServerProcessLogin(mes *Message.Message)(err error){
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

	tf:=&utils.Transfer{
		Conn:this.Conn,
	}

	err=tf.WritePkg(data)
	return 

}


