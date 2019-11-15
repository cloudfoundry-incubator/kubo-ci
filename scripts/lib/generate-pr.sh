create_pr_payload() {
  title="$1 upgrade $2"
  body="This is an auto generated PR created for $1 upgrade to $2"
  echo '{"title":"'"$title"'","body":"'"$body"'","head":"'"$3"'","base":"'"$4"'"}'
}

# Needs to be called from the directory where PR needs to be generated
generate_pull_request() {
  local component=$1
  local tag=$2
  local repo=$3
  local base_branch=$4

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

  git config --global user.email "cfcr+cibot@pivotal.io"
  git config --global user.name "CFCR CI BOT"

  branch_name="upgrade/${component}${tag}"
  git checkout -b $branch_name
  git add .
  git commit -m "Upgrade $component to $tag"
  git push origin $branch_name

  # create a PR here
  payload=$(create_pr_payload "$component" "$tag" "$branch_name" "$base_branch")
  curl -u "cfcr:${CFCR_USER_TOKEN}" -H "Content-Type: application/json" -X POST -d "$payload" "https://api.github.com/repos/pivotal-cf/${repo}/pulls" --fail
}
