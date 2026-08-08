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
  imageUrl TEXT NOT NULL, annotations BLOB, latestAnnotations BLOB,
  FOREIGN KEY(reference) REFERENCES raw_images(reference) ON DELETE CASCADE
);
CREATE VIRTUAL TABLE images_fts USING FTS5(
  content='images',
  reference,
  description
)
/* images_fts(reference,description) */;
CREATE TABLE IF NOT EXISTS 'images_fts_data'(id INTEGER PRIMARY KEY, block BLOB);
CREATE TABLE IF NOT EXISTS 'images_fts_idx'(segid, term, pgno, PRIMARY KEY(segid, term)) WITHOUT ROWID;
CREATE TABLE IF NOT EXISTS 'images_fts_docsize'(id INTEGER PRIMARY KEY, sz BLOB);
CREATE TABLE IF NOT EXISTS 'images_fts_config'(k PRIMARY KEY, v) WITHOUT ROWID;
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
CREATE TABLE images_workflow_runs (
  reference TEXT NOT NULL,
  started DATETIME NOT NULL,
  result TEXT NOT NULL,
  blob BLOB NOT NULL,
  FOREIGN KEY(reference) REFERENCES images(reference) ON DELETE CASCADE
);
CREATE TABLE images_changes (
  reference TEXT NOT NULL,
  time DATETIME NOT NULL,
  type TEXT NOT NULL,

  changedBasic BOOLEAN NOT NULL DEFAULT FALSE,
  changedLinks BOOLEAN NOT NULL DEFAULT FALSE,
  changedReleaseNotes BOOLEAN NOT NULL DEFAULT FALSE,
  changedDescription BOOLEAN NOT NULL DEFAULT FALSE,
  changedGraph BOOLEAN NOT NULL DEFAULT FALSE,
  changedVulnerabilities BOOLEAN NOT NULL DEFAULT FALSE, changedScorecard BOOLEAN NOT NULL DEFAULT FALSE, changedProvenance BOOLEAN NOT NULL DEFAULT FALSE, changedSBOM BOOLEAN NOT NULL DEFAULT FALSE,

  FOREIGN KEY(reference) REFERENCES images(reference) ON DELETE CASCADE
);
CREATE TRIGGER images_fts_insert AFTER INSERT ON images BEGIN
  INSERT INTO images_fts(rowid, reference, description) VALUES (new.rowid, new.reference, new.description);
END;
CREATE TRIGGER images_fts_delete AFTER DELETE ON images BEGIN
  INSERT INTO images_fts(images_fts, rowid, reference, description) VALUES('delete', old.rowid, old.reference, old.description);
END;
CREATE TRIGGER images_fts_update AFTER UPDATE ON images BEGIN
  INSERT INTO images_fts(images_fts, rowid, reference, description) VALUES('delete', old.rowid, old.reference, old.description);
  INSERT INTO images_fts(rowid, reference, description) VALUES (new.rowid, new.reference, new.description);
END;
CREATE TRIGGER images_changes_images_insert AFTER INSERT ON images BEGIN
  INSERT INTO images_changes(
    reference,
    time,
    type,

    changedBasic
  ) VALUES (
    new.reference,
    datetime('now', 'subsecond'),
    "insert",

    TRUE
  );
END;
CREATE TRIGGER images_changes_images_update AFTER UPDATE ON images WHEN
    old.latestReference <> new.latestReference
    OR old.description <> new.description
    OR old.imageUrl <> new.imageUrl
  BEGIN
  INSERT INTO images_changes(
    reference,
    time,
    type,

    changedBasic
  ) VALUES (
    new.reference,
    datetime('now', 'subsecond'),
    "update",

    TRUE
  );
END;
CREATE TRIGGER images_changes_images_links_insert AFTER INSERT ON images_linksv2 BEGIN
  INSERT INTO images_changes(
    reference,
    time,
    type,

    changedLinks
  ) VALUES (
    new.reference,
    datetime('now', 'subsecond'),
    "insert",

    TRUE
  );
END;
CREATE TRIGGER images_changes_images_links_update AFTER UPDATE ON images_linksv2 WHEN
    old.links <> new.links
  BEGIN
  INSERT INTO images_changes(
    reference,
    time,
    type,

    changedLinks
  ) VALUES (
    new.reference,
    datetime('now', 'subsecond'),
    "update",

    TRUE
  );
END;
CREATE TRIGGER images_changes_images_release_notes_insert AFTER INSERT ON images_release_notes BEGIN
  INSERT INTO images_changes(
    reference,
    time,
    type,

    changedReleaseNotes
  ) VALUES (
    new.reference,
    datetime('now', 'subsecond'),
    "insert",

    TRUE
  );
END;
CREATE TRIGGER images_changes_images_release_notes_update AFTER UPDATE ON images_release_notes WHEN
    old.title <> new.title
    OR old.html <> new.html
    OR old.markdown <> new.markdown
    OR old.released <> new.released
  BEGIN
  INSERT INTO images_changes(
    reference,
    time,
    type,

    changedReleaseNotes
  ) VALUES (
    new.reference,
    datetime('now', 'subsecond'),
    "update",

    TRUE
  );
END;
CREATE TRIGGER images_changes_images_descriptions_insert AFTER INSERT ON images_descriptions BEGIN
  INSERT INTO images_changes(
    reference,
    time,
    type,

    changedDescription
  ) VALUES (
    new.reference,
    datetime('now', 'subsecond'),
    "insert",

    TRUE
  );
END;
CREATE TRIGGER images_changes_images_descriptions_update AFTER UPDATE ON images_descriptions WHEN
    old.html <> new.html
    OR old.markdown <> new.markdown
  BEGIN
  INSERT INTO images_changes(
    reference,
    time,
    type,

    changedDescription
  ) VALUES (
    new.reference,
    datetime('now', 'subsecond'),
    "update",

    TRUE
  );
END;
CREATE TRIGGER images_changes_images_graphs_insert AFTER INSERT ON images_graphs BEGIN
  INSERT INTO images_changes(
    reference,
    time,
    type,

    changedGraph
  ) VALUES (
    new.reference,
    datetime('now', 'subsecond'),
    "insert",

    TRUE
  );
END;
CREATE TRIGGER images_changes_images_graphs_update AFTER UPDATE ON images_graphs WHEN
    old.graph <> new.graph
  BEGIN
  INSERT INTO images_changes(
    reference,
    time,
    type,

    changedGraph
  ) VALUES (
    new.reference,
    datetime('now', 'subsecond'),
    "update",

    TRUE
  );
END;
CREATE TABLE revision (
  -- Ensure only a single row can exist
  id INTEGER PRIMARY KEY CHECK (id = 0),
  revision INT NOT NULL
);
CREATE TABLE images_scorecards (
  reference TEXT NOT NULL,
  score REAL NOT NULL,
  scorecard BLOB NOT NULL,
  PRIMARY KEY (reference),
  FOREIGN KEY(reference) REFERENCES images(reference) ON DELETE CASCADE
);
CREATE TRIGGER images_changes_images_scorecards_insert AFTER INSERT ON images_scorecards BEGIN
  INSERT INTO images_changes(
    reference,
    time,
    type,

    changedScorecard
  ) VALUES (
    new.reference,
    datetime('now', 'subsecond'),
    "insert",

    TRUE
  );
END;
CREATE TRIGGER images_changes_images_scorecards_update AFTER UPDATE ON images_scorecards WHEN
    old.scorecard <> new.scorecard
  BEGIN
  INSERT INTO images_changes(
    reference,
    time,
    type,

    changedScorecard
  ) VALUES (
    new.reference,
    datetime('now', 'subsecond'),
    "update",

    TRUE
  );
END;
CREATE TABLE images_provenance (
  reference TEXT NOT NULL,
  provenance BLOB NOT NULL,
  PRIMARY KEY (reference),
  FOREIGN KEY(reference) REFERENCES images(reference) ON DELETE CASCADE
);
CREATE TRIGGER images_changes_images_provenance_insert AFTER INSERT ON images_provenance BEGIN
  INSERT INTO images_changes(
    reference,
    time,
    type,

    changedProvenance
  ) VALUES (
    new.reference,
    datetime('now', 'subsecond'),
    "insert",

    TRUE
  );
END;
CREATE TRIGGER images_changes_images_provenance_update AFTER UPDATE ON images_provenance WHEN
    old.provenance <> new.provenance
  BEGIN
  INSERT INTO images_changes(
    reference,
    time,
    type,

    changedProvenance
  ) VALUES (
    new.reference,
    datetime('now', 'subsecond'),
    "update",

    TRUE
  );
END;
CREATE TABLE images_sbom (
  reference TEXT NOT NULL,
  sbom BLOB NOT NULL,
  PRIMARY KEY (reference),
  FOREIGN KEY(reference) REFERENCES images(reference) ON DELETE CASCADE
);
CREATE TRIGGER images_changes_images_sbom_insert AFTER INSERT ON images_sbom BEGIN
  INSERT INTO images_changes(
    reference,
    time,
    type,

    changedSBOM
  ) VALUES (
    new.reference,
    datetime('now', 'subsecond'),
    "insert",

    TRUE
  );
END;
CREATE TRIGGER images_changes_images_sbom_update AFTER UPDATE ON images_sbom WHEN
    old.sbom <> new.sbom
  BEGIN
  INSERT INTO images_changes(
    reference,
    time,
    type,

    changedSBOM
  ) VALUES (
    new.reference,
    datetime('now', 'subsecond'),
    "update",

    TRUE
  );
END;
CREATE TABLE images_vulnerabilitiesv3 (
  reference TEXT NOT NULL,
  count INT NOT NULL,
  vulnerabilities BLOB NOT NULL,
  PRIMARY KEY (reference),
  FOREIGN KEY(reference) REFERENCES images(reference) ON DELETE CASCADE
);
CREATE TRIGGER images_changes_images_vulnerabilities_insert AFTER INSERT ON images_vulnerabilitiesv3 BEGIN
  INSERT INTO images_changes(
    reference,
    time,
    type,

    changedVulnerabilities
  ) VALUES (
    new.reference,
    datetime('now', 'subsecond'),
    "insert",

    TRUE
  );
END;
CREATE TRIGGER images_changes_images_vulnerabilities_update AFTER UPDATE ON images_vulnerabilitiesv3 WHEN
    old.vulnerabilities <> new.vulnerabilities
  BEGIN
  INSERT INTO images_changes(
    reference,
    time,
    type,

    changedVulnerabilities
  ) VALUES (
    new.reference,
    datetime('now', 'subsecond'),
    "update",

    TRUE
  );
END;
CREATE TABLE images_tags (
  reference TEXT NOT NULL,
  tag TEXT NOT NULL,
  PRIMARY KEY (reference, tag),
  FOREIGN KEY(reference) REFERENCES images(reference) ON DELETE CASCADE
);
CREATE TABLE images_updates (
  newReference TEXT PRIMARY KEY NOT NULL,
  newAnnotations BLOB,
  oldReference TEXT NOT NULL,
  oldAnnotations BLOB,
  versionDiffSortable INT NOT NULL,
  identified DATETIME NOT NULL,
  released DATETIME
);
CREATE TRIGGER images_changes_images_insert_version AFTER INSERT ON images WHEN
    new.latestReference <> ""
    AND new.latestReference <> new.reference
  BEGIN
  INSERT INTO images_updates(
    newReference,
    newAnnotations,
    oldReference,
    oldAnnotations,
    versionDiffSortable,
    identified,
    released
  ) VALUES (
    new.latestReference,
    new.latestAnnotations,
    new.reference,
    new.annotations,
    new.versionDiffSortable,
    datetime('now', 'subsecond'),
    new.latestCreated
  );
END;
CREATE TRIGGER images_changes_images_update_version AFTER UPDATE ON images WHEN
    new.latestReference <> ""
    AND new.latestReference <> old.latestReference
    AND new.latestReference <> new.reference
  BEGIN
  INSERT INTO images_updates(
    newReference,
    newAnnotations,
    oldReference,
    oldAnnotations,
    versionDiffSortable,
    identified,
    released
  ) VALUES (
    new.latestReference,
    new.latestAnnotations,
    new.reference,
    new.annotations,
    new.versionDiffSortable,
    datetime('now', 'subsecond'),
    new.latestCreated
  );
END;
