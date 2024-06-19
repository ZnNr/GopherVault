create table if not exists credentials (
                                           id SERIAL PRIMARY KEY,
                                           user_name TEXT UNIQUE NOT NULL,
                                           login TEXT UNIQUE NOT NULL,
                                           password TEXT NOT NULL,
                                           metadata TEXT
)