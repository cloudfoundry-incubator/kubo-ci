#!/bin/bash

set -exu -o pipefail

REPOS=${REPO:-target-repos}

FILE=slack-notification/text

echo "$MESSAGE" > "$FILE"

for REPO in $REPOS/*; do
    # .git/ref is provided by concourse resource
    REF=$(git -C "$REPO" show -s --format=%h $(cat "$REPO/.git/ref"))

    COMMITTER=$(git -C "$REPO" show -s --format="%ce" "$REF")
    COMMITTER_SLACK_NAME=$(bosh int git-kubo-home/slackers "--path=/$COMMITTER" || echo "$COMMITTER")
    COMMITTER_SLACK_NAME=$(echo "$COMMITTER_SLACK_NAME" | sed '/^$/d')

    AUTHOR=$(git -C "$REPO" show -s --format="%ae" "$REF")
    AUTHOR_SLACK_NAME=$(bosh int git-kubo-home/slackers "--path=/$AUTHOR" || echo "$AUTHOR")
    AUTHOR_SLACK_NAME=$(echo "$AUTHOR_SLACK_NAME" | sed '/^$/d')

    COMMIT_LINK=$(echo "<https://$(git -C "$REPO" remote get-url origin | cut -d@ -f2 | sed -e 's|:|/|' -e 's|.git$||')/commit/$REF|$REF>")
    echo "<@$COMMITTER_SLACK_NAME> and <@$AUTHOR_SLACK_NAME> committed in $REPO (commit $COMMIT_LINK)" >> $FILE
    if [[ "$COMMITTER_SLACK_NAME" == "$COMMITTER" ]] || [[ "$AUTHOR_SLACK_NAME" == "$AUTHOR" ]]; then
        echo "<!subteam^S7V8MPT6U> There is an unknown email id in this commit!" >> "$FILE"
    fi
done

if [ ! -z "${LOCK_NAME}" ]; then
    echo "Lock: $LOCK_NAME" >> "$FILE"
fi
