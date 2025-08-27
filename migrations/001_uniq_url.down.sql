BEGIN;

ALTER TABLE links
    ADD COLUMN url text;

UPDATE links l
SET url = u.url
FROM urls u
WHERE l.url_fk = u.id;

ALTER TABLE links
    DROP CONSTRAINT IF EXISTS links_url_id_fkey;

DROP INDEX IF EXISTS idx_links_url_id;

ALTER TABLE links
    DROP COLUMN url_fk;

DROP TABLE IF EXISTS urls;

COMMIT