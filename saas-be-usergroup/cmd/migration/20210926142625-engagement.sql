
-- +migrate Up
CREATE TABLE engagements (
    id BIGSERIAL PRIMARY KEY,
    action VARCHAR(255) NOT NULL,
    action_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    browser VARCHAR(255) NOT NULL,
    browser_version VARCHAR(255) NOT NULL,
    platform VARCHAR(255) NOT NULL,
    identifier TEXT NOT NULL,
    host VARCHAR(255) NOT NULL,
    path VARCHAR(255) NOT NULL,
    full_url VARCHAR(255) NOT NULL,
    view_port VARCHAR(255) NOT NULL,
    os VARCHAR(255) NOT NULL,
    os_version VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NULL,
    updated_at TIMESTAMP NULL
);

-- +migrate Down
DROP TABLE engagements;