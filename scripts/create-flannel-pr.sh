#!/bin/bash

set -exu -o pipefail

create_pr_payload() {
  title="Flannel upgrade $1"
  body="This is an auto generated PR created for flannel upgrade to $1"
  echo '{"title":"'"$title"'","body":"'"$body"'","head":"'"$2"'","base":"master"}'
}

main() {
  local tag=$(cat flannel-release/tag)

  mkdir -p ~/.ssh
  cat > ~/.ssh/config <<EOF
StrictHostKeyChecking no
LogLevel quiet
EOF
  chmod 0600 ~/.ssh/config

  cat > ~/.ssh/id_rsa <<EOF
${GIT_SSH_KEY}
EOF
  chmod 0600 ~/.ssh/id_rsa
  eval $(ssh-agent) >/dev/null 2>&1
  trap "kill $SSH_AGENT_PID" 0
  ssh-add ~/.ssh/id_rsa

  pushd git-kubo-release-output

  git config --global user.email "cfcr+cibot@pivotal.io"
  git config --global user.name "CFCR CI BOT"

  branch_name="upgrade/flannel${tag}"
  git checkout -b $branch_name
  git add .
  if git diff-index --quiet HEAD --; then
    echo "No changes detected"
    exit 0
  fi

  git commit -m "Upgrade flannel to $tag"
  git push origin $branch_name

  # create a PR here
  token=${CFCR_USER_TOKEN}
  payload=$(create_pr_payload $tag $branch_name)
  curl -u cfcr:${CFCR_USER_TOKEN} -H "Content-Type: application/json" -X POST -d "$payload" https://api.github.com/repos/cloudfoundry-incubator/kubo-release/pulls

popd
}

main

