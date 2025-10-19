-- Source revision: 7
-- Target revision: 8
-- Summary: Track annotations

-- TODO: Remove in v1

ALTER TABLE images ADD annotations BLOB;
ALTER TABLE images ADD latestAnnotations BLOB;

INSERT INTO revision (id, revision) VALUES (0, 8) ON CONFLICT DO UPDATE SET revision=excluded.revision;
