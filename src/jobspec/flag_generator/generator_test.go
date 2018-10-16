package flag_generator_test

import (
	"io/ioutil"
	"jobspec/flag_generator"
	"path/filepath"

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

	Context("spec file", func() {
		It("reads spec from file", func() {
			spec, err := flag_generator.ReadSpecFile("./fixtures/oldSpec")
			Expect(err).NotTo(HaveOccurred())
			Expect(spec.Templates["bin/script.erb"]).To(Equal("bin/script"))
			Expect(spec.Packages).To(ContainElement("somePackage"))
			Expect(spec.Properties).To(HaveKey("boshProperty"))
			Expect(spec.Consumes[0]["name"]).To(Equal("some-link"))
			Expect(spec.Provides[0]["name"]).To(Equal("test-job"))
		})

		It("writes spec to file", func() {
			expected := `consumes: []
name: ""
packages: []
properties:
  my-prop:
    description: awesome
provides: []
templates: {}`
			test := &flag_generator.JobSpec{
				Properties: map[string]flag_generator.Property{
					"my-prop": flag_generator.Property{
						Description: "awesome",
					},
				},
			}

			tfile, err := ioutil.TempFile("", "")
			Expect(err).NotTo(HaveOccurred())
			defer tfile.Close()

			tfilepath, err := filepath.Abs(tfile.Name())
			Expect(err).NotTo(HaveOccurred())
			flag_generator.WriteSpecFile(tfilepath, test)

			actual, err := ioutil.ReadFile(tfilepath)
			Expect(err).NotTo(HaveOccurred())
			Expect(actual).To(MatchYAML(expected))
		})
	})
})
