BEGIN;

CREATE TABLE users
(
    id SERIAL PRIMARY KEY
);

ALTER TABLE links
    ADD COLUMN user_id    INT,
    ADD COLUMN is_deleted BOOLEAN NOT NULL DEFAULT false;

-- ALTER TABLE links
--     ADD CONSTRAINT links_user_id_fkey
--         FOREIGN KEY (user_id) REFERENCES users (id);

-- todo: enable after merge https://github.com/Yandex-Practicum/go-autotests/pull/89/files

COMMIT;