#!/bin/bash

set -eEou pipefail

RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo "====================================================="
echo "|                   Housekeeper                     |"
echo "====================================================="

echo "               ~~~~~Pipelines~~~~~                   "
echo "The following pipelines do not match what is set on concourse,"
echo "we recommend that you run ./set_pipeline on each file to debug."
echo "This may take some time while we check all pipelines in the current directory."

NUM_EXPECTED_LINES=7
for pipeline in $(find . -maxdepth 1 -name '*.yml'); do
  if [[ ${NUM_EXPECTED_LINES} -ne $(echo "n" | ./set_pipeline "${pipeline}" | wc -l) ]]; then
    echo -e "${RED}$(basename "${pipeline}")${NC}"
  fi
done

echo "                ~~~~~Tasks~~~~~                      "
echo "The following tasks may not be referenced by any pipeline"
echo "Please manually check each and delete if unused."

for task in $(find ./tasks -name '*.yml'); do
  if ! grep -q "$(basename "${task}")" *.yml; then
    echo -e "${YELLOW}${task}${NC}"
  fi
done

echo "               ~~~~~Scripts~~~~~~                    "
echo "The following scripts may not be referenced by any task"
echo "they may however be called directly by a pipeline or other script, which is not checked"
echo "Please manually check each and delete if unused."

for script in $(find ./scripts -name '*.sh'); do
  if ! grep -qr "$(basename "${script}")" ./tasks; then
    if ! grep -qr "$(basename "${script}")" *.yml; then
      if ! grep -qr "$(basename "${script}")" ./scripts; then
        echo -e "${YELLOW}${script}${NC}"
      fi
    fi
  fi
done

echo "============ Housekeeper has finished ==============="
