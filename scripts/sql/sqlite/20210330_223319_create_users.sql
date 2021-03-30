-- up
CREATE TABLE IF NOT EXISTS users (
    id CHAR(36) NOT NULL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    firstname VARCHAR(100) NOT NULL,
    lastname VARCHAR(100) NOT NULL,
    password TEXT,
    externalid TEXT,
    type VARCHAR(50) NOT NULL DEFAULT 'standard',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users (email);


CREATE TRIGGER IF NOT EXISTS users_after_update AFTER UPDATE ON users BEGIN
    UPDATE users SET modified_at = DATETIME('now') WHERE id = NEW.id;
end;


-- down
DROP INDEX IF EXISTS idx_users_email;

DROP TRIGGER IF EXISTS users_after_update;

DROP TABLE IF EXISTS users;
