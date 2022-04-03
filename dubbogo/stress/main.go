package main

import (
	"dubbo-go-pixiu-benchmark/dubbogo/pkg"
	"dubbo-go-pixiu-benchmark/stats"
	"dubbo.apache.org/dubbo-go/v3/common/logger"
	"flag"
	"sync"
	"time"
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

				//TODO 编写dubbo proxy服务调用

				elapsed := time.Since(start)
				if start.After(warmDeadline) {
					hist.Add(elapsed.Nanoseconds())
				}
			}
		}()
	}
}

func main() {

	flag.Parse()
	if *testName == "" {
		logger.Fatalf("test_name not set")
	}

	req := &pkg.StressRequest{
		ResponseType: 0,
		ResponseSize: int32(*rspSize),
		Payload: &pkg.Payload{
			Type: pkg.PayloadType_COMPRESSABLE,
			Body: make([]byte, *rqSize),
		},
	}

	//FIXME 压测器解析dubbo协议请求连接 执行请求
	r := req
	if r == nil {

	}

	//warmDeadline := time.Now().Add(time.Duration(*warmupDur) * time.Second)
	//endDeadline := warmDeadline.Add(time.Duration(*duration) * time.Second)

}
