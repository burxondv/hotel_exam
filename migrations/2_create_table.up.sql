CREATE TABLE "bookings"(
    "id" SERIAL PRIMARY KEY,
    "hotel_id" INTEGER NOT NULL,
    "room_id" INTEGER NOT NULL,
    "guest_id" INTEGER NOT NULL,
    "from" DATE NOT NULL,
    "to" DATE NOT NULL,
    "total_price" INTEGER NOT NULL
);