package main

import (
	"context"
	"os"
)

import (
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	hessian "github.com/apache/dubbo-go-hessian2"

	"github.com/dubbogo/gost/log"
)

type UserProvider struct {
	GetContext func(ctx context.Context) (rsp *ContextContent, err error)
}

type ContextContent struct {
	Path              string
	InterfaceName     string
	DubboVersion      string
	LocalAddr         string
	RemoteAddr        string
	UserDefinedStrVal string
	CtxStrVal         string
	CtxIntVal         int64
}

func (c *ContextContent) JavaClassName() string {
	return "org.apache.dubbo.ContextContent"
}

// need to setup environment variable "CONF_CONSUMER_FILE_PATH" to "conf/client.yml" before run
func main() {
	var userProvider = &UserProvider{}
	config.SetConsumerService(userProvider)
	hessian.RegisterPOJO(&ContextContent{})
	if err := config.Load(); err != nil {
		panic(err)
	}
	gxlog.CInfo("\n\n\nstart to test dubbo")

	atta := make(map[string]interface{})
	atta["string-value"] = "string-demo"
	atta["int-value"] = 1231242
	atta["user-defined-value"] = &ContextContent{InterfaceName: "test.interface.name"}
	reqContext := context.WithValue(context.Background(), constant.DubboCtxKey("attachment"), atta)
	rspContent, err := userProvider.GetContext(reqContext)
	if err != nil {
		gxlog.CError("error: %v\n", err)
		os.Exit(1)
		return
	}
	gxlog.CInfo("response result: %+v\n", rspContent)
}
