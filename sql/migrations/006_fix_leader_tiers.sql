-- +goose Up
-- Cap leader tiers at 5.0 (F tier). Some seed values were set above the valid range.
UPDATE leaders
SET tier = 5.0
WHERE tier > 5.0;

-- +goose Down
-- No reliable reverse: previous values were invalid and are not recorded.
