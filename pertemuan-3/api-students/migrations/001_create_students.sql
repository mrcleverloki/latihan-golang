CREATE TABLE IF NOT EXISTS students (
    id SERIAL PRIMARY KEY,
    nim VARCHAR(20) NOT NULL,
    name VARCHAR(100) NOT NULL,
    grade VARCHAR(5) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Batasan NIM unik tanpa membedakan huruf besar/kecil
CREATE UNIQUE INDEX IF NOT EXISTS students_nim_lower_key 
ON students (LOWER(nim));

-- Indeks tambahan untuk mempercepat pencarian nama (ILIKE)
CREATE INDEX IF NOT EXISTS students_name_lower_idx 
ON students (LOWER(name));