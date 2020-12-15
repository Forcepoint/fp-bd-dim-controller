create table IF NOT EXISTS module_element_type
(
    id bigint unsigned auto_increment
        primary key,
    created_at datetime(3) null,
    updated_at datetime(3) null,
    deleted_at datetime(3) null,
    element_type longtext null,
    module_id bigint unsigned null
);