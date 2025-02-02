#!/usr/bin/env bash

# This script will go through all compose files, looking for images defined like
# so:
# services:
#   cupdate:
#     image: ghcr.io/alexgustafsson/cupdate:0.16.0
#
# For all such images, the script will perform a lookup towards Cupdate's API to
# find the latest image and update the manifest in-place if the user whishes to
# do so.

# CUPDATE_API must be set to the URL where your Cupdate instance is reachable
CUPDATE_API="https://cupdate.home.local/api/v1/image"
# COMPOSE_FILE_NAME is the name of compose files. All files by this name
# in current directory will be checked
COMPOSE_FILE_NAME="compose.yaml"

function confirm() {
  while true; do
    echo -n "$1 [Y/n] "
    read -rsn 1 result </dev/tty
    echo
    case "$result" in
    'y' | 'Y' | '')
      return 0
      ;;
    'n' | 'N')
      return 1
      ;;
    *)
      continue
      ;;
    esac
  done
}

# For each compose file found in the current directory
while read -r composeFile; do
  # For each image reference tag replacement found in the compose file, such as
  # services:
  #   cupdate:
  #     image: ghcr.io/alexgustafsson/cupdate:0.16.0
  while read -r reference; do
    if [[ -z "$reference" ]]; then
      continue
    fi

    # Use the Cupdate API to find the latest reference
    latestReference="$(curl --silent --get "$CUPDATE_API" --data-urlencode "reference=$reference" | jq -rc .latestReference 2>/dev/null)"
    if [[ ! $! -eq 0 ]] || [[ -z "$latestReference" ]] || [[ "$latestReference" = "$reference" ]]; then
      continue
    fi

    # Print a header and ask for confirmation
    echo "$composeFile"
    echo -e "\e[31m$reference\e[0m -> \e[32m$latestReference\e[0m"
    confirm "Update?"

    # If update is confirmed
    if [[ $? -eq 0 ]]; then
      # Create a diff using yq as yq doesn't support keeping whitespace or
      # comments around
      diff="$(diff -u --ignore-blank-lines -L "$composeFile" -L "$composeFile" <(grep -v '^$' "$composeFile") <(yq ".services[] |= select(.image == \"$reference\").image = \"$latestReference\"" "$composeFile"))"
      # Print and apply the diff
      echo "$diff"
      echo "$diff" | patch --silent --unified --posix --fuzz 3 "$composeFile" -i -
    fi

    # Separate runs by a newline
    echo
  done <<<"$(yq '.services[] | select(.image|type == "!!str") | .image' "$composeFile" | grep -v '^$')"
done <<<"$(find . -type f -name "$COMPOSE_FILE_NAME")"
