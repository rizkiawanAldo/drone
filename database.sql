-- Create extension for UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create estate table
CREATE TABLE IF NOT EXISTS estates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    width INTEGER NOT NULL CHECK (width BETWEEN 1 AND 50000),
    length INTEGER NOT NULL CHECK (length BETWEEN 1 AND 50000),
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
);

-- Create tree table
CREATE TABLE IF NOT EXISTS trees (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    estate_id UUID NOT NULL REFERENCES estates(id) ON DELETE CASCADE,
    x INTEGER NOT NULL CHECK (x >= 1),
    y INTEGER NOT NULL CHECK (y >= 1),
    height INTEGER NOT NULL CHECK (height BETWEEN 1 AND 30),
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    -- Ensure tree coordinates are valid for the estate dimensions
    CONSTRAINT valid_coordinates CHECK (x <= (SELECT width FROM estates WHERE id = estate_id) AND y <= (SELECT length FROM estates WHERE id = estate_id)),
    -- Ensure only one tree per plot
    UNIQUE (estate_id, x, y)
); 