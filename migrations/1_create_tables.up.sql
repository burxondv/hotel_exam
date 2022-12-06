CREATE TABLE "hotels"(
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR NOT NULL,
    "address" VARCHAR NOT NULL,
    "stars_popular" INTEGER NOT NULL,
    "image_url" VARCHAR NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "rooms"(
    "id" SERIAL PRIMARY KEY,
    "hotel_id" INTEGER NOT NULL REFERENCES hotels(id),
    "status" BOOLEAN NOT NULL,
    "image_url" VARCHAR NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "guests"(
    "id" SERIAL PRIMARY KEY,
    "hotel_id" INTEGER NOT NULL REFERENCES hotels(id),
    "room_id" INTEGER NOT NULL REFERENCES rooms(id),
    "first_name" VARCHAR NOT NULL,
    "last_name" VARCHAR NOT NULL,
    "phone_number" VARCHAR UNIQUE NOT NULL,
    "email" VARCHAR UNIQUE NOT NULL,
    "password" VARCHAR NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "owners"(
    "id" SERIAL PRIMARY KEY,
    "hotel_id" INTEGER NOT NULL REFERENCES hotels(id),
    "first_name" VARCHAR NOT NULL,
    "last_name" VARCHAR NOT NULL,
    "email" VARCHAR NOT NULL,
    "password" VARCHAR NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
