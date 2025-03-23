-- Track changes
CREATE TABLE images_changes (
  reference TEXT NOT NULL,
  time DATETIME NOT NULL,
  type TEXT NOT NULL,

  changedBasic BOOLEAN NOT NULL DEFAULT FALSE,
  changedLinks BOOLEAN NOT NULL DEFAULT FALSE,
  changedReleaseNotes BOOLEAN NOT NULL DEFAULT FALSE,
  changedDescription BOOLEAN NOT NULL DEFAULT FALSE,
  changedGraph BOOLEAN NOT NULL DEFAULT FALSE,
  changedVulnerabilities BOOLEAN NOT NULL DEFAULT FALSE,
  changedScorecard BOOLEAN NOT NULL DEFAULT FALSE,
  changedProvenance BOOLEAN NOT NULL DEFAULT FALSE,

  FOREIGN KEY(reference) REFERENCES images(reference) ON DELETE CASCADE
);

-- NOTE: There are no ON DELETE triggers as they would insert data on a cascade
-- delete, at which point we already remove the images_changes entry

-- Update changes on INSERT to images table
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

-- Update images_updates on UPDATE to images table
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

-- Update changes on INSERT to images_links table
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

-- Update images_updates on UPDATE to images_links table
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

-- Update changes on INSERT to images_release_notes table
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

-- Update images_updates on UPDATE to images_release_notes table
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

-- Update changes on INSERT to images_descriptions table
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

-- Update images_updates on UPDATE to images_descriptions table
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

-- Update changes on INSERT to images_graphs table
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

-- Update images_updates on UPDATE to images_graphs table
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

-- Update changes on INSERT to images_vulnerabilities table
CREATE TRIGGER images_changes_images_vulnerabilities_insert AFTER INSERT ON images_vulnerabilitiesv2 BEGIN
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

-- Update images_updates on UPDATE to images_vulnerabilities table
CREATE TRIGGER images_changes_images_vulnerabilities_update AFTER UPDATE ON images_vulnerabilitiesv2 WHEN
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

-- Update changes on INSERT to images_scorecards table
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

-- Update images_updates on UPDATE to images_scorecards table
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

-- Update changes on INSERT to images_provenance table
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

-- Update images_updates on UPDATE to images_provenance table
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
