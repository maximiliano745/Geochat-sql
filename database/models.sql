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

-- DROP TABLE if exists grupo_miembros cascade;
-- DROP TABLE if exists user_groups cascade;
--  DROP TABLE IF EXISTS user_groups;

CREATE TABLE IF NOT EXISTS user_groups (
-- CREATE TABLE user_groups (
    id serial NOT NULL,
    iddueño INT NOT NULL,         -- ID del dueño del grupo
    group_name VARCHAR(150) NOT NULL,
    created_at timestamp DEFAULT now(),
    updated_at timestamp NOT NULL,
    CONSTRAINT pk_user_groups PRIMARY KEY(id),
    FOREIGN KEY (iddueño) REFERENCES users(id)  -- Restricción de clave externa
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

-- Crear la tabla para la relación entre grupos y miembros
CREATE TABLE IF NOT EXISTS grupo_miembros (
    id_grupo serial NOT NULL,
    id_miembro serial NOT NULL,
    CONSTRAINT pk_grupo_miembros PRIMARY KEY (id_grupo, id_miembro),
    FOREIGN KEY (id_grupo) REFERENCES user_groups(id),
    FOREIGN KEY (id_miembro) REFERENCES users(id)
);


-- CREATE TABLE IF NOT EXISTS ROLES (
--     id SERIAL PRIMARY KEY,
--     name VARCHAR(255) NOT NULL
-- );

CREATE TABLE IF NOT EXISTS ROLES (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL
);


CREATE TABLE IF NOT EXISTS ROLE_USER (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    role_id INT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (role_id) REFERENCES roles(id)
);

--ALTER TABLE ROLES ADD CONSTRAINT unique_role_name UNIQUE (name);
-- ALTER TABLE ROLES DROP CONSTRAINT unique_role_name;

INSERT INTO ROLES (name) VALUES ('maxi') ON CONFLICT DO NOTHING;
INSERT INTO ROLES (name) VALUES ('dueño_grupo') ON CONFLICT DO NOTHING;
INSERT INTO ROLES (name) VALUES ('usuario_grupo') ON CONFLICT DO NOTHING;
