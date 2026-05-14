
create table sk (
    id bigserial primary key not null,
    name varchar(128) not null,
    is_active boolean not null,
    created_at timestamp with time zone default now() not null,
    updated_at timestamp with time zone default now() not null
);

insert into sk (id, name, is_active, created_at, updated_at)
values
    (1, 'Альфастрахование', true, now(), now()),
    (3, 'Ингосстрах', true, now(), now()),
    (5, 'Росгосстрах', true, now(), now()),
    (7, 'Согласие', true, now(), now()),
    (10, 'МАКС', false, now(), now()),
    (25, 'Югория', false, now(), now()),
    (27, 'Зетта Страхование', false, now(), now()),
    (32, 'РЕСО', false, now(), now()),
    (33, 'ВСК', false, now(), now()),
    (36, 'Ренессанс Страхование', false, now(), now()),
    (38, 'Совкомбанк Страхование', false, now(), now()),
    (39, 'Интач', false, now(), now()),
    (46, 'СОГАЗ', false, now(), now()),
    (47, 'ГАЙДЕ', false, now(), now()),
    (52, 'Гелиос', false, now(), now()),
    (61, 'ОСК', false, now(), now()),
    (72, 'АСКО', false, now(), now()),
    (75, 'Абсолют', false, now(), now()),
    (98, 'Астро-Волга', false, now(), now()),
    (107, 'Т-Страхование', false, now(), now()),
    (125, 'ВЕРНА', false, now(), now()),
    (127, 'ЕВРОИНС', false, now(), now()),
    (142, 'Мафин', false, now(), now()),
    (144, 'Сбербанк cтрахование', false, now(), now()),
    (145, 'Энергогарант', false, now(), now()),
    (200, 'Инсапп', false, now(), now())
;
