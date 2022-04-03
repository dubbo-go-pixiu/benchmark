package main

import (
	"context"
	"dubbo-go-pixiu-benchmark/dubbogo/main/pkg"
	"dubbo.apache.org/dubbo-go/v3/config"
	"fmt"
	hessian "github.com/apache/dubbo-go-hessian2"
)

var (
	userProvider = &pkg.UserProvider1{}
)

func init() {
	config.SetConsumerService(userProvider)
	hessian.RegisterPOJO(&pkg.User{})
}

func main() {

	path := "/Users/windwheel/Documents/gitrepo/dubbo-go-benchmark/3.0/dubbogo/stress/dubbogo.yml"
	err := config.Load(config.WithPath(path))

	if err != nil {
		panic(err)
	}

	user, err := userProvider.GetUser(context.TODO(), &pkg.User{ID: "1113333", Name: "chengxingyuan", Code: 111, Age: 111})

	if err != nil {
		fmt.Println("打印当前客户端调用结果: %s", user)
	}

}
