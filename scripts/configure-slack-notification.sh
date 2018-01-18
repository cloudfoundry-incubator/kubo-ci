#!/bin/bash

set -exu -o pipefail

REPO=${REPO:-target-repo}

# .git/ref is provided by concourse resource
REF=$(cat "$REPO/.git/ref")

COMMITTER=$(git -C "$REPO" show -s --format="%ce" "$REF")
COMMITTER_SLACK_NAME=$(bosh int git-kubo-home/slackers "--path=/$COMMITTER")

AUTHOR=$(git -C "$REPO" show -s --format="%ae" "$REF")
AUTHOR_SLACK_NAME=$(bosh int git-kubo-home/slackers "--path=/$AUTHOR")

message="$MESSAGE
Committer: $COMMITTER
Author: $AUTHOR
Repo: $REPO
Ref: https://github.com/$REPO/commit/$REF
Slack Usernames: <@$COMMITTER_SLACK_NAME> <@$AUTHOR_SLACK_NAME>"

if [ ! -z "${LOCK_NAME}" ]; then
    message+="
Lock: $LOCK_NAME"
fi

echo "$message" > slack-notification/text
