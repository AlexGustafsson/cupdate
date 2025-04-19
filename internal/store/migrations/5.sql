-- Source revision: 5
-- Target revision: 6
-- Summary: Cascade delete tags (#143)

-- TODO: Remove in v1

-- Clean up
DELETE FROM images_tags WHERE reference NOT IN (SELECT reference FROM images);

-- Cascade can't be added in ALTER in sqlite
ALTER TABLE images_tags RENAME TO images_tags_old;

CREATE TABLE images_tags (
  reference TEXT NOT NULL,
  tag TEXT NOT NULL,
  PRIMARY KEY (reference, tag),
  FOREIGN KEY(reference) REFERENCES images(reference) ON DELETE CASCADE
);

INSERT INTO images_tags SELECT * FROM images_tags_old;

INSERT INTO revision (id, revision) VALUES (0, 6) ON CONFLICT DO UPDATE SET revision=excluded.revision;
