#!/bin/bash

set -exu -o pipefail

COMMITTER=$(cat $REPO/.git/committer)

SLACK_NAME=$(echo $COMMITTER | cut -d@ -f1)

REF=$(cat $REPO/.git/ref)

message="$MESSAGE\nCommitter: $COMMITTER\nRepo: $REPO\nRef: $REF\nSlack Username (guess): <@$SLACK_NAME>"

if [ -d "kubo-lock" ]; then
    kubo_lock_name=$(cat kubo-lock/name)
    message+="\nLock: $kubo_lock_name"
fi

echo $message > slack-notification/text

echo "@$SLACK_NAME" > slack-notification/channel
