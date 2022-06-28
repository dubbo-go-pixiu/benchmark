package main

import (
	"dubbo-go-pixiu-benchmark/api"
)
import (
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/common/logger"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
)

type GreeterProvider struct {
	api.UnimplementedGreeterServer
}

func (s *GreeterProvider) SayHelloStream(svr api.Greeter_SayHelloStreamServer) error {

	//args := make(chan interface{}(constant.AttachmentKey))

	//type attachCtxType string
	//var attachmentKey = attachCtxType(constant.AttachmentKey)
	//
	//
	//fmt.Println("当前上下文: %s",ctx)
	attachments := svr.Context().Value(constant.DubboCtxKey("attachment"))
	logger.Infof("get triple attachment %v", attachments)
	c, err := svr.Recv()
	if err != nil {
		return err
	}
	logger.Infof("Dubbo-go3 GreeterProvider recv 1 user, name = %s\n", c.Name)
	err = svr.Send(&api.User{
		Name: "hello " + c.Name,
		Age:  18,
		Id:   "123456789",
	})
	if err != nil {
		return err
	}
	return nil
}

func main() {
	config.SetProviderService(&GreeterProvider{})

	path := "github.com/dubbo-go-pixiu/benchmark/dubbogo/server/conf/dubbogo.yml"
	if err := config.Load(config.WithPath(path)); err != nil {
		logger.GetLogger().Error(err)
	}

	select {}

}
