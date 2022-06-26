package main

import (
	"context"
	"dubbo-go-pixiu-benchmark/api"
	"dubbo-go-pixiu-benchmark/helpers"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"sync"
	"time"
)

import (
	"dubbo-go-pixiu-benchmark/dubbogo/pkg"

	"dubbo-go-pixiu-benchmark/stats"

	"dubbo.apache.org/dubbo-go/v3/common"
	dubboConstant "dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/common/logger"
	"dubbo.apache.org/dubbo-go/v3/protocol/dubbo"
	invocationImpl "dubbo.apache.org/dubbo-go/v3/protocol/invocation"
)

var (
	port      = flag.String("port", "50051", "Localhost port to connect to.")
	numRPC    = flag.Int("r", 1, "The number of concurrent RPCs on each connection.")
	numConn   = flag.Int("c", 1, "The number of parallel connections.")
	warmupDur = flag.Int("w", 10, "Warm-up duration in seconds")
	duration  = flag.Int("d", 60, "Benchmark duration in seconds")
	rqSize    = flag.Int("req", 1, "Request message size in bytes.")
	rspSize   = flag.Int("resp", 1, "Response message size in bytes.")
	rpcType   = flag.String("rpc_type", "unary",
		`Configure different stress rpc type. Valid options are:
		   unary;
		   streaming.`)
	testName = flag.String("test_name", "", "Name of the test used for creating profiles.")
	wg       sync.WaitGroup
	hopts    = stats.HistogramOptions{
		NumBuckets:   2495,
		GrowthFactor: .01,
	}
	mu    sync.Mutex
	hists []*stats.Histogram
	failInterceptor   *helpers.FailInterceptor
	banyandSession    *gexec.Session
	rootPath          string
	deferRootPathFunc func()
)

func runWithConn(req *pkg.StressRequest, warmDeadline, endDeadline time.Time) {
	for i := 0; i < *numRPC; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

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

	req := &pkg.StressRequest{
		ResponseType: 0,
		ResponseSize: int32(*rspSize),
		Payload: &pkg.Payload{
			Type: pkg.PayloadType_COMPRESSABLE,
			Body: make([]byte, *rqSize),
		},
	}
	banyandbBinary, err := gexec.Build("github.com/dubbo-go-pixiu/benchmark/dubbogo/server/cmd/server.go");
	Expect(err).ShouldNot(HaveOccurred())
	cmd := exec.Command(banyandBinary)
	banyandSession, err = gexec.Start(cmd, os.Stdout, os.Stdout)
	Expect(err).ShouldNot(HaveOccurred())


	url, err := common.NewURL("127.0.0.1:20000/org.apache.dubbo.sample.UserProvider",
		common.WithProtocol(dubbo.DUBBO), common.WithParamsValue(dubboConstant.SerializationKey, dubboConstant.Hessian2Serialization),
		common.WithParamsValue(dubboConstant.GenericFilterKey, "true"),
		common.WithParamsValue(dubboConstant.InterfaceKey, ""),
		common.WithParamsValue(dubboConstant.ReferenceFilterKey, "generic,filter"),
		// dubboAttachment must contains group and version info
		common.WithParamsValue(dubboConstant.GroupKey, ""),
		common.WithParamsValue(dubboConstant.VersionKey, ""),
		common.WithPath(dubboConstant.InterfaceKey))

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
	}

	r := req
	if r == nil {

	}

	warmDeadline := time.Now().Add(time.Duration(*warmupDur) * time.Second)
	endDeadline := warmDeadline.Add(time.Duration(*duration) * time.Second)

	if endDeadline != time.Now() {

	}
	runWithConn(req,warmDeadline,endDeadline)

}
