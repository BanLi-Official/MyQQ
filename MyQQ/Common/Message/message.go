package Message

import(
	_"fmt"
)


const(
	LoginMesType= "LoginMes"//客户端发送给服务器端的登录信息
	LoginResMesType="LoginResMes" //服务器端发给客户端的登录结果信息
	RegisterMesType="RegisterMes"//注册信息
	RegisterResMesType="RegisterResMes"//服务器返回给客户端的的注册结果信息
	NotifyUserStatusMesType="NotifyUserStatusMes"//配合服务器推送用户状态信息
	SmsMesType="SmsMes"//用户发送给服务器的群发信息
	PppMesType="PppMes"//用户发送给特定用户的信息
	GetAllUserType="GetAllUser" //用于向服务器请求所有用户
	GetAllUserResType ="GetAllUserRes"//用于存放服务器找到的所有用户
	PppMes_OffLineType="PppMes_OffLine"//用于发送离线信息的信息内容
)


//定义几个展示用户状态的常量
const(
	UserOnline=iota
	UserOffline
	UserBusyStatus
)



type Message struct{//客户端与服务器端之间进行相互传递信息的载体
	Type string  `json:"type"` //消息类型
	Data string `json:"data"` //消息内容
}


type LoginMes struct{//客户端发送给服务器端的登录信息
	UserId string `json:"userId"` //用户id
	UserPSW string `json:"userPSW"`  //用户密码
	UserName string  `json:"userName"` //用户名字
}

type LoginResMes struct{//服务器端发给客户端的登录结果信息
	Code int `json:"code"` //返回状态码，500表示用户未注册 ，200表示登陆成功
	Error string `json:"error"` //返回错误信息
	UsersId []string `json:"usersId"`  //记录当前在线用户
}


type RegisterMes struct{//注册信息
	User User `json:"user"`
}

type RegisterResMes struct{//注册信息
	Code int `json:"code"` //返回状态码，500表示用户未注册 ，200表示登陆成功
	Error string `json:"error"` //返回错误信息
}


//配合服务器推送用户状态信息
type NotifyUserStatusMes struct{
	UserId string `json:"userId"`
	Status int `json:"status"`
}



//增加一个SmsMes 发送信息
type SmsMes struct{
	Content string `json:"content"`//内容
	User //继承
}


//点对点发送消息
type PppMes struct{
	Content string `json:"content"`//内容
	FromUser string `json:"fromUser"`  //来源于用户
	ToUser string `json:"toUser"`   //发送目标
}



//查询所有用户信息请求
type GetAllUser struct{
	User string `json:"user"`  //请求的用户
}

type GetAllUserRes struct{
	AllUser map[string] *User `json:"allUser"` //返回的用户数组
}


//点对点发送消息
type PppMes_OffLine struct{
	Content string `json:"content"`//内容
	FromUser string `json:"fromUser"`  //来源于用户
	ToUser string `json:"toUser"`   //发送目标
}