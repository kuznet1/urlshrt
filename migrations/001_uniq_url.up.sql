BEGIN;

CREATE TABLE urls
(
    id  SERIAL PRIMARY KEY,
    url text NOT NULL UNIQUE
);

INSERT INTO urls (url)
SELECT DISTINCT url
FROM links;

ALTER TABLE links
    ADD COLUMN url_fk int;

UPDATE links l
SET url_fk = u.id
FROM urls u
WHERE l.url = u.url;


ALTER TABLE links
    ALTER COLUMN url_fk SET NOT NULL;
ALTER TABLE links
    ADD CONSTRAINT links_url_id_fkey FOREIGN KEY (url_fk)
        REFERENCES urls (id);

ALTER TABLE links
    DROP COLUMN url;

CREATE INDEX idx_links_url_id ON links(url_fk);

COMMIT;