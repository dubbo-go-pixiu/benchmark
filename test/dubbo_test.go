package test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"dubbo-go-pixiu-benchmark/protocol/dubbo/go-client/pkg"

	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	hessian "github.com/apache/dubbo-go-hessian2"
	. "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/gmeasure"
)

var (
	curPath      string
	userProvider = &pkg.UserProvider{}
	sampleConfig = gmeasure.SamplingConfig{
		N:           100,
		Duration:    120 * time.Second,
		NumParallel: 10,
	}
)

func TestCases(t *testing.T) {
	gomega.RegisterFailHandler(Fail)
	RunSpecs(t, "test")
}

var dubboServerSession, pixiuSession *gexec.Session

var _ = Describe("test", Ordered, func() {
	BeforeAll(func() {
		var err error
		curPath, err = os.Getwd()
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		dubboServerSession = prepareDubboServer()
		pixiuSession = preparePixiu()
		time.Sleep(3 * time.Second)

		prepareDubboClient()
	})

	//TODO(kenwaycai): output to external file

	It("dubbo protocol performance test", func() {
		defer GinkgoRecover()

		experiment := gmeasure.NewExperiment("dubbo protocol performance test")
		AddReportEntry(experiment.Name, experiment)

		experiment.Sample(func(idx int) {
			reqUser := &pkg.User{}
			reqUser.ID = "003"
			experiment.MeasureDuration("GetUser", func() {
				user, err := userProvider.GetUser(context.TODO(), reqUser)
				gomega.Expect(err).To(gomega.Succeed())
				fmt.Printf("consumer:%+v", user)
			})
		}, sampleConfig)

		experiment.Sample(func(idx int) {
			experiment.MeasureDuration("GetGender", func() {
				gender, err := userProvider.GetGender(context.TODO(), 1)
				gomega.Expect(err).To(gomega.Succeed())
				fmt.Printf("consumer:%+v", gender)
			})
		}, sampleConfig)

		experiment.Sample(func(idx int) {
			experiment.MeasureDuration("GetUser0", func() {
				ret, err := userProvider.GetUser0("003", "Moorse")
				gomega.Expect(err).To(gomega.Succeed())
				fmt.Printf("consumer:%+v", ret)
			})
		}, sampleConfig)

		experiment.Sample(func(idx int) {
			experiment.MeasureDuration("GetUsers", func() {
				ret1, err := userProvider.GetUsers([]string{"002", "003"})
				gomega.Expect(err).To(gomega.Succeed())
				fmt.Printf("consumer:%+v", ret1)
			})
		}, sampleConfig)

		experiment.Sample(func(idx int) {
			experiment.MeasureDuration("GetUser2", func() {
				var i int32 = 1
				user, err := userProvider.GetUser2(context.TODO(), i)
				gomega.Expect(err).To(gomega.Succeed())
				fmt.Printf("consumer:%+v", user)
			})
		}, sampleConfig)

		experiment.Sample(func(idx int) {
			experiment.MeasureDuration("GetErr", func() {
				reqUser := &pkg.User{}
				reqUser.ID = "003"
				_, err := userProvider.GetErr(context.TODO(), reqUser)
				gomega.Expect(err).To(gomega.HaveOccurred())
				fmt.Printf("consumer:%+v", err.Error())
			})
		}, sampleConfig)
	})

	It("pixiu to dubbo protocol performance test", func() {

		urlPrefix := "http://localhost:8881/dubbo.io/org.apache.dubbo.sample.UserProvider/%s"

		experiment := gmeasure.NewExperiment("pixiu to dubbo protocol performance test")
		AddReportEntry(experiment.Name, experiment)

		experiment.Sample(func(idx int) {
			experiment.MeasureDuration("GetUser", func() {
				defer GinkgoRecover()

				url := fmt.Sprintf(urlPrefix, "GetUser")
				data := "{\"types\":\"object\",\"values\":{\"id\":\"003\"}}"

				resp, err := http.Post(url, "application/json", strings.NewReader(data))
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(resp.Status, 200)
				_, err = ioutil.ReadAll(resp.Body)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})

			experiment.MeasureDuration("GetGender", func() {
				defer GinkgoRecover()

				url := fmt.Sprintf(urlPrefix, "GetGender")
				data := "{\"types\":\"int\",\"values\":{\"1\"}}"

				resp, err := http.Post(url, "application/json", strings.NewReader(data))
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(resp.Status, 200)
				_, err = ioutil.ReadAll(resp.Body)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})

			experiment.MeasureDuration("GetUser0", func() {
				defer GinkgoRecover()

				url := fmt.Sprintf(urlPrefix, "GetUser0")
				data := "{\"types\":\"string,string\",\"values\":[\"003\", \"Moorse\"]}"

				resp, err := http.Post(url, "application/json", strings.NewReader(data))
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(resp.Status, 200)
				_, err = ioutil.ReadAll(resp.Body)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})

			experiment.MeasureDuration("GetUsers", func() {
				defer GinkgoRecover()

				url := fmt.Sprintf(urlPrefix, "GetUsers")
				data := "{\"types\":\"[]string\",\"values\":[\"002\", \"003\"]}"

				resp, err := http.Post(url, "application/json", strings.NewReader(data))
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(resp.Status, 200)
				_, err = ioutil.ReadAll(resp.Body)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})

		}, sampleConfig)
	})

	AfterAll(func() {
		dubboServerSession.Terminate().Wait()
		pixiuSession.Terminate().Wait(5 * time.Second)
	})

})

func prepareDubboServer() *gexec.Session {
	serverProcess, err := gexec.Build("dubbo-go-pixiu-benchmark/protocol/dubbo/go-server/cmd")
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	command := exec.Command(serverProcess)
	session, err := gexec.Start(command, ioutil.Discard, ioutil.Discard)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	return session
}

func prepareDubboClient() {
	hessian.RegisterJavaEnum(pkg.Gender(pkg.MAN))
	hessian.RegisterJavaEnum(pkg.Gender(pkg.WOMAN))
	hessian.RegisterPOJO(&pkg.User{})

	config.SetConsumerService(userProvider)

	curPath = curPath + "/../protocol/dubbo/go-client/conf/dubbogo.yml"
	err := config.Load(config.WithPath(curPath))
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func preparePixiu() *gexec.Session {
	command := exec.Command("../dist/pixiu", "gateway", "start", "-c", curPath+"/../protocol/dubbo/pixiu/conf/config.yaml")
	//session, err := gexec.Start(command, ioutil.Discard, ioutil.Discard)
	session, err := gexec.Start(command, ioutil.Discard, ioutil.Discard)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	return session
}
