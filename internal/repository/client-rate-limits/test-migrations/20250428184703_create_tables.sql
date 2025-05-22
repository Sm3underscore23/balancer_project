-- +goose Up
-- +goose StatementBegin
CREATE TABLE clients_limits (
    id SERIAL PRIMARY KEY,
    client_id VARCHAR UNIQUE NOT NULL,
    capacity FLOAT NOT NULL,
    rate_per_sec FLOAT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE clients_limits;
-- +goose StatementEnd
