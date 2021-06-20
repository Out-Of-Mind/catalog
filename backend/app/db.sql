DROP DATABASE IF EXISTS catalog_db;
CREATE DATABASE catalog_db;

CREATE USER catalog_user WITH PASSWORD 'password';
ALTER ROLE catalog_user SET client_encoding TO 'utf8';
ALTER ROLE catalog_user SET default_transaction_isolation TO 'read committed';
ALTER ROLE catalog_user SET timezone TO 'UTC';
GRANT ALL PRIVILEGES ON DATABASE catalog_db TO catalog_user;

# \c catalog_db catalog_user

DROP TABLE IF EXISTS groups;
CREATE TABLE groups(
   group_id INT GENERATED ALWAYS AS IDENTITY,
   group_name VARCHAR(255) NOT NULL,
   PRIMARY KEY(group_id)
);

DROP TABLE IF EXISTS users;
CREATE TABLE users(
   user_id INT GENERATED ALWAYS AS IDENTITY,
   user_name VARCHAR(255) NOT NULL,
   group_id INT,
   PRIMARY KEY(user_id),
   CONSTRAINT fk_group
      FOREIGN KEY(group_id) 
	  REFERENCES groups(group_id)
	  ON DELETE SET NULL
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
   PRIMARY KEY(item_id),
   CONSTRAINT fk_category
      FOREIGN KEY(category_id) 
	  REFERENCES categories(category_id)
	  ON DELETE CASCADE
);

SELECT category_name, item_name FROM categories, items, users WHERE users.user_id=1 AND categories.group_id=users.group_id AND items.category_id=categories.category_id;