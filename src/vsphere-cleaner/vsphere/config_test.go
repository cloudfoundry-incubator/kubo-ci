package vsphere_test

import (
	"errors"
	"vsphere-cleaner/vsphere"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func getFunctionAssigningAndReturning(assign string, err error) func(interface{}) error {
	return func(s interface{}) error {
		str := s.(*string)
		*str = assign
		return err
	}
}

var _ = Describe("Config", func() {
	Describe("IP", func() {
		It("unmarshals valid IP", func() {
			ip := new(vsphere.IP)
			err := ip.UnmarshalYAML(getFunctionAssigningAndReturning("10.1.1.1", nil))
			Expect(err).ToNot(HaveOccurred())
			Expect(*ip).To(Equal(vsphere.IP("10.1.1.1")))
		})

		Context("when internal ip has invalid format", func() {
			It("throws error", func() {
				ip := new(vsphere.IP)
				err := ip.UnmarshalYAML(getFunctionAssigningAndReturning("foo", nil))
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when internal function returns error", func() {
			It("throws error", func() {
				ip := new(vsphere.IP)
				err := ip.UnmarshalYAML(getFunctionAssigningAndReturning("10.1.1.1", errors.New("I am error")))
				Expect(err).To(HaveOccurred())
			})
		})
	})

})
