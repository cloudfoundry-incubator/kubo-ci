#!/bin/bash

cp gcs-shipable-version/* gcs-shipable-version-output
echo | cat kubo-version/number - >> gcs-shipable-version-output/shipable
