package generic_test

import (
	"tests/test_helpers"

    "regexp"

	. "github.com/onsi/ginkgo"
    . "github.com/onsi/ginkgo/extensions/table"
    . "github.com/onsi/gomega"
)

var _ = Describe("API Versions", func() {
	var (
		runner *test_helpers.KubectlRunner
	)

	BeforeEach(func() {
		runner = test_helpers.NewKubectlRunner()
	})

    DescribeTable("api-versions", func(api_regex string) {
        regex, err := regexp.Compile(api_regex)
        Expect(err).NotTo(HaveOccurred())

		output := runner.GetOutput("api-versions")
	    for i := 0; i < len(output); i++ {
		    if regex.MatchString(output[i]) {
			    return
		    }
	    }
        Fail("Unable to find api-version: '" + api_regex + "'")
    },
        Entry("RBAC v1alpha1 is enabled", "^rbac.*/v1alpha1"),
        Entry("RBAC v1beta1 is enabled", "^rbac.*/v1beta1"),
        Entry("RBAC v1 is enabled", "^rbac.*/v1"),
    )

})
