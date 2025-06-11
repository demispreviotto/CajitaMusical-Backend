-- Habilita la extensión pgcrypto para gen_random_uuid() si no está ya (requiere privilegios de superusuario)
-- Si estás en PostgreSQL 13+, gen_random_uuid() es nativo y no necesitas pgcrypto.
-- CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Paso 1: Eliminar las tablas existentes si existen (SOLO PARA DESARROLLO/TABLAS VACÍAS)
DROP TABLE IF EXISTS authentication;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS songs; -- Si también quieres que las canciones sean migradas, asegúrate de que GORM recree esta tabla también.
-- ¡CUIDADO! Las líneas de DROP TABLE eliminan tus datos existentes.

-- Paso 2: Recrear la tabla 'users' con 'id' como UUID
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Usa gen_random_uuid() para generar UUIDs
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Paso 3: Recrear la tabla 'authentication' con 'user_id' referenciando a UUID de 'users'
CREATE TABLE authentication (
    id SERIAL PRIMARY KEY,
    user_id UUID UNIQUE NOT NULL, -- Ahora UUID
    password_hash VARCHAR(255) NOT NULL,
    -- Add any other authentication fields if needed
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- Paso 4: Recrear la tabla 'sessions' con 'user_id' y 'session_id' como UUIDs
CREATE TABLE sessions (
    session_id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- UUID para SessionID
    user_id UUID NOT NULL, -- Ahora UUID
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    user_agent TEXT,
    ip_address VARCHAR(45),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);