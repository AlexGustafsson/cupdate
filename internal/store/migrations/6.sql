-- Source revision: 6
-- Target revision: 7
-- Summary: Remove unused temporary table from 5.sql

-- TODO: Remove in v1

-- Clean up
DROP TABLE images_tags_old;

INSERT INTO revision (id, revision) VALUES (0, 7) ON CONFLICT DO UPDATE SET revision=excluded.revision;
