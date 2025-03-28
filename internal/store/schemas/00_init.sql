CREATE TABLE revision (
  -- Ensure only a single row can exist
  id INTEGER PRIMARY KEY CHECK (id = 0),
  revision INT NOT NULL
);
INSERT INTO revision (id, revision) VALUES (0, 4);

CREATE TABLE raw_images (
  reference TEXT PRIMARY KEY NOT NULL,
  tags BLOB NOT NULL,
  graph BLOB NOT NULL,
  lastProcessed DATETIME
);

CREATE TABLE images (
  reference TEXT PRIMARY KEY NOT NULL,
  created DATETIME,
  latestReference TEXT,
  latestCreated DATETIME,
  versionDiffSortable INT NOT NULL,
  description TEXT NOT NULL,
  lastModified DATETIME NOT NULL,
  imageUrl TEXT NOT NULL,
  FOREIGN KEY(reference) REFERENCES raw_images(reference) ON DELETE CASCADE
);

CREATE TABLE images_tags (
  reference TEXT NOT NULL,
  tag TEXT NOT NULL,
  PRIMARY KEY (reference, tag),
  FOREIGN KEY(reference) REFERENCES images(reference) ON DELETE CASCADE
);
-- TODO: Remove in v1.
-- There was previously an issue with a missing cascade on images_tags, which
-- could cause images_tags to have rows not related to an entry in images, which
-- would result in counts like "vulnerable images" or "failed images" being off.
-- To remediate this issue, remove all entries without references.
DELETE FROM images_tags WHERE reference NOT IN (SELECT reference FROM images);

-- TODO: Rename in v1. This was done as an easy way to migrate somewhat
-- gracefully without having to drop the entire database
DROP TABLE IF EXISTS images_links;
CREATE TABLE images_linksv2 (
  reference TEXT NOT NULL,
  links BLOB NOT NULL,
  PRIMARY KEY (reference),
  FOREIGN KEY(reference) REFERENCES images(reference) ON DELETE CASCADE
);

CREATE TABLE images_release_notes (
  reference TEXT NOT NULL,
  title TEXT NOT NULL,
  html TEXT NOT NULL,
  markdown TEXT NOT NULL,
  released DATETIME NOT NULL,
  PRIMARY KEY (reference),
  FOREIGN KEY(reference) REFERENCES images(reference) ON DELETE CASCADE
);

CREATE TABLE images_descriptions (
  reference TEXT NOT NULL,
  html TEXT NOT NULL,
  markdown TEXT NOT NULL,
  PRIMARY KEY (reference),
  FOREIGN KEY(reference) REFERENCES images(reference) ON DELETE CASCADE
);

CREATE TABLE images_graphs (
  reference TEXT NOT NULL,
  graph BLOB NOT NULL,
  PRIMARY KEY (reference),
  FOREIGN KEY(reference) REFERENCES images(reference) ON DELETE CASCADE
);

-- TODO: Rename in v1. This was done as an easy way to migrate somewhat
-- gracefully without having to drop the entire database
DROP TABLE IF EXISTS images_vulnerabilities;
CREATE TABLE images_vulnerabilitiesv2 (
  reference TEXT NOT NULL,
  count INT NOT NULL,
  vulnerabilities BLOB NOT NULL,
  PRIMARY KEY (reference),
  FOREIGN KEY(reference) REFERENCES images(reference) ON DELETE CASCADE
);

CREATE TABLE images_scorecards (
  reference TEXT NOT NULL,
  score REAL NOT NULL,
  scorecard BLOB NOT NULL,
  PRIMARY KEY (reference),
  FOREIGN KEY(reference) REFERENCES images(reference) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS images_workflow_runs (
  reference TEXT NOT NULL,
  started DATETIME NOT NULL,
  result TEXT NOT NULL,
  blob BLOB NOT NULL,
  FOREIGN KEY(reference) REFERENCES images(reference) ON DELETE CASCADE
);

CREATE TABLE images_provenance (
  reference TEXT NOT NULL,
  provenance BLOB NOT NULL,
  PRIMARY KEY (reference),
  FOREIGN KEY(reference) REFERENCES images(reference) ON DELETE CASCADE
);

CREATE TABLE images_sbom (
  reference TEXT NOT NULL,
  sbom BLOB NOT NULL,
  PRIMARY KEY (reference),
  FOREIGN KEY(reference) REFERENCES images(reference) ON DELETE CASCADE
);
