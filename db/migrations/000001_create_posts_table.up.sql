create table if not exists posts (
    id uuid primary key default gen_random_uuid(),
    title varchar(100) not null,
    content text,
    author varchar(100) not null,
    is_comments_allowed boolean default true,
    created_at timestamp default now()
);