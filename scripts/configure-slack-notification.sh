#!/bin/bash

set -exu -o pipefail

# .git/ref is provided by concourse resource
REF=$(cat "$REPO/.git/ref")

COMMITTER=$(git -C "$REPO" show -s --format="%ce" "$REF")
AUTHOR=$(git -C "$REPO" show -s --format="%ae" "$REF")


COMMITTER_SLACK_NAME=$(bosh int git-kubo-home/slackers "--path=/$COMMITTER")
AUTHOR_SLACK_NAME=$(bosh int git-kubo-home/slackers "--path=/$AUTHOR")

message="$MESSAGE\nCommitter: $COMMITTER\nAuthor: $AUTHOR\nRepo: $REPO\nRef: $REF\nSlack Usernames: <@$COMMITTER_SLACK_NAME> <@$AUTHOR_SLACK_NAME>"

if [ ! -z "${LOCK_NAME}" ]; then
    message+="\nLock: $LOCK_NAME"
fi

echo "$message" > slack-notification/text

echo "#cfcr-ci" > slack-notification/channel
