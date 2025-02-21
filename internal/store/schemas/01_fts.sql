CREATE VIRTUAL TABLE images_fts USING FTS5(
  content='images',
  reference,
  description
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
