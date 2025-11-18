CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE "asset_status" AS ENUM (
    'offline',
    'online'
);
CREATE TABLE IF NOT EXISTS "assets" (
    "ID"      UUID PRIMARY KEY DEFAULT UUID_GENERATE_V4(),
    "name"          VARCHAR(255) NOT NULL UNIQUE,
    "status"        "asset_status" NOT NULL DEFAULT 'offline',
    "locationID"   UUID NOT NULL REFERENCES "locations"("ID") ON DELETE SET NULL,
    "createdAtUTC"    TIMESTAMP(3) NOT NULL DEFAULT NOW(),
    "lastUpdatedAtUTC"  TIMESTAMP(3) NOT NULL DEFAULT NOW()
);