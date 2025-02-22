#!/usr/bin/env bash

goRevision="$(grep -e 'const Revision = \([0-9]\+\)' internal/store/init.go | cut -d' ' -f 4)"
sqlRevision="$(grep -e 'INSERT INTO revision (id, revision) VALUES (0, \([0-9]\+\));' internal/store/schemas/00_init.sql | sed 's/[^0-9 ]//g' | cut -d' ' -f 8)"

# Revision must be found
if [[ -z "$goRevision" ]] || [[ -z "$sqlRevision" ]]; then
  echo "Cannot identify the store revision"
  echo "Go revision: $goRevision"
  echo "SQL revision: $sqlRevision"
  exit 1
fi

# Revisions must match
if [[ ! "$goRevision" = "$sqlRevision" ]]; then
  echo "The revision in the init scripts does not match the one of the go code ($sqlRevision vs $goRevision)"
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
