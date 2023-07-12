package Message

import(
	_"fmt"
)


const(
	LoginMesType= "LoginMes"//客户端发送给服务器端的登录信息
	LoginResMesType="LoginResMes" //服务器端发给客户端的登录结果信息
	RegisterMesType="RegisterMes"//注册信息
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
}


type RegisterMes struct{//注册信息

	
}