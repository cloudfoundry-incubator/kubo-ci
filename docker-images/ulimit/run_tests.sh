#!/bin/bash

if [ $# -ne 1 ]; then
  echo "Usage:"
  echo "  ${0} <docker image id>"
  exit 1
fi

output=$(docker run ${1})

re='^[0-9]+$'
if ! [[ $output =~ $re ]] ; then
   echo "Expect ulimit value to be a number" >&2; exit 1
fi

