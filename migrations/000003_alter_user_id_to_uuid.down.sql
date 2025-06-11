-- Revierte los cambios eliminando las tablas.
DROP TABLE IF EXISTS authentication;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS users;
-- DROP TABLE IF EXISTS songs; -- Solo si la creaste en la migración UP y quieres eliminarla aquí.