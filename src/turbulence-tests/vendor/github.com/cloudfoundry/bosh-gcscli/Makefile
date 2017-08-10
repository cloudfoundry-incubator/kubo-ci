# Copyright 2017 Google Inc.
# 
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
# 
#    http://www.apache.org/licenses/LICENSE-2.0
# 
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

default: test-int

# build the binary
build:
	go build

# Fetch base dependencies as well as testing packages
get-deps:
	go get
	# Ginkgo and omega test tools
	go get github.com/onsi/ginkgo/ginkgo
	go get github.com/onsi/gomega

# Cleans up directory and source code with gofmt
clean:
	go clean ./...

# Run gofmt on all code
fmt:
	gofmt -l -w .

# Run linter with non-strict checking
lint:
	ls -d */ | grep -v vendor | xargs -L 1 golint

# Vet code
vet:
	go tool vet $$(ls -d */ | grep -v vendor)

# Generate a $StorageClass.lock which contains our bucket name
# used for testing. Buckets must be unique among all in GCS,
# we cannot simply hardcode a bucket.
.PHONY: FORCE
regional.lock:
	@test -s "regional.lock" || \
	{ echo -n "bosh-gcs"; \
	cat /dev/urandom | tr -dc 'a-z0-9' | fold -w 40 | head -n 1 ;} > regional.lock

# Create a bucket using the name located in $StorageClass.lock with
# a sane location.
regional-bucket: regional.lock
	@gsutil ls | grep "$$(cat regional.lock)"&> /dev/null; if [ $$? -ne 0 ]; then \
		gsutil mb -c REGIONAL -l us-east1 "gs://$$(cat regional.lock)"; \
	fi

.PHONY: FORCE
multiregional.lock:
	@test -s "multiregional.lock" || \
	{ echo -n "bosh-gcs"; \
	cat /dev/urandom | tr -dc 'a-z0-9' | fold -w 40 | head -n 1 ;} > multiregional.lock

multiregional-bucket: multiregional.lock
	@gsutil ls | grep "$$(cat multiregional.lock)"&> /dev/null; if [ $$? -ne 0 ]; then \
		gsutil mb -c MULTI_REGIONAL -l us "gs://$$(cat multiregional.lock)"; \
	fi

# Create all buckets necessary for the test.
prep-gcs: regional-bucket multiregional-bucket

# Remove all buckets listed in $StorageClass.lock files.
clean-gcs:
	test -s "multiregional.lock" && test -s "regional.lock"
	@gsutil rb "gs://$$(cat regional.lock)"
	rm regional.lock
	@gsutil rb "gs://$$(cat multiregional.lock)"
	rm multiregional.lock

# Perform only unit tests
test-unit:
	ginkgo -r -skipPackage integration

# Perform all tests, including integration tests.
test-int: get-deps clean fmt lint vet build prep-gcs
	 export MULTIREGIONAL_BUCKET_NAME="$$(cat multiregional.lock)" && \
	 export REGIONAL_BUCKET_NAME="$$(cat regional.lock)" && \
	 ginkgo -r

# Perform all non-long tests, including integration tests.
test-fast-int: get-deps clean fmt lint vet build prep-gcs
	 export MULTIREGIONAL_BUCKET_NAME="$$(cat multiregional.lock)" && \
	 export REGIONAL_BUCKET_NAME="$$(cat regional.lock)" && \
	 export SKIP_LONG_TESTS="yes" && \
	 ginkgo -r