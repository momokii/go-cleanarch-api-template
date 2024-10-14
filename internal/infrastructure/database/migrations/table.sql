CREATE TABLE users (
    id SERIAL NOT NULL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    role BIGINT NULL DEFAULT 2,
    is_deleted BOOLEAN DEFAULT FALSE,
    CONSTRAINT users_role_fkey FOREIGN KEY (role) REFERENCES role(id)
);

CREATE TABLE role (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

insert into table role (name) values ('Admin');
insert into table role (name) values ('User');
insert into table role (name) values ('SuperAdmin');

-- password: Admin123
insert into table users (username, password, role) values ('admin1', '$2a$16$M7vqg6tCH.2oGkD7ePaelupK.jEQfdkkhihGatKb.OlUfCkluOMh6', 3);