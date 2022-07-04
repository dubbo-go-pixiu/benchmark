package client

import (
	"dubbo-go-pixiu-benchmark/pixiu/dubbo/go-client/pkg"
)

import (
	"context"
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3/common/logger"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	hessian "github.com/apache/dubbo-go-hessian2"
)

var (
	userProvider = &pkg.UserProvider{}
)

type FenMan struct {
	id  string "识别的身份ID"
	age int64  "当前年龄"
}

// need to setup environment variable "DUBBO_GO_CONFIG_PATH" to "conf/dubbogo.yml" before run
func main() {
	hessian.RegisterJavaEnum(pkg.Gender(pkg.MAN))
	hessian.RegisterJavaEnum(pkg.Gender(pkg.WOMAN))

	//ID   string `hessian:"id"`
	//Name string
	//Age  int32
	//Time time.Time
	//Sex  Gender // not
	hessian.RegisterPOJO(&pkg.User{
		ID: "123", Name: "windwheel",
		Age: 24, Time: time.Now(), Sex: pkg.Gender(pkg.MAN)})

	config.SetConsumerService(userProvider)

	path := "github.com/dubbo-go-pixiu/benchmark/pixiu/dubbo/go-client/conf/dubbogo.yml"

	err := config.Load(config.WithPath(path))
	if err != nil {
		logger.Error(err)
	}

	logger.Infof("\n\ntest")
	test()
}

func test() {
	logger.Infof("\n\n\nstart to test dubbo")
	reqUser := &pkg.User{}
	reqUser.ID = "003"
	user, err := userProvider.GetUser(context.TODO(), reqUser)
	if err != nil {
		logger.Error(err)
	}
	logger.Infof("response result: %v", user)

	logger.Infof("\n\n\nstart to test dubbo - enum")
	gender, err := userProvider.GetGender(context.TODO(), 1)
	if err != nil {
		logger.Error(err)
	}
	logger.Infof("response result: %v", gender)

	logger.Infof("\n\n\nstart to test dubbo - GetUser0")
	ret, err := userProvider.GetUser0("003", "Moorse")
	if err != nil {
		logger.Error(err)
	}
	logger.Infof("response result: %v", ret)

	logger.Infof("\n\n\nstart to test dubbo - GetUsers")
	ret1, err := userProvider.GetUsers([]string{"002", "003"})
	if err != nil {
		logger.Error(err)
	}
	logger.Infof("response result: %v", ret1)

	logger.Infof("\n\n\nstart to test dubbo - getUser")

	var i int32 = 1
	user, err = userProvider.GetUser2(context.TODO(), i)
	if err != nil {
		logger.Error(err)
	}
	logger.Infof("response result: %v", user)

	logger.Infof("\n\n\nstart to test dubbo - getErr")
	reqUser.ID = "003"
	_, err = userProvider.GetErr(context.TODO(), reqUser)
	if err == nil {
		logger.GetLogger().Errorf("err is nil")
	}
	logger.Infof("getErr - error: %v", err)
}
