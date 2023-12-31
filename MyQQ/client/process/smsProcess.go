package process

import(
	"fmt"
	"gocode/MyQQ/Common/Message"
	"gocode/MyQQ/Common/Errors"
	"encoding/json"
	"gocode/MyQQ/client/utils"


)

type SmsProcess struct{

}
//发送一个p2p的信息
func (this *SmsProcess) SendPppMes (content string,toUser string)(err error){
	//先判断是否在线
	isOnline :=checkIsOnline(toUser)
	if isOnline == false{
		err=Errors.ERROR_USER_NOT_ONLINE
		return
	}

	//创建一个Mes
	var mes Message.Message
	mes.Type=Message.PppMesType
	
	//创建一个PppMes，赋值并序列化
	var pppMes Message.PppMes
	pppMes.Content=content
	pppMes.FromUser=CurUser.UserId
	pppMes.ToUser=toUser

	data,err :=json.Marshal(pppMes)
	if err !=nil{
		fmt.Printf("将smsMes序列化失败，err=%v\n",err)
		return
	}
	//PppMes赋值给Mes
	mes.Data=string(data)

	//序列化Mes
	data,err =json.Marshal(mes)
	if err !=nil{
		fmt.Printf("将smsMes序列化失败，err=%v\n",err)
		return
	}
	

	//发送信息
	tf :=&utils.Transfer{
		Conn:CurUser.Conn,
	}
	err =tf.WritePkg(data)
	if err !=nil{
		fmt.Printf("WritePkg(单发信息)失败，err=%v\n",err)
		return
	}


	return


}




//发送一个p2p的信息,这里用于实现离线信息
func (this *SmsProcess) SendPppMes_OffLine (content string,toUser string)(err error){
	//创建一个Mes
	var mes Message.Message
	mes.Type=Message.PppMes_OffLineType
	
	//创建一个pppMes_OffLine，赋值并序列化
	var pppMes_OffLine Message.PppMes_OffLine
	pppMes_OffLine.Content=content
	pppMes_OffLine.FromUser=CurUser.UserId
	pppMes_OffLine.ToUser=toUser

	data,err := json.Marshal(pppMes_OffLine)
	if err!=nil{
		fmt.Printf("将pppMes_OffLine序列化过程中发生错误，err=%v\n",err)
	}

	//将序列化后的内容传入mes中
	mes.Data=string(data)

	//将整个mes序列化
	data ,err=json.Marshal(mes)
	if err!=nil{
		fmt.Printf("将包含pppMes_OffLine的mes序列化过程中发生错误，err=%v\n",err)
	}

	//将离线信息发送到服务器中
	tf := &utils.Transfer{
		Conn:CurUser.Conn,
	}

	//发送！
	err =tf.WritePkg(data)
	if err !=nil {
		fmt.Printf("发送离线信息到服务器发生错误，err=%v",err)
	}

	return 



}



func (this *SmsProcess) SendGroupMes (content string)(err error){
	//创建一个Mes
	var mes Message.Message
	mes.Type=Message.SmsMesType

	//创建一个smsMes的结构体
	var smsMes Message.SmsMes
	smsMes.UserId=CurUser.UserId
	smsMes.Content=content
	smsMes.UserStatus=CurUser.UserStatus

	//4.将smsMes序列化，并赋值给mes.data
	data,err :=json.Marshal(smsMes)
	if err !=nil{
		fmt.Printf("将smsMes序列化失败，err=%v\n",err)
		return
	}
	mes.Data=string(data)

	//5.将mes序列化
	data,err = json.Marshal(mes)
	if err !=nil{
		fmt.Printf("将mes序列化失败，err=%v\n",err)
		return
	}


	//先声明一个用于传送数据的变量
	tf:=&utils.Transfer{
		Conn:CurUser.Conn,
	}
	//6.发送序列化之后的mes
	err =tf.WritePkg(data)
	if err !=nil{
		fmt.Printf("WritePkg(群发信息)失败，err=%v\n",err)
		return
	}


	return


}