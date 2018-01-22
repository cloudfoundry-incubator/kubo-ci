#!/bin/sh

set -exu -o pipefail

REPOS=${REPO:-target-repos}

FILE=slack-notification/text

echo "{" > "$FILE"
printf '"text": "%s"' "Pipeline: https://ci.kubo.sh/teams/$BUILD_TEAM_NAME/pipelines/$BUILD_PIPELINE_NAME" >> "$FILE"

echo '"attachments": [' >> "$FILE"

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


    echo '{ "color": "#ff0000",' >> "$FILE"
    printf '"title": "%s",' "$REPO (commit $REF)" >> "$FILE"
    printf '"title_link": "%s"' "$COMMIT_LINK" >> "$FILE"
    printf '"fields": [' >> "$FILE"
    printf '{"title": "Author", "short": true, "value": "%s"}' "$AUTHOR_SLACK_NAME" >> "$FILE"
    printf '{"title": "Committer", "short": true, "value": "%s"}' "$COMMITTER_SLACK_NAME" >> "$FILE"

    echo '}' >> "$FILE"
done

echo ']' >> "$FILE"
echo "}" >> $FILE

#
#{
#	"text": "Build Failed. <https://ci.kubo.sh|Pipeline Job>",
#	"attachments": [
#        { "color": "#ff0000",
#           "title": "Kubo-CI (commit foo)", "title_link": "http://github.com",
#            "fields": [
#			{ "title": "Author","value": "<@akshay>","short": true},
#			{ "title": "Committer", "value": "<@akshay>", "short": true }]
#        }
#    ]
#}