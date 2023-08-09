package model

import(
	"fmt"
	"github.com/garyburd/redigo/redis"  //引入redis包
	"gocode/MyQQ/Common/Message"

	"encoding/json"
	
)
//在服务器启动之后就初始化一个userDAO实例，把他做成全局变量，在需要和redis进行操作的时候就直接使用即可
var (
	MyUserDao *UserDao
)


//定义一个userDao的结构体
//完成对user的各种操作

type UserDao struct{
	pool *redis.Pool
}



//使用工厂模式创建一个userDao实例
func NewUserDao(pool *redis.Pool)(userDao *UserDao){
	userDao = &UserDao{
		pool :pool,
	}
	return
}




//1.根据用户id，返回一个用户实例＋errors
func (this *UserDao)GetUserById(conn redis.Conn,id string)(user *User,err error){
	//通过给定的id去查找用户
	res,err:=redis.String(conn.Do("hget","user",id))
	if err !=nil{
		fmt.Printf("conn.Do(hget,users,id)失败 err=%v，id=%v",err,id)
		if err ==redis.ErrNil{//表示在redis库中没有找到该id的用户
			err = ERROR_USER_NOTEXITS
		}
		return 
	}

	user =&User{}
	//将查到的对象反序列化
	err= json.Unmarshal([]byte(res),user)
	if err !=nil{
		fmt.Printf("在库中查到的对象反序列化失败 err=%v",err)
		return 
	}
	return 
}


//2.检查传来的用户名与密码是否正确
func (this *UserDao) CheckLogin(userId string,userPsw string)(user *User,err error){
	//先从数据库连接池中取出一个连接
	conn := this.pool.Get()
	defer conn.Close()

	//获取数据库中的信息
	user ,err =this.GetUserById(conn,userId)

	if err!= nil{
		fmt.Printf("获取数据库信息错误 GetUserById Error=%v\n",err)
		return
	}
	if user.UserPSW!=userPsw{
		fmt.Printf("CheckLogin（）密码错误\n")
		err=ERROR_USER_PSW
		return
	}
	return
}



//3.将传来的信息入库
func (this *UserDao) Register(user *Message.User)(err error){
	//先从数据库连接池中取出一个连接
	conn := this.pool.Get()
	defer conn.Close()

	//在数据库中查这个id，看看有没有已经被注册
	_,err = this.GetUserById(conn,user.UserId)
	if err == nil{
		fmt.Println("客户端传来的id已经被注册了")
		err=ERROR_USER_EXITS
		return 
	}

	//将user序列化
	data ,err :=json.Marshal(user)
	if err!= nil{
		fmt.Printf("user序列化错误 Error=%v\n",err)
		return
	}
	
	//将序列化后的信息传入redis中
	_,err=conn.Do("Hset","user",user.UserId,string(data))
	if err!= nil{
		fmt.Printf("将序列化后的用户注册信息传入redis中错误 Error=%v\n",err)
		return
	}
	return

}

//4.查找数据库中的所有用户
func (this *UserDao) GetAllUser_From_Redis()(err error,AllUsers map[string] *Message.User){
	//先从数据库连接池中取出一个连接
	conn := this.pool.Get()
	defer conn.Close()

	//通过给定的id去查找用户
	res,err:=redis.Strings(conn.Do("hgetall","user"))
	//把这些内容存到一个map中去
	//var AllUsers map[string] *Message.User
	AllUsers=make(map[string] *Message.User)
	var str string

	for i,v:= range  res{
		//fmt.Printf("res[%v]=%v\n",i,v)
		if i%2==0{
			str=v
		}else{
			
			var user Message.User
			err=json.Unmarshal([]byte(v),&user)
			//fmt.Printf("v=%v\n",v) 
			//fmt.Printf("user_t=%v\n",user) 


			AllUsers[str]=&user
		}


	}
	//fmt.Println(AllUsers)

	//返回给上一级


	
	return
}





//3.将传来的信息入库
func (this *UserDao) InsertOffLineMassage(from string ,to string ,contain string)(err error){
	//先从数据库连接池中取出一个连接
	conn := this.pool.Get()
	defer conn.Close()

	//
	
	

}