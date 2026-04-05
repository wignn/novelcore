-- Genres lookup table
CREATE TABLE IF NOT EXISTS genres (
  id SERIAL PRIMARY KEY,
  name VARCHAR(64) NOT NULL UNIQUE,
  slug VARCHAR(64) NOT NULL UNIQUE
);

-- Tags lookup table
CREATE TABLE IF NOT EXISTS tags (
  id SERIAL PRIMARY KEY,
  name VARCHAR(64) NOT NULL UNIQUE,
  slug VARCHAR(64) NOT NULL UNIQUE
);

-- Authors table
CREATE TABLE IF NOT EXISTS authors (
  id CHAR(27) PRIMARY KEY,
  name VARCHAR(128) NOT NULL,
  bio TEXT DEFAULT '',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Translation Groups table
CREATE TABLE IF NOT EXISTS translation_groups (
  id CHAR(27) PRIMARY KEY,
  name VARCHAR(128) NOT NULL UNIQUE,
  website_url TEXT DEFAULT '',
  description TEXT DEFAULT '',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Main novels table
CREATE TABLE IF NOT EXISTS novels (
  id CHAR(27) PRIMARY KEY,
  title VARCHAR(256) NOT NULL,
  alternative_title TEXT DEFAULT '',
  description TEXT DEFAULT '',
  cover_image_url TEXT DEFAULT '',
  author_id CHAR(27) REFERENCES authors(id),
  status VARCHAR(32) DEFAULT 'ongoing',
  novel_type VARCHAR(32) DEFAULT 'web_novel',
  country_of_origin VARCHAR(64) DEFAULT '',
  year_published INT,
  total_chapters INT DEFAULT 0,
  rating_avg DECIMAL(3,2) DEFAULT 0.00,
  rating_count INT DEFAULT 0,
  view_count BIGINT DEFAULT 0,
  bookmark_count INT DEFAULT 0,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Novel-Genre many-to-many
CREATE TABLE IF NOT EXISTS novel_genres (
  novel_id CHAR(27) REFERENCES novels(id) ON DELETE CASCADE,
  genre_id INT REFERENCES genres(id) ON DELETE CASCADE,
  PRIMARY KEY (novel_id, genre_id)
);

-- Novel-Tag many-to-many
CREATE TABLE IF NOT EXISTS novel_tags (
  novel_id CHAR(27) REFERENCES novels(id) ON DELETE CASCADE,
  tag_id INT REFERENCES tags(id) ON DELETE CASCADE,
  PRIMARY KEY (novel_id, tag_id)
);

-- Chapters table
CREATE TABLE IF NOT EXISTS chapters (
  id CHAR(27) PRIMARY KEY,
  novel_id CHAR(27) REFERENCES novels(id) ON DELETE CASCADE,
  chapter_number DECIMAL(10,1) NOT NULL,
  title VARCHAR(256) DEFAULT '',
  translator_group_id CHAR(27) REFERENCES translation_groups(id),
  source_url TEXT DEFAULT '',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Novel daily stats for ranking
CREATE TABLE IF NOT EXISTS novel_daily_stats (
  id SERIAL PRIMARY KEY,
  novel_id CHAR(27) REFERENCES novels(id) ON DELETE CASCADE,
  stat_date DATE NOT NULL,
  views INT DEFAULT 0,
  bookmarks_added INT DEFAULT 0,
  reviews_added INT DEFAULT 0,
  UNIQUE(novel_id, stat_date)
);

-- Indexes
CREATE INDEX idx_chapters_novel_id ON chapters(novel_id);
CREATE INDEX idx_chapters_novel_number ON chapters(novel_id, chapter_number);
CREATE INDEX idx_novels_status ON novels(status);
CREATE INDEX idx_novels_country ON novels(country_of_origin);
CREATE INDEX idx_novels_rating ON novels(rating_avg DESC);
CREATE INDEX idx_novels_views ON novels(view_count DESC);
CREATE INDEX idx_novels_bookmarks ON novels(bookmark_count DESC);
CREATE INDEX idx_novels_updated ON novels(updated_at DESC);
CREATE INDEX idx_daily_stats_date ON novel_daily_stats(stat_date);

-- Trigger for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = CURRENT_TIMESTAMP;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_novels_updated_at BEFORE UPDATE ON novels
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER set_chapters_updated_at BEFORE UPDATE ON chapters
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Seed genres
INSERT INTO genres (name, slug) VALUES
  ('Action', 'action'), ('Adventure', 'adventure'), ('Comedy', 'comedy'),
  ('Drama', 'drama'), ('Fantasy', 'fantasy'), ('Horror', 'horror'),
  ('Martial Arts', 'martial-arts'), ('Mystery', 'mystery'),
  ('Romance', 'romance'), ('Sci-fi', 'sci-fi'), ('Slice of Life', 'slice-of-life'),
  ('Supernatural', 'supernatural'), ('Tragedy', 'tragedy'),
  ('Wuxia', 'wuxia'), ('Xianxia', 'xianxia'), ('Xuanhuan', 'xuanhuan'),
  ('Harem', 'harem'), ('Historical', 'historical'), ('Josei', 'josei'),
  ('Mature', 'mature'), ('Mecha', 'mecha'), ('Psychological', 'psychological'),
  ('School Life', 'school-life'), ('Seinen', 'seinen'), ('Shoujo', 'shoujo'),
  ('Shounen', 'shounen'), ('Gender Bender', 'gender-bender'), ('Ecchi', 'ecchi'),
  ('Sports', 'sports'), ('Smut', 'smut')
ON CONFLICT (name) DO NOTHING;

-- Seed tags
INSERT INTO tags (name, slug) VALUES
  ('Male Protagonist', 'male-protagonist'), ('Female Protagonist', 'female-protagonist'),
  ('Overpowered Protagonist', 'overpowered-protagonist'),
  ('Reincarnation', 'reincarnation'), ('Transmigration', 'transmigration'),
  ('System', 'system'), ('Level System', 'level-system'),
  ('Cultivation', 'cultivation'), ('Revenge', 'revenge'),
  ('Second Chance', 'second-chance'), ('Weak to Strong', 'weak-to-strong'),
  ('Dense Protagonist', 'dense-protagonist'),
  ('Academy', 'academy'), ('Monsters', 'monsters'),
  ('Dungeons', 'dungeons'), ('Isekai', 'isekai'),
  ('Villainess', 'villainess'), ('Regression', 'regression'),
  ('Kingdom Building', 'kingdom-building'), ('Game Elements', 'game-elements')
ON CONFLICT (name) DO NOTHING;
