create table if not exists comments (
    id uuid primary key default gen_random_uuid(),
    author varchar(100),
    content text not null check(length(content) <= 2000),
    post_id uuid references posts(id) on delete cascade,
    parent_comment_id uuid references comments(id) on delete cascade,
    created_at timestamp default now()
);