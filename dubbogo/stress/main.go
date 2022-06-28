package main

import (
	"dubbo-go-pixiu-benchmark/api"
	"dubbo-go-pixiu-benchmark/helpers"
	"dubbo-go-pixiu-benchmark/stats"
)

import (
	"context"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"sync"
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3/common"
	dubboConstant "dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/common/logger"
	"dubbo.apache.org/dubbo-go/v3/protocol/dubbo"
	invocationImpl "dubbo.apache.org/dubbo-go/v3/protocol/invocation"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var (
	port      = flag.String("port", "50051", "Localhost port to connect to.")
	numRPC    = flag.Int("r", 1, "The number of concurrent RPCs on each connection.")
	warmupDur = flag.Int("w", 10, "Warm-up duration in seconds")
	duration  = flag.Int("d", 60, "Benchmark duration in seconds")
	wg        sync.WaitGroup
	hopts     = stats.HistogramOptions{
		NumBuckets:   2495,
		GrowthFactor: .01,
	}
	mu              sync.Mutex
	hists           []*stats.Histogram
	failInterceptor *helpers.FailInterceptor
	serverSession   *gexec.Session
)

func runWithConn(invoker protocol.Invoker, invCtx context.Context, invoc *invocationImpl.RPCInvocation, warmDeadline, endDeadline time.Time) {
	for i := 0; i < *numRPC; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			res := invoker.Invoke(invCtx, invoc)

			if res.Error() != nil {
				urlStrs := res.Result().([]string)
				ret := make([]*common.URL, 0, len(urlStrs))
				for _, v := range urlStrs {
					tempURL, err := common.NewURL(v)
					if err != nil {
						url := common.URL{}
						logger.Infof("current instance url :", url)
					}
					ret = append(ret, tempURL)
				}
				logger.Infof("final result: ", ret)
			}

			hist := stats.NewHistogram(hopts)
			for {
				start := time.Now()
				if start.After(endDeadline) {
					mu.Lock()
					hists = append(hists, hist)
					mu.Unlock()
					return
				}

				elapsed := time.Since(start)
				if start.After(warmDeadline) {
					hist.Add(elapsed.Nanoseconds())
				}
			}
		}()
	}

}

func main() {

	serverCommand, err := gexec.Build("github.com/dubbo-go-pixiu/benchmark/dubbogo/server/cmd/server")
	Expect(err).ShouldNot(HaveOccurred())
	cmd := exec.Command(serverCommand)
	serverSession, err = gexec.Start(cmd, os.Stdout, os.Stdout)
	Expect(err).ShouldNot(HaveOccurred())

	address := "127.0.0.1:20000"

	url, err := common.NewURL(address+"/org.apache.dubbo.sample.UserProvider",
		common.WithProtocol(dubbo.DUBBO), common.WithParamsValue(dubboConstant.SerializationKey, dubboConstant.Hessian2Serialization),
		common.WithParamsValue(dubboConstant.GenericFilterKey, "true"),
		common.WithParamsValue(dubboConstant.InterfaceKey, ""),
		common.WithParamsValue(dubboConstant.ReferenceFilterKey, "generic,filter"),
		// dubboAttachment must contains group and version info
		common.WithParamsValue(dubboConstant.GroupKey, ""),
		common.WithParamsValue(dubboConstant.VersionKey, ""),
		common.WithPath(dubboConstant.InterfaceKey))

	Eventually(helpers.HealthCheck(address, 30*time.Second, time.Second)).Should(Succeed())

	if err != nil {
		fmt.Println("current url: ", url)
	}
	dubboProtocol := dubbo.NewDubboProtocol()

	args := []reflect.Value{}
	args = append(args, reflect.ValueOf(&api.HelloRequest{Name: "request name"}))
	bizReply := &api.HelloReply{}
	params := []interface{}{"1", "username"}

	invoker := dubboProtocol.Refer(url)
	invoc := invocationImpl.NewRPCInvocationWithOptions(invocationImpl.WithMethodName("GetUser"),
		invocationImpl.WithArguments(params),
		invocationImpl.WithReply(bizReply),
		invocationImpl.WithCallBack(nil), invocationImpl.WithParameterValues(args))

	invCtx := context.Background()
	if err != nil {
		if invoker == nil {
			logger.Errorf("can't connect to upstream server %s with address %s")
		}

		warmDeadline := time.Now().Add(time.Duration(*warmupDur) * time.Second)
		endDeadline := warmDeadline.Add(time.Duration(*duration) * time.Second)

		if endDeadline != time.Now() {

		}
		runWithConn(invoker, invCtx, invoc, warmDeadline, endDeadline)
	}

}
