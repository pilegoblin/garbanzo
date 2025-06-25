--
-- We'll start with the user table, the foundation of the application.
--
CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar UNIQUE NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "auth_id" varchar UNIQUE NOT NULL,
  "avatar_url" varchar,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "user_color" varchar NOT NULL DEFAULT '000000'
);

--
-- Pods are the main communities where users gather.
-- Each pod is owned by a single user.
--
CREATE TABLE "pods" (
  "id" bigserial PRIMARY KEY,
  "owner_id" bigint NOT NULL,
  "name" varchar NOT NULL,
  "invite_code" varchar,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

--
-- This is a "join table" that connects users to the pods they are members of.
-- This allows for a many-to-many relationship between users and pods.
--
CREATE TABLE "pod_members" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "pod_id" bigint NOT NULL,
  "joined_at" timestamptz NOT NULL DEFAULT (now())
);

--
-- Beans exist within a pod. They can be for text or voice chat.
--
CREATE TABLE "beans" (
  "id" bigserial PRIMARY KEY,
  "pod_id" bigint NOT NULL,
  "name" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

--
-- Messages are posted by users within a specific bean.
--
CREATE TABLE "messages" (
  "id" text PRIMARY KEY,
  "bean_id" bigint NOT NULL,
  "author_id" bigint NOT NULL,
  "content" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz -- To track if a message has been edited
);

--
-- Foreign Key Constraints ensure data integrity. For example, a pod's owner_id
-- must correspond to an actual user in the "users" table.
--
ALTER TABLE "pods" ADD FOREIGN KEY ("owner_id") REFERENCES "users" ("id") ON DELETE CASCADE;
ALTER TABLE "pod_members" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;
ALTER TABLE "pod_members" ADD FOREIGN KEY ("pod_id") REFERENCES "pods" ("id") ON DELETE CASCADE;
ALTER TABLE "beans" ADD FOREIGN KEY ("pod_id") REFERENCES "pods" ("id") ON DELETE CASCADE;
ALTER TABLE "messages" ADD FOREIGN KEY ("bean_id") REFERENCES "beans" ("id") ON DELETE CASCADE;
ALTER TABLE "messages" ADD FOREIGN KEY ("author_id") REFERENCES "users" ("id") ON DELETE SET NULL; -- If a user is deleted, their messages can remain but without an author.

--
-- Indexes are crucial for query performance, especially on foreign keys and
-- frequently queried columns. A unique index on pod_members prevents
-- a user from joining the same pod more than once.
--
CREATE INDEX ON "pods" ("owner_id");
CREATE INDEX ON "pod_members" ("user_id");
CREATE INDEX ON "pod_members" ("pod_id");
CREATE UNIQUE INDEX ON "pod_members" ("user_id", "pod_id");
CREATE INDEX ON "beans" ("pod_id");
CREATE INDEX ON "messages" ("bean_id");
CREATE INDEX ON "messages" ("author_id");