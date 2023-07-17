package main

import(

	"github.com/garyburd/redigo/redis"  //引入redis包
	"time"
)

//定义一个全局的pool
var pool *redis.Pool


func InitPool(address string,maxIdle int,maxActive int,idleTimeout time.Duration){
	pool=&redis.Pool{
		MaxIdle:maxIdle,//最大空闲链接数
		MaxActive:maxActive,//表示数据库最大的连接数，0表示没有限制
		IdleTimeout:idleTimeout,  //最大空闲时间
		Dial :func()(redis.Conn,error){//初始化连接的代码，指出连接的ip和端口
			return redis.Dial("tcp",address)
		},
	}
}