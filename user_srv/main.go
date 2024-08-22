package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"net"

	"mxshop_srvs/user_srv/handler"
	"mxshop_srvs/user_srv/proto"
)

func main() {
	//接收命令行输入
	IP := flag.String("ip", "127.0.0.1", "IP 地址")
	Port := flag.Int("port", 50051, "端口号")
	//解析
	flag.Parse()
	fmt.Println("ip:", *IP) //注意是指针类型
	fmt.Println("port:", *Port)
	server := grpc.NewServer()
	//把用户服务注册到grpc的server中
	proto.RegisterUserServer(server, &handler.UserService{})
	//启动监听
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic("failed to listen:" + err.Error())
	}
	err = server.Serve(lis)
	if err != nil {
		panic("failed to start grpc:" + err.Error())
	}

	//service := handler.UserService{}
	//for i := 0; i < 10; i++ {
	//	user, err := service.CreateUser(context.Background(), &proto.CreatUserInfo{
	//		Nickname: "name_" + strconv.Itoa(i),
	//		PassWord: "12345" + strconv.Itoa(i),
	//		Mobile:   "1937211791" + strconv.Itoa(i),
	//	})
	//	if err != nil {
	//		fmt.Println("CreateUser error:", err)
	//		continue
	//	}
	//	fmt.Println("userId:", user.Id)
	//}

}
