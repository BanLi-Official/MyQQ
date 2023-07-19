package Message



//定义一个用户的结构体

type User struct{
	//确定字段信息
	//为了保证序列化和反序列化的成功，我们必须保证json中的字段与本结构体中的字段一致
	UserId string `json:"userId"` //用户id
	UserPSW string `json:"userPSW"`  //用户密码
	UserName string  `json:"userName"` //用户名字
	UserStatus int `json:"userStatus"`//用户状态

}