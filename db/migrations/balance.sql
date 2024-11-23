INSERT INTO users (username, balance) VALUES ('user1', 100.0), ('user2', 100.0) ON CONFLICT (username) DO NOTHING;
