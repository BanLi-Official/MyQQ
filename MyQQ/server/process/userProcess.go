package processes

import(
	"fmt"
	"net"
	"gocode/MyQQ/Common/Message"
	"gocode/MyQQ/server/utils"
	"encoding/json"
	"gocode/MyQQ/server/model"


)

//写一个结构体将这些函数串起来
type UserProcess struct{
	//将连接conn整合到结构体中
	Conn net.Conn
	//添加一个标识，用于区分用户
	UserId string
}


//编写一个通知所有在线用户的方法，userId通知其他用户：我上线了
func (this *UserProcess) NotifyOtherOnlineUser(userId string){
	//遍历onlineusers，然后一个个的发送
	for id ,up :=range this.onlineUsers{
		//过滤掉自己
		if id ==userId{
			continue
		}

		//开始通知，单独写一个方法
		up.NotifyMeOnline(userId)
	}
}



func (this *UserProcess) NotifyMeOnline(userId string){
	//组装NotifyUserStatusMes
	//先声明一个M.M，用于存储发过去的信息
	var mes Message.Message
	mes.Type=Message.RegisterResMesType

	//返回消息给客户端的话，必须要声明一个NotifyUserStatusMesType来存消息
	var notifyUserStatusMes Message.NotifyUserStatusMes
	notifyUserStatusMes.UserId=userId
	notifyUserStatusMes.Status=Message.UserOnline

	//将notifyUserStatusMes序列化
	data, err :=json.Marshal(notifyUserStatusMes)
	if err!=nil{
		fmt.Printf("notifyUserStatusMes序列化失败！，err=%v\n",err)
	}
	//将序列化后的信息存入mes
	mes.Data=string(data)

	//对mes再次进行序列化准备发送
	data, err =json.Marshal(mes)
	if err!=nil{
		fmt.Printf("mes序列化失败！，err=%v\n",err)
	}

	tf:=&utils.Transfer{
		Conn:this.Conn,
	}

	err=tf.WritePkg(data)
	if err!=nil{
		fmt.Printf("mes发送失败！，err=%v\n",err)
	}

	return 


}

//编写一个ServerProcessRegister函数用来专门处理登录请求
func (this *UserProcess) ServerProcessRegister(mes *Message.Message)(err error){
	//核心代码
	//先从mes中取出data，并反序列化为LoginMes
	var RegisterMes Message.RegisterMes
	err= json.Unmarshal([]byte(mes.Data),&RegisterMes)
	if err!=nil{
		fmt.Printf("mes中的data反序列化失败！，err=%v\n",err)
	}

	//先声明一个M.M，用于存储发过去的信息
	var resMes Message.Message
	resMes.Type=Message.RegisterResMesType

	//返回消息给客户端的话，必须要声明一个RegisterResMes来存消息
	var RegisterResMes  Message.RegisterResMes

	err =model.MyUserDao.Register(&RegisterMes.User)
	if err!=nil{
		if err==model.ERROR_USER_EXITS{
			RegisterResMes.Code=505
			RegisterResMes.Error=model.ERROR_USER_EXITS.Error()

		}else{
			RegisterResMes.Code=506
			RegisterResMes.Error="注册发生未知错误"
		}

	}else{
		RegisterResMes.Code=200
	}


	//将RegisterResMes序列化
	data, err :=json.Marshal(RegisterResMes)
	if err!=nil{
		fmt.Printf("RegisterResMes序列化失败！，err=%v\n",err)
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


	//使用model.MyUserDao到redis中去验证
	user ,err:=model.MyUserDao.CheckLogin(loginMes.UserId,loginMes.UserPSW)

	if err!=nil{
		if err==model.ERROR_USER_NOTEXITS{
			//不合法，用户不存在
			loginResMes.Code=500
			loginResMes.Error=err.Error()
		}else if err==model.ERROR_USER_PSW{
			//不合法，密码错误
			loginResMes.Code=505
			loginResMes.Error=err.Error()
		}else{
			//不合法
			loginResMes.Code=404
			loginResMes.Error="未知错误"
		}
	}else{
		//合法
		loginResMes.Code=200
		//将本次登录成功的用户数据放入userId
		this.UserId=loginMes.UserId
		//将本次成功登录的用户放入userMgr
		userMgr.AddOnlineUsers(this)
		//通知其他用户我上线了
		this.NotifyOtherOnlineUser(loginMes.UserId)
		//将目前在线用户信息用循环的方式放入loginResMes
		for id,_ :=range userMgr.onlineUsers{
			loginResMes.UsersId=append(loginResMes.UsersId,id)
		}
		fmt.Printf("合法用户登录，user=%v\n",user)
	}
	

	//假设id=100, psw=123
	// if loginMes.UserId=="100" && loginMes.UserPSW=="123"{
	// 	//合法
	// 	loginResMes.Code=200
	// }else{
	// 	//不合法，用户不存在
	// 	loginResMes.Code=500
	// 	loginResMes.Error="此账号不存在，请注册后使用"
	// }

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


