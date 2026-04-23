ALTER TABLE schools DROP CONSTRAINT IF EXISTS schools_status_sekolah_check;
ALTER TABLE schools DROP CONSTRAINT IF EXISTS schools_jenjang_check;

ALTER TABLE schools
  DROP COLUMN IF EXISTS contact,
  DROP COLUMN IF EXISTS status_sekolah,
  DROP COLUMN IF EXISTS jenjang;
