CREATE TABLE IF NOT EXISTS github_advisories (
  id TEXT PRIMARY KEY NOT NULL,
  repository NOT NULL,
  published DATETIME NOT NULL,
  severity TEXT NOT NULL,
  introduced_version TEXT NOT NULL,
  fixed_version TEXT
);
