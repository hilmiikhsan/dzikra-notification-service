-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_push_notifications (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    detail VARCHAR(255) NULL,
    url VARCHAR(255) NULL,
    user_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Fungsi untuk memperbarui kolom `updated_at`
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger untuk tabel `user_push_notifications`
CREATE TRIGGER set_updated_at_user_push_notifications
BEFORE UPDATE ON user_push_notifications
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_push_notifications;
-- +goose StatementEnd