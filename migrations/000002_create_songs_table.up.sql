-- migrations/000002_create_songs_table.up.sql

CREATE TABLE IF NOT EXISTS songs (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    artist VARCHAR(255) NOT NULL,
    album VARCHAR(255),
    genre VARCHAR(255),
    year INTEGER,
    track_number INTEGER,
    duration_seconds INTEGER,
    file_path VARCHAR(512) NOT NULL UNIQUE,
    filename VARCHAR(255) NOT NULL,
    album_art_path VARCHAR(512),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Opcional: Índices para mejorar el rendimiento de las búsquedas
CREATE INDEX IF NOT EXISTS idx_songs_title ON songs (title);
CREATE INDEX IF NOT EXISTS idx_songs_artist ON songs (artist);
CREATE INDEX IF NOT EXISTS idx_songs_album ON songs (album);