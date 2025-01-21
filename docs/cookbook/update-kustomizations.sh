#!/usr/bin/env bash

# This script will go through all kustomization files, looking for image tag
# changes such as the following:
# images:
#   - name: ghcr.io/alexgustafsson/cupdate
#     newTag: 0.15.0
#
# For all such images, the script will perform a lookup towards Cupdate's API to
# find the latest image and update the manifest in-place if the user whishes to
# do so.

# CUPDATE_API must be set to the URL where your Cupdate instance is reachable
CUPDATE_API="https://cupdate.home.local/api/v1/image"
# MANIFEST_FILE_NAME is the name of kustomization files. All files by this name
# in current directory will be checked
MANIFEST_FILE_NAME="kustomization.yaml"

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

# For each kustomization file found in the current directory
while read -r kustomization; do
  # For each image reference tag replacement found in the kustomization, such as
  # images:
  #   - name: ghcr.io/alexgustafsson/cupdate
  #     newTag: 0.15.0
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
    echo "$kustomization"
    echo -e "\e[31m$reference\e[0m -> \e[32m$latestReference\e[0m"
    confirm "Update?"

    # If update is confirmed
    if [[ $? -eq 0 ]]; then
      # Parse reference
      name=${reference%%:*}
      version=${reference#*:}
      latestVersion=${latestReference#*:}

      # Create a diff using yq as yq doesn't support keeping whitespace or
      # comments around
      diff="$(diff -u --ignore-blank-lines -L "$kustomization" -L "$kustomization" <(grep -v '^$' "$kustomization") <(yq ".images[] |= select(.name == \"$name\" and .newTag == \"$version\").newTag = \"$latestVersion\"" "$kustomization"))"
      # Print and apply the diff
      echo "$diff"
      echo "$diff" | patch --silent --unified --posix --fuzz 3 "$kustomization" -i -
    fi

    # Separate runs by a newline
    echo
  done <<<"$(yq '.images[] | select((.name|type == "!!str") and (.newTag|type == "!!str")) | "\(.name):\(.newTag)"' "$kustomization" | grep -v '^:$')"
done <<<"$(find . -type f -name "$MANIFEST_FILE_NAME")"
