#!/bin/sh -e

cp -r kubo-lock/* kubo-lock-with-proxy/

proxy_setting=$(bosh-cli int - --path /proxy_setting < proxy-tf/metadata)
echo >> kubo-lock-with-proxy/metadata
echo "http_proxy: $proxy_setting" >> kubo-lock-with-proxy/metadata
echo "https_proxy: $proxy_setting" >> kubo-lock-with-proxy/metadata
