CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username TEXT
        NOT NULL
        UNIQUE
        CONSTRAINT name_length CHECK (char_length(username) <= 255),
    password_hash TEXT
        NOT NULL
        CONSTRAINT password_hash_length CHECK (char_length(password_hash) <= 511),
    create_time TIMESTAMP
        NOT NULL,
    is_admin BOOLEAN DEFAULT('false')
        NOT NULL
   
);
CREATE TABLE IF NOT EXISTS feature (
    id BIGSERIAL PRIMARY KEY,
    feature_data TEXT

        CONSTRAINT feature_data_length CHECK (char_length(feature_data) <= 255)
    
);

CREATE TABLE IF NOT EXISTS banner (
      id BIGSERIAL PRIMARY KEY,
      feature_id BIGSERIAL REFERENCES feature (id)
	NOT NULL,
      title TEXT
        NOT NULL
        CONSTRAINT title_length CHECK (char_length(title) <= 255),
      banner_data TEXT
      CONSTRAINT banner_data_length CHECK (char_length(title) <= 3000),
      url TEXT
        NOT NULL
        CONSTRAINT banner_title_length CHECK (char_length(title) <= 255),
      is_active BOOLEAN DEFAULT ('true')
        NOT NULL,
      create_time TIMESTAMP
        NOT NULL,
      update_time TIMESTAMP
        NOT NULL
      
      
        

);

CREATE TABLE IF NOT EXISTS tag (
    id BIGSERIAL PRIMARY KEY,
    tag_data TEXT

        CONSTRAINT tag_data_length CHECK (char_length(tag_data) <= 255)

);
CREATE TABLE IF NOT EXISTS banner_tag (
    id BIGSERIAL PRIMARY KEY,
    tag_id BIGSERIAL REFERENCES tag (id),
    banner_id BIGSERIAL REFERENCES banner(id),
    feature_id BIGSERIAL REFERENCES feature(id),
	  UNIQUE (tag_id, feature_id)
);

--password: testuser--
INSERT INTO users(username, password_hash, create_time, is_admin) 
        VALUES ('testuser', 'ae5deb822e0d71992900471a7199d0d95b8e7c9d05c40a8245a281fd2c1d6684', CURRENT_TIMESTAMP, 'false'),
               ('testadmin', 'ae5deb822e0d71992900471a7199d0d95b8e7c9d05c40a8245a281fd2c1d6684', CURRENT_TIMESTAMP, 'true');

do $$
begin
for r in 1..20 loop
insert into tag(id) values(DEFAULT);
insert into feature(id) values (DEFAULT);
end loop;
end;
$$;

