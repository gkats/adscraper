CREATE TABLE ads (
  id SERIAL PRIMARY KEY,
  headline1 VARCHAR(50) NOT NULL,
  headline2 VARCHAR(50) NOT NULL,
  description VARCHAR(100) NOT NULL,
  path VARCHAR NOT NULL,
  rest TEXT,
  raw TEXT,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX ads_h1_h2_desc_index ON ads (headline1, headline2, description);