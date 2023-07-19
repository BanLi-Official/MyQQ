package process

import(
	"fmt"
	"gocode/MyQQ/Common/Message"
	"encoding/json"
	"gocode/MyQQ/client/utils"


)

type SmsProcess struct{

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