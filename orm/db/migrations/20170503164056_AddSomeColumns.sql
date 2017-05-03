
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE aman.verticals ADD my_column VARCHAR(255);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE aman.verticals DROP COLUMN my_column;