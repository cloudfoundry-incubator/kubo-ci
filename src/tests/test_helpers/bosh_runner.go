package test_helpers

import (
	"fmt"
	"os/exec"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

const EtcdHostname = "master-0.etcd.cfcr.internal"

func RunEtcdCommandFromWorker(deployment, workerID string, args ...string) string {
	remoteArgs := []string{
		"/var/vcap/packages/etcdctl/etcdctl",
		fmt.Sprintf("--endpoints https://%s:2379", EtcdHostname),
		"--cert-file", "/var/vcap/jobs/flanneld/config/etcd-client.crt",
		"--key-file", "/var/vcap/jobs/flanneld/config/etcd-client.key",
		"--ca-file", "/var/vcap/jobs/flanneld/config/etcd-ca.crt",
	}
	remoteArgs = append(remoteArgs, args...)
	s := RunSSHWithDeployment(deployment, "worker/"+workerID, fmt.Sprintf("sudo su -c '%s'", strings.Join(remoteArgs, " ")))
	Eventually(s, "20s", "1s").Should(gexec.Exit())
	ss := string(s.Out.Contents())
	return ss
}

func RunEtcdCommandFromMasterWithFullPrivilege(deployment, masterID string, args ...string) string {
	remoteArgs := []string{
		"/var/vcap/packages/etcdctl/etcdctl",
		fmt.Sprintf("--endpoints https://%s:2379", EtcdHostname),
		"--cert-file", "/var/vcap/jobs/etcd/config/etcdctl.crt",
		"--key-file", "/var/vcap/jobs/etcd/config/etcdctl.key",
		"--ca-file", "/var/vcap/jobs/etcd/config/etcdctl-ca.crt",
	}
	remoteArgs = append(remoteArgs, args...)
	s := RunSSHWithDeployment(deployment, "master/"+masterID, fmt.Sprintf("sudo su -c '%s'", strings.Join(remoteArgs, " ")))
	Eventually(s, "20s", "1s").Should(gexec.Exit())
	ss := string(s.Out.Contents())
	return ss
}

func RunSSHWithDeployment(deploymentName, instance string, args ...string) *gexec.Session {
	nargs := []string{"-d", deploymentName, "ssh", "--opts=-q",
		instance,
	}
	nargs = append(nargs, args...)

	return RunCommand("bosh", nargs...)
}

func RunCommand(cmd string, args ...string) *gexec.Session {
	c1 := exec.Command(cmd, args...)
	session, err := gexec.Start(c1, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())

	return session
}
