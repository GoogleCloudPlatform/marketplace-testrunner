#!/bin/bash

# Usage:
#   tag_latest.sh <repo name> <local tag>
# For example:
#   tag_latest.sh gcr.io/$PROJECT_ID/testrunner "0.1.3"
#   tag_latest.sh gcr.io/$PROJECT_ID/testrunner $TAG_NAME

set -eu

latest_remote_version=$( \
  gcloud container images list-tags "$1" \
    --format "value(tags)" --filter 'tags~\d+\.\d+\.\d+'\
    | tr ',' '\n' \
    | egrep '[0-9]+\.[0-9]+\.[0-9]+' \
    | sort -t '.' -k 1,1 -k 2,2 -k 3,3 -k 4,4 -g \
    | tail -n 1)
echo "Latest remote version: $latest_remote_version"
latest_version=$( \
  printf '%s\n%s' "$latest_remote_version" "$2" \
    | sort -t '.' -k 1,1 -k 2,2 -k 3,3 -k 4,4 -g \
    | tail -n 1)
echo "Latest version should be: $latest_version"
if [[ "$2" != "$latest_version" ]]; then
  docker pull "$1:$latest_version"
fi
docker tag "$1:$latest_version" "$1:latest"
