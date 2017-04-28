#!/bin/sh -e

cp -r kubo-lock/* kubo-lock-with-proxy/
bosh-cli interpolate kubo-lock/metadata --ops-file=git-kubo-ci/ops-files/internetless.yml > kubo-lock-with-proxy/metadata

proxy_setting=$(bosh-cli int - --path /proxy_setting < proxy-tf/metadata)
echo "http_proxy: $proxy_setting" >> kubo-lock-with-proxy/metadata
echo "https_proxy: $proxy_setting" >> kubo-lock-with-proxy/metadata
