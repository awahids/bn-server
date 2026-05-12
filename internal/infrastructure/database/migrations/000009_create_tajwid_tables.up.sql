CREATE TABLE IF NOT EXISTS tajwid_rules (
    id              VARCHAR(50)  PRIMARY KEY,
    name            VARCHAR(100) NOT NULL,
    arabic_name     VARCHAR(200),
    category        VARCHAR(50)  NOT NULL,
    description     TEXT         NOT NULL,
    trigger_letters VARCHAR(200),
    examples        JSONB        NOT NULL DEFAULT '[]',
    audio_url       VARCHAR(500),
    sort_order      INT          NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_tajwid_rules_category   ON tajwid_rules (category);
CREATE INDEX IF NOT EXISTS idx_tajwid_rules_sort_order ON tajwid_rules (sort_order);

DELETE FROM tajwid_rules;

INSERT INTO tajwid_rules (id, name, arabic_name, category, description, trigger_letters, examples, audio_url, sort_order) VALUES
-- NUN SUKUN / TANWIN (5 rules)
('izhar-halqi',       'Izhar Halqi',         'إظهار حلقي',    'nun_sukun', 'Nun sukun atau tanwin dibaca jelas tanpa dengung ketika bertemu salah satu huruf halqi (tenggorokan)', 'ء ه ع ح غ خ',
 '[{"surah_name":"Al-Baqarah","ayah_number":136,"full_text":"مِنْ أَحَدٍ","highlighted_text":"مِنْ أَ","translation":"dari seorangpun"},{"surah_name":"Al-Insan","ayah_number":4,"full_text":"إِنَّا أَعْتَدْنَا","highlighted_text":"أَعْتَدْنَا","translation":"sesungguhnya Kami menyediakan"}]',
 '/audio/tajwid/izhar-halqi.mp3', 1),

('idgham-bighunnah',  'Idgham Bighunnah',    'إدغام بغنة',    'nun_sukun', 'Nun sukun atau tanwin dimasukkan ke huruf berikutnya disertai dengung (ghunnah) ketika bertemu ي ن م و', 'ي ن م و',
 '[{"surah_name":"Al-Zalzalah","ayah_number":7,"full_text":"مَن يَعْمَلْ","highlighted_text":"مَن يَ","translation":"barangsiapa mengerjakan"}]',
 '/audio/tajwid/idgham-bighunnah.mp3', 2),

('idgham-bilaghunnah','Idgham Bilaghunnah',  'إدغام بلا غنة', 'nun_sukun', 'Nun sukun atau tanwin dimasukkan ke huruf berikutnya TANPA dengung ketika bertemu ر atau ل', 'ر ل',
 '[{"surah_name":"Al-Baqarah","ayah_number":2,"full_text":"مِن رَّبِّهِمْ","highlighted_text":"مِن رَّ","translation":"dari Tuhan mereka"}]',
 '/audio/tajwid/idgham-bilaghunnah.mp3', 3),

('iqlab',             'Iqlab',               'إقلاب',         'nun_sukun', 'Nun sukun atau tanwin berubah menjadi mim ketika bertemu huruf ب, disertai dengung', 'ب',
 '[{"surah_name":"Al-Humazah","ayah_number":4,"full_text":"كَلَّا لَيُنبَذَنَّ","highlighted_text":"لَيُنبَ","translation":"sekali-kali tidak, pasti dia akan dilemparkan"}]',
 '/audio/tajwid/iqlab.mp3', 4),

('ikhfa-haqiqi',      'Ikhfa Haqiqi',        'إخفاء حقيقي',   'nun_sukun', 'Nun sukun atau tanwin dibaca samar-samar dengan dengung ketika bertemu 15 huruf ikhfa', 'ت ث ج د ذ ز س ش ص ض ط ظ ف ق ك',
 '[{"surah_name":"Al-Baqarah","ayah_number":2,"full_text":"وَمِنْ كُلِّ","highlighted_text":"مِنْ كُ","translation":"dan dari setiap"}]',
 '/audio/tajwid/ikhfa-haqiqi.mp3', 5),

-- MAD (4 rules)
('mad-thobi-i',       'Mad Thobi''i',        'مد طبيعي',      'mad', 'Mad asli/alami, dibaca 2 harakat. Terjadi pada alif sesudah fathah, wau sesudah dhammah, ya sesudah kasrah', 'ا و ي',
 '[{"surah_name":"Al-Fatihah","ayah_number":1,"full_text":"الرَّحِيمِ","highlighted_text":"حِيمِ","translation":"Maha Penyayang"}]',
 '/audio/tajwid/mad-thobi-i.mp3', 6),

('mad-jaiz-munfasil', 'Mad Jaiz Munfasil',   'مد جائز منفصل', 'mad', 'Mad yang terjadi bila huruf mad berada di akhir kata dan huruf hamzah berada di awal kata berikutnya. Dibaca 4-5 harakat', 'ا و ي',
 '[{"surah_name":"Al-Baqarah","ayah_number":5,"full_text":"إِنَّا أَعْطَيْنَاكَ","highlighted_text":"إِنَّا أَ","translation":"Sesungguhnya Kami telah memberimu"}]',
 '/audio/tajwid/mad-jaiz-munfasil.mp3', 7),

('mad-wajib-muttasil','Mad Wajib Muttasil',  'مد واجب متصل',  'mad', 'Mad yang wajib dibaca 4-5 harakat karena huruf mad dan hamzah berada dalam satu kata', 'ا و ي',
 '[{"surah_name":"Al-Baqarah","ayah_number":14,"full_text":"جَاءُوا","highlighted_text":"جَاءُ","translation":"mereka datang"}]',
 '/audio/tajwid/mad-wajib-muttasil.mp3', 8),

('mad-lazim-kilmi',   'Mad Lazim Kilmi',     'مد لازم كلمي',  'mad', 'Mad yang harus dibaca 6 harakat karena huruf mad bertemu huruf sukun dalam satu kata', 'ا و ي',
 '[{"surah_name":"Al-An''am","ayah_number":143,"full_text":"آلذَّكَرَيْنِ","highlighted_text":"آل","translation":"Apakah dua yang jantan"}]',
 '/audio/tajwid/mad-lazim-kilmi.mp3', 9),

-- QALQALAH (2 rules)
('qalqalah-sughra',   'Qalqalah Sughra',     'قلقلة صغرى',    'qalqalah', 'Qalqalah kecil: huruf qalqalah (ق ط ب ج د) sukun di tengah kata, dibaca dengan pantulan ringan', 'ق ط ب ج د',
 '[{"surah_name":"Al-Falaq","ayah_number":2,"full_text":"مِن شَرِّ مَا خَلَقَ","highlighted_text":"خَلَقَ","translation":"dari kejahatan makhluk-Nya"}]',
 '/audio/tajwid/qalqalah-sughra.mp3', 10),

('qalqalah-kubra',    'Qalqalah Kubra',      'قلقلة كبرى',    'qalqalah', 'Qalqalah besar: huruf qalqalah sukun karena waqf (berhenti), dibaca dengan pantulan kuat', 'ق ط ب ج د',
 '[{"surah_name":"Al-Ikhlas","ayah_number":1,"full_text":"قُلْ هُوَ اللَّهُ أَحَدٌ","highlighted_text":"أَحَدٌ","translation":"Katakanlah: Dialah Allah Yang Maha Esa"}]',
 '/audio/tajwid/qalqalah-kubra.mp3', 11),

-- GHUNNAH
('ghunnah',           'Ghunnah',             'غنة',           'ghunnah', 'Dengung yang wajib dibaca pada nun atau mim yang bertasydid, selama 2 harakat', 'ن م',
 '[{"surah_name":"Al-Fatihah","ayah_number":7,"full_text":"وَلَا الضَّالِّينَ","highlighted_text":"الضَّالِّ","translation":"dan bukan orang-orang yang sesat"}]',
 '/audio/tajwid/ghunnah.mp3', 12),

-- ALIF LAM (2 rules)
('alif-lam-syamsiyah','Alif Lam Syamsiyah',  'ال شمسية',      'al', 'Alif lam dibaca lebur/diidghamkan ke huruf berikutnya (huruf syams). Tidak ada bunyi "L"', 'ت ث د ذ ر ز س ش ص ض ط ظ ل ن',
 '[{"surah_name":"Al-Fatihah","ayah_number":2,"full_text":"الرَّحْمَٰنِ","highlighted_text":"الرَّ","translation":"Yang Maha Pengasih"}]',
 '/audio/tajwid/alif-lam-syamsiyah.mp3', 13),

('alif-lam-qamariyah','Alif Lam Qamariyah',  'ال قمرية',      'al', 'Alif lam dibaca jelas/izhar ke huruf berikutnya (huruf qamar). Bunyi "L" terdengar jelas', 'ا ب ج ح خ ع غ ف ق ك م و ه ي',
 '[{"surah_name":"Al-Fatihah","ayah_number":2,"full_text":"الْحَمْدُ","highlighted_text":"الْحَ","translation":"Segala puji"}]',
 '/audio/tajwid/alif-lam-qamariyah.mp3', 14),

-- WAQF
('waqf-tam',          'Waqf Tam',            'وقف تام',       'waqf', 'Berhenti sempurna pada akhir kalimat yang maknanya sudah lengkap. Tanda: قلى atau ط', 'وقف',
 '[{"surah_name":"Al-Fatihah","ayah_number":5,"full_text":"إِيَّاكَ نَعْبُدُ وَإِيَّاكَ نَسْتَعِينُ","highlighted_text":"نَسْتَعِينُ","translation":"hanya kepada-Mu kami memohon pertolongan"}]',
 '/audio/tajwid/waqf-tam.mp3', 15);
