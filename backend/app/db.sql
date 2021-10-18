CREATE USER catalog_user WITH ENCRYPTED PASSWORD 'password';
ALTER ROLE catalog_user SET client_encoding TO 'utf8';
ALTER ROLE catalog_user SET default_transaction_isolation TO 'read committed';
ALTER ROLE catalog_user SET timezone TO 'UTC';

DROP DATABASE IF EXISTS catalog_db;
CREATE DATABASE catalog_db WITH OWNER catalog_user;

# \c catalog_db catalog_user

DROP TABLE IF EXISTS groups;
CREATE TABLE groups(
   group_id INT GENERATED ALWAYS AS IDENTITY,
   group_name VARCHAR(255) NOT NULL,
   UNIQUE(group_name),
   invite_link VARCHAR(32) NOT NULL,
   UNIQUE(invite_link),
   select_link VARCHAR(32) NOT NULL,
   UNIQUE(select_link),
   PRIMARY KEY(group_id)
);

DROP TABLE IF EXISTS users;
CREATE TABLE users(
   user_id INT GENERATED ALWAYS AS IDENTITY,
   user_name VARCHAR(255) NOT NULL,
   email TEXT NOT NULL,
   password TEXT NOT NULL,
   group_id INT,
   PRIMARY KEY(user_id),
   UNIQUE(user_name),
   CONSTRAINT fk_group
      FOREIGN KEY(group_id) 
	  REFERENCES groups(group_id)
	  ON DELETE SET NULL
);

DROP TABLE IF EXISTS groups_users;
CREATE TABLE groups_users(
   user_id INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
   group_id INT NOT NULL REFERENCES groups(group_id) ON DELETE CASCADE,
   CONSTRAINT pk_group_user PRIMARY KEY (user_id, group_id)
);

DROP TABLE IF EXISTS owned_groups;
CREATE TABLE owned_groups(
   user_id INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
   group_id INT NOT NULL REFERENCES groups(group_id) ON DELETE CASCADE,
   CONSTRAINT pk_owned_group PRIMARY KEY (user_id, group_id)
);

DROP TABLE IF EXISTS categories;
CREATE TABLE categories(
   category_id INT GENERATED ALWAYS AS IDENTITY,
   category_name TEXT NOT NULL,
   group_id INT NOT NULL,
   PRIMARY KEY(category_id),
   UNIQUE(category_name),
   CONSTRAINT fk_group
      FOREIGN KEY(group_id) 
	  REFERENCES groups(group_id)
	  ON DELETE CASCADE
);

DROP TABLE IF EXISTS items;
CREATE TABLE items(
   item_id INT GENERATED ALWAYS AS IDENTITY,
   item_name TEXT NOT NULL,
   category_id INT NOT NULL,
   UNIQUE(item_name),
   PRIMARY KEY(item_id),
   CONSTRAINT fk_category
      FOREIGN KEY(category_id) 
	  REFERENCES categories(category_id)
	  ON DELETE CASCADE
);

CREATE INDEX users_id_idx ON users(user_id);
CREATE INDEX users_group_id_idx ON users(group_id);
CREATE INDEX items_category_id_idx ON items(category_id);

insert into groups(group_name) values('test_group');
insert into categories(category_name, group_id) values('test_category', 1);
insert into items(item_name, category_id) values('test_item', 1);
insert into items(item_name, category_id) values('test_item', 1);
insert into items(item_name, category_id) values('test_item', 1);

SELECT categories.category_name, categories.category_id FROM categories, users
WHERE users.user_id=2 AND categories.group_id=users.group_id;

SELECT items.item_name, items.category_id FROM items, categories, users
WHERE users.user_id=2 AND categories.group_id=users.group_id AND items.category_id=categories.category_id;

INSERT INTO items(item_name, category_id)
VALUES('a', (SELECT category_id FROM categories WHERE categories.category_name='Что-то'));

DELETE FROM items WHERE item_name='a' AND category_id=(SELECT category_id FROM categories WHERE category_name='Что-то');

INSERT INTO categories(category_name, group_id) VALUES('категория', (SELECT group_id FROM users WHERE user_id=2));

DELETE FROM categories WHERE category_name='категория' AND group_id=(SELECT group_id FROM users WHERE user_id=2);

SELECT groups.group_name FROM groups_users, groups WHERE groups_users.group_id=groups.group_id AND groups_users.user_id=1;
SELECT groups.group_name, owned_groups.link FROM owned_groups, groups WHERE owned_groups.group_id=groups.group_id AND owned_groups.user_id=1;

UPDATE users SET group_id=(SELECT group_id FROM groups WHERE invite_link=$1) WHERE user_id=$2;
SELECT groups_users.group_id FROM groups_users WHERE groups_users.group_id=(SELECT group_id FROM groups WHERE invite_link=$1);
UPDATE users SET group_id=(SELECT groups_users.group_id FROM groups_users WHERE groups_users.group_id=(SELECT group_id FROM groups WHERE select_link='test')) WHERE user_id=1;

UPDATE users SET group_id=(SELECT owned_groups.group_id FROM owned_groups WHERE owned_groups.link='qqqqqq') WHERE user_id=1;

SELECT * FROM owned_groups WHERE user_id=1 AND group_id=4;