package main

import (
	"fmt"
	"net"
	"time"

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"mxshop_srvs/user_srv/global"
	"mxshop_srvs/user_srv/handler"
	"mxshop_srvs/user_srv/initialize"
	"mxshop_srvs/user_srv/proto"
)

func startGRPCServer() {
	//启动监听
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", global.ServerConfig.Host, global.ServerConfig.Port))
	if err != nil {
		panic("failed to listen:" + err.Error())
	}

	server := grpc.NewServer()
	//把用户服务注册到grpc的server中
	proto.RegisterUserServer(server, &handler.UserService{})

	//
	healthCheckServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(server, healthCheckServer)
	healthCheckServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	err = server.Serve(lis)
	if err != nil {
		panic("failed to start grpc:" + err.Error())
	}
}
func registerToConsul() {
	//注册grpc服务健康检查(默认自带的proto，可以去看文档)
	//将grpc服务注册到consul中
	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%d",
		global.ServerConfig.ConsulInfo.Host,
		global.ServerConfig.ConsulInfo.Port)
	client, err := api.NewClient(config)
	if err != nil {
		zap.S().Fatal("consul client init failed")
		panic(err)
	}
	registration := &api.AgentServiceRegistration{
		ID:      global.ServerConfig.Name + "1",
		Name:    global.ServerConfig.Name,
		Address: global.ServerConfig.Host,
		Port:    global.ServerConfig.Port,
		Tags:    []string{"csm", "grpc", "user-srv", "test"},
		Check: &api.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d", global.ServerConfig.Host, global.ServerConfig.Port),
			Timeout:                        "5s",
			Interval:                       "5s",
			DeregisterCriticalServiceAfter: "5s",
		},
	}
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}
}
func main() {
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()
	zap.S().Info(global.ServerConfig)
	//启动grpc服务
	go startGRPCServer()
	time.Sleep(time.Second)
	//等待1s，grpc服务启动后注册到consul中
	registerToConsul()
	//主线程阻塞
	select {}
}
