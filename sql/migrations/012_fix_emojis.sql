-- +goose Up
UPDATE leaders SET discord_emoji_string = '<:Vercingetorix_Civ6:1520172393769861130>' WHERE id = 78;
UPDATE leaders SET discord_emoji_string = '<:Olympias_Civ6:1520172394726293514>' WHERE id = 79;
UPDATE leaders SET discord_emoji_string = '<:Te_KinichII_Civ6:1520172395699376411>' WHERE id = 80;
UPDATE leaders SET discord_emoji_string = '<:Ahiram_Civ6:1520172396815061032>' WHERE id = 81;
UPDATE leaders SET discord_emoji_string = '<:AlHasan_ibn_Sulaiman_Civ6:1520172398836453397>' WHERE id = 82;
UPDATE leaders SET discord_emoji_string = '<:Spearthrower_Owl_Civ6:1520172397758644326>' WHERE id = 83;
UPDATE leaders SET discord_emoji_string = '<:Kiviuq_Civ6:1520172400199864401>' WHERE id = 84;
UPDATE leaders SET discord_emoji_string = '<:Trisong_Detsen_Civ6:1520172401105571921>' WHERE id = 85;
UPDATE leaders SET discord_emoji_string = '<:Anacaona_Civ6:1520172403127222283>' WHERE id = 86;
UPDATE leaders SET discord_emoji_string = '<:Stanislaw_II_Civ6:1520172011022713012>' WHERE id = 87;
UPDATE leaders SET discord_emoji_string = '<:Maria_Theresa_Civ6:1520172402099617875>' WHERE id = 88;
UPDATE leaders SET discord_emoji_string = '<:Theodoric_Civ6:1520172392599654431>' WHERE id = 89;

-- +goose Down
UPDATE leaders SET discord_emoji_string = '<:Ambiorix_Civ6:1229388711087702046>' WHERE id = 78;
UPDATE leaders SET discord_emoji_string = '<:Alexander_Civ6:1229381348251406417>' WHERE id = 79;
UPDATE leaders SET discord_emoji_string = '<:Lady_Six_Sky_Civ6:1229393212645572618>' WHERE id = 80;
UPDATE leaders SET discord_emoji_string = '<:Dido_Civ6:1229388795565309984>' WHERE id = 81;
UPDATE leaders SET discord_emoji_string = '<:Suleiman_Civ6:1229599564277874720>' WHERE id = 82;
UPDATE leaders SET discord_emoji_string = '<:Montezuma_Civ6:1229393522717622313>' WHERE id = 83;
UPDATE leaders SET discord_emoji_string = '<:Gandhi_Civ6:1229388944014049300>' WHERE id = 84;
UPDATE leaders SET discord_emoji_string = '<:Gandhi_Civ6:1229388944014049300>' WHERE id = 85;
UPDATE leaders SET discord_emoji_string = NULL WHERE id = 86;
UPDATE leaders SET discord_emoji_string = NULL WHERE id = 87;
UPDATE leaders SET discord_emoji_string = NULL WHERE id = 88;
UPDATE leaders SET discord_emoji_string = NULL WHERE id = 89;
