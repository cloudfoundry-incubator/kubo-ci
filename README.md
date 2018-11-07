# Continuous Integration (CI) for Cloud Foundry Container Runtime (CFCR)

## Coverage

Repos covered by the pipelines:
- [Kubo release](https://www.github.com/cloudfoundry-incubator/kubo-release) - a BOSH release for Kubernetes
- [Kubo deployment](https://www.github.com/cloudfoundry-incubator/kubo-deployment) - BOSH deployment manifest and ops-files for deploying a single kubernetes cluster using CFCR.
- [CFCR Etcd](https://github.com/cloudfoundry-incubator/cfcr-etcd-release) - A BOSH release of Etcd used in CFCR
- [Docker BOSH release](https://github.com/cloudfoundry-incubator/docker-boshrelease) - A BOSH release of Docker used in CFCR

Other related repos:
- [Kubo disaster recovery acceptance tests
  (kdrats)](https://github.com/cloudfoundry-incubator/kubo-disaster-recovery-acceptance-tests) - Provides tests used by
  pipelines contained within this repo.
- [Kubo locks] - provides configuration information for the CI environments
  used by the pipelines.  Private repo for CFCR team use.

## Pipelines
To view the pipelines visit [ci.kubo.sh](https://ci.kubo.sh)

### Testing
#### Kubo-deployment
The main pipeline, runs on every commit in kubo-release and kubo-deployment. Is also used to create Github releases.
Runs tests on GCP, AWS and vSphere.

#### Long-lived-openstack
Runs tests on Openstack with a long-lived director instead of a freshly-deployed one for each run.
[Configuring OpenStack and Creating New Environments](openstack.md)

#### custom-CIDRs
Tests the CIDRS suite of the integration tests.  It tests changing the Pod and Service CIDRs for kubernetes.

#### Istio
Tests whether Istio specifications can be deployed on CFCR, and the Istio integration tests are run.

#### Gaffer
Attempts to update the CFCR deployment on which the gaffer application is run.  The gaffer application backs https://gaffer.kubo.sh

#### kubo-release-pr & kubo-deployment-pr
Test pull requests to the repos.  Have to be triggered manually.

#### pod-security-policy
Tests that the PodSecurityPolicy admission controller and policies work as expected.

#### CFCR-etcd
Tests the etcd release used in CFCR with different configurations. Creates and destroys the cluster at each run.

#### CFCR-etcd-long-running
Uses a long-lived etcd cluster and specifically focuses on potential “split brain” failure scenarios, centred around network partitions and VM restarts.

### Bumping components
The bump-* pipelines are used to automatically update components of CFCR or the CI infrastructure. The docker-boshrelease pipeline tests and releases new versions of the docker release.

### Building docker images
This pipeline builds all the images used in the other pipelines based on the [dockerfiles](https://github.com/cloudfoundry-incubator/kubo-ci/tree/master/docker-images) located in kubo-ci.

### Maintaining CI IaaS environments

The recycle pipeline is used to clean up environments.  It is triggered
automatically whenever a lock is released by the pipeline using that
environment.
The cleanup-gcp is used to clean up leftover LoadBalancers and Disks in GCP.
The vsphere-cleaner pipeline generates a binary which is used by the recycle pipeline to clean vsphere environments.

