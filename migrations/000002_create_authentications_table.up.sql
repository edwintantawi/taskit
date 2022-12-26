CREATE TABLE authentications (
  id          VARCHAR(64)   PRIMARY KEY,
  user_id     VARCHAR(255)  NOT NULL,
  token       TEXT          NOT NULL,
  expires_at  TIMESTAMP     NOT NULL,

  CONSTRAINT fk_authentications_users FOREIGN KEY(user_id) REFERENCES users(id)
);