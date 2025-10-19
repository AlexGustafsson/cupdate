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

# Revision must match the last migration
# shellcheck disable=SC2012
last_migration_file="$(ls -v1 internal/store/migrations | tail -n1 | cut -d'.' -f1)"
grep "internal/store/migrations/$last_migration_file.sql" -e "^INSERT INTO revision (id, revision) VALUES (0, $goRevision);$" -e "^INSERT INTO revision (id, revision) VALUES (0, $goRevision) ON CONFLICT DO UPDATE SET revision=excluded.revision;$" &>/dev/null
if [[ $? -gt 0 ]]; then
	echo "Go revision ($goRevision) not matched by the last migration: $last_migration_file"
fi
