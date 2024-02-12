create table messages
(
    message_id varchar(512) not null,
    header varchar(128),
    body varchar(512),
    created_at date,
    updated_at date,
    status varchar(128)
);

create table message_histories
(
    message_id varchar(512) not null,
    created_at date,
    from_status varchar(128),
    to_status varchar(128)
);