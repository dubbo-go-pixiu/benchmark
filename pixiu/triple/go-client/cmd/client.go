package main

import (
	"context"
	"sync"
)

import (
	"dubbo-go-pixiu-benchmark/api"

	"dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/common/logger"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
)

var grpcGreeterImpl = new(api.GreeterClientImpl)

func init() {
	config.SetConsumerService(grpcGreeterImpl)
}

// export DUBBO_GO_CONFIG_PATH= PATH_TO_SAMPLES/helloworld/go-client/conf/dubbogo.yml
func main() {
	path := "/Users/windwheel/Documents/gitrepo/dubbo-go-triple-demo/context/triple/go-client/conf/dubbogo.yml"

	config.SetConsumerService(grpcGreeterImpl)
	err := config.Load(config.WithPath(path))
	if err != nil {
		panic(err)
	}

	logger.Info("start to test dubbo")

	// set user defined context attachment, map value can be string or []string, otherwise it is not to be transferred
	userDefinedValueMap := make(map[string]interface{})
	userDefinedValueMap["key1"] = "user defined value 1"
	userDefinedValueMap["key2"] = "user defined value 2"
	userDefinedValueMap["key3"] = []string{"user defined value 3.1", "user defined value 3.2"}
	userDefinedValueMap["key4"] = []string{"user defined value 4.1", "user defined value 4.2"}

	ctx := context.WithValue(context.Background(), constant.DubboCtxKey("attachment"), userDefinedValueMap)
	//测试流式 rpc
	logger.Info("start to test dubbo stream context")
	request := &api.HelloRequest{
		Name: "laurence",
	}
	stream, err := grpcGreeterImpl.SayHelloStream(ctx)
	if err != nil {
		logger.Error(err)
	}
	// stream grpc双向流式发送
	err = stream.Send(request)
	if err != nil {
		logger.Error(err)
	}
	logger.Infof("client stream send request: %v\n", request)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		reply, err := stream.Recv()
		if err != nil {
			logger.Error(err)
		}
		logger.Infof("client stream received result: %v\n", reply)
	}()
	wg.Wait()
}
