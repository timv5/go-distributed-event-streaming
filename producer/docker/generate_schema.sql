create table messages
(
    id varchar(512) not null,
    header varchar(128),
    body varchar(512),
    created_at date,
    updated_at date,
    status varchar(128),
    PRIMARY KEY (id)
);