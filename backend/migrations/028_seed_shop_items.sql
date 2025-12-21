-- Migration: Seed shop items
-- File: 028_seed_shop_items.sql

-- Healing Items
INSERT INTO shop_items (name, description, category, effect_type, effect_value, gtk_cost, icon_url) VALUES
('Potion', 'Restores 20 HP to one character', 'healing', 'heal_hp', 20, 50, '/assets/items/potion.png'),
('Super Potion', 'Restores 50 HP to one character', 'healing', 'heal_hp', 50, 150, '/assets/items/super_potion.png'),
('Hyper Potion', 'Restores 100 HP to one character', 'healing', 'heal_hp', 100, 300, '/assets/items/hyper_potion.png'),
('Max Potion', 'Fully restores HP to one character', 'healing', 'heal_full_hp', 100, 500, '/assets/items/max_potion.png'),
('Revive', 'Revives a fainted character with 50% HP', 'healing', 'revive', 50, 400, '/assets/items/revive.png'),
('Max Revive', 'Revives a fainted character with full HP', 'healing', 'revive', 100, 800, '/assets/items/max_revive.png')
ON CONFLICT DO NOTHING;

-- Status Cure Items
INSERT INTO shop_items (name, description, category, effect_type, effect_value, gtk_cost, icon_url) VALUES
('Antidote', 'Cures poison status', 'status', 'cure_poison', 0, 80, '/assets/items/antidote.png'),
('Burn Heal', 'Cures burn status', 'status', 'cure_burn', 0, 80, '/assets/items/burn_heal.png'),
('Ice Heal', 'Cures freeze status', 'status', 'cure_freeze', 0, 80, '/assets/items/ice_heal.png'),
('Paralyze Heal', 'Cures paralysis status', 'status', 'cure_paralysis', 0, 80, '/assets/items/paralyze_heal.png'),
('Awakening', 'Wakes up a sleeping character', 'status', 'cure_sleep', 0, 80, '/assets/items/awakening.png'),
('Full Heal', 'Cures all status effects', 'status', 'cure_all_status', 0, 300, '/assets/items/full_heal.png')
ON CONFLICT DO NOTHING;

-- PP Recovery Items
INSERT INTO shop_items (name, description, category, effect_type, effect_value, gtk_cost, icon_url) VALUES
('PP Potion', 'Restores 5 PP to one ability', 'pp', 'restore_pp', 5, 100, '/assets/items/pp_potion.png'),
('PP Max', 'Fully restores PP to one ability', 'pp', 'restore_pp_full', 0, 300, '/assets/items/pp_max.png'),
('Elixir', 'Restores 5 PP to all abilities', 'pp', 'restore_pp_all', 5, 500, '/assets/items/elixir.png'),
('Max Elixir', 'Fully restores PP to all abilities', 'pp', 'restore_pp_all_full', 0, 1000, '/assets/items/max_elixir.png')
ON CONFLICT DO NOTHING;

-- Battle Buff Items
INSERT INTO shop_items (name, description, category, effect_type, effect_value, gtk_cost, icon_url) VALUES
('X Attack', 'Raises Attack by 50% for 3 turns', 'battle', 'buff_attack', 150, 200, '/assets/items/x_attack.png'),
('X Defense', 'Raises Defense by 50% for 3 turns', 'battle', 'buff_defense', 150, 200, '/assets/items/x_defense.png'),
('X Speed', 'Raises Speed by 50% for 3 turns', 'battle', 'buff_speed', 150, 200, '/assets/items/x_speed.png'),
('Guard Spec', 'Prevents stat reduction for 5 turns', 'battle', 'buff_guard', 0, 300, '/assets/items/guard_spec.png'),
('Dire Hit', 'Greatly increases critical hit rate for 3 turns', 'battle', 'buff_critical', 100, 300, '/assets/items/dire_hit.png')
ON CONFLICT DO NOTHING;

-- Egg Incubation Items
INSERT INTO shop_items (name, description, category, effect_type, effect_value, gtk_cost, is_consumable, icon_url) VALUES
('Egg Scanner', 'Reveals egg stats before hatching', 'egg', 'reveal_stats', 0, 2000, true, '/assets/items/egg_scanner.png'),
('Basic Nest', 'Reduces incubation time by 25%', 'egg', 'accelerate', 25, 500, true, '/assets/items/basic_nest.png'),
('Advanced Incubator', 'Reduces incubation time by 50%', 'egg', 'accelerate', 50, 1000, true, '/assets/items/incubator.png'),
('Solar Heat Lamp', 'Reduces incubation time by 75%', 'egg', 'accelerate', 75, 1500, true, '/assets/items/solar_lamp.png'),
('Instant Hatch Serum', 'Hatches egg immediately', 'egg', 'instant_hatch', 100, 5000, true, '/assets/items/instant_hatch.png')
ON CONFLICT DO NOTHING;

-- Evolution Items (Future use)
INSERT INTO shop_items (name, description, category, effect_type, effect_value, gtk_cost, icon_url) VALUES
('Rare Candy', 'Instantly increases character level by 1', 'evolution', 'level_up', 1, 1000, '/assets/items/rare_candy.png'),
('Stat Booster', 'Permanently increases one stat by 5', 'evolution', 'boost_stat', 5, 2000, '/assets/items/stat_booster.png')
ON CONFLICT DO NOTHING;
