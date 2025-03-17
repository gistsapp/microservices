CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS user_entity (
    user_id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(255) NOT NULL,
    email VARCHAR(320) NOT NULL,
    picture TEXT
);

CREATE TABLE IF NOT EXISTS token(
  token_id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id uuid NOT NULL,
  token VARCHAR(255) NOT NULL,
  expires_at TIMESTAMP NOT NULL
);

ALTER TABLE user_entity ADD CONSTRAINT user_entity_username_key UNIQUE (username); --username is actually the user handle

ALTER TABLE token ADD CONSTRAINT token_user_id_fkey FOREIGN KEY (user_id) REFERENCES user_entity(user_id);

CREATE TABLE IF NOT EXISTS federated_identity(
  federated_identity_id VARCHAR(255) PRIMARY KEY, -- local is email, oidc is the provider id
  user_id uuid NOT NULL,
  provider VARCHAR(255) NOT NULL,
  data JSONB NOT NULL
);

ALTER TABLE federated_identity ADD CONSTRAINT federated_identity_user_id_fkey FOREIGN KEY (user_id) REFERENCES user_entity(user_id);

CREATE TABLE IF NOT EXISTS verification_token(
  token VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL
);

ALTER TABLE verification_token ADD CONSTRAINT verification_token_email_key UNIQUE (email);

ALTER TABLE verification_token ADD CONSTRAINT pkey PRIMARY KEY (token, email);
