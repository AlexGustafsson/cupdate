CREATE TABLE IF NOT EXISTS raw_images (
  reference TEXT PRIMARY KEY NOT NULL,
  tags BLOB NOT NULL,
  graph BLOB NOT NULL,
  lastProcessed DATETIME
);

CREATE TABLE IF NOT EXISTS images (
  reference TEXT PRIMARY KEY NOT NULL,
  latestReference TEXT NOT NULL,
  description TEXT NOT NULL,
  lastModified DATETIME NOT NULL,
  imageUrl TEXT NOT NULL,
  FOREIGN KEY(reference) REFERENCES raw_images(reference) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS tags (
  name TEXT PRIMARY KEY NOT NULL,
  color TEXT NOT NULL,
  description TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS images_tags (
  reference TEXT NOT NULL,
  tag TEXT NOT NULL,
  PRIMARY KEY (reference, tag),
  FOREIGN KEY(reference) REFERENCES images(reference) ON DELETE CASCADE
  FOREIGN KEY(tag) REFERENCES tags(name) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS images_links (
  reference TEXT NOT NULL,
  url TEXT NOT NULL,
  type TEXT NOT NULL,
  PRIMARY KEY (reference, url),
  FOREIGN KEY(reference) REFERENCES images(reference) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS images_release_notes (
  reference TEXT NOT NULL,
  title TEXT NOT NULL,
  html TEXT NOT NULL,
  markdown TEXT NOT NULL,
  released DATETIME NOT NULL,
  PRIMARY KEY (reference),
  FOREIGN KEY(reference) REFERENCES images(reference) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS images_descriptions (
  reference TEXT NOT NULL,
  html TEXT NOT NULL,
  markdown TEXT NOT NULL,
  PRIMARY KEY (reference),
  FOREIGN KEY(reference) REFERENCES images(reference) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS images_graphs (
  reference TEXT NOT NULL,
  graph BLOB NOT NULL,
  PRIMARY KEY (reference),
  FOREIGN KEY(reference) REFERENCES images(reference) ON DELETE CASCADE
);