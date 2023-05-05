CREATE TYPE status AS ENUM('-1', '0', '1', '2');

ALTER TABLE comments
    ADD COLUMN Processed_Slug text,
    ADD COLUMN Processed_Body text,
    ADD COLUMN Processed_Author text,
    ADD COLUMN Process_Status status;