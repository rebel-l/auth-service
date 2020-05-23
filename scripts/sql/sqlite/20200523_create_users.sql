-- up
CREATE TABLE users (
    id CHAR(50) PRIMARY KEY NOT NULL,
    email VARCHAR(255) NOT NULL
);

-- down
DROP TABLE users;
