-- DROP TABLE IF EXISTS users;
CREATE TABLE IF NOT EXISTS users (
-- CREATE TABLE users (
    id serial NOT NULL,
    username VARCHAR(150) NOT NULL ,
    pasword varchar(256) NOT NULL,
    email VARCHAR(150) NOT NULL UNIQUE,
    created_at timestamp DEFAULT now(),
    updated_at timestamp NOT NULL,
    hash varchar(256) ,
    CONSTRAINT pk_users PRIMARY KEY(id)
);

-- DROP TABLE IF EXISTS user_groups;
CREATE TABLE IF NOT EXISTS user_groups (
-- CREATE TABLE user_groups (
    id serial NOT NULL,
    group_name VARCHAR(150) NOT NULL,
    description TEXT,
    created_at timestamp DEFAULT now(),
    updated_at timestamp NOT NULL,
    CONSTRAINT pk_user_groups PRIMARY KEY(id)
);

-- DROP TABLE IF EXISTS user_group_membership;
CREATE TABLE IF NOT EXISTS user_group_membership (
-- CREATE TABLE  user_group_membership (
    user_id serial NOT NULL,
    group_id serial NOT NULL,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT fk_group FOREIGN KEY (group_id) REFERENCES user_groups (id),
    CONSTRAINT pk_user_group_membership PRIMARY KEY (user_id, group_id)
);

-- DROP TABLE IF EXISTS PedidosContactos;
CREATE TABLE IF NOT EXISTS PedidosContactos (
-- CREATE TABLE PedidosContactos (
    id serial NOT NULL,
    idusuarioofrece serial NOT NULL,
    idusuarioacepta serial NOT NULL,
    estado boolean DEFAULT false,
    CONSTRAINT pk_PedidosContactos PRIMARY KEY(id),
    CONSTRAINT fk_usuario_ofrece FOREIGN KEY (idusuarioofrece) REFERENCES users (id),
    CONSTRAINT fk_usuario_acepta FOREIGN KEY (idusuarioacepta) REFERENCES users (id)
);

