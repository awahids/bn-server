UPDATE tajwid_rules
SET
    examples = $$[
      {"surah_name":"Al-Baqarah","ayah_number":109,"full_text":"وَدَّ كَثِيرٌ مِنْ أَهْلِ الْكِتَابِ لَوْ يَرُدُّونَكُم","highlighted_text":"مِنْ أَهْ","translation":"banyak dari Ahli Kitab ingin mengembalikan kalian"},
      {"surah_name":"Al-Insan","ayah_number":4,"full_text":"إِنَّا أَعْتَدْنَا لِلْكَافِرِينَ سَلَاسِلَا وَأَغْلَالًا وَسَعِيرًا","highlighted_text":"أَعْتَدْنَا","translation":"sesungguhnya Kami menyediakan bagi orang-orang kafir rantai dan belenggu serta neraka yang menyala"}
    ]$$::jsonb,
    updated_at = NOW()
WHERE id = 'izhar-halqi';

UPDATE tajwid_rules
SET
    examples = $$[
      {"surah_name":"Az-Zalzalah","ayah_number":7,"full_text":"فَمَن يَعْمَلْ مِثْقَالَ ذَرَّةٍ خَيْرًا يَرَهُ","highlighted_text":"مَن يَ","translation":"maka barangsiapa mengerjakan kebaikan seberat zarrah, niscaya dia akan melihatnya"}
    ]$$::jsonb,
    updated_at = NOW()
WHERE id = 'idgham-bighunnah';

UPDATE tajwid_rules
SET
    examples = $$[
      {"surah_name":"Al-Baqarah","ayah_number":5,"full_text":"أُولَٰئِكَ عَلَىٰ هُدًى مِّن رَّبِّهِمْ وَأُولَٰئِكَ هُمُ الْمُفْلِحُونَ","highlighted_text":"مِّن رَّ","translation":"mereka itulah yang mendapat petunjuk dari Tuhan mereka, dan mereka itulah orang-orang yang beruntung"}
    ]$$::jsonb,
    updated_at = NOW()
WHERE id = 'idgham-bilaghunnah';

UPDATE tajwid_rules
SET
    examples = $$[
      {"surah_name":"Al-Humazah","ayah_number":4,"full_text":"كَلَّا لَيُنبَذَنَّ فِي الْحُطَمَةِ","highlighted_text":"لَيُنبَ","translation":"sekali-kali tidak, pasti dia benar-benar akan dilemparkan ke dalam Huthamah"}
    ]$$::jsonb,
    updated_at = NOW()
WHERE id = 'iqlab';

UPDATE tajwid_rules
SET
    examples = $$[
      {"surah_name":"Qaf","ayah_number":7,"full_text":"وَأَنْبَتْنَا فِيهَا مِنْ كُلِّ زَوْجٍ بَهِيجٍ","highlighted_text":"مِنْ كُ","translation":"dan Kami tumbuhkan di sana segala macam pasangan yang indah"}
    ]$$::jsonb,
    updated_at = NOW()
WHERE id = 'ikhfa-haqiqi';

UPDATE tajwid_rules
SET
    examples = $$[
      {"surah_name":"Al-Fatihah","ayah_number":3,"full_text":"الرَّحْمَٰنِ الرَّحِيمِ","highlighted_text":"حِيمِ","translation":"Yang Maha Pengasih, Maha Penyayang"}
    ]$$::jsonb,
    updated_at = NOW()
WHERE id = 'mad-thobi-i';

UPDATE tajwid_rules
SET
    examples = $$[
      {"surah_name":"Al-Kafirun","ayah_number":1,"full_text":"قُلْ يَا أَيُّهَا الْكَافِرُونَ","highlighted_text":"يَا أَ","translation":"Katakanlah, wahai orang-orang kafir"}
    ]$$::jsonb,
    updated_at = NOW()
WHERE id = 'mad-jaiz-munfasil';

UPDATE tajwid_rules
SET
    examples = $$[
      {"surah_name":"An-Nasr","ayah_number":1,"full_text":"إِذَا جَاءَ نَصْرُ اللَّهِ وَالْفَتْحُ","highlighted_text":"جَاءَ","translation":"apabila telah datang pertolongan Allah dan kemenangan"}
    ]$$::jsonb,
    updated_at = NOW()
WHERE id = 'mad-wajib-muttasil';

UPDATE tajwid_rules
SET
    examples = $$[
      {"surah_name":"Al-An''am","ayah_number":143,"full_text":"قُلْ آلذَّكَرَيْنِ حَرَّمَ أَمِ الْأُنثَيَيْنِ","highlighted_text":"آل","translation":"Katakanlah, apakah dua yang jantan yang diharamkan ataukah dua yang betina"}
    ]$$::jsonb,
    updated_at = NOW()
WHERE id = 'mad-lazim-kilmi';

UPDATE tajwid_rules
SET
    examples = $$[
      {"surah_name":"Al-Falaq","ayah_number":2,"full_text":"قُلْ أَعُوذُ بِرَبِّ الْفَلَقِ مِن شَرِّ مَا خَلَقَ","highlighted_text":"خَلَقَ","translation":"Aku berlindung kepada Tuhan yang menguasai subuh dari kejahatan makhluk-Nya"}
    ]$$::jsonb,
    updated_at = NOW()
WHERE id = 'qalqalah-sughra';

UPDATE tajwid_rules
SET
    examples = $$[
      {"surah_name":"Al-Masad","ayah_number":1,"full_text":"تَبَّتْ يَدَا أَبِي لَهَبٍ وَتَبَّ","highlighted_text":"وَتَبَّ","translation":"binasalah kedua tangan Abu Lahab, dan benar-benar binasa dia"}
    ]$$::jsonb,
    updated_at = NOW()
WHERE id = 'qalqalah-kubra';

UPDATE tajwid_rules
SET
    examples = $$[
      {"surah_name":"Al-Kawthar","ayah_number":1,"full_text":"إِنَّا أَعْطَيْنَاكَ الْكَوْثَرَ","highlighted_text":"إِنَّا","translation":"sesungguhnya Kami telah memberikan kepadamu nikmat yang banyak"}
    ]$$::jsonb,
    updated_at = NOW()
WHERE id = 'ghunnah';

UPDATE tajwid_rules
SET
    examples = $$[
      {"surah_name":"Al-Fatihah","ayah_number":3,"full_text":"بِسْمِ اللَّهِ الرَّحْمَٰنِ الرَّحِيمِ","highlighted_text":"الرَّ","translation":"dengan nama Allah Yang Maha Pengasih, Maha Penyayang"}
    ]$$::jsonb,
    updated_at = NOW()
WHERE id = 'alif-lam-syamsiyah';

UPDATE tajwid_rules
SET
    examples = $$[
      {"surah_name":"Al-Fatihah","ayah_number":2,"full_text":"الْحَمْدُ لِلَّهِ رَبِّ الْعَالَمِينَ","highlighted_text":"الْحَ","translation":"segala puji bagi Allah, Tuhan seluruh alam"}
    ]$$::jsonb,
    updated_at = NOW()
WHERE id = 'alif-lam-qamariyah';

UPDATE tajwid_rules
SET
    examples = $$[
      {"surah_name":"Al-Fatihah","ayah_number":5,"full_text":"إِيَّاكَ نَعْبُدُ وَإِيَّاكَ نَسْتَعِينُ","highlighted_text":"نَسْتَعِينُ","translation":"hanya kepada-Mu kami menyembah dan hanya kepada-Mu kami memohon pertolongan"}
    ]$$::jsonb,
    updated_at = NOW()
WHERE id = 'waqf-tam';
