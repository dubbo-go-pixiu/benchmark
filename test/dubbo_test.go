package test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
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
	userProvider = &pkg.UserProvider{}
)

func TestCases(t *testing.T) {
	gomega.RegisterFailHandler(Fail)
	RunSpecs(t, "test")
}

var session *gexec.Session

var _ = Describe("test", Ordered, func() {
	BeforeAll(func() {
		session = prepareDubboServer()
		time.Sleep(5 * time.Second)
		prepareDubboClient()
	})

	It("test-get-user", func() {
		experiment := gmeasure.NewExperiment("end-to-end web-server performance")
		AddReportEntry(experiment.Name, experiment)

		reqUser := &pkg.User{}
		reqUser.ID = "003"
		experiment.MeasureDuration("GetUser", func() {
			user, err := userProvider.GetUser(context.TODO(), reqUser)
			gomega.Expect(err).To(gomega.Succeed())
			fmt.Printf("consumer:%+v", user)
		})

		// TODO(kenwaycai): output to external file
		//fmt.Println(experiment.String())

	})

	AfterAll(func() {
		session.Terminate().Wait()
	})

})

func prepareDubboServer() *gexec.Session {
	serverProcess, err := gexec.Build("dubbo-go-pixiu-benchmark/protocol/dubbo/go-server/cmd")
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	command := exec.Command(serverProcess)
	session, err := gexec.Start(command, os.Stdout, os.Stderr)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	return session
}

func prepareDubboClient() {
	hessian.RegisterJavaEnum(pkg.Gender(pkg.MAN))
	hessian.RegisterJavaEnum(pkg.Gender(pkg.WOMAN))
	hessian.RegisterPOJO(&pkg.User{})

	config.SetConsumerService(userProvider)

	err := config.Load(config.WithPath("/mnt/d/Workspace/benchmark/protocol/dubbo/go-client/conf/dubbogo.yml"))
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}
