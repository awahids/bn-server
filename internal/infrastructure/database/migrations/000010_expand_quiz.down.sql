DELETE FROM quiz_questions
WHERE id IN (
    'hm1','hm2','hm3','hm4','hm5','hm6','hm7','hm8','hm9','hm10',
    'hh1','hh2','hh3','hh4','hh5','hh6','hh7','hh8','hh9','hh10',
    'tm1','tm2','tm3','tm4','tm5','tm6','tm7','tm8','tm9','tm10',
    'th1','th2','th3','th4','th5','th6','th7','th8','th9','th10',
    'wm1','wm2','wm3','wm4','wm5','wm6','wm7','wm8','wm9','wm10',
    'wh1','wh2','wh3','wh4','wh5','wh6','wh7','wh8','wh9','wh10',
    'gm1','gm2','gm3','gm4','gm5','gm6','gm7','gm8','gm9','gm10',
    'gh1','gh2','gh3','gh4','gh5','gh6','gh7','gh8','gh9','gh10',
    'fe1','fe2','fe3','fe4','fe5','fe6','fe7','fe8','fe9','fe10',
    'fm1','fm2','fm3','fm4','fm5','fm6','fm7','fm8','fm9','fm10',
    'fh1','fh2','fh3','fh4','fh5','fh6','fh7','fh8','fh9','fh10',
    'se1','se2','se3','se4','se5','se6','se7','se8','se9','se10',
    'sm1','sm2','sm3','sm4','sm5','sm6','sm7','sm8','sm9','sm10',
    'sh1','sh2','sh3','sh4','sh5','sh6','sh7','sh8','sh9','sh10'
);

DELETE FROM quiz_categories
WHERE id IN ('fiqih', 'sirah');

ALTER TABLE quiz_questions
DROP COLUMN IF EXISTS difficulty_order;
