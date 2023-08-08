package processes

import (

	"fmt"
	"net"
	"gocode/MyQQ/Common/Message"
	"gocode/MyQQ/server/utils"
	"encoding/json"


)

type SmsProcess struct{

}


//转发消息
func (this *SmsProcess)SendGroupMes(mes *Message.Message){
	//遍历服务器端的onlineUsers
	//将消息发出

	//取出mes中的内容SmsMes
	var smsMes Message.SmsMes
	err:=json.Unmarshal([]byte(mes.Data),&smsMes)
	if err!= nil{
		fmt.Printf("json.Unmarshal Error=%v\n",err)
		return
	}

	data ,err := json.Marshal(mes)
	if err!= nil{
		fmt.Printf("json.Marshal Error=%v\n",err)
		return
	}

	for id,up :=range userMgr.onlineUsers{
		//在这里过滤掉自己
		if id==smsMes.UserId{
			continue
		}

		this.SendMesToEachOnlineUser(data,up.Conn)
	}
}


func (this *SmsProcess)SendMesToEachOnlineUser(data []byte,conn net.Conn){
	//创建一个transfer实例，发送data
	tf:=&utils.Transfer{
		Conn:conn,
	}
	err:=tf.WritePkg(data)
	if err!=nil{
		fmt.Println("转发消息失败 err=",err)
	}

}


//转发消息给某人
func (this *SmsProcess)SendPppMesToSb(mes *Message.Message){
	//取出mes中的内容PppMes
	var pppMes Message.PppMes
	err:=json.Unmarshal([]byte(mes.Data),&pppMes)
	if err!= nil{
		fmt.Printf("json.Unmarshal Error=%v\n",err)
		return
	}

	//fmt.Println(pppMes)

	//发送给特定的人

	//找到这个人
	for id,up :=range userMgr.onlineUsers{
		//在这里过滤掉自己
		if id==pppMes.ToUser{
			//创建一个transfer实例，发送data
			tf:=&utils.Transfer{
				Conn:up.Conn,
			}
			//将loginMes序列化
			data, err :=json.Marshal(mes)
			if err!=nil{
				fmt.Printf("mes序列化失败！，err=%v\n",err)
			}

			err=tf.WritePkg(data)
			if err!=nil{
				fmt.Println("转发消息失败 err=",err)
				return
			}

			fmt.Println("转发消息成功")
		}
	}

	//找到这个人的conn
}