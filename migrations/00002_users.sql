-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL
);

INSERT INTO users (email, password) VALUES ('dylan@pixelparade', '$2a$14$geSrruTlnk29Pt4G9RRzHuPS5I2OnZzkTxObagiIrkU0putlv8YPe')
ON CONFLICT (email) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
