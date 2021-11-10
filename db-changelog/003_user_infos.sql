create table if not exists user_infos
(
    uuid       text not null
        constraint user_infos_pkey
            primary key,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    name       text,
    phone      text,
    email      text,
    link       text,
    question   text
);
