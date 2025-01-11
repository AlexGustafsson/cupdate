# Vulndb

Vulndb is a tiny sqlite file that contains information useful to statically look
up known vulnerabilities in container images based on their source repositories.
For now it uses recent GitHub's advisory database.

The database is updated daily and published as an OCI artifact used by Cupdate.
The artifact is available here:
<https://github.com/AlexGustafsson/cupdate/pkgs/container/cupdate%2Fvulndb>.

For more advanced scanning requirements, use something like
[Trivy](https://github.com/aquasecurity/trivy) or
[Grype](https://github.com/anchore/grype).

## Running

When run, vulndb will download its source data, compile it and push an OCI
artifactory. It is intended to run as a GitHub action.

```shell
INPUT_GITHUB_ACTOR="..." INPUT_GITHUB_TOKEN="..." go run tools/vulndb/*.go
```

## Schema

See [internal/db/createTablesIfNotExist.sql](internal/db/createTablesIfNotExist.sql).

## Data sources

- GitHub Advisory Database: <https://github.com/github/advisory-database>
