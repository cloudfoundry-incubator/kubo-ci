package flag_generator_test

import (
	"jobspec/flag_generator"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/pflag"
)

type TestK8sFlags struct {
	flag *pflag.Flag
}

func (tkf TestK8sFlags) AddFlags(fs *pflag.FlagSet) {
	if tkf.flag != nil {
		fs.AddFlag(tkf.flag)
	}
}

var _ = Describe("Generator", func() {
	Context("k8s flag", func() {
		It("returns populated Property", func() {
			prop := flag_generator.GenerateArgsFromFlags(TestK8sFlags{
				flag: &pflag.Flag{
					Name: "my-flag",
				},
			}, []string{})
			Expect(prop).NotTo(BeNil())
			Expect(prop.Properties["my-flag"]).NotTo(BeNil())
		})

		It("skips blacklisted flags", func() {
			prop := flag_generator.GenerateArgsFromFlags(TestK8sFlags{
				flag: &pflag.Flag{
					Name: "my-flag",
				},
			}, []string{"my-flag"})
			Expect(prop).NotTo(BeNil())
			Expect(prop.Properties).NotTo(HaveKey("my-flag"))
		})

		It("includes global 'v' flag", func() {
			prop := flag_generator.GenerateArgsFromFlags(TestK8sFlags{}, []string{})
			Expect(prop).NotTo(BeNil())
			Expect(prop.Properties["v"]).NotTo(BeNil())
		})
	})
})
