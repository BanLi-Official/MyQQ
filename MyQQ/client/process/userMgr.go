package process

import(
	"fmt"
	"gocode/MyQQ/Common/Message"
	"gocode/MyQQ/client/model"
	"encoding/json"
	"gocode/MyQQ/client/utils"

)

//客户端维护的map
var onlineUsers map[string] *Message.User=make(map[string] *Message.User,10)

//当前用户结构体
var CurUser model.CurUser//在用户登录完成后，完成对curuser的初始化

//在客户端显示当前在线的用户
func outputOnlineUser(){
	//遍历onlineUsers
	fmt.Println("当前在线的用户列表(All)：")
	for id,v:=range onlineUsers{
		if v.UserStatus==1{continue} 
		fmt.Printf("用户id:%v\t  用户状态：%v\n",id,v.UserStatus)
	}

}


//在客户端显示除自己以外当前在线的用户
func outputOnlineUser_outSelf(){
	//遍历onlineUsers
	fmt.Println("当前在线的用户列表：")
	for id,v:=range onlineUsers{
		if !checkIsOnline(id){continue} 
		if id==CurUser.UserId{continue}
		fmt.Printf("用户id:%v\t  用户状态：%v\n",id,v.UserStatus)
	}
}







//编写一个方法，处理返回的NotifyUserStatus
func updateUserStatus(notifyUserStatusMes *Message.NotifyUserStatusMes){
	//适当优化，如果存在该账号，那么则只是改变账号状态
	user ,ok:=onlineUsers[notifyUserStatusMes.UserId]
	if !ok{
		user=&Message.User{
			UserId:notifyUserStatusMes.UserId,
		}
	}
	user.UserStatus=notifyUserStatusMes.Status
	onlineUsers[notifyUserStatusMes.UserId]=user

	outputOnlineUser()
}


//查找目前这个用户是否在线
func checkIsOnline(user string)(isOnline bool){
	//遍历onlineUsers
	isOnline=false
	for _,v:=range onlineUsers{
		if v.UserStatus==Message.UserOnline{
			isOnline=true
		}
		
	}

	return
}



//在客户端显示所有用户
func GetAllUser(){

	//先向服务器请求所有用户信息
	//先做一个外包装Mes
	var mes Message.Message
	mes.Type=Message.GetAllUserType

	//在做一个内容物
	var getAllUser Message.GetAllUser
	getAllUser.User=CurUser.UserId

	//将这个内容物序列化存入mes
	data,err:=json.Marshal(getAllUser)
	if err !=nil{
		fmt.Printf("将getAllUser序列化失败，err=%v\n",err)
		return
	}
	mes.Data=string(data)

	//将这个信息序列化
	data ,err =json.Marshal(mes)
	if err !=nil{
		fmt.Printf("将带getAllUser的mes序列化失败，err=%v\n",err)
		return
	}

	//将这个信息发送到服务器
	//先声明一个用于传送数据的变量
	tf:=&utils.Transfer{
		Conn:CurUser.Conn,
	}
	//发送
	err =tf.WritePkg(data)
	if err !=nil{
		fmt.Printf("将带getAllUser的mes发送到服务器失败，err=%v\n",err)
		return
	}


	//处理返回信息
	



	//展示所有用户（包括其用户状态）
}

