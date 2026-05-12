UPDATE tajwid_rules
SET
    examples = $$[
      {"surah_name":"Al-Baqarah","ayah_number":109,"full_text":"وَدَّ كَثِيرٌ مِّنْ أَهْلِ ٱلْكِتَـٰبِ لَوْ يَرُدُّونَكُم مِّنۢ بَعْدِ إِيمَـٰنِكُمْ كُفَّارًا حَسَدًا مِّنْ عِندِ أَنفُسِهِم مِّنۢ بَعْدِ مَا تَبَيَّنَ لَهُمُ ٱلْحَقُّ ۖ فَٱعْفُوا۟ وَٱصْفَحُوا۟ حَتَّىٰ يَأْتِىَ ٱللَّهُ بِأَمْرِهِۦٓ ۗ إِنَّ ٱللَّهَ عَلَىٰ كُلِّ شَىْءٍ قَدِيرٌ","highlighted_text":"مِّنْ أَهْلِ","translation":"Banyak dari Ahli Kitab menginginkan agar kalian kembali kafir setelah beriman."},
      {"surah_name":"Al-Insan","ayah_number":4,"full_text":"إِنَّآ أَعْتَدْنَا لِلْكَـٰفِرِينَ سَلَـٰسِلَا۟ وَأَغْلَـٰلًا وَسَعِيرًا","highlighted_text":"أَعْتَدْنَا","translation":"Sesungguhnya Kami menyediakan bagi orang-orang kafir rantai, belenggu, dan neraka yang menyala."}
    ]$$::jsonb,
    updated_at = NOW()
WHERE id = 'izhar-halqi';

UPDATE tajwid_rules
SET
    examples = $$[
      {"surah_name":"Az-Zalzalah","ayah_number":7,"full_text":"فَمَن يَعْمَلْ مِثْقَالَ ذَرَّةٍ خَيْرًا يَرَهُۥ","highlighted_text":"مَن يَعْ","translation":"Siapa yang mengerjakan kebaikan seberat zarrah, niscaya dia akan melihatnya."}
    ]$$::jsonb,
    updated_at = NOW()
WHERE id = 'idgham-bighunnah';

UPDATE tajwid_rules
SET
    examples = $$[
      {"surah_name":"Al-Baqarah","ayah_number":5,"full_text":"أُو۟لَـٰٓئِكَ عَلَىٰ هُدًى مِّن رَّبِّهِمْ ۖ وَأُو۟لَـٰٓئِكَ هُمُ ٱلْمُفْلِحُونَ","highlighted_text":"مِّن رَّبِّ","translation":"Mereka berada di atas petunjuk dari Tuhan mereka, dan mereka itulah orang-orang yang beruntung."}
    ]$$::jsonb,
    updated_at = NOW()
WHERE id = 'idgham-bilaghunnah';

UPDATE tajwid_rules
SET
    examples = $$[
      {"surah_name":"Al-Humazah","ayah_number":4,"full_text":"كَلَّا ۖ لَيُنۢبَذَنَّ فِى ٱلْحُطَمَةِ","highlighted_text":"لَيُنۢبَ","translation":"Sekali-kali tidak, dia benar-benar akan dilemparkan ke dalam Huthamah."}
    ]$$::jsonb,
    updated_at = NOW()
WHERE id = 'iqlab';

UPDATE tajwid_rules
SET
    examples = $$[
      {"surah_name":"Qaf","ayah_number":7,"full_text":"وَٱلْأَرْضَ مَدَدْنَـٰهَا وَأَلْقَيْنَا فِيهَا رَوَٰسِىَ وَأَنۢبَتْنَا فِيهَا مِن كُلِّ زَوْجٍۭ بَهِيجٍ","highlighted_text":"مِن كُلِّ","translation":"Dan bumi Kami hamparkan, lalu Kami tumbuhkan di dalamnya segala macam pasangan yang indah."}
    ]$$::jsonb,
    updated_at = NOW()
WHERE id = 'ikhfa-haqiqi';

UPDATE tajwid_rules
SET
    examples = $$[
      {"surah_name":"Al-Fatihah","ayah_number":3,"full_text":"ٱلرَّحْمَـٰنِ ٱلرَّحِيمِ","highlighted_text":"ٱلرَّحِيمِ","translation":"Yang Maha Pengasih, Maha Penyayang."}
    ]$$::jsonb,
    updated_at = NOW()
WHERE id = 'mad-thobi-i';

UPDATE tajwid_rules
SET
    examples = $$[
      {"surah_name":"Al-Kafirun","ayah_number":1,"full_text":"قُلْ يَـٰٓأَيُّهَا ٱلْكَـٰفِرُونَ","highlighted_text":"يَـٰٓأَيُّهَا","translation":"Katakanlah, wahai orang-orang kafir."}
    ]$$::jsonb,
    updated_at = NOW()
WHERE id = 'mad-jaiz-munfasil';

UPDATE tajwid_rules
SET
    examples = $$[
      {"surah_name":"An-Nasr","ayah_number":1,"full_text":"إِذَا جَآءَ نَصْرُ ٱللَّهِ وَٱلْفَتْحُ","highlighted_text":"جَآءَ","translation":"Apabila telah datang pertolongan Allah dan kemenangan."}
    ]$$::jsonb,
    updated_at = NOW()
WHERE id = 'mad-wajib-muttasil';

UPDATE tajwid_rules
SET
    examples = $$[
      {"surah_name":"Al-An'am","ayah_number":143,"full_text":"ثَمَـٰنِيَةَ أَزْوَٰجٍ ۖ مِّنَ ٱلضَّأْنِ ٱثْنَيْنِ وَمِنَ ٱلْمَعْزِ ٱثْنَيْنِ ۗ قُلْ ءَآلذَّكَرَيْنِ حَرَّمَ أَمِ ٱلْأُنثَيَيْنِ أَمَّا ٱشْتَمَلَتْ عَلَيْهِ أَرْحَامُ ٱلْأُنثَيَيْنِ ۖ نَبِّـُٔونِى بِعِلْمٍ إِن كُنتُمْ صَـٰدِقِينَ","highlighted_text":"ءَآلذَّ","translation":"Katakanlah, apakah dua yang jantan yang diharamkan, atau dua yang betina?"}
    ]$$::jsonb,
    updated_at = NOW()
WHERE id = 'mad-lazim-kilmi';

UPDATE tajwid_rules
SET
    examples = $$[
      {"surah_name":"Al-Falaq","ayah_number":2,"full_text":"مِن شَرِّ مَا خَلَقَ","highlighted_text":"خَلَقَ","translation":"Dari kejahatan makhluk-Nya."}
    ]$$::jsonb,
    updated_at = NOW()
WHERE id = 'qalqalah-sughra';

UPDATE tajwid_rules
SET
    examples = $$[
      {"surah_name":"Al-Masad","ayah_number":1,"full_text":"تَبَّتْ يَدَآ أَبِى لَهَبٍ وَتَبَّ","highlighted_text":"وَتَبَّ","translation":"Binasalah kedua tangan Abu Lahab, dan benar-benar binasa dia."}
    ]$$::jsonb,
    updated_at = NOW()
WHERE id = 'qalqalah-kubra';

UPDATE tajwid_rules
SET
    examples = $$[
      {"surah_name":"Al-Kawthar","ayah_number":1,"full_text":"إِنَّآ أَعْطَيْنَـٰكَ ٱلْكَوْثَرَ","highlighted_text":"إِنَّآ","translation":"Sesungguhnya Kami telah memberimu nikmat yang banyak."}
    ]$$::jsonb,
    updated_at = NOW()
WHERE id = 'ghunnah';

UPDATE tajwid_rules
SET
    examples = $$[
      {"surah_name":"Al-Fatihah","ayah_number":1,"full_text":"بِسْمِ ٱللَّهِ ٱلرَّحْمَـٰنِ ٱلرَّحِيمِ","highlighted_text":"ٱلرَّحْ","translation":"Dengan nama Allah Yang Maha Pengasih lagi Maha Penyayang."}
    ]$$::jsonb,
    updated_at = NOW()
WHERE id = 'alif-lam-syamsiyah';

UPDATE tajwid_rules
SET
    examples = $$[
      {"surah_name":"Al-Fatihah","ayah_number":2,"full_text":"ٱلْحَمْدُ لِلَّهِ رَبِّ ٱلْعَـٰلَمِينَ","highlighted_text":"ٱلْحَمْ","translation":"Segala puji bagi Allah, Tuhan seluruh alam."}
    ]$$::jsonb,
    updated_at = NOW()
WHERE id = 'alif-lam-qamariyah';

UPDATE tajwid_rules
SET
    examples = $$[
      {"surah_name":"Al-Fatihah","ayah_number":5,"full_text":"إِيَّاكَ نَعْبُدُ وَإِيَّاكَ نَسْتَعِينُ","highlighted_text":"نَسْتَعِينُ","translation":"Hanya kepada-Mu kami menyembah dan hanya kepada-Mu kami memohon pertolongan."}
    ]$$::jsonb,
    updated_at = NOW()
WHERE id = 'waqf-tam';
