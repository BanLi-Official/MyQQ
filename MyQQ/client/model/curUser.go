package model

import(
	"net"
	"gocode/MyQQ/Common/Message"
)

//在客户端很多地方会用到curUser，所以把它变成一个全局变量
type CurUser struct{
	Conn net.Conn
	Message.User
}