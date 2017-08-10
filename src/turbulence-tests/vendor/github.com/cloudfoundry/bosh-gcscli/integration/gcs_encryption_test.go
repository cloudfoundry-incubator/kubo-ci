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

	"crypto/sha256"

	"github.com/cloudfoundry/bosh-gcscli/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

// encryptionKeyBytes are used as the key in tests requiring encryption.
var encryptionKeyBytes = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}

// encryptionKeyBytesHash is the has of the encryptionKeyBytes
//
// Typical usage is ensuring the encryption key is actually used by GCS.
var encryptionKeyBytesHash = sha256.Sum256(encryptionKeyBytes)

var _ = Describe("Integration", func() {
	Context("general (Default Applicaton Credentials) configuration", func() {
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

		encryptedConfigs := []TableEntry{
			Entry("MultiRegional bucket, default StorageClass, encrypted", &config.GCSCli{
				BucketName:    multiRegional,
				EncryptionKey: encryptionKeyBytes,
			}),
			Entry("Regional bucket, default StorageClass, encrypted", &config.GCSCli{
				BucketName:    regional,
				EncryptionKey: encryptionKeyBytes,
			}),
		}

		DescribeTable("Get with correct encryption_key works",
			func(config *config.GCSCli) {
				ctx.AddConfig(config)
				AssertEncryptionWorks(gcsCLIPath, ctx)
			},
			encryptedConfigs...)

		DescribeTable("Get with wrong encryption_key should fail",
			func(config *config.GCSCli) {
				ctx.AddConfig(config)
				AssertWrongKeyEncryptionFails(gcsCLIPath, ctx)
			},
			encryptedConfigs...)

		DescribeTable("Get with no encryption_key should fail",
			func(config *config.GCSCli) {
				ctx.AddConfig(config)
				AssertNoKeyEncryptionFails(gcsCLIPath, ctx)
			},
			encryptedConfigs...)
	})
})
