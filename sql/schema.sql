CREATE TABLE IF NOT EXISTS users(
    id serial PRIMARY KEY,
    email text NOT NULL UNIQUE,
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP,
    picture text,
    username text NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS pods(
    id serial PRIMARY KEY,
    pod_name text NOT NULL,
    pod_owner integer REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS pod_users(
    user_id integer REFERENCES users(id),
    pod_id integer REFERENCES pods(id),
    PRIMARY KEY (user_id, pod_id)
);

CREATE TABLE IF NOT EXISTS channels(
    id serial PRIMARY KEY,
    channel_name text NOT NULL,
    pod_id integer NOT NULL REFERENCES pods(id)
);

CREATE TABLE IF NOT EXISTS posts(
    id serial PRIMARY KEY,
    content text NOT NULL,
    user_id integer NOT NULL REFERENCES users(id),
    channel_id integer NOT NULL REFERENCES channels(id),
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP
);

