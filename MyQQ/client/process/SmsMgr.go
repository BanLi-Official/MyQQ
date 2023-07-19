package process

import(
	"fmt"
	"gocode/MyQQ/Common/Message"
	"encoding/json"
)

func outputGroupMes(mes *Message.Message){
	//显示该信息,反序列化mes.Data

	var smsMes Message.SmsMes
	err:=json.Unmarshal([]byte(mes.Data),&smsMes)
	if err!= nil{
		fmt.Printf("json.Unmarshal Error=%v\n",err)
		return
	}


	//显示信息
	info :=fmt.Sprintf("用户id：\t%s 对大家说：\t%s",smsMes.UserId,smsMes.Content)
	fmt.Println(info)
	fmt.Println()

}