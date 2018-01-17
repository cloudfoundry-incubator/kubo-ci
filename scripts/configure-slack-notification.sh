#!/bin/bash

set -exu -o pipefail


COMMITTER=$(cat git-kubo-ci/.git/committer)
export COMMITTER

SLACK_NAME=$(echo $COMMITTER | cut -d@ -f1)

REF=$(cat git-kubo-ci/.git/ref)
export REF

echo "Committer: $COMMITTER\nRef: $REF" > slack-notification/text

echo "@$SLACK_NAME" > slack-notification/channel
