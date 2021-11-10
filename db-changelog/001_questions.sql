create table if not exists questions
(
    uuid             text not null
        constraint questions_pkey
            primary key,
    created_at       timestamp with time zone,
    updated_at       timestamp with time zone,
    text             text,
    difficulty_level smallint
);