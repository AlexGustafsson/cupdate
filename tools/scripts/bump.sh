#!/usr/bin/env bash

# usage: ./tools/scripts/bump.sh <old> <new>
from="$1"
to="$2"

# Files to edit:
# .github/ISSUE_TEMPLATE/server-bug.yaml
# .github/ISSUE_TEMPLATE/ui-bug.yaml
# deploy/base/deployment.yaml
# docs/docker/compose.yaml
# docs/docker/README.md
# docs/cookbook/README.md
# docs/cookbook/update-compose-files.sh
# docs/cookbook/update-kustomizations.sh
# docs/kubernetes/kustomization.yaml
# docs/kubernetes/README.md
# docs/podman/README.md
# README.md

if [[ -z "$from" ]] || [[ -z "$to" ]]; then
	echo "usage: $0 <old> <new>"
	exit 1
fi

from=${from//./\\.}
to=${to//./\\.}

files_to_edit="$(sed '1,/# Files to edit:/d' "$0" | sed '/^$/q' | sed '/^$/d' | cut -c3-)"

while read -r file; do
	sed -i "s/$from/$to/g" "$file"
done <<<"$files_to_edit"
