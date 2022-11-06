CREATE TABLE company (
	id uuid NOT NULL,
	name varchar(80) NOT NULL UNIQUE,
	description varchar(3000),
    total_employees int NOT NULL,
	is_registered BOOLEAN NOT NULL,
	type varchar(255) NOT NULL,
	CONSTRAINT "company_pk" PRIMARY KEY ("id")
);
CREATE INDEX company_id_idx ON company USING HASH (id);