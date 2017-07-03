CREATE TABLE ad_keywords (
  id SERIAL PRIMARY KEY,
  ad_id INTEGER REFERENCES ads (id),
  keyword_id INTEGER REFERENCES keywords (id),
  position INTEGER NOT NULL,
  position_count INTEGER NOT NULL DEFAULT 1,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX ad_keywords_ad_id_keyword_id_position_index ON ad_keywords (ad_id, keyword_id, position);