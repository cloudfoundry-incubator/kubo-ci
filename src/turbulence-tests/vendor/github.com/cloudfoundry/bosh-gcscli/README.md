## GCS CLI

A CLI for uploading, fetching and deleting content to/from
the [GCS blobstore](https://cloud.google.com/storage/). This is **not**
an official Google Product.

## Installation

```
go get github.com/cloudfoundry/bosh-gcscli
```

## Usage

Given a JSON config file (`config.json`)...

``` json
{
  "bucket_name":            "<string> (required)",

  "credentials_source":     "<string> [path_to_service_account_file|none]",
  "storage_class":          "<string> (explicit_storage_class|none)",
}
```


Empty `credentials_source` implies attempting to use Application Default
Credentials.

Empty `storage_class` implies using the default for the bucket.

``` bash
# Usage
bosh-gcscli --help

# Command: "put"
# Upload a blob to the GCS blobstore.
bosh-gcscli -c config.json put <path/to/file> <remote-blob>

# Command: "get"
# Fetch a blob from the GCS blobstore.
# Destination file will be overwritten if exists.
bosh-gcscli -c config.json get <remote-blob> <path/to/file>

# Command: "delete"
# Remove a blob from the GCS blobstore.
bosh-gcscli -c config.json delete <remote-blob>

# Command: "exists"
# Checks if blob exists in the GCS blobstore.
bosh-gcscli -c config.json exists <remote-blob>
```

Alternatively, this package's underlying client can be used to access GCS,
see the [godoc](https://godoc.org/github.com/cloudfoundry/bosh-gcscli)
for more information.

## Tooling

A Makefile is provided for ease of development. Targets are annotated
with descriptions.

Integration tests expect to be run from a host with [Application Default
Credentials](https://developers.google.com/identity/protocols/application-default-credentials)
available which has permissions to create and delete buckets.
Application Default Credentials are present on any GCE instance and inherit
the permisions of the [service account](https://cloud.google.com/iam/docs/service-accounts)
assigned to the instance.

## License

This library is licensed under Apache 2.0. Full license text is
available in [LICENSE](LICENSE).