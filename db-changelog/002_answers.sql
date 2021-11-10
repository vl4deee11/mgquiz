create table if not exists answers
(
    uuid       text not null
        constraint answers_pkey
            primary key,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    text       text,
    question_uuid text
        constraint fk_questions_answers
            references questions,
    is_right   boolean
);