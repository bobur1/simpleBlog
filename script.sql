create table users
(
  id       serial       not null
    constraint users_pkey
    primary key,
  email    varchar(200) not null,
  password varchar(200) not null
);

alter table users
  owner to postgres;

create unique index users_id_uindex
  on users (id);

create unique index users_email_uindex
  on users (email);

create table articles
(
  id      serial       not null
    constraint articles_pkey
    primary key,
  title   varchar(200) not null,
  context text,
  user_id integer
    constraint articles_users_id_fk
    references users
    on delete cascade
);

alter table articles
  owner to postgres;

create unique index articles_id_uindex
  on articles (id);

create table comments
(
  id         serial not null
    constraint comments_pkey
    primary key,
  article_id integer
    constraint comments_articles_id_fk
    references articles
    on delete cascade,
  user_id    integer
    constraint comments_users_id_fk
    references users
    on delete cascade,
  parent_id  integer,
  context    text
);

alter table comments
  owner to postgres;

create unique index comments_id_uindex
  on comments (id);

create table uploads
(
  id         serial       not null
    constraint uploads_pkey
    primary key,
  article_id integer      not null
    constraint uploads_articles_id_fk
    references articles
    on delete cascade,
  path       varchar(200) not null
);

alter table uploads
  owner to postgres;

create unique index uploads_id_uindex
  on uploads (id);

create table vote_article
(
  mark       integer,
  id         serial  not null
    constraint vote_article_pk
    primary key,
  user_id    integer
    constraint vote_article_users_id_fk
    references users
    on delete cascade,
  article_id integer not null
    constraint vote_article_articles_id_fk
    references articles
    on delete cascade
);

alter table vote_article
  owner to postgres;

create unique index vote_article_id_uindex
  on vote_article (id);

