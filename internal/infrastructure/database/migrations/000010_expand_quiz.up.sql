ALTER TABLE quiz_questions
ADD COLUMN IF NOT EXISTS difficulty_order INT
GENERATED ALWAYS AS (
    CASE difficulty
        WHEN 'easy' THEN 1
        WHEN 'medium' THEN 2
        ELSE 3
    END
) STORED;

INSERT INTO quiz_categories (id, name, description, icon, color)
VALUES
    ('fiqih', 'Fiqih Islam', 'Hukum dan tata cara ibadah', 'Scale', 'chart-5'),
    ('sirah', 'Sirah Nabawiyah', 'Kisah Nabi dan sejarah Islam', 'Star', 'chart-6')
ON CONFLICT (id) DO UPDATE
SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    icon = EXCLUDED.icon,
    color = EXCLUDED.color,
    updated_at = NOW();

INSERT INTO quiz_questions (id, question, options, correct_answer, explanation, material, category_id, difficulty) VALUES
-- Hijaiyah medium
('hm1', 'Huruf hijaiyah manakah yang memiliki tiga titik di atas?', '["ب","ت","ث","ن"]', 2, 'Huruf ث memiliki tiga titik di atas.', 'Huruf tsa (ث) dibedakan dari ba dan ta melalui jumlah titiknya: tiga titik di atas.', 'hijaiyah', 'medium'),
('hm2', 'Huruf manakah yang dibaca "Kh"?', '["ح","خ","ج","غ"]', 1, 'Huruf خ dibaca kha atau kh.', 'Kha (خ) keluar dari tenggorokan dan memiliki satu titik di atas.', 'hijaiyah', 'medium'),
('hm3', 'Huruf manakah yang letak titiknya di bawah?', '["ف","ق","ب","ن"]', 2, 'Huruf ب memiliki satu titik di bawah.', 'Ba (ب) dibedakan dari ta dan tsa dengan satu titik di bawah.', 'hijaiyah', 'medium'),
('hm4', 'Huruf manakah yang dibaca "Syin"?', '["س","ش","ص","ض"]', 1, 'Huruf ش adalah syin.', 'Syin (ش) berbentuk mirip sin tetapi memiliki tiga titik di atas.', 'hijaiyah', 'medium'),
('hm5', 'Pasangan huruf yang bentuk dasarnya sama tetapi titiknya berbeda adalah?', '["د dan ر","ب dan ت","ع dan ل","ك dan م"]', 1, 'Ba dan ta berbagi bentuk dasar yang serupa.', 'Beberapa huruf hijaiyah dibedakan oleh jumlah dan posisi titik, contohnya ba, ta, dan tsa.', 'hijaiyah', 'medium'),
('hm6', 'Huruf manakah yang dibaca "Dzal"?', '["ذ","ز","د","ض"]', 0, 'Huruf ذ dibaca dzal.', 'Dzal (ذ) mirip dal tetapi memiliki satu titik di atas.', 'hijaiyah', 'medium'),
('hm7', 'Huruf manakah yang tidak memiliki titik?', '["ز","ش","ح","ن"]', 2, 'Huruf ح tidak memiliki titik.', 'Ha (ح) berbeda dari jim dan kha karena tidak memiliki titik.', 'hijaiyah', 'medium'),
('hm8', 'Huruf manakah yang dibaca "Qaf"?', '["ف","ق","ك","غ"]', 1, 'Huruf ق dibaca qaf.', 'Qaf (ق) memiliki dua titik di atas dan makhraj yang lebih tebal.', 'hijaiyah', 'medium'),
('hm9', 'Huruf apakah yang berada setelah ز dalam urutan hijaiyah?', '["س","ر","ش","ص"]', 0, 'Setelah ز adalah س.', 'Urutan hijaiyah membantu mengenali letak huruf dalam alfabet Arab.', 'hijaiyah', 'medium'),
('hm10', 'Huruf manakah yang dibaca "Ain"?', '["غ","ع","ء","ه"]', 1, 'Huruf ع dibaca ain.', 'Ain (ع) adalah huruf tenggorokan yang khas dalam bahasa Arab.', 'hijaiyah', 'medium'),

-- Hijaiyah hard
('hh1', 'Huruf manakah yang tepat untuk transliterasi "Zh"?', '["ض","ظ","ذ","ز"]', 1, 'Huruf ظ dibaca zha atau zho.', 'Zha (ظ) termasuk huruf isti''la yang dibaca tebal.', 'hijaiyah', 'hard'),
('hh2', 'Huruf mana yang merupakan pasangan bertitik dari ح?', '["خ","ج","ه","ع"]', 0, 'Huruf خ adalah versi bertitik dari bentuk dasar ح.', 'Ha (ح) dan kha (خ) memiliki bentuk dasar yang serupa, tetapi kha bertitik di atas.', 'hijaiyah', 'hard'),
('hh3', 'Huruf manakah yang termasuk huruf halki (tenggorokan)?', '["س","ع","ل","ف"]', 1, 'Huruf ع keluar dari tenggorokan.', 'Huruf halki meliputi ء ه ع ح غ خ dan keluar dari area tenggorokan.', 'hijaiyah', 'hard'),
('hh4', 'Huruf yang memiliki dua titik di bawah adalah?', '["ي","ن","ق","ث"]', 0, 'Huruf ي memiliki dua titik di bawah saat bentuk terpisah.', 'Ya (ي) dikenali dari dua titik di bawah pada bentuk terpisahnya.', 'hijaiyah', 'hard'),
('hh5', 'Huruf manakah yang dibaca paling dekat dengan suara "emphatic s"?', '["س","ص","ش","ز"]', 1, 'Huruf ص dibaca shad, yaitu s tebal.', 'Shad (ص) adalah huruf tebal yang berbeda dari sin (س).', 'hijaiyah', 'hard'),
('hh6', 'Huruf yang bentuk akhirnya sering turun ke bawah garis adalah?', '["ر","و","ن","ل"]', 2, 'Nun dalam beberapa posisi akhir turun ke bawah garis.', 'Bentuk sambung huruf dapat berubah tergantung posisinya di awal, tengah, atau akhir.', 'hijaiyah', 'hard'),
('hh7', 'Huruf manakah yang dibaca "Dhad"?', '["ض","ظ","د","ذ"]', 0, 'Huruf ض dibaca dhad.', 'Dhad (ض) adalah huruf tebal yang menjadi ciri khas bahasa Arab.', 'hijaiyah', 'hard'),
('hh8', 'Huruf apakah yang tepat untuk bunyi glottal stop atau hamzah?', '["ع","ء","ا","ه"]', 1, 'Simbol ء adalah hamzah.', 'Hamzah (ء) melambangkan hentian suara di pangkal tenggorokan.', 'hijaiyah', 'hard'),
('hh9', 'Huruf yang mirip ف tetapi memiliki dua titik di atas adalah?', '["ق","ك","ث","ت"]', 0, 'Huruf ق adalah qaf.', 'Qaf (ق) mirip fa dalam sebagian bentuk, namun qaf memiliki dua titik di atas.', 'hijaiyah', 'hard'),
('hh10', 'Huruf manakah yang berada sebelum ي dalam urutan hijaiyah?', '["ه","و","ن","م"]', 1, 'Sebelum ya adalah waw.', 'Waw (و) berada tepat sebelum ya (ي) dalam urutan hijaiyah.', 'hijaiyah', 'hard'),

-- Tajwid medium
('tm1', 'Apa hukum nun sukun bertemu huruf ب?', '["Ikhfa","Iqlab","Izhar","Idgham"]', 1, 'Nun sukun atau tanwin yang bertemu ب dibaca iqlab.', 'Iqlab mengubah bunyi nun sukun atau tanwin menjadi mim samar ketika bertemu ba.', 'tajwid', 'medium'),
('tm2', 'Berapa harakat bacaan Mad Thobi''i?', '["2","3","4","6"]', 0, 'Mad thobi''i dibaca 2 harakat.', 'Mad thobi''i adalah mad asli yang panjangnya dua harakat.', 'tajwid', 'medium'),
('tm3', 'Huruf mana yang termasuk huruf qalqalah?', '["س","ط","ل","م"]', 1, 'Huruf ط termasuk huruf qalqalah.', 'Huruf qalqalah ada lima: ق ط ب ج د.', 'tajwid', 'medium'),
('tm4', 'Apa hukum alif lam pada kata الشَّمْسُ?', '["Qamariyah","Syamsiyah","Ghunnah","Mad"]', 1, 'Lam pada kata tersebut tidak dibaca jelas sehingga termasuk syamsiyah.', 'Alif lam syamsiyah terjadi ketika lam lebur ke huruf syamsiyah sesudahnya.', 'tajwid', 'medium'),
('tm5', 'Nun bertasydid wajib dibaca dengan?', '["Mad","Qalqalah","Ghunnah","Waqf"]', 2, 'Nun bertasydid wajib dibaca ghunnah.', 'Ghunnah adalah dengung dua harakat pada nun atau mim bertasydid.', 'tajwid', 'medium'),
('tm6', 'Mad yang terjadi karena huruf mad bertemu hamzah dalam satu kata disebut?', '["Mad jaiz munfasil","Mad wajib muttasil","Mad aridh lissukun","Mad lazim"]', 1, 'Mad wajib muttasil terjadi dalam satu kata.', 'Jika huruf mad dan hamzah berada dalam satu kata, hukumnya mad wajib muttasil.', 'tajwid', 'medium'),
('tm7', 'Jika tanwin bertemu huruf ل, maka hukumnya adalah?', '["Idgham bighunnah","Idgham bilaghunnah","Iqlab","Izhar"]', 1, 'Tanwin bertemu lam dibaca idgham tanpa dengung.', 'Idgham bilaghunnah berlaku ketika nun sukun atau tanwin bertemu lam atau ra.', 'tajwid', 'medium'),
('tm8', 'Huruf ikhfa berjumlah?', '["5","10","15","17"]', 2, 'Huruf ikhfa haqiqi berjumlah 15.', 'Ikhfa haqiqi terjadi ketika nun sukun atau tanwin bertemu salah satu dari 15 huruf ikhfa.', 'tajwid', 'medium'),
('tm9', 'Apa nama hukum berhenti sempurna pada akhir kalimat yang maknanya telah lengkap?', '["Waqf Tam","Waqf Hasan","Qalqalah Kubra","Mad Lazim"]', 0, 'Berhenti sempurna disebut waqf tam.', 'Waqf tam dilakukan pada akhir makna yang utuh dan tidak terkait erat dengan sesudahnya.', 'tajwid', 'medium'),
('tm10', 'Alif lam dibaca jelas pada kata الْقَمَرُ karena termasuk?', '["Syamsiyah","Qamariyah","Ghunnah","Mad thobi''i"]', 1, 'Lam dibaca jelas karena termasuk alif lam qamariyah.', 'Huruf qamariyah membuat lam pada alif lam dibaca jelas.', 'tajwid', 'medium'),

-- Tajwid hard
('th1', 'Perbedaan utama Idgham Bighunnah dan Bilaghunnah adalah?', '["Posisi nun","Adanya dengung","Jumlah huruf mad","Jenis waqf"]', 1, 'Perbedaannya terletak pada adanya dengung.', 'Idgham bighunnah dibaca dengan dengung, sedangkan bilaghunnah tanpa dengung.', 'tajwid', 'hard'),
('th2', 'Pada lafaz مِنْ خَيْرٍ, hukum nun sukun adalah?', '["Izhar Halqi","Ikhfa","Iqlab","Idgham"]', 0, 'Nun sukun bertemu kha termasuk izhar halqi.', 'Huruf kha (خ) termasuk huruf halqi sehingga nun sukun dibaca jelas.', 'tajwid', 'hard'),
('th3', 'Mad jaiz munfasil terjadi saat?', '["Mad bertemu sukun dalam satu kata","Mad di akhir kata bertemu hamzah di awal kata berikutnya","Mad bertemu tasydid","Mad karena waqf"]', 1, 'Mad jaiz munfasil terjadi pada dua kata yang terpisah.', 'Ciri mad jaiz munfasil adalah huruf mad di akhir kata dan hamzah di awal kata setelahnya.', 'tajwid', 'hard'),
('th4', 'Qalqalah kubra biasanya terdengar saat?', '["Huruf qalqalah berharakat fathah","Huruf qalqalah di akhir kata ketika waqf","Huruf mim bertasydid","Lam qamariyah"]', 1, 'Qalqalah kubra muncul ketika berhenti pada huruf qalqalah di akhir kata.', 'Pantulan qalqalah kubra lebih kuat karena terjadi saat waqf.', 'tajwid', 'hard'),
('th5', 'Pada lafaz مِنْ كُلِّ, hukum nun sukun adalah?', '["Ikhfa","Iqlab","Izhar","Idgham bilaghunnah"]', 0, 'Nun sukun bertemu kaf dibaca ikhfa.', 'Ikhfa terjadi ketika nun sukun atau tanwin bertemu salah satu dari 15 huruf ikhfa, termasuk huruf kaf.', 'tajwid', 'hard'),
('th6', 'Mana pasangan huruf untuk Idgham Bilaghunnah?', '["ي dan و","م dan ن","ر dan ل","ق dan ك"]', 2, 'Idgham bilaghunnah hanya berlaku untuk ra dan lam.', 'Ra (ر) dan lam (ل) adalah huruf idgham tanpa dengung.', 'tajwid', 'hard'),
('th7', 'Berapa harakat bacaan Mad Lazim Kilmi?', '["2","4","5","6"]', 3, 'Mad lazim kilmi dibaca 6 harakat.', 'Mad lazim adalah salah satu mad far''i yang dibaca panjang enam harakat.', 'tajwid', 'hard'),
('th8', 'Apa hukum tanwin pada frasa سَمِيعٌ بَصِيرٌ?', '["Izhar","Ikhfa","Iqlab","Idgham"]', 2, 'Tanwin bertemu ba dibaca iqlab.', 'Iqlab mengubah nun/tanwin menjadi bunyi mim samar sebelum ba.', 'tajwid', 'hard'),
('th9', 'Huruf syamsiyah berfungsi membuat bunyi?', '["Lam dibaca jelas","Lam lebur pada huruf sesudahnya","Mad menjadi pendek","Nun menjadi mim"]', 1, 'Pada alif lam syamsiyah, bunyi lam tidak terdengar jelas.', 'Huruf syamsiyah menyebabkan lam diidghamkan ke huruf setelahnya.', 'tajwid', 'hard'),
('th10', 'Contoh ghunnah wajib adalah ketika terdapat?', '["Nun mati bertemu alif","Mim atau nun bertasydid","Lam sukun","Huruf mad di akhir kata"]', 1, 'Ghunnah wajib terjadi pada mim atau nun bertasydid.', 'Ghunnah adalah dengung dua harakat pada mim dan nun yang bertasydid.', 'tajwid', 'hard'),

-- Worship medium
('wm1', 'Rukun wudhu yang harus dilakukan setelah membasuh wajah adalah?', '["Membasuh kaki","Mengusap kepala","Membasuh tangan sampai siku","Membaca niat keras-keras"]', 2, 'Setelah wajah, wudhu dilanjutkan dengan membasuh tangan sampai siku.', 'Urutan rukun wudhu yang umum dipelajari: niat, wajah, tangan, kepala, kaki, dan tertib.', 'worship', 'medium'),
('wm2', 'Shalat witir paling sedikit terdiri dari berapa rakaat?', '["1","2","3","5"]', 0, 'Witir paling sedikit satu rakaat.', 'Shalat witir adalah shalat sunnah penutup malam dan minimal dilakukan satu rakaat.', 'worship', 'medium'),
('wm3', 'Saat sujud, anggota badan yang menempel ke lantai berjumlah?', '["5","6","7","8"]', 2, 'Ada tujuh anggota sujud.', 'Tujuh anggota sujud adalah dahi, dua telapak tangan, dua lutut, dan dua ujung kaki.', 'worship', 'medium'),
('wm4', 'Apa hukum puasa Ramadan bagi muslim yang baligh dan mampu?', '["Sunnah","Wajib","Makruh","Mubah"]', 1, 'Puasa Ramadan hukumnya wajib.', 'Puasa Ramadan adalah salah satu rukun Islam yang diwajibkan bagi muslim yang memenuhi syarat.', 'worship', 'medium'),
('wm5', 'Zakat fitrah wajib ditunaikan sebelum?', '["Shalat Id","Tarawih","Subuh 1 Syawal","Berbuka puasa"]', 0, 'Zakat fitrah sebaiknya ditunaikan sebelum shalat Id.', 'Waktu terbaik membayar zakat fitrah adalah sebelum pelaksanaan shalat Idulfitri.', 'worship', 'medium'),
('wm6', 'Jumlah takbir pada shalat jenazah adalah?', '["2","3","4","5"]', 2, 'Shalat jenazah memiliki empat kali takbir.', 'Shalat jenazah dilakukan dengan empat takbir tanpa ruku dan sujud.', 'worship', 'medium'),
('wm7', 'Apa nama mandi wajib setelah selesai haid?', '["Wudhu","Tayamum","Ghusl","Istinja"]', 2, 'Mandi wajib disebut ghusl.', 'Ghusl adalah mandi besar untuk menghilangkan hadas besar.', 'worship', 'medium'),
('wm8', 'Jika tidak ada air dan tidak bisa memakainya, bersuci dilakukan dengan?', '["Istinja","Qadha","Tayamum","Tasbih"]', 2, 'Pengganti wudhu atau mandi saat tidak ada air adalah tayamum.', 'Tayamum menggunakan debu suci sebagai keringanan saat tidak bisa menggunakan air.', 'worship', 'medium'),
('wm9', 'Rukun Islam yang berkaitan dengan kemampuan finansial dan fisik ke Baitullah adalah?', '["Zakat","Haji","Puasa","Syahadat"]', 1, 'Rukun tersebut adalah haji.', 'Haji diwajibkan sekali seumur hidup bagi muslim yang mampu.', 'worship', 'medium'),
('wm10', 'Bacaan tahiyat akhir dibaca pada posisi?', '["Berdiri","Ruku","Duduk","Sujud"]', 2, 'Tahiyat akhir dibaca sambil duduk.', 'Tahiyat akhir merupakan bagian penutup sebelum salam dalam shalat.', 'worship', 'medium'),

-- Worship hard
('wh1', 'Mana yang termasuk syarat sah shalat?', '["Baligh","Islam","Menutup aurat","Laki-laki"]', 2, 'Menutup aurat termasuk syarat sah shalat.', 'Syarat sah shalat mencakup suci dari hadas, suci tempat/pakaian, masuk waktu, menutup aurat, dan menghadap kiblat.', 'worship', 'hard'),
('wh2', 'Jika seseorang lupa rakaat lalu melakukan sujud sahwi, itu dilakukan karena?', '["Meninggalkan rukun dengan sengaja","Kelupaan dalam shalat","Tidak berwudhu","Tidak membaca iftitah"]', 1, 'Sujud sahwi dilakukan untuk menutup kekurangan karena lupa.', 'Sujud sahwi disyariatkan ketika terjadi lupa dalam shalat, misalnya ragu jumlah rakaat atau meninggalkan sunnah ab''adh.', 'worship', 'hard'),
('wh3', 'Nisab zakat emas secara umum setara dengan?', '["20 gram","40 gram","85 gram","100 gram"]', 2, 'Nisab emas umumnya sekitar 85 gram.', 'Nisab zakat emas yang umum dipakai adalah 85 gram emas murni.', 'worship', 'hard'),
('wh4', 'Ibadah umrah tidak mencakup salah satu amalan berikut:', '["Ihram","Wukuf di Arafah","Sa''i","Tahallul"]', 1, 'Wukuf di Arafah hanya ada dalam haji.', 'Rukun umrah meliputi ihram, tawaf, sa''i, dan tahallul.', 'worship', 'hard'),
('wh5', 'Seseorang yang sedang sakit dan tidak mampu berdiri saat shalat dapat?', '["Meninggalkan shalat","Shalat sesuai kemampuan","Hanya berdzikir","Menunggu sembuh"]', 1, 'Shalat tetap wajib sesuai kemampuan.', 'Syariat memberi keringanan: shalat dapat dilakukan duduk atau berbaring bila tidak mampu berdiri.', 'worship', 'hard'),
('wh6', 'Apa tujuan utama niat dalam ibadah?', '["Memperkeras bacaan","Membedakan ibadah satu dengan lainnya","Memperindah gerakan","Menambah rakaat"]', 1, 'Niat membedakan jenis ibadah dan menujukannya kepada Allah.', 'Niat adalah amalan hati yang membedakan ibadah wajib, sunnah, dan kebiasaan.', 'worship', 'hard'),
('wh7', 'Puasa seseorang batal jika dengan sengaja?', '["Tidur siang","Muntah tak sengaja","Makan dan minum","Bersiwak"]', 2, 'Makan dan minum dengan sengaja membatalkan puasa.', 'Hal yang membatalkan puasa antara lain makan, minum, dan hubungan suami istri dengan sengaja pada siang Ramadan.', 'worship', 'hard'),
('wh8', 'Tertib dalam wudhu berarti?', '["Mengulang dua kali","Mengikuti urutan rukun","Membaca doa panjang","Membasuh kaki dulu"]', 1, 'Tertib berarti mengikuti urutan yang benar.', 'Tertib adalah melakukan rukun-rukun wudhu sesuai urutan yang ditetapkan.', 'worship', 'hard'),
('wh9', 'Salam pertama dalam shalat berfungsi untuk?', '["Membatalkan wudhu","Menutup shalat","Mengganti tahiyat","Memulai shalat"]', 1, 'Salam menutup shalat.', 'Salam adalah penanda keluarnya seseorang dari ibadah shalat.', 'worship', 'hard'),
('wh10', 'Yang wajib dibayar oleh orang yang sengaja membatalkan puasa Ramadan dengan hubungan suami istri adalah?', '["Qadha saja","Fidyah saja","Kafarat dan qadha","Tidak ada"]', 2, 'Dalam fikih, ada kewajiban kafarat berat dan qadha.', 'Beberapa pelanggaran puasa Ramadan mewajibkan qadha dan kafarat, terutama jima'' di siang hari Ramadan.', 'worship', 'hard'),

-- General medium
('gm1', 'Kitab suci yang diturunkan kepada Nabi Musa AS adalah?', '["Injil","Taurat","Zabur","Al-Qur''an"]', 1, 'Nabi Musa menerima Taurat.', 'Allah menurunkan Taurat kepada Nabi Musa AS sebagai petunjuk bagi Bani Israil.', 'general', 'medium'),
('gm2', 'Siapakah ayah dari Nabi Ismail AS?', '["Nabi Yaqub","Nabi Ibrahim","Nabi Nuh","Nabi Dawud"]', 1, 'Ayah Nabi Ismail adalah Nabi Ibrahim AS.', 'Nabi Ibrahim AS dikenal sebagai ayah para nabi, termasuk Nabi Ismail AS.', 'general', 'medium'),
('gm3', 'Malaikat yang bertugas menyampaikan wahyu adalah?', '["Mikail","Israfil","Jibril","Izrail"]', 2, 'Jibril menyampaikan wahyu kepada para nabi.', 'Malaikat Jibril memiliki tugas utama menyampaikan wahyu dari Allah.', 'general', 'medium'),
('gm4', 'Peristiwa turunnya Al-Qur''an pertama kali terjadi di?', '["Gua Tsur","Madinah","Gua Hira","Thaif"]', 2, 'Wahyu pertama turun di Gua Hira.', 'Nabi Muhammad SAW menerima wahyu pertama di Gua Hira saat berkhalwat.', 'general', 'medium'),
('gm5', 'Bulan hijriyah yang digunakan untuk puasa wajib adalah?', '["Muharram","Syawal","Ramadan","Dzulqa''dah"]', 2, 'Puasa wajib dilaksanakan pada bulan Ramadan.', 'Ramadan adalah bulan diturunkannya Al-Qur''an dan bulan diwajibkannya puasa.', 'general', 'medium'),
('gm6', 'Apa nama tahun kelahiran Nabi Muhammad SAW yang terkenal dalam sejarah?', '["Tahun Gajah","Tahun Hijrah","Tahun Badar","Tahun Fathu Makkah"]', 0, 'Nabi lahir pada Tahun Gajah.', 'Tahun Gajah dikenal karena upaya penyerangan Ka''bah oleh pasukan bergajah.', 'general', 'medium'),
('gm7', 'Siapakah khalifah pertama setelah wafatnya Nabi Muhammad SAW?', '["Umar bin Khattab","Utsman bin Affan","Ali bin Abi Thalib","Abu Bakar Ash-Shiddiq"]', 3, 'Khalifah pertama adalah Abu Bakar.', 'Abu Bakar Ash-Shiddiq diangkat sebagai khalifah pertama kaum muslimin.', 'general', 'medium'),
('gm8', 'Arah kiblat umat Islam menghadap ke?', '["Masjid Nabawi","Baitul Maqdis","Ka''bah","Gunung Uhud"]', 2, 'Kiblat umat Islam adalah Ka''bah.', 'Ka''bah di Masjidil Haram menjadi arah kiblat dalam shalat.', 'general', 'medium'),
('gm9', 'Lailatul Qadar terjadi pada bulan?', '["Rajab","Ramadan","Safar","Muharram"]', 1, 'Lailatul Qadar terjadi di bulan Ramadan.', 'Lailatul Qadar adalah malam yang lebih baik dari seribu bulan dan dicari pada akhir Ramadan.', 'general', 'medium'),
('gm10', 'Apa arti hijrah Nabi ke Madinah bagi dakwah Islam?', '["Akhir kenabian","Awal dakwah sembunyi-sembunyi","Permulaan masyarakat Islam yang kuat","Turunnya wahyu pertama"]', 2, 'Hijrah membuka fase baru pembangunan masyarakat Islam.', 'Setelah hijrah ke Madinah, dakwah Islam berkembang dengan landasan sosial dan politik yang lebih kuat.', 'general', 'medium'),

-- General hard
('gh1', 'Nama lain surat Al-Fatihah yang berarti "tujuh ayat yang diulang-ulang" adalah?', '["Al-Mulk","As-Sab''ul Matsani","Yasin","An-Nur"]', 1, 'Al-Fatihah disebut As-Sab''ul Matsani.', 'Al-Fatihah memiliki banyak nama, salah satunya As-Sab''ul Matsani.', 'general', 'hard'),
('gh2', 'Perang yang pertama kali terjadi antara kaum muslimin dan Quraisy adalah?', '["Uhud","Khandaq","Badar","Hunain"]', 2, 'Perang Badar adalah perang besar pertama.', 'Perang Badar terjadi pada tahun kedua Hijriyah dan menjadi titik penting dalam sejarah Islam.', 'general', 'hard'),
('gh3', 'Shahabat yang dijuluki Al-Faruq adalah?', '["Abu Bakar","Umar bin Khattab","Ali bin Abi Thalib","Bilal bin Rabah"]', 1, 'Al-Faruq adalah gelar Umar bin Khattab.', 'Umar dikenal tegas dalam membedakan kebenaran dan kebatilan.', 'general', 'hard'),
('gh4', 'Apa nama piagam yang mengatur kehidupan bersama di Madinah?', '["Piagam Aqabah","Piagam Hudaibiyah","Piagam Madinah","Piagam Tabuk"]', 2, 'Dokumen itu dikenal sebagai Piagam Madinah.', 'Piagam Madinah adalah salah satu contoh awal tata kelola masyarakat majemuk dalam Islam.', 'general', 'hard'),
('gh5', 'Nabi yang terkenal dengan mukjizat mampu memahami bahasa burung adalah?', '["Nabi Musa","Nabi Sulaiman","Nabi Yusuf","Nabi Zakaria"]', 1, 'Mukjizat itu diberikan kepada Nabi Sulaiman AS.', 'Nabi Sulaiman AS dikaruniai kerajaan besar dan kemampuan memahami bahasa hewan.', 'general', 'hard'),
('gh6', 'Surah terpanjang dalam Al-Qur''an adalah?', '["Ali Imran","Al-Baqarah","An-Nisa","Yusuf"]', 1, 'Surah Al-Baqarah adalah yang terpanjang.', 'Al-Baqarah memiliki 286 ayat dan menjadi surah terpanjang dalam Al-Qur''an.', 'general', 'hard'),
('gh7', 'Tokoh yang pertama kali mengumandangkan adzan dalam Islam adalah?', '["Bilal bin Rabah","Abu Hurairah","Muadz bin Jabal","Salman Al-Farisi"]', 0, 'Bilal bin Rabah dikenal sebagai muadzin pertama.', 'Bilal bin Rabah RA memiliki suara yang merdu dan diberi amanah untuk adzan.', 'general', 'hard'),
('gh8', 'Bai''at Aqabah terjadi sebelum peristiwa?', '["Isra Mi''raj","Hijrah ke Madinah","Fathu Makkah","Perang Uhud"]', 1, 'Bai''at Aqabah membuka jalan menuju hijrah.', 'Bai''at Aqabah memperkuat dukungan penduduk Yatsrib terhadap dakwah Rasulullah.', 'general', 'hard'),
('gh9', 'Apa nama induk hadis yang disusun Imam Bukhari?', '["Al-Muwaththa","Sahih al-Bukhari","Sunan Abu Dawud","Musnad Ahmad"]', 1, 'Karya terkenal Imam Bukhari adalah Sahih al-Bukhari.', 'Sahih al-Bukhari adalah salah satu kitab hadis paling otoritatif dalam Islam.', 'general', 'hard'),
('gh10', 'Rukun iman yang berkaitan dengan percaya adanya hari pembalasan adalah?', '["Iman kepada kitab","Iman kepada qada dan qadar","Iman kepada hari akhir","Iman kepada rasul"]', 2, 'Hari pembalasan termasuk iman kepada hari akhir.', 'Iman kepada hari akhir mencakup keyakinan tentang kebangkitan, hisab, surga, dan neraka.', 'general', 'hard'),

-- Fiqih easy
('fe1', 'Apa hukum shalat lima waktu bagi muslim yang baligh dan berakal?', '["Sunnah","Wajib","Makruh","Mubah"]', 1, 'Shalat lima waktu hukumnya wajib.', 'Shalat lima waktu merupakan kewajiban pokok yang harus dijaga oleh setiap muslim yang memenuhi syarat.', 'fiqih', 'easy'),
('fe2', 'Sebelum shalat, seseorang harus berada dalam keadaan?', '["Lapar","Suci","Mengantuk","Diam"]', 1, 'Suci dari hadas adalah syarat sah shalat.', 'Bersuci dari hadas kecil dan besar menjadi syarat penting sebelum melaksanakan shalat.', 'fiqih', 'easy'),
('fe3', 'Bersuci dengan debu suci saat tidak ada air disebut?', '["Istinja","Ghusl","Tayamum","Wudhu"]', 2, 'Pengganti bersuci saat tidak ada air adalah tayamum.', 'Tayamum adalah keringanan dari Allah saat air tidak tersedia atau tidak bisa digunakan.', 'fiqih', 'easy'),
('fe4', 'Zakat fitrah biasanya dikeluarkan pada bulan?', '["Muharram","Ramadan","Rabiul Awal","Safar"]', 1, 'Zakat fitrah ditunaikan di akhir Ramadan.', 'Zakat fitrah berkaitan dengan penyucian puasa dan dibayar menjelang Idulfitri.', 'fiqih', 'easy'),
('fe5', 'Puasa Ramadan dilaksanakan pada waktu?', '["Dari subuh sampai maghrib","Dari zuhur sampai isya","Dari asar sampai subuh","Sepanjang hari tanpa berbuka"]', 0, 'Puasa dimulai dari terbit fajar sampai terbenam matahari.', 'Puasa Ramadan menahan diri dari hal-hal yang membatalkan sejak fajar hingga maghrib.', 'fiqih', 'easy'),
('fe6', 'Rukun Islam yang dilakukan di Makkah bagi yang mampu adalah?', '["Zakat","Puasa","Haji","Shalat"]', 2, 'Ibadah tersebut adalah haji.', 'Haji menjadi kewajiban sekali seumur hidup bagi muslim yang mampu.', 'fiqih', 'easy'),
('fe7', 'Air yang suci dan menyucikan boleh digunakan untuk?', '["Bermain","Wudhu","Masak saja","Mandi biasa saja"]', 1, 'Air suci dan menyucikan dapat dipakai untuk wudhu.', 'Dalam fikih thaharah, air mutlak dapat digunakan untuk bersuci.', 'fiqih', 'easy'),
('fe8', 'Jumlah rakaat shalat Maghrib adalah?', '["2","3","4","5"]', 1, 'Shalat Maghrib terdiri dari tiga rakaat.', 'Maghrib adalah shalat wajib setelah matahari terbenam dengan jumlah tiga rakaat.', 'fiqih', 'easy'),
('fe9', 'Membayar zakat mal diwajibkan ketika harta telah mencapai?', '["Niat","Nisab","Umur","Perjalanan"]', 1, 'Salah satu syarat zakat mal adalah mencapai nisab.', 'Nisab adalah batas minimal harta yang mewajibkan seseorang mengeluarkan zakat.', 'fiqih', 'easy'),
('fe10', 'Arah yang dihadapi saat shalat adalah?', '["Madinah","Gunung Arafah","Ka''bah","Masjid Aqsa"]', 2, 'Saat shalat muslim menghadap kiblat yaitu Ka''bah.', 'Menghadap kiblat merupakan salah satu syarat sah shalat bagi yang mampu.', 'fiqih', 'easy'),

-- Fiqih medium
('fm1', 'Salah satu hal yang membatalkan wudhu adalah?', '["Tidur nyenyak","Membaca Al-Qur''an","Berdzikir","Tersenyum"]', 0, 'Tidur nyenyak termasuk pembatal wudhu menurut pelajaran fikih dasar.', 'Pembatal wudhu mencakup keluarnya sesuatu dari dua jalan, hilang akal, dan beberapa keadaan lainnya.', 'fiqih', 'medium'),
('fm2', 'Apabila seseorang tertinggal shalat wajib karena lupa, ia wajib?', '["Meninggalkannya","Mengqadha","Membayar fidyah saja","Mengganti dengan sedekah"]', 1, 'Shalat yang tertinggal karena lupa harus diqadha.', 'Qadha shalat dilakukan sesegera mungkin ketika seseorang ingat.', 'fiqih', 'medium'),
('fm3', 'Najis yang dianggap ringan dan cara menyucikannya dengan dipercik air disebut?', '["Mughallazhah","Mutawassithah","Mukhaffafah","Ma''fu"]', 2, 'Najis mukhaffafah termasuk najis ringan.', 'Contoh najis ringan yang dikenal di pelajaran fikih adalah air kencing bayi laki-laki yang belum makan selain ASI.', 'fiqih', 'medium'),
('fm4', 'Salah satu syarat wajib zakat mal adalah kepemilikan harta selama?', '["1 pekan","1 bulan","1 haul","10 hari"]', 2, 'Haul berarti kepemilikan selama satu tahun hijriyah.', 'Sebagian jenis zakat mal mensyaratkan harta dimiliki selama satu haul.', 'fiqih', 'medium'),
('fm5', 'Orang yang sedang safar mendapat keringanan untuk?', '["Membatalkan semua ibadah","Menjama'' dan mengqashar shalat tertentu","Tidak berpuasa selamanya","Mengganti zakat dengan sedekah"]', 1, 'Musafir mendapat rukhsah untuk jamak dan qashar.', 'Fikih memberi keringanan kepada musafir dalam kondisi tertentu, termasuk jamak dan qashar.', 'fiqih', 'medium'),
('fm6', 'Talak yang masih memungkinkan rujuk selama masa iddah disebut?', '["Bain kubra","Bain sughra","Raj''i","Fasakh"]', 2, 'Talak raj''i memungkinkan rujuk selama iddah.', 'Dalam fikih keluarga, talak raj''i berbeda dengan talak bain yang memutus hubungan lebih jauh.', 'fiqih', 'medium'),
('fm7', 'Bagian waris untuk suami jika istri wafat tanpa anak adalah?', '["1/2","1/4","1/8","2/3"]', 0, 'Dalam faraidh, suami mendapat setengah bila istri tidak punya anak.', 'Pembagian waris Islam memiliki ketentuan tetap untuk beberapa ahli waris.', 'fiqih', 'medium'),
('fm8', 'Menyentuh mushaf Al-Qur''an menurut pelajaran fikih dasar sebaiknya dalam keadaan?', '["Lapar","Suci","Safar","Sedih"]', 1, 'Dianjurkan dan diajarkan menyentuh mushaf dalam keadaan suci.', 'Adab terhadap Al-Qur''an termasuk menjaga kesucian saat menyentuh mushaf.', 'fiqih', 'medium'),
('fm9', 'Salah satu rukun nikah adalah adanya?', '["Panggung","Wali","Musik","Mahar besar"]', 1, 'Wali termasuk rukun nikah.', 'Rukun nikah meliputi calon suami, calon istri, wali, dua saksi, dan ijab kabul.', 'fiqih', 'medium'),
('fm10', 'Hewan ternak yang wajib dizakati jika memenuhi syarat disebut harta?', '["Perdagangan","Pertanian","An''am","Rikaz"]', 2, 'Harta ternak dikenal dengan istilah an''am.', 'Zakat hewan ternak berlaku pada unta, sapi, dan kambing dengan syarat tertentu.', 'fiqih', 'medium'),

-- Fiqih hard
('fh1', 'Jual beli yang sah menurut fikih mensyaratkan adanya?', '["Riba","Akad yang jelas","Paksaan","Barang haram"]', 1, 'Akad yang jelas menjadi syarat penting dalam muamalah.', 'Fikih muamalah menekankan kerelaan kedua pihak dan kejelasan objek akad.', 'fiqih', 'hard'),
('fh2', 'Dalam pembagian waris, anak laki-laki dibanding anak perempuan mendapat bagian?', '["Sama rata","Dua banding satu","Setengahnya","Seperempatnya"]', 1, 'Secara umum bagian anak laki-laki dua kali anak perempuan.', 'Ketentuan faraidh menetapkan bagian tertentu bagi para ahli waris sesuai nas.', 'fiqih', 'hard'),
('fh3', 'Masa tunggu bagi perempuan yang ditalak disebut?', '["Haul","Iddah","Kifarat","Nisab"]', 1, 'Masa tunggu setelah perceraian disebut iddah.', 'Iddah memiliki hikmah menjaga kejelasan nasab dan memberi masa renungan.', 'fiqih', 'hard'),
('fh4', 'Salah satu larangan dalam ihram saat haji atau umrah adalah?', '["Bertalbiyah","Memakai wewangian","Tawaf","Sa''i"]', 1, 'Orang berihram dilarang memakai wewangian.', 'Larangan ihram bertujuan menjaga kekhusyukan dan adab selama manasik.', 'fiqih', 'hard'),
('fh5', 'Riba fadhl terjadi ketika?', '["Pinjam meminjam tanpa saksi","Pertukaran barang ribawi sejenis tidak sama takaran","Jual beli kredit","Zakat ditunda"]', 1, 'Riba fadhl berkaitan dengan pertukaran barang ribawi sejenis yang tidak setara.', 'Fikih melarang riba dalam bentuk tambahan yang tidak dibenarkan pada komoditas tertentu.', 'fiqih', 'hard'),
('fh6', 'Qashar shalat berlaku untuk shalat fardu yang berjumlah?', '["2 rakaat","3 rakaat","4 rakaat","Semua rakaat"]', 2, 'Qashar hanya berlaku pada shalat empat rakaat.', 'Musafir dapat mengqashar zuhur, asar, dan isya menjadi dua rakaat.', 'fiqih', 'hard'),
('fh7', 'Nazar yang isinya maksiat seharusnya?', '["Ditunaikan","Dibatalkan dan tidak boleh dilaksanakan","Diganti puasa","Diwariskan"]', 1, 'Nazar maksiat tidak boleh ditunaikan.', 'Ketaatan tidak bisa dibangun di atas pelanggaran terhadap syariat.', 'fiqih', 'hard'),
('fh8', 'Jika imam lupa lalu melakukan sujud sahwi, makmum yang mengikutinya seharusnya?', '["Keluar dari shalat","Tetap mengikuti imam","Diam saja tanpa sujud","Membatalkan wudhu"]', 1, 'Makmum mengikuti gerakan imam dalam sujud sahwi.', 'Kaedah umum shalat berjamaah adalah mengikuti imam selama bukan pada kemaksiatan.', 'fiqih', 'hard'),
('fh9', 'Wakaf berarti?', '["Meminjamkan uang dengan bunga","Menahan pokok harta dan mengalirkan manfaatnya","Menjual tanah murah","Membagikan waris sebelum wafat"]', 1, 'Wakaf menahan pokok harta dan menyerahkan manfaatnya untuk kebaikan.', 'Wakaf menjadi salah satu amal jariyah yang pahalanya terus mengalir.', 'fiqih', 'hard'),
('fh10', 'Dalam fikih, darah haid membuat seorang wanita tidak wajib melakukan?', '["Dzikir","Shalat","Sedekah","Doa"]', 1, 'Wanita haid tidak diwajibkan shalat sampai suci kembali.', 'Syariat memberi keringanan khusus bagi wanita haid terkait beberapa ibadah tertentu.', 'fiqih', 'hard'),

-- Sirah easy
('se1', 'Siapakah nabi terakhir dalam Islam?', '["Nabi Isa","Nabi Musa","Nabi Muhammad","Nabi Ibrahim"]', 2, 'Nabi Muhammad SAW adalah nabi terakhir.', 'Dalam aqidah Islam, Nabi Muhammad SAW adalah penutup para nabi.', 'sirah', 'easy'),
('se2', 'Di kota manakah Nabi Muhammad SAW dilahirkan?', '["Madinah","Makkah","Thaif","Yaman"]', 1, 'Rasulullah lahir di Makkah.', 'Makkah adalah kota kelahiran Nabi Muhammad SAW.', 'sirah', 'easy'),
('se3', 'Siapakah istri pertama Nabi Muhammad SAW?', '["Aisyah","Hafshah","Khadijah","Zainab"]', 2, 'Istri pertama beliau adalah Khadijah RA.', 'Khadijah RA sangat berjasa mendukung dakwah Nabi sejak awal.', 'sirah', 'easy'),
('se4', 'Ke mana Nabi hijrah dari Makkah?', '["Syam","Yaman","Madinah","Mesir"]', 2, 'Hijrah dilakukan dari Makkah ke Madinah.', 'Hijrah ke Madinah menjadi titik balik besar dalam sejarah Islam.', 'sirah', 'easy'),
('se5', 'Siapakah sahabat yang menemani Nabi di gua saat hijrah?', '["Umar bin Khattab","Ali bin Abi Thalib","Abu Bakar Ash-Shiddiq","Utsman bin Affan"]', 2, 'Abu Bakar menemani Rasulullah saat hijrah.', 'Dalam perjalanan hijrah, Rasulullah dan Abu Bakar sempat bersembunyi di Gua Tsur.', 'sirah', 'easy'),
('se6', 'Nama kakek Nabi Muhammad SAW adalah?', '["Abu Thalib","Abdul Muthalib","Abdullah","Hamzah"]', 1, 'Kakek beliau adalah Abdul Muthalib.', 'Abdul Muthalib adalah tokoh Quraisy yang merawat Nabi setelah wafatnya ibu beliau.', 'sirah', 'easy'),
('se7', 'Perang pertama yang diikuti Rasulullah melawan Quraisy adalah?', '["Uhud","Khandaq","Badar","Hunain"]', 2, 'Perang Badar adalah perang besar pertama.', 'Badar menjadi kemenangan penting bagi kaum muslimin pada awal Islam.', 'sirah', 'easy'),
('se8', 'Masjid pertama yang dibangun Nabi saat tiba di Madinah adalah?', '["Masjid Nabawi","Masjid Quba","Masjidil Haram","Masjid Aqsa"]', 1, 'Masjid Quba dibangun lebih dahulu.', 'Masjid Quba memiliki kedudukan istimewa dalam sejarah hijrah.', 'sirah', 'easy'),
('se9', 'Siapakah paman Nabi yang sangat memusuhi dakwah Islam?', '["Hamzah","Abu Lahab","Abbas","Abu Thalib"]', 1, 'Abu Lahab dikenal keras menentang dakwah Nabi.', 'Penentangan terhadap dakwah Rasulullah datang bahkan dari sebagian kerabat dekat.', 'sirah', 'easy'),
('se10', 'Perjanjian damai antara kaum muslimin dan Quraisy yang terkenal adalah?', '["Aqabah","Hudaibiyah","Badar","Tabuk"]', 1, 'Perjanjian itu adalah Hudaibiyah.', 'Perjanjian Hudaibiyah menjadi langkah strategis bagi perkembangan dakwah Islam.', 'sirah', 'easy'),

-- Sirah medium
('sm1', 'Nama ayah Nabi Muhammad SAW adalah?', '["Abdul Muthalib","Abdullah","Abu Talib","Harits"]', 1, 'Ayah beliau bernama Abdullah.', 'Abdullah wafat sebelum Nabi Muhammad SAW lahir.', 'sirah', 'medium'),
('sm2', 'Peristiwa Isra Mi''raj menunjukkan perjalanan Nabi dari Masjidil Haram ke?', '["Masjid Quba","Masjid Nabawi","Masjidil Aqsa","Gua Hira"]', 2, 'Isra berlangsung dari Masjidil Haram ke Masjidil Aqsa.', 'Isra Mi''raj adalah mukjizat besar yang juga membawa perintah shalat.', 'sirah', 'medium'),
('sm3', 'Suku Nabi Muhammad SAW adalah?', '["Aus","Khazraj","Quraisy","Tsaqif"]', 2, 'Rasulullah berasal dari suku Quraisy.', 'Quraisy adalah suku besar yang memegang peranan penting di Makkah.', 'sirah', 'medium'),
('sm4', 'Siapakah budak yang dimerdekakan dan diangkat menjadi anak asuh Nabi?', '["Bilal","Zaid bin Haritsah","Salman","Ammar"]', 1, 'Zaid bin Haritsah sangat dekat dengan Rasulullah.', 'Zaid bin Haritsah memiliki posisi khusus dalam kehidupan Rasulullah SAW.', 'sirah', 'medium'),
('sm5', 'Tahun kesedihan disebut demikian karena wafatnya Khadijah dan?', '["Hamzah","Abu Talib","Ali","Umar"]', 1, 'Abu Talib wafat pada tahun yang sama dengan Khadijah.', 'Tahun tersebut sangat berat bagi Rasulullah karena kehilangan dua pendukung besar.', 'sirah', 'medium'),
('sm6', 'Siapa sahabat yang pertama kali masuk Islam dari kalangan laki-laki dewasa?', '["Ali bin Abi Thalib","Abu Bakar Ash-Shiddiq","Umar bin Khattab","Utsman bin Affan"]', 1, 'Abu Bakar adalah lelaki dewasa pertama yang masuk Islam.', 'Abu Bakar memiliki peran awal yang sangat besar dalam dakwah.', 'sirah', 'medium'),
('sm7', 'Perang Uhud terjadi setelah?', '["Badar","Khandaq","Tabuk","Hunain"]', 0, 'Perang Uhud terjadi setelah Badar.', 'Urutan peristiwa penting membantu memahami perkembangan masyarakat Madinah.', 'sirah', 'medium'),
('sm8', 'Bai''at Aqabah dilakukan oleh penduduk?', '["Makkah","Habasyah","Yatsrib","Thaif"]', 2, 'Penduduk Yatsrib memberikan bai''at kepada Nabi.', 'Bai''at Aqabah membuka jalan hijrah dan pembentukan komunitas muslim di Madinah.', 'sirah', 'medium'),
('sm9', 'Siapa panglima muslim pada Perang Mu''tah yang pertama ditunjuk?', '["Khalid bin Walid","Zaid bin Haritsah","Sa''d bin Abi Waqqash","Amr bin Ash"]', 1, 'Zaid bin Haritsah menjadi panglima pertama dalam urutan komando Mu''tah.', 'Perang Mu''tah menunjukkan kedisiplinan dan strategi dalam pasukan muslim.', 'sirah', 'medium'),
('sm10', 'Fathu Makkah berarti?', '["Hijrah ke Habasyah","Penaklukan Makkah","Perang pertama","Perjanjian damai"]', 1, 'Fathu Makkah berarti pembebasan atau penaklukan Makkah.', 'Peristiwa Fathu Makkah menunjukkan kemenangan Islam yang disertai sikap pemaaf Rasulullah.', 'sirah', 'medium'),

-- Sirah hard
('sh1', 'Dalam strategi Perang Khandaq, siapa yang mengusulkan pembuatan parit?', '["Umar bin Khattab","Salman Al-Farisi","Ali bin Abi Thalib","Abu Ubaidah"]', 1, 'Usulan menggali parit berasal dari Salman Al-Farisi.', 'Ide strategi ini belum umum di Arabia dan terbukti efektif mempertahankan Madinah.', 'sirah', 'hard'),
('sh2', 'Hijrah pertama sebagian sahabat sebelum ke Madinah dilakukan ke?', '["Syam","Habasyah","Mesir","Bahrain"]', 1, 'Sebagian sahabat hijrah ke Habasyah.', 'Habasyah dipilih karena di sana ada raja yang dikenal adil.', 'sirah', 'hard'),
('sh3', 'Siapakah yang menggantikan Nabi tidur di tempat beliau pada malam hijrah?', '["Abu Bakar","Ali bin Abi Thalib","Utsman","Zubair"]', 1, 'Ali bin Abi Thalib tidur di tempat Rasulullah.', 'Peran Ali menunjukkan keberanian dan pengorbanan besar dalam mendukung hijrah.', 'sirah', 'hard'),
('sh4', 'Perjanjian Hudaibiyah secara lahir tampak berat, tetapi membawa dampak?', '["Melemahkan dakwah","Mempercepat penyebaran Islam","Mengakhiri shalat","Memecah kaum muslimin"]', 1, 'Hudaibiyah justru membuka jalan bagi penyebaran Islam yang lebih luas.', 'Masa damai membuat interaksi sosial meningkat dan dakwah berkembang pesat.', 'sirah', 'hard'),
('sh5', 'Perang Tabuk terkenal sebagai ekspedisi yang berat karena?', '["Lawan tidak hadir","Musim dingin panjang","Jarak jauh dan cuaca panas","Terjadi di malam hari"]', 2, 'Tabuk ditempuh dalam kondisi sangat panas dan jauh.', 'Ekspedisi Tabuk menguji kesiapan iman dan pengorbanan kaum muslimin.', 'sirah', 'hard'),
('sh6', 'Tokoh Quraisy yang terkenal sangat cerdas bernegosiasi pada Hudaibiyah adalah?', '["Abu Sufyan","Suhail bin Amr","Walid bin Mughirah","Amr bin Hisham"]', 1, 'Suhail bin Amr menjadi wakil Quraisy dalam perjanjian Hudaibiyah.', 'Namanya sering dikaitkan dengan detail negosiasi Hudaibiyah.', 'sirah', 'hard'),
('sh7', 'Setelah wafat Nabi, siapakah yang memimpin pengumpulan mushaf Al-Qur''an pertama kali secara resmi?', '["Utsman bin Affan","Abu Bakar Ash-Shiddiq","Ali bin Abi Thalib","Muawiyah"]', 1, 'Pengumpulan awal mushaf resmi dimulai pada masa Abu Bakar.', 'Inisiatif ini dilakukan demi menjaga Al-Qur''an setelah banyak penghafal wafat.', 'sirah', 'hard'),
('sh8', 'Siapakah komandan Quraisy pada Perang Uhud?', '["Abu Jahl","Khalid bin Walid","Abu Sufyan","Amr bin Ash"]', 2, 'Pemimpin Quraisy dalam Uhud adalah Abu Sufyan.', 'Memahami tokoh-tokoh yang terlibat membantu membaca dinamika sirah.', 'sirah', 'hard'),
('sh9', 'Piagam Madinah menunjukkan bahwa Nabi membangun masyarakat dengan prinsip?', '["Pemaksaan","Keadilan dan tanggung jawab bersama","Diskriminasi suku","Kekuasaan tanpa aturan"]', 1, 'Piagam Madinah menekankan hidup bersama dan tanggung jawab sosial.', 'Sirah Nabi menunjukkan kepemimpinan yang adil dan terstruktur di Madinah.', 'sirah', 'hard'),
('sh10', 'Peristiwa bai''at Ridwan terjadi karena rumor wafatnya?', '["Ali bin Abi Thalib","Abu Bakar","Utsman bin Affan","Bilal bin Rabah"]', 2, 'Bai''at Ridwan terjadi setelah tersebar rumor terbunuhnya Utsman.', 'Peristiwa ini menunjukkan loyalitas para sahabat kepada Rasulullah SAW.', 'sirah', 'hard')
ON CONFLICT (id) DO UPDATE
SET
    question = EXCLUDED.question,
    options = EXCLUDED.options,
    correct_answer = EXCLUDED.correct_answer,
    explanation = EXCLUDED.explanation,
    material = EXCLUDED.material,
    category_id = EXCLUDED.category_id,
    difficulty = EXCLUDED.difficulty,
    updated_at = NOW();
