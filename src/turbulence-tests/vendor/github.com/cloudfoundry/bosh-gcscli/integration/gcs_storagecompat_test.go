/*
 * Copyright 2017 Google Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package integration

import (
	"fmt"
	"os"

	"github.com/cloudfoundry/bosh-gcscli/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Integration", func() {
	Context("invalid storage_class for bucket (Default Applicaton Credentials) configuration", func() {
		regional := os.Getenv(RegionalBucketEnv)
		multiRegional := os.Getenv(MultiRegionalBucketEnv)

		var ctx AssertContext
		BeforeEach(func() {
			Expect(regional).ToNot(BeEmpty(),
				fmt.Sprintf(NoBucketMsg, RegionalBucketEnv))
			Expect(multiRegional).ToNot(BeEmpty(),
				fmt.Sprintf(NoBucketMsg, MultiRegionalBucketEnv))

			ctx = NewAssertContext()
		})
		AfterEach(func() {
			ctx.Cleanup()
		})

		configurations := []TableEntry{
			Entry("MultiRegional bucket, regional StorageClass", &config.GCSCli{
				BucketName:   multiRegional,
				StorageClass: "REGIONAL",
			}),
			Entry("Regional bucket, multiregional StorageClass", &config.GCSCli{
				BucketName:   regional,
				StorageClass: "MULTI_REGIONAL",
			}),
		}

		DescribeTable("Invalid Put should fail",
			func(config *config.GCSCli) {
				ctx.AddConfig(config)
				AssertPutFails(gcsCLIPath, ctx)
			},
			configurations...)
	})
})
