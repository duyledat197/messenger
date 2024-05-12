CREATE TABLE IF NOT EXISTS users(
  id text PRIMARY KEY,
  username text UNIQUE,
  password text,
  email text,
  facebook text,
  gmail text,
  discord text
);

