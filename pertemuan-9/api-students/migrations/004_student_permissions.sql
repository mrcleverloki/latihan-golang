-- ---------------------------------------------------------------
-- 1. Tambah permission untuk entity students
-- ---------------------------------------------------------------
INSERT INTO permissions (name, description) VALUES
    ('student:list',       'Melihat daftar seluruh student'),
    ('student:read:any',   'Melihat data student milik siapa pun'),
    ('student:create',     'Mendaftarkan student baru'),
    ('student:update:any', 'Mengubah data student mana pun'),
    ('student:delete',     'Menghapus data student')
ON CONFLICT (name) DO NOTHING;

-- ---------------------------------------------------------------
-- 2. Pasangkan ke role masing-masing
-- ---------------------------------------------------------------
INSERT INTO role_permissions (role_name, permission_name) VALUES
    ('admin', 'student:list'),
    ('admin', 'student:read:any'),
    ('admin', 'student:create'),
    ('admin', 'student:update:any'),
    ('admin', 'student:delete'),
    ('staff', 'student:list'),
    ('staff', 'student:read:any'),
    ('staff', 'student:create')
ON CONFLICT DO NOTHING;

-- ---------------------------------------------------------------
-- 3. Tambah kolom owner_id pada tabel students
-- ---------------------------------------------------------------
ALTER TABLE students ADD COLUMN IF NOT EXISTS owner_id INTEGER;

-- Tangani data lama: jika ada baris yang owner_id-nya NULL, pasang ke user pertama
UPDATE students 
SET owner_id = (SELECT id FROM users ORDER BY id ASC LIMIT 1)
WHERE owner_id IS NULL AND EXISTS (SELECT 1 FROM users);

-- Pasang foreign key constraint
ALTER TABLE students DROP CONSTRAINT IF EXISTS students_owner_id_fkey;
ALTER TABLE students
    ADD CONSTRAINT students_owner_id_fkey
    FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS students_owner_id_idx ON students (owner_id);
