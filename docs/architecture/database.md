# Database

Cupdate uses sqlite to persist data.

<!--
| Meaning      | Left  | Right |
| ------------ | ----- | ----- |
| Zero or one  |   |o  |   o|  |
| Exactly      |   ||  |   ||  |
| Zero or more |   }o  |   o{  |
| One or more  |   }|  |   |{  |
 -->

```mermaid
erDiagram
  raw_images {
    TEXT reference PK "OCI reference"
    BLOB graph "JSON-encoded graph"
    BLOB tags "JSON-encoded list of tags"
  }

  images ||--|| raw_images : maps
  images {
    TEXT reference PK, FK "OCI reference"
    TEXT latestReference "OCI refrence of latest version"
    TEXT description "Image description"
    TEXT imageUrl "URL to a image"
    DATETIME lastModified "When the entry was last modified"
  }

  tags {
    TEXT name PK "Tag name"
    TEXT color "Tag's CSS color"
    TEXT description "Tag description"
  }

  images ||--o{ images_tags : has
  tags ||--|| images_tags : maps
  images_tags {
    TEXT reference PK, FK "OCI reference"
    TEXT tag PK, FK "Tag name"
  }

  images ||--o{ images_links : has
  images_links {
    TEXT reference PK, FK "OCI reference"
    TEXT url PK "Link URL"
    TEXT type "Link type, such as 'docker'"
  }

  images ||--|| images_release_notes : maps
  images_release_notes {
    TEXT reference PK, FK "OCI reference"
    TEXT title "Title of release"
    TEXT html "HTML body content"
    TEXT markdown "Markdown body content"
    DATETIME release "Time of release"
  }

  images ||--|| images_descriptions : maps
  images_descriptions {
    TEXT reference PK, FK "OCI reference"
    TEXT html "HTML body content"
    TEXT markdown "Markdown body content"
  }

  images ||--|| images_graphs : maps
  images_graphs {
    TEXT reference PK, FK "OCI reference"
    BLOB graph "JSON-encoded graph"
  }
```
