UPDATE tajwid_rules
SET
    examples = $$[
      {"surah_name":"Al-Baqarah","ayah_number":136,"full_text":"مِنْ أَحَدٍ","highlighted_text":"مِنْ أَ","translation":"dari seorangpun"},
      {"surah_name":"Al-Insan","ayah_number":4,"full_text":"إِنَّا أَعْتَدْنَا","highlighted_text":"أَعْتَدْنَا","translation":"sesungguhnya Kami menyediakan"}
    ]$$::jsonb,
    updated_at = NOW()
WHERE id = 'izhar-halqi';

UPDATE tajwid_rules
SET
    examples = $$[
      {"surah_name":"Al-Zalzalah","ayah_number":7,"full_text":"مَن يَعْمَلْ","highlighted_text":"مَن يَ","translation":"barangsiapa mengerjakan"}
    ]$$::jsonb,
    updated_at = NOW()
WHERE id = 'idgham-bighunnah';

UPDATE tajwid_rules
SET
    examples = $$[
      {"surah_name":"Al-Baqarah","ayah_number":2,"full_text":"مِن رَّبِّهِمْ","highlighted_text":"مِن رَّ","translation":"dari Tuhan mereka"}
    ]$$::jsonb,
    updated_at = NOW()
WHERE id = 'idgham-bilaghunnah';

UPDATE tajwid_rules
SET
    examples = $$[
      {"surah_name":"Al-Humazah","ayah_number":4,"full_text":"كَلَّا لَيُنبَذَنَّ","highlighted_text":"لَيُنبَ","translation":"sekali-kali tidak, pasti dia akan dilemparkan"}
    ]$$::jsonb,
    updated_at = NOW()
WHERE id = 'iqlab';

UPDATE tajwid_rules
SET
    examples = $$[
      {"surah_name":"Al-Baqarah","ayah_number":2,"full_text":"وَمِنْ كُلِّ","highlighted_text":"مِنْ كُ","translation":"dan dari setiap"}
    ]$$::jsonb,
    updated_at = NOW()
WHERE id = 'ikhfa-haqiqi';

UPDATE tajwid_rules
SET
    examples = $$[
      {"surah_name":"Al-Fatihah","ayah_number":1,"full_text":"الرَّحِيمِ","highlighted_text":"حِيمِ","translation":"Maha Penyayang"}
    ]$$::jsonb,
    updated_at = NOW()
WHERE id = 'mad-thobi-i';

UPDATE tajwid_rules
SET
    examples = $$[
      {"surah_name":"Al-Baqarah","ayah_number":5,"full_text":"إِنَّا أَعْطَيْنَاكَ","highlighted_text":"إِنَّا أَ","translation":"Sesungguhnya Kami telah memberimu"}
    ]$$::jsonb,
    updated_at = NOW()
WHERE id = 'mad-jaiz-munfasil';

UPDATE tajwid_rules
SET
    examples = $$[
      {"surah_name":"Al-Baqarah","ayah_number":14,"full_text":"جَاءُوا","highlighted_text":"جَاءُ","translation":"mereka datang"}
    ]$$::jsonb,
    updated_at = NOW()
WHERE id = 'mad-wajib-muttasil';

UPDATE tajwid_rules
SET
    examples = $$[
      {"surah_name":"Al-An''am","ayah_number":143,"full_text":"آلذَّكَرَيْنِ","highlighted_text":"آل","translation":"Apakah dua yang jantan"}
    ]$$::jsonb,
    updated_at = NOW()
WHERE id = 'mad-lazim-kilmi';

UPDATE tajwid_rules
SET
    examples = $$[
      {"surah_name":"Al-Falaq","ayah_number":2,"full_text":"مِن شَرِّ مَا خَلَقَ","highlighted_text":"خَلَقَ","translation":"dari kejahatan makhluk-Nya"}
    ]$$::jsonb,
    updated_at = NOW()
WHERE id = 'qalqalah-sughra';

UPDATE tajwid_rules
SET
    examples = $$[
      {"surah_name":"Al-Ikhlas","ayah_number":1,"full_text":"قُلْ هُوَ اللَّهُ أَحَدٌ","highlighted_text":"أَحَدٌ","translation":"Katakanlah: Dialah Allah Yang Maha Esa"}
    ]$$::jsonb,
    updated_at = NOW()
WHERE id = 'qalqalah-kubra';

UPDATE tajwid_rules
SET
    examples = $$[
      {"surah_name":"Al-Fatihah","ayah_number":7,"full_text":"وَلَا الضَّالِّينَ","highlighted_text":"الضَّالِّ","translation":"dan bukan orang-orang yang sesat"}
    ]$$::jsonb,
    updated_at = NOW()
WHERE id = 'ghunnah';

UPDATE tajwid_rules
SET
    examples = $$[
      {"surah_name":"Al-Fatihah","ayah_number":2,"full_text":"الرَّحْمَٰنِ","highlighted_text":"الرَّ","translation":"Yang Maha Pengasih"}
    ]$$::jsonb,
    updated_at = NOW()
WHERE id = 'alif-lam-syamsiyah';

UPDATE tajwid_rules
SET
    examples = $$[
      {"surah_name":"Al-Fatihah","ayah_number":2,"full_text":"الْحَمْدُ","highlighted_text":"الْحَ","translation":"Segala puji"}
    ]$$::jsonb,
    updated_at = NOW()
WHERE id = 'alif-lam-qamariyah';

UPDATE tajwid_rules
SET
    examples = $$[
      {"surah_name":"Al-Fatihah","ayah_number":5,"full_text":"إِيَّاكَ نَعْبُدُ وَإِيَّاكَ نَسْتَعِينُ","highlighted_text":"نَسْتَعِينُ","translation":"hanya kepada-Mu kami memohon pertolongan"}
    ]$$::jsonb,
    updated_at = NOW()
WHERE id = 'waqf-tam';
