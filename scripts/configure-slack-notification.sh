#!/bin/bash

set -exu -o pipefail

COMMITTER=$(cat $REPO/.git/committer)

SLACK_NAME=$(echo $COMMITTER | cut -d@ -f1)

REF=$(cat git-kubo-ci/.git/ref)

echo "$MESSAGE\nCommitter: $COMMITTER\nRef: $REF\nSlack Username (Guess): @$SLACK_NAME" > slack-notification/text

SLACK_NAME=tony
echo "@$SLACK_NAME @muchhals" > slack-notification/channel
