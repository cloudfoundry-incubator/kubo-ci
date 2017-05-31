package ipcalc_test

import (
	"vsphere-cleaner/ipcalc"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Calculator", func() {
	Describe("GetIPsFromCIDR", func() {
		It("should calculate a list of IPs", func() {
			ips, err := ipcalc.GetIPsFromCIDR("10.1.1.0/31")
			Expect(err).ToNot(HaveOccurred())
			Expect(ips).To(Equal([]string{"10.1.1.1"}))
		})

		It("should not include broadcast IPs", func() {
			ips, err := ipcalc.GetIPsFromCIDR("10.1.1.254/31")
			Expect(err).ToNot(HaveOccurred())
			Expect(ips).To(Equal([]string{"10.1.1.254"}))
		})

		It("should return error if the CIDR isn't valid", func() {
			_, err := ipcalc.GetIPsFromCIDR("Where is the seismic c-beam?")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("GetIPsFromRange", func() {
		It("should return the IP if it is not a range", func() {
			ips, _ := ipcalc.GetIPsFromRange("10.1.1.1")
			Expect(ips).To(Equal([]string{"10.1.1.1"}))
		})

		It("should calculate a list of two IPs", func() {
			ips, err := ipcalc.GetIPsFromRange("10.1.1.1-10.1.1.2")
			Expect(err).ToNot(HaveOccurred())
			Expect(ips).To(Equal([]string{"10.1.1.1", "10.1.1.2"}))
		})

		It("should calculate a list of more than two IPs", func() {
			ips, err := ipcalc.GetIPsFromRange("10.1.1.1-10.1.1.3")
			Expect(err).ToNot(HaveOccurred())
			Expect(ips).To(Equal([]string{"10.1.1.1", "10.1.1.2", "10.1.1.3"}))
		})

		table.DescribeTable("errors", func(ipRange string) {
			_, err := ipcalc.GetIPsFromRange(ipRange)
			Expect(err).To(HaveOccurred())

		},
			table.Entry("Single IP is invalid", "foo"),
			table.Entry("First IP in range is invalid", "foo-10.1.1.1"),
			table.Entry("Last IP in range is invalid", "10.1.1.1-foo"),
			table.Entry("First IP is bigger than last IP", "10.1.1.2-10.1.1.1"),
		)
	})

	Describe("Difference", func() {
		It("should find difference between two empty arrays", func() {
			Expect(ipcalc.Difference([]string{}, []string{})).To(Equal([]string{}))
		})

		It("should find difference when only second array is empty", func() {
			Expect(ipcalc.Difference([]string{"foo"}, []string{})).To(Equal([]string{"foo"}))
		})

		It("should find difference when both arrays are not empty", func() {
			Expect(ipcalc.Difference([]string{"foo", "bar"}, []string{"bar"})).To(Equal([]string{"foo"}))
			Expect(ipcalc.Difference([]string{"bar", "foo", "bar"}, []string{"bar"})).To(Equal([]string{"foo"}))
		})
	})
})
