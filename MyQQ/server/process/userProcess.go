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
	for id ,up :=range userMgr.onlineUsers{
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
	mes.Type=Message.NotifyUserStatusMesType

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




//编写一个方法用于删除用户,并通知用户已下线
func (this *UserProcess) ServerProcessDelete(conn net.Conn)(err error){
	var userid_string string
	userid_string,err=userMgr.GetOnlineUserId(conn)
	if err !=nil{
		fmt.Printf("服务器端删除角色的时候GetOnlineUserId错误 ，err=%v\n",err)
		return
	}
	userMgr.DelOnlineUsers(userid_string)
	fmt.Printf("离线客户已经删除，当前在线客户如下：\n")
	userMgr.ShowOnlineUser()


	//通知所有人
	//遍历onlineusers，然后一个个的发送
	for _,up :=range userMgr.onlineUsers{
		//开始通知，单独写一个方法
		up.NotifySBOffline(userid_string)
	}
	return

}


//通知有人下线
func (this *UserProcess) NotifySBOffline(userId string){
	//组装NotifyUserStatusMes
	//先声明一个M.M，用于存储发过去的信息
	var mes Message.Message
	mes.Type=Message.NotifyUserStatusMesType

	//返回消息给客户端的话，必须要声明一个NotifyUserStatusMesType来存消息
	var notifyUserStatusMes Message.NotifyUserStatusMes
	notifyUserStatusMes.UserId=userId
	notifyUserStatusMes.Status=Message.UserOffline
	fmt.Println("notifyUserStatusMes=",notifyUserStatusMes)

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



//编写一个GetUsers函数专门用来处理获取所有用户的请求
func (this *UserProcess) GetUsers(mes *Message.Message)(err error){
	//提取mes中的data信息
	var getAllUser Message.GetAllUser
	err=json.Unmarshal([]byte(mes.Data),&getAllUser)
	//展示data信息
	fmt.Println(getAllUser)

	//获取所有用户
	var AllUsers map[string] *Message.User
	err,AllUsers=model.MyUserDao.GetAllUser_From_Redis()
	if err!=nil{
		fmt.Printf("model.MyUserDao.GetAllUser_From_Redis()失败！，err=%v\n",err)
	}
	//fmt.Println(AllUsers)

	//发送给对面
	//2023年8月8日16:35:25  
	//	完成了客户端通知服务器端进行查找所有用户的行为，完成了服务器端从redis中查找所有用户，并展示找到的用户，
	//  下面应该将这个结果发送给对面，客户端展示所有用户与其状态（目前还有一个问题：所有用户的状态还没有写，所有人的状态都是0，后面在服务器端判断它的用户状态），用户选择
	//发送对象，将数据发送到服务器，服务器判断是否在线，不在线则进行信息存储


	//先向服务器请求所有用户信息
	//先做一个外包装Mes
	var mesRes Message.Message
	mesRes.Type=Message.GetAllUserResType

	//在做一个内容物
	var getAllUserRes Message.GetAllUserRes
	getAllUserRes.AllUser=AllUsers

	//将这个内容物序列化存入mes
	data,err:=json.Marshal(getAllUserRes)
	if err !=nil{
		fmt.Printf("将getAllUserRes序列化失败，err=%v\n",err)
		return
	}
	mesRes.Data=string(data)


	//将这个信息序列化
	data ,err =json.Marshal(mesRes)
	if err !=nil{
		fmt.Printf("将带getAllUserRes的mesRes序列化失败，err=%v\n",err)
		return
	}


	//将这个信息发送到客户端
	//先声明一个用于传送数据的变量
	tf:=&utils.Transfer{
		Conn:this.Conn,
	}
	//发送
	err =tf.WritePkg(data)
	if err !=nil{
		fmt.Printf("将带getAllUserRes的mes发送到客户端失败，err=%v\n",err)
		return
	}







	return
}

