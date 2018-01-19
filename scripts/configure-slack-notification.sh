#!/bin/bash

set -exu -o pipefail

REPOS=${REPO:-target-repos}

FILE=slack-notification/text

echo "$MESSAGE" > $FILE

for REPO in $REPOS/*; do
    # .git/ref is provided by concourse resource
    REF=$(git -C "$REPO" show -s --format=%t $(cat "$REPO/.git/ref"))

    COMMITTER=$(git -C "$REPO" show -s --format="%ce" "$REF")
    COMMITTER_SLACK_NAME=$(bosh int git-kubo-home/slackers "--path=/$COMMITTER" || echo "$COMMITTER")
    COMMITTER_SLACK_NAME=$(echo "$COMMITTER_SLACK_NAME" | sed '/^$/d')

    AUTHOR=$(git -C "$REPO" show -s --format="%ae" "$REF")
    AUTHOR_SLACK_NAME=$(bosh int git-kubo-home/slackers "--path=/$AUTHOR" || echo "$AUTHOR")
    AUTHOR_SLACK_NAME=$(echo "$AUTHOR_SLACK_NAME" | sed '/^$/d')

    echo "<@$COMMITTER_SLACK_NAME> and <@$AUTHOR_SLACK_NAME> committed in $REPO (commit $REF)" >> $FILE
done

if [ ! -z "${LOCK_NAME}" ]; then
    echo "Lock: $LOCK_NAME" >> $FILE
fi
