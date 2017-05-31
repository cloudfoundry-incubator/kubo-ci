package vsphere_test

import (
	"net/url"
	"vsphere-cleaner/vsphere"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	Describe("UsedIPs", func() {
		It("should get IPs in InternalCIDR but not those in ReservedIPs", func() {
			config := vsphere.Config{InternalCIDR: "10.1.1.0/30", ReservedIPs: []string{"10.1.1.3-10.1.1.4"}}
			ips, err := config.UsedIPs()
			Expect(err).NotTo(HaveOccurred())
			Expect(ips).To(Equal([]string{"10.1.1.1", "10.1.1.2"}))
		})

		It("should return error if InternalCIDR is not valid", func() {
			config := vsphere.Config{InternalCIDR: "10.1.1.0/130", ReservedIPs: []string{"10.1.1.3-10.1.1.4"}}
			_, err := config.UsedIPs()
			Expect(err).To(HaveOccurred())
		})

		It("should return error if ReservedIPs is not valid", func() {
			config := vsphere.Config{InternalCIDR: "10.1.1.0/30", ReservedIPs: []string{"a"}}
			_, err := config.UsedIPs()
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Url", func() {
		It("should build url from config", func() {
			vSphereUrl := vsphere.Config{IP: "host", User: "user", Password: "password"}.BuildUrl()
			expectedUrl, _ := url.Parse("https://user:password@host/sdk")
			Expect(vSphereUrl).To(Equal(expectedUrl))
		})
	})

})
