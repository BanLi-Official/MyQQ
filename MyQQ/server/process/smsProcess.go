package processes

import (

	"fmt"
	"net"
	"gocode/MyQQ/Common/Message"
	"gocode/MyQQ/server/utils"
	"encoding/json"
	"gocode/MyQQ/server/model"

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



//转发消息到redis里面
func (this *SmsProcess)SendPppMesToRedis(mes *Message.Message){
	//取出mes中的内容PppMes
	var pppMes_OffLine Message.PppMes_OffLine
	err:=json.Unmarshal([]byte(mes.Data),&pppMes_OffLine)
	if err!= nil{
		fmt.Printf("json.Unmarshal Error=%v\n",err)
		return
	}


	//fmt.Println(pppMes_OffLine)
	//fmt.Println("From:",pppMes_OffLine.FromUser)
	//fmt.Println("To  :",pppMes_OffLine.ToUser)
	//fmt.Println("Contain:",pppMes_OffLine.Content)

	//存入数据库

	model.MyUserDao.InsertOffLineMassage(pppMes_OffLine.FromUser,pppMes_OffLine.ToUser,pppMes_OffLine.Content)


	//2023年8月9日15:44:21
	//      今天完成了从服务器发送所有用户的功能，同时客户端也可以分辨各个用户的在线状态，根据在线状态分别发送在线信息和离线信息；
	// 服务器端接收到了来自客户端的离线信息请求，并且能够展示出来这个信息的详细内容
	// 	 明天的任务：将客户端发过来的离线系信息传到redis当中，然后当用户登录成功的时候系统自动发送离线信息给用户

}