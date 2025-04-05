-- Source revision: 4
-- Target revision: 5
-- Summary: Use osv vulnerabilities

-- TODO: Remove in v1
DROP TABLE IF EXISTS images_vulnerabilities;
DROP TABLE IF EXISTS images_vulnerabilitiesv2;

-- TODO: Rename in v1
CREATE TABLE images_vulnerabilitiesv3 (
  reference TEXT NOT NULL,
  count INT NOT NULL,
  vulnerabilities BLOB NOT NULL,
  PRIMARY KEY (reference),
  FOREIGN KEY(reference) REFERENCES images(reference) ON DELETE CASCADE
);

-- Update changes on INSERT to images_vulnerabilities table
-- TODO: Always drop in v1
DROP TRIGGER IF EXISTS images_changes_images_vulnerabilities_insert;
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

-- Update images_updates on UPDATE to images_vulnerabilities table
-- TODO: Always drop in v1
DROP TRIGGER IF EXISTS images_changes_images_vulnerabilities_update;
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

INSERT INTO revision (id, revision) VALUES (0, 5) ON CONFLICT DO UPDATE SET revision=excluded.revision;
