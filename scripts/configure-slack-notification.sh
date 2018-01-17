#!/bin/bash

set -exu -o pipefail


COMMITTER=$(cat git-kubo-ci/.git/committer)
export COMMITTER

REF=$(cat git-kubo-ci/.git/ref)
export REF

echo "Committer: $COMMITTER\nRef: $REF" > slack-notification/text
echo "@tony" > slack-notification/channel
