package persistent_volume_test

import (
	"fmt"
	"os"

	"github.com/go-redis/redis"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"integration-tests/test_helpers"
)

const (
	testKey   = "Earth"
	testValue = "The gold grows greed like a proud scallywag."
)

var _ = Describe("Pod Storage", func() {
	It("should not be erased on recreation", func() {

		runner := test_helpers.NewKubectlRunner()
		runner.RunKubectlCommand("create", "namespace", runner.Namespace()).Wait("60s")

		By("Installing redis")
		runner.ExpectEventualSuccess("apply", "-f", test_helpers.PathFromRoot("specs/redis-master-storage.yml"))
		runner.ExpectEventualSuccess("apply", "-f", test_helpers.PathFromRoot("specs/redis-master.yml"))
		runner.ExpectEventualSuccess("rollout", "status", "deployment/redis-master", "-w")

		workerIP := os.Getenv("WORKER_IP")
		client := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:30404", workerIP),
			Password: "",
			DB:       0,
		})

		By("Writing the value")
		expectEventualSuccess(client.Ping().Result)
		expectImmediateSuccess(client.Set(testKey, testValue, 0).Result)
		expectImmediateSuccess(client.Save().Result)

		By("Re-installing redis")
		runner.ExpectEventualSuccess("delete", "-f", test_helpers.PathFromRoot("specs/redis-master.yml"))
		runner.ExpectEventualSuccess("apply", "-f", test_helpers.PathFromRoot("specs/redis-master.yml"))
		runner.ExpectEventualSuccess("rollout", "status", "deployment/redis-master", "-w")

		By("Verifying the main assertion")
		expectEventualSuccess(client.Ping().Result)
		value, err := client.Get(testKey).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(Equal(testValue))

		By("Tearing down")
		runner.ExpectEventualSuccess("delete", "-f", test_helpers.PathFromRoot("specs/redis-master.yml"))
	})
})

func expectEventualSuccess(f func() (string, error)) {
	Eventually(func() error { _, err := f(); return err }).ShouldNot(HaveOccurred())
}

func expectImmediateSuccess(f func() (string, error)) {
	_, err := f()
	Expect(err).ToNot(HaveOccurred())
}
