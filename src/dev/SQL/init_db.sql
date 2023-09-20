CREATE TYPE role_type AS ENUM ('user', 'admin');

CREATE TABLE IF NOT EXISTS users (
	id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	name VARCHAR(35) NOT NULL,
	email VARCHAR(254) NOT NULL UNIQUE,
	about TEXT DEFAULT '',
	role role_type DEFAULT 'user',
	password VARCHAR(128) NOT NULL
);

CREATE TABLE IF NOT EXISTS tag (
	id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	name VARCHAR(35) NOT NULL,
	about TEXT DEFAULT '',
	color VARCHAR(10) NOT NULL
);

CREATE TABLE IF NOT EXISTS project (
	id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	name VARCHAR(35) NOT NULL,
	about TEXT DEFAULT '',
	color VARCHAR(10) NOT NULL,
	is_private boolean NOT NULL,
	total_count_hours FLOAT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS goal (
	id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	hours_count FLOAT NOT NULL,
	name VARCHAR(35) NOT NULL,
	project_id INT NOT NULL REFERENCES project(id) ON DELETE CASCADE,
	description TEXT DEFAULT '',
	time_start TIMESTAMP NOT NULL,
	time_end TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS entry (
	id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	project_id INT REFERENCES project(id) ON DELETE CASCADE,
	description TEXT DEFAULT '',
	time_start TIMESTAMP NOT NULL,
	time_end TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS tag_entry (
	tag_id INT NOT NULL REFERENCES tag(id) ON DELETE CASCADE,
	entry_id INT NOT NULL REFERENCES entry(id) ON DELETE CASCADE,
	PRIMARY KEY (tag_id, entry_id)
);

CREATE TABLE IF NOT EXISTS friend_relation (
	subscriber_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	PRIMARY KEY (subscriber_id, user_id)
);

INSERT INTO
	users (name, email, about, role, password)
VALUES
	('test', 'test', 'test', 'user', '');

-- insert into friend_relation (subscriber_id, user_id) VALUES (1, 2);
-- insert into friend_relation (subscriber_id, user_id) VALUES (2, 1);
-- select f1.subscriber_id, f1.user_id from friend_relation f1
-- join friend_relation f2 on f2.user_id = f1.subscriber_id and f2.subscriber_id = f1.user_id
-- where f1.user_id = 2;
CREATE OR REPLACE FUNCTION update_total_count_hours() 
RETURNS TRIGGER AS $$ 
BEGIN
	IF (TG_OP = 'INSERT') THEN
		UPDATE project
		SET total_count_hours = total_count_hours + (EXTRACT(EPOCH FROM (NEW.time_end - NEW.time_start))) / 3600 
		WHERE id = NEW.project_id;
	ELSIF (TG_OP = 'UPDATE') THEN
		UPDATE project
		SET total_count_hours = total_count_hours + (EXTRACT(EPOCH FROM (NEW.time_end - NEW.time_start)) - EXTRACT(EPOCH FROM (OLD.time_end - OLD.time_start))) / 3600
		WHERE id = NEW.project_id;
	ELSIF (TG_OP = 'DELETE') THEN
		UPDATE project
		SET total_count_hours = total_count_hours - (EXTRACT(EPOCH FROM (OLD.time_end - OLD.time_start))) / 3600
		WHERE id = OLD.project_id;
	END IF;
	RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_total_count_hours_trigger
AFTER INSERT OR UPDATE OR DELETE
	ON entry FOR EACH ROW EXECUTE FUNCTION update_total_count_hours();