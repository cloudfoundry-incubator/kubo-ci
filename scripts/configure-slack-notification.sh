#!/bin/bash

set -exu -o pipefail

ROOT="$(pwd)"
REPOS=${REPO:-target-repos}
SLACK_TEXT_TEMPLATE='{
    "text": $pipeline,
    "attachments": $attachments
}'

SLACK_ATTACHMENT_TEMPLATE='{
    "color": "#ff0000",
    "title": $title,
    "fields": [
        {"title": "Author", "short": true, "value": $author},
        {"title": "Committer", "short": true, "value": $committer}
    ]
}'

function main() {
  local attachments="[]"
  for repo in ${REPOS}/*; do
    local attachment="$(jq -n \
      --arg title "$(basename "${repo}") (commit $(get_commit_link "${repo}"))" \
      --arg author "$(get_author_name "${repo}")" \
      --arg committer "$(get_committer_name "${repo}")" \
      "${SLACK_ATTACHMENT_TEMPLATE}")"

    attachments="$(echo "${attachments}" | jq \
        --argjson attachment "${attachment}" \
        '. += [$attachment]')"
  done

  jq -n \
    --arg pipeline "Build Failed. <https://ci.kubo.sh/teams/\$BUILD_TEAM_NAME/pipelines/\$BUILD_PIPELINE_NAME|Pipeline Job>" \
    --argjson attachments "$attachments" \
    "${SLACK_TEXT_TEMPLATE}" > "${ROOT}/slack-notification/text"
}

function get_repo_ref() {
  local repo="${1}"
  git -C "${repo}" show -s --format=%h $(cat "${repo}/.git/ref")
}

function get_commit_link() {
  local repo="${1}"
  local ref="$(get_repo_ref ${repo})"
  echo "<https://$(git -C "${repo}" remote get-url origin | cut -d@ -f2 | sed -e 's|:|/|' -e 's|.git$||')/commit/${ref}|${ref}>"
}

function get_author_name() {
  local repo="${1}"
  local author=$(git -C "${repo}" show -s --format="%ae" "$(get_repo_ref "${repo}")")

  get_slacker_name "${author}"
}

function get_committer_name() {
  local repo="${1}"
  local committer=$(git -C "${repo}" show -s --format="%ce" "$(get_repo_ref "${repo}")")

  get_slacker_name "${committer}"
}

function get_slacker_name() {
  local lookup_name="${1}"
  local slack_name="$(bosh int git-kubo-home/slackers "--path=/${lookup_name}" || echo "${lookup_name}")"

  echo "${slack_name}" | sed '/^$/d'
}

main
