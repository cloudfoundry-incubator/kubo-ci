#!/bin/bash

set -exu -o pipefail

COMMITTER=$(cat $REPO/.git/committer)
export COMMITTER

SLACK_NAME=$(echo $COMMITTER | cut -d@ -f1)

REF=$(cat git-kubo-ci/.git/ref)
export REF

echo "$MESSAGE\nCommitter: $COMMITTER\nRef: $REF\nSlack Username (Guess): $SLACK_NAME" > slack-notification/text

echo "@$SLACK_NAME" > slack-notification/channel
