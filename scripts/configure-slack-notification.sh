#!/bin/bash

set -exu -o pipefail

COMMITTER=$(cat $REPO/.git/committer)

SLACK_NAME=$(echo $COMMITTER | cut -d@ -f1)

REF=$(cat $REPO/.git/ref)

echo "$MESSAGE\nCommitter: $COMMITTER\nRepo: $REPO\nRef: $REF\nSlack Username (guess): <@$SLACK_NAME>" > slack-notification/text

echo "@$SLACK_NAME" > slack-notification/channel
