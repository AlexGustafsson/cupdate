#!/usr/bin/env bash

goRevision="$(grep -e 'const Revision = \([0-9]\+\)' internal/store/init.go | cut -d' ' -f 4)"

# Revision must be found
if [[ -z "$goRevision" ]]; then
	echo "Cannot identify the store revision"
	echo "Go revision: $goRevision"
	exit 1
fi

# Migrations must exist
for i in $(seq 0 "$((goRevision - 1))"); do
	path="internal/store/migrations/$i.sql"
	if [[ ! -f "$path" ]]; then
		echo "Missing migration: $path"
		exit 1
	fi
done

# Revision must be updated by migrations
for i in $(seq 0 "$((goRevision - 1))"); do
	path="internal/store/migrations/$i.sql"
	grep "$path" -e "^INSERT INTO revision (id, revision) VALUES (0, $((i + 1)));$" -e "^INSERT INTO revision (id, revision) VALUES (0, $((i + 1))) ON CONFLICT DO UPDATE SET revision=excluded.revision;$" &>/dev/null
	if [[ $? -gt 0 ]]; then
		echo "Missing revision update in migration script: $path"
	fi
done
