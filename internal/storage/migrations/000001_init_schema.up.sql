create table users(
    id BIGSERIAL primary key, 
    username varchar(20) not null unique,
    password varchar(255) not null,
    email varchar(255) not null unique,
    status varchar(20) not null default 'active',
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp,
    last_login timestamp
);

create table newsletters(
    id BIGSERIAL primary key,
    title varchar(255) not null,
    author BIGINT not null references users(id),
    description text not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);

create table subscribers(
    id BIGSERIAL primary key,
    user_id BIGINT not null references users(id),
    newsletter_id BIGINT not null references newsletters(id),
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp,
    unique(user_id, newsletter_id)
);

create table letters(
    id BIGSERIAL primary key,
    newsletter_id BIGINT not null references newsletters(id),
    title varchar(255) not null,
    content text not null,
    status varchar(20) not null default 'draft',
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);

create table views(
    id BIGSERIAL primary key,
    letter_id BIGINT not null references letters(id),
    user_id BIGINT not null references users(id),
    ip_address inet,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);

-- Create indexes for better query performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_newsletters_author ON newsletters(author);
CREATE INDEX idx_subscribers_user ON subscribers(user_id);
CREATE INDEX idx_subscribers_newsletter ON subscribers(newsletter_id);
CREATE INDEX idx_letters_newsletter ON letters(newsletter_id);
CREATE INDEX idx_views_letter ON views(letter_id);
CREATE INDEX idx_views_user ON views(user_id);

