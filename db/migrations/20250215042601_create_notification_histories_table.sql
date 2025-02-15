-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS notification_histories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recipient VARCHAR(255) NOT NULL,
    template_id UUID NOT NULL, 
    status VARCHAR(10) NOT NULL,
    error_message TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

ALTER TABLE notification_histories ADD CONSTRAINT fk_notification_histories_template_id FOREIGN KEY (template_id) REFERENCES notification_templates(id) ON DELETE CASCADE ON UPDATE CASCADE;

-- Fungsi untuk memperbarui kolom `updated_at`
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger untuk tabel `notification_histories`
CREATE TRIGGER set_updated_at_notification_histories
BEFORE UPDATE ON notification_histories
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE notification_histories DROP CONSTRAINT IF EXISTS fk_notification_histories_template_id;
DROP TABLE IF EXISTS notification_histories;
-- +goose StatementEnd