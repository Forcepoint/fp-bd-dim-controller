create table IF NOT EXISTS element_batches
(
    id         bigint unsigned auto_increment
        primary key,
    created_at datetime(3) null,
    updated_at datetime(3) null,
    deleted_at datetime(3) null
);

create table IF NOT EXISTS list_elements
(
    id              bigint unsigned auto_increment
        primary key,
    created_at      datetime(3)     null,
    updated_at      datetime(3)     null,
    deleted_at      datetime(3)     null,
    source          longtext        null,
    service_name    varchar(191)    null,
    type            varchar(191)    null,
    value           varchar(500)    not null,
    safe            tinyint(1)      null,
    update_batch_id bigint unsigned null,
    constraint value
        unique (value)
);

create index IF NOT EXISTS batchid
    on list_elements (update_batch_id);

create index IF NOT EXISTS idx_list_elements_value
    on list_elements (value);

create index IF NOT EXISTS svcname
    on list_elements (service_name);

create index IF NOT EXISTS type
    on list_elements (type);

create table IF NOT EXISTS log_entries
(
    id          bigint unsigned auto_increment
        primary key,
    created_at  datetime(3) null,
    updated_at  datetime(3) null,
    deleted_at  datetime(3) null,
    module_name longtext    null,
    level       longtext    null,
    message     longtext    null,
    caller      longtext    null,
    time        datetime(3) null
);

create table IF NOT EXISTS module_endpoints
(
    id                 bigint unsigned auto_increment
        primary key,
    created_at         datetime(3)     null,
    updated_at         datetime(3)     null,
    deleted_at         datetime(3)     null,
    secure             tinyint(1)      null,
    endpoint           longtext        null,
    module_metadata_id bigint unsigned null
);

create table IF NOT EXISTS module_metadata
(
    id                  bigint unsigned auto_increment
        primary key,
    created_at          datetime(3)  null,
    updated_at          datetime(3)  null,
    deleted_at          datetime(3)  null,
    module_service_name varchar(191) null,
    module_display_name longtext     null,
    module_type         varchar(191) null,
    module_description  text         null,
    inbound_route       varchar(191) null,
    internal_ip         longtext     null,
    internal_port       longtext     null,
    icon_url            longtext     null,
    configured          tinyint(1)   null,
    configurable        tinyint(1)   null,
    last_ping           datetime(3)  null,
    constraint inbound_route
        unique (inbound_route)
);

create index IF NOT EXISTS sname
    on module_metadata (module_service_name);

create index IF NOT EXISTS type
    on module_metadata (module_type);

create table IF NOT EXISTS update_statuses
(
    id                 bigint unsigned auto_increment
        primary key,
    created_at         datetime(3)     null,
    updated_at         datetime(3)     null,
    deleted_at         datetime(3)     null,
    service_name       longtext        null,
    status             longtext        null,
    update_batch_id    bigint unsigned null,
    module_metadata_id bigint unsigned null
);

create table IF NOT EXISTS users
(
    id bigint unsigned auto_increment
        primary key,
    created_at datetime(3) null,
    updated_at datetime(3) null,
    deleted_at datetime(3) null,
    name longtext null,
    email varchar(100) null,
    password longtext null,
    admin tinyint(1) default 0 not null,
    constraint value
        unique (email)
);

