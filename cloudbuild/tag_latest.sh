#!/bin/bash

# Determine the appropriate image to tag "latest".
# It pulls down the tags that are currently available remotely.
# Among those and the new local tag, the highest version is picked
# to be tagged as "latest".
# We're only considering SemVer major.minor.patch tags.
# Usage:
#   tag_latest.sh <repo name> <local tag>
# For example:
#   tag_latest.sh gcr.io/$PROJECT_ID/testrunner "0.1.3"
#   tag_latest.sh gcr.io/$PROJECT_ID/testrunner $TAG_NAME

set -eu

readonly repo="$1"
readonly tag="$2"

latest_remote_version=$( \
  gcloud container images list-tags "$repo" \
    --format "value(tags)" --filter 'tags~\d+\.\d+\.\d+'\
    | tr ',' '\n' \
    | egrep '[0-9]+\.[0-9]+\.[0-9]+' \
    | sort -t '.' -k 1,1 -k 2,2 -k 3,3 -g \
    | tail -n 1)
echo "Latest remote version: $latest_remote_version"
latest_version=$( \
  printf '%s\n%s' "$latest_remote_version" "$tag" \
    | sort -t '.' -k 1,1 -k 2,2 -k 3,3 -g \
    | tail -n 1)
echo "Latest version should be: $latest_version"

# If the tag is not the local tag, we need to pull it down
# before tagging "latest" locally.
if [[ "$tag" != "$latest_version" ]]; then
  docker pull "$repo:$latest_version"
fi
docker tag "$repo:$latest_version" "$repo:latest"
