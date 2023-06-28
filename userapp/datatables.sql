DROP TABLE IF EXISTS newusers;
CREATE TABLE newusers (
  id         INT AUTO_INCREMENT NOT NULL,
  username VARCHAR(128) NOT NULL,
  userpass VARCHAR(128) NOT NULL,
  PRIMARY KEY (`id`)
);

DROP TABLE IF EXISTS newtodos;
CREATE TABLE newtodos (
  id         INT AUTO_INCREMENT NOT NULL,
  user_id    INT NOT NULL,
  title VARCHAR(128) NOT NULL,
  curr_status VARCHAR(128) NOT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (user_id) REFERENCES newusers(id) ON DELETE CASCADE
);
INSERT INTO newusers (username, userpass) VALUES ('John','123');
INSERT INTO newusers (username, userpass) VALUES ('Jane', 'app123');
INSERT INTO newusers (username, userpass) VALUES ('David','ten5');

INSERT INTO newtodos (user_id, title, curr_status) VALUES (1, 'Buy groceries', 'Not done');
INSERT INTO newtodos (user_id, title, curr_status) VALUES (1, 'Finish homework', 'In progress');
INSERT INTO newtodos (user_id, title, curr_status) VALUES (2, 'Buy clothes', 'Not done');
INSERT INTO newtodos (user_id, title, curr_status) VALUES (3, 'Finish work', 'Done');
