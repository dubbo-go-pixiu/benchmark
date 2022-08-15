package test

import (
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/gmeasure"
	"io/ioutil"
	"os/exec"
	"time"
)

var (
	CurPath      string
	SampleConfig = gmeasure.SamplingConfig{
		N:           100,
		Duration:    120 * time.Second,
		NumParallel: 10,
	}
)

func PreparePixiu(path string) *gexec.Session {
	command := exec.Command("../dist/pixiu", "gateway", "start", "-c", path)
	session, err := gexec.Start(command, ioutil.Discard, ioutil.Discard)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	return session
}
