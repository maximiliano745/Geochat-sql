--DROP TABLE IF EXISTS users;
CREATE TABLE IF NOT EXISTS users (
    id serial NOT NULL,
    username VARCHAR(150) NOT NULL ,
    pasword varchar(256) NOT NULL,
    email VARCHAR(150) NOT NULL UNIQUE,
    created_at timestamp DEFAULT now(),
    updated_at timestamp NOT NULL,
    hash varchar(256) ,
    CONSTRAINT pk_users PRIMARY KEY(id)
);
-- DROP TABLE IF EXISTS posts;
/* 
CREATE TABLE IF NOT EXISTS posts (
    id serial NOT NULL,
    user_id int NOT NULL,
    body text NOT NULL,
    created_at timestamp DEFAULT now(),
    updated_at timestamp NOT NULL,
    CONSTRAINT pk_notes PRIMARY KEY(id),
    CONSTRAINT fk_posts_users FOREIGN KEY(user_id) REFERENCES users(id)
);
 */