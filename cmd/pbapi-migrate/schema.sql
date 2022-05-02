pragma foreign_keys = on;

create table ssh_keys (
    id integer not null,
    body text not null,
    primary key (id)
);

create table ssh_key_resources (
    id integer not null,
    ssh_key_id integer not null,
    cid text not null,
    primary key (id),
    foreign key (ssh_key_id) references ssh_keys(id) on update cascade on delete cascade
);