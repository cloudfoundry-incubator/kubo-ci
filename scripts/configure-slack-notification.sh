#!/bin/bash

set -exu -o pipefail

COMMITTER=$(cat $REPO/.git/committer)

SLACK_NAME=$(bosh int git-kubo-home/slackers "--path=/$COMMITER")

REF=$(cat $REPO/.git/ref)

message="$MESSAGE\nCommitter: $COMMITTER\nRepo: $REPO\nRef: $REF\nSlack Username: <@$SLACK_NAME>"

if [ ! -z "${LOCK_NAME}" ]; then
    message+="\nLock: $LOCK_NAME"
fi

echo $message > slack-notification/text

echo "@$SLACK_NAME" > slack-notification/channel
