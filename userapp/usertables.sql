DROP TABLE IF EXISTS users;
CREATE TABLE users (
  id         INT AUTO_INCREMENT NOT NULL,
  username VARCHAR(128) NOT NULL,
  PRIMARY KEY (`id`)
);

DROP TABLE IF EXISTS todos;
CREATE TABLE todos (
  id         INT AUTO_INCREMENT NOT NULL,
  user_id    INT NOT NULL,
  title VARCHAR(128) NOT NULL,
  curr_status VARCHAR(128) NOT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
INSERT INTO users (username) VALUES ('John');
INSERT INTO users (username) VALUES ('Jane');
INSERT INTO users (username) VALUES ('David');

INSERT INTO todos (user_id, title, curr_status) VALUES (1, 'Buy groceries', 'Not done');
INSERT INTO todos (user_id, title, curr_status) VALUES (1, 'Finish homework', 'In progress');
INSERT INTO todos (user_id, title, curr_status) VALUES (2, 'Buy clothes', 'Not done');
INSERT INTO todos (user_id, title, curr_status) VALUES (3, 'Finish work', 'Done');
