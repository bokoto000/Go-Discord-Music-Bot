CREATE TABLE IF NOT EXISTS guilds (
    id serial NOT NULL,
    guild_id text NOT NULL,
    CONSTRAINT guilds_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS keywords(    
    id serial NOT NULL,
    guild_id text NOT NULL,
    keyword_key UNIQUE text,
    keyword_value text,
    CONSTRAINT keywords_pkey PRIMARY KEY (id)
);