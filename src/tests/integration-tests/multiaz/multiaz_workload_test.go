package multiaz

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"gopkg.in/yaml.v2"
)

var _ = Describe("Multi-AZ workload deployment", func() {
	BeforeEach(func() {
		deployNginx := runner.RunKubectlCommand("create", "-f", nginxSpec)
		Eventually(deployNginx, "60s").Should(gexec.Exit(0))

		rolloutWatch := runner.RunKubectlCommand("rollout", "status", "deployment/nginx", "-w")
		Eventually(rolloutWatch, "120s").Should(gexec.Exit(0))
	})

	AfterEach(func() {
		session := runner.RunKubectlCommand("delete", "-f", nginxSpec)
		session.Wait("60s")
	})

	FIt("deploys three pods across three azs", func() {
		nodeList := runner.RunKubectlCommand("get", "nodes", "-o", "yaml")
		Eventually(nodeList, "60s").Should(gexec.Exit(0))
		nodeZones, err := getNodeZones(nodeList.Out.Contents())
		Expect(err).NotTo(HaveOccurred())

		podList := runner.RunKubectlCommand("get", "pods", "-l", "app=nginx", "-o", "yaml")
		Eventually(podList, "60s").Should(gexec.Exit(0))
		podNodes, err := getPodInstanceNodes(podList.Out.Contents())
		Expect(err).NotTo(HaveOccurred())
		Expect(podNodes).To(HaveLen(3))

		pod1AZ := nodeZones[podNodes[0]]
		pod2AZ := nodeZones[podNodes[1]]
		pod3AZ := nodeZones[podNodes[2]]

		Expect(pod1AZ).NotTo(Equal(pod2AZ))
		Expect(pod1AZ).NotTo(Equal(pod3AZ))
		Expect(pod2AZ).NotTo(Equal(pod3AZ))
	})
})

func getNodeZones(nodeDescriptionYAML []byte) (map[string]string, error) {
	var nodeDescription struct {
		Items []struct {
			Metadata struct {
				Name   string            `yaml:"name"`
				Labels map[string]string `yaml:"labels"`
			} `yaml:"metadata"`
		} `yaml:"items"`
	}

	err := yaml.Unmarshal(nodeDescriptionYAML, &nodeDescription)
	if err != nil {
		return nil, err
	}

	nodeZoneMap := map[string]string{}
	for _, item := range nodeDescription.Items {
		for labelKey, labelValue := range item.Metadata.Labels {
			if labelKey == "failure-domain.beta.kubernetes.io/zone" {
				nodeZoneMap[item.Metadata.Name] = labelValue
			}
		}
	}

	return nodeZoneMap, nil
}

func getPodInstanceNodes(podDescriptionYAML []byte) ([]string, error) {
	var podDescription struct {
		Items []struct {
			Spec struct {
				NodeName string `yaml:"nodeName"`
			} `yaml:"spec"`
		} `yaml:"items"`
	}

	err := yaml.Unmarshal(podDescriptionYAML, &podDescription)
	if err != nil {
		return nil, err
	}

	var podNodes []string
	for _, item := range podDescription.Items {
		podNodes = append(podNodes, item.Spec.NodeName)
	}

	return podNodes, nil
}
