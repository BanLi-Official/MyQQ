package processes
import(
	"fmt"
	"net"
)

//一个服务器只需要一个UserMgr，所以可以把它处理成一个全局变量
var (
	userMgr *UserMgr
)


//声明一个结构体，存储一个map用于维护这个map
type UserMgr struct{
	onlineUsers map[string] *UserProcess
}

//完成对UserMgr的初始化工作,程序开始就会自动调用init()
func init() {
	userMgr=&UserMgr{
		onlineUsers : make(map[string] *UserProcess,1024),
	}
}

//完成对onlineUsers的添加
func (this *UserMgr)AddOnlineUsers(up *UserProcess){
	fmt.Printf("up=%v",up)
	this.onlineUsers[up.UserId] = up
}

//完成对onlineUsers的删除
func (this *UserMgr)DelOnlineUsers(userId string){
	delete(this.onlineUsers,userId)
}

//返回所有在线用户
func (this *UserMgr)GetAllOnlineUsers()( map[string] *UserProcess){
	return this.onlineUsers
}

//按照id返回用户
func (this *UserMgr)GetOnlineUserById(userId string)(up  *UserProcess,err error){
	up,ok:=this.onlineUsers[userId]
	if !ok{
		err=fmt.Errorf("用户%v不存在",userId)
		return
	}
	return
}


//按照conn返回用户
func (this *UserMgr)GetOnlineUserId(conn net.Conn)(userid string,err error){
	for _,up :=range this.onlineUsers{
		if conn==up.Conn{
			userid=up.UserId
			return 
		}
	}
	err=fmt.Errorf("无法通过conn找到该用户")
	return
}


func (this *UserMgr)ShowOnlineUser()(){
	fmt.Println("当前在线客户如下：")
	for _,up :=range this.onlineUsers{
		fmt.Printf("用户id：%s\n",up.UserId)
	}

}