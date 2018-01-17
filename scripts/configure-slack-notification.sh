#!/bin/bash

set -exu -o pipefail

echo "I am a message" > slack-notification/text
echo "@tony" > slack-notification/channel
