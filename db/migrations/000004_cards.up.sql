CREATE TABLE IF NOT EXISTS cards (
                                     id SERIAL PRIMARY KEY,
                                     user_name TEXT NOT NULL,
                                     bank_name TEXT NOT NULL,
                                     number TEXT,
                                     cv TEXT,
                                     password TEXT,
                                     card_type TEXT,
                                     metadata TEXT
);