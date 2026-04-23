ALTER TABLE schools
  ADD COLUMN IF NOT EXISTS jenjang VARCHAR(32) NOT NULL DEFAULT 'Lainnya',
  ADD COLUMN IF NOT EXISTS status_sekolah VARCHAR(20) NOT NULL DEFAULT 'swasta',
  ADD COLUMN IF NOT EXISTS contact VARCHAR(100) NOT NULL DEFAULT '';

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM pg_constraint
    WHERE conname = 'schools_jenjang_check'
  ) THEN
    ALTER TABLE schools
      ADD CONSTRAINT schools_jenjang_check
      CHECK (jenjang IN ('TK', 'SD', 'SMP', 'SMA', 'SMK', 'MI', 'MTs', 'MA', 'Lainnya'));
  END IF;
END $$;

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM pg_constraint
    WHERE conname = 'schools_status_sekolah_check'
  ) THEN
    ALTER TABLE schools
      ADD CONSTRAINT schools_status_sekolah_check
      CHECK (status_sekolah IN ('negeri', 'swasta'));
  END IF;
END $$;
