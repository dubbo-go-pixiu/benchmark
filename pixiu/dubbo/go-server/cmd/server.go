package server

import (
	"dubbo-go-pixiu-benchmark/pixiu/dubbo/go-server/pkg"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3/common/logger"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	hessian "github.com/apache/dubbo-go-hessian2"
)

var (
	survivalTimeout = int(3e9)
)

// need to setup environment variable "DUBBO_GO_CONFIG_PATH" to "conf/dubbogo.yml" before run
func main() {

	// ------for hessian2------
	hessian.RegisterJavaEnum(pkg.Gender(pkg.MAN))
	hessian.RegisterJavaEnum(pkg.Gender(pkg.WOMAN))
	hessian.RegisterPOJO(&pkg.User{})
	config.SetProviderService(&pkg.UserProvider{})
	//config.SetProviderService(&pkg.UserProvider1{})
	//config.SetProviderService(&pkg.UserProvider2{})
	//config.SetProviderService(&pkg.ComplexProvider{})
	//config.SetProviderService(&pkg.WrapperArrayClassProvider{})
	// ------------
	path := "/Users/windwheel/Documents/gitrepo/dubbo-go-triple-demo/pixiu/dubbo/go-server/conf/dubbogo.yml"

	if err := config.Load(config.WithPath(path)); err != nil {
		panic(err)
	}

	initSignal()
}

func initSignal() {
	signals := make(chan os.Signal, 1)
	// It is not possible to block SIGKILL or syscall.SIGSTOP
	signal.Notify(signals, os.Interrupt, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		sig := <-signals
		logger.Infof("get signal %s", sig.String())
		switch sig {
		case syscall.SIGHUP:
			// reload()
		default:
			time.AfterFunc(time.Duration(survivalTimeout), func() {
				logger.Warnf("app exit now by force...")
				os.Exit(1)
			})

			// The program exits normally or timeout forcibly exits.
			fmt.Println("provider app exit now...")
			return
		}
	}
}
