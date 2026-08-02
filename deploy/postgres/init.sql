-- PostgreSQL initialization script for VEDO EduTrack.
--
-- Applied once when the PostgreSQL container starts for the first time.
-- Production-grade extensions and defaults are enabled here.
--
-- See ADR-DES.DATA.storage-strategy (PostgreSQL + sqlc/pgx).

-- UUID generation (used for aggregate IDs, idempotency keys).
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Full-text search for resource catalog (F3).
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- Schema creation placeholder — actual schemas are managed by Atlas migrations.
-- See ADR-IMPL.PROCESS.development-tooling §3 (Atlas).
CREATE SCHEMA IF NOT EXISTS edutrack;
COMMENT ON SCHEMA edutrack IS 'VEDO EduTrack application schema — managed by Atlas migrations.';
