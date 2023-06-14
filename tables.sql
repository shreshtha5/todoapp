DROP TABLE IF EXISTS todoapp;
CREATE TABLE todoapp (
  id         INT AUTO_INCREMENT NOT NULL,
  title VARCHAR(128) NOT NULL,
  curr_status VARCHAR(128) NOT NULL,
  PRIMARY KEY (`id`)
);

INSERT INTO todoapp
  (title, curr_status)
VALUES
  ('Return Booking', 'Not Done'),
  ('GSleave', 'Done'),
  ('Fruits', 'Not Done');
  