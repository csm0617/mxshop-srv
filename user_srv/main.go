package main

import (
	"fmt"
	"google.golang.org/grpc/health"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"

	"mxshop_srvs/user_srv/global"
	"mxshop_srvs/user_srv/handler"
	"mxshop_srvs/user_srv/initialize"
	"mxshop_srvs/user_srv/proto"
)

func main() {
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitConfig()

	server := grpc.NewServer()
	//把用户服务注册到grpc的server中
	proto.RegisterUserServer(server, &handler.UserService{})
	//启动监听
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", global.ServerConfig.Host, global.ServerConfig.Port))
	if err != nil {
		panic("failed to listen:" + err.Error())
	}
	//注册grpc服务健康检查(默认自带的proto，可以去看文档)
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	err = server.Serve(lis)
	if err != nil {
		panic("failed to start grpc:" + err.Error())
	}
}
