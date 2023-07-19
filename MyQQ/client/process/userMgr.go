package process

import(
	"fmt"
	"gocode/MyQQ/Common/Message"
	"gocode/MyQQ/client/model"
)

//客户端维护的map
var onlineUsers map[string] *Message.User=make(map[string] *Message.User,10)
//当前用户结构体
var CurUser model.CurUser//在用户登录完成后，完成对curuser的初始化

//在客户端显示当前在线的用户
func outputOnlineUser(){
	//遍历onlineUsers
	fmt.Println("当前在线的用户列表：")
	for id,_:=range onlineUsers{
		fmt.Println("用户id:\t",id)
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