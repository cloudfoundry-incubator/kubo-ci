package master_failure_test

import (
	"crypto/tls"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"tests/test_helpers"

	"io/ioutil"
	"net/http"

	"github.com/cppforlife/turbulence/incident"
	"github.com/cppforlife/turbulence/incident/selector"
	"github.com/cppforlife/turbulence/tasks"
)

var _ = Describe("A single master failure", func() {
	Specify("Master VM is resurrected within 10 minutes", func() {
		boshDirector := test_helpers.NewDirector()
		deployment, err := boshDirector.FindDeployment(test_helpers.DeploymentName)
		Expect(err).NotTo(HaveOccurred())
		countRunningApiServerOnMaster := test_helpers.CountProcessesOnVmsOfType(deployment, test_helpers.MasterVmType, "kubernetes-api", test_helpers.VmRunningState)

		Expect(countRunningApiServerOnMaster()).To(Equal(1))

		By("Deleting the Master VM")
		hellRaiser := test_helpers.TurbulenceClient()
		killOneMaster := incident.Request{
			Selector: selector.Request{
				Deployment: &selector.NameRequest{
					Name: test_helpers.DeploymentName,
				},
				Group: &selector.NameRequest{
					Name: test_helpers.MasterVmType,
				},
				ID: &selector.IDRequest{
					Limit: selector.MustNewLimitFromString("1"),
				},
			},
			Tasks: tasks.OptionsSlice{
				tasks.KillOptions{},
			},
		}
		incident := hellRaiser.CreateIncident(killOneMaster)
		incident.Wait()
		Expect(countRunningApiServerOnMaster()).Should(Equal(0))

		By("Waiting for resurrection")
		Eventually(countRunningApiServerOnMaster, "10m", "20s").Should(Equal(1))

		By("Setting up SSH")
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		result, err := client.Get(fmt.Sprintf("https://%s:8443/healthz", test_helpers.GetMasterIP(deployment)))
		Expect(err).ToNot(HaveOccurred())
		Expect(result.StatusCode).To(Equal(http.StatusOK))
		response, err := ioutil.ReadAll(result.Body)
		Expect(err).ToNot(HaveOccurred())
		Expect(string(response)).To(Equal("ok"))
	})
})
