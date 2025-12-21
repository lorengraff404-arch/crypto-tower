-- Phase 15.4: Add effect fields to moves table
ALTER TABLE moves ADD COLUMN IF NOT EXISTS effect_type VARCHAR(30);
ALTER TABLE moves ADD COLUMN IF NOT EXISTS effect_chance INT DEFAULT 0 CHECK (effect_chance >= 0 AND effect_chance <= 100);
ALTER TABLE moves ADD COLUMN IF NOT EXISTS effect_duration INT DEFAULT 0;
ALTER TABLE moves ADD COLUMN IF NOT EXISTS effect_magnitude INT DEFAULT 0;

-- Update some existing moves with effects (examples)
UPDATE moves SET 
    effect_type = 'burn',
    effect_chance = 30,
    effect_duration = 3,
    effect_magnitude = 5
WHERE name LIKE '%Fire%' OR name LIKE '%Flame%' AND effect_type IS NULL;

UPDATE moves SET 
    effect_type = 'freeze',
    effect_chance = 20,
    effect_duration = 2,
    effect_magnitude = 50
WHERE name LIKE '%Ice%' OR name LIKE '%Frost%' AND effect_type IS NULL;

UPDATE moves SET 
    effect_type = 'poison',
    effect_chance = 40,
    effect_duration = 4,
    effect_magnitude = 7
WHERE name LIKE '%Poison%' OR name LIKE '%Toxic%' AND effect_type IS NULL;

-- Create some new buff/debuff moves
INSERT INTO moves (name, class, element, power, accuracy, max_pp, current_pp, effect_type, effect_chance, effect_duration, effect_magnitude, description) VALUES
('Power Up', 'Support', 'NORMAL', 0, 100, 10, 10, 'atk_up', 100, 2, 30, 'Increases Attack by 30% for 2 turns'),
('Iron Defense', 'Support', 'STEEL', 0, 100, 10, 10, 'def_up', 100, 2, 30, 'Increases Defense by 30% for 2 turns'),
('Agility', 'Support', 'NORMAL', 0, 100, 10, 10, 'spd_up', 100, 2, 50, 'Increases Speed by 50% for 2 turns'),
('Intimidate', 'Support', 'DARK', 0, 100, 10, 10, 'atk_down', 100, 2, 30, 'Decreases enemy Attack by 30% for 2 turns'),
('Screech', 'Support', 'NORMAL', 0, 85, 10, 10, 'def_down', 100, 2, 30, 'Decreases enemy Defense by 30% for 2 turns')
ON CONFLICT DO NOTHING;
