CREATE TABLE IF NOT EXISTS "users" (
    "id" SERIAL PRIMARY KEY,
    "callsign" VARCHAR(10) NOT NULL,
    "hangar" INT NOT NULL
);