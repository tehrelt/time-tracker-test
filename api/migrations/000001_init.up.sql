CREATE TABLE IF NOT EXISTS "users" (
  "id" VARCHAR NOT NULL PRIMARY KEY,
  "surname" VARCHAR NOT NULL,
  "name" VARCHAR NOT NULL,
  "patronymic" VARCHAR NOT NULL,
  "address" VARCHAR NOT NULL,
  "passport_serie" VARCHAR NOT NULL,
  "passport_number" VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS "activity" (
  "id" SERIAL NOT NULL PRIMARY KEY,
  "user_id" VARCHAR NOT NULL REFERENCES "users"("id"),
  "start_time" TIMESTAMP NOT NULL,
  "end_time" TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS "users_passport_uindex" ON "users"("passport_serie", "passport_number");