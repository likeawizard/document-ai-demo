CREATE TABLE receipts(
    id uuid NOT NULL,
    filename text,
    status text,
    mime_type text,
    path text,
    PRIMARY KEY(id)
);

CREATE TABLE tags(
    id SERIAL NOT NULL,
    name text NOT NULL,
    PRIMARY KEY(id)
);

CREATE TABLE tags_to_receipts(
    id SERIAL NOT NULL,
    tag_id INT NOT NULL,
    receipt_id UUID NOT NULL,
    PRIMARY KEY(id),
    CONSTRAINT fk_tag FOREIGN KEY (tag_id) REFERENCES tags(id),
    CONSTRAINT fk_receipt FOREIGN KEY (receipt_id) REFERENCES receipts(id)
);