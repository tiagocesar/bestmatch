drop database if exists bestmatch;

create database bestmatch;

\c bestmatch;

-- Enabling distance calculations
create extension cube;
create extension earthdistance;

create table if not exists materials
(
    id   uuid
        primary key,
    name varchar(100) not null
);

alter table materials
    owner to root;

create table if not exists partners
(
    id      uuid
        primary key,
    name    varchar(100) not null,
    address point        not null,
    radius  integer      not null,
    rating  real         not null
);

alter table partners
    owner to root;

create table if not exists partners_materials
(
    partner_id  uuid
        constraint partners_materials_partners_id_fk
            references partners,
    material_id uuid
        constraint partners_materials_materials_id_fk
            references materials,
    constraint partners_materials_pk
        unique (partner_id, material_id)
);

alter table partners_materials
    owner to root;

-- Sample data
INSERT INTO materials (id, name)
VALUES ('07cab731-d981-4915-9444-cc997eec351f', 'wood'),
       ('1606f175-3502-4028-9501-6b591c00f1f3', 'carpet'),
       ('ac47d822-ffc9-48b7-8492-4d49e921d4df', 'tiles'),
       ('e802e6d4-9559-4218-b17b-2c199907e867', 'marble'),
       ('95c74044-953d-4a4e-b226-b10b7553c525', 'granite');

INSERT INTO partners (id, name, address, radius, rating)
VALUES ('b276cb54-ac52-4f8c-adb1-afce5ced67c4', 'Acme Inc.',      point(4.8986299, 52.3706706), 10, 5.0),
       ('6360d1e7-ccd0-43d0-8bf5-d7bc807213d3', 'De Twee Broers', point(4.8964412, 52.3706706), 10, 4.3),
       ('785bfeb5-b982-4c56-88ca-bbd7a49cee4c', 'Really far',     point(5.073058, 52.421221),   12, 4.5);

INSERT INTO partners_materials (partner_id, material_id)
VALUES ('b276cb54-ac52-4f8c-adb1-afce5ced67c4', '07cab731-d981-4915-9444-cc997eec351f'),
       ('b276cb54-ac52-4f8c-adb1-afce5ced67c4', '1606f175-3502-4028-9501-6b591c00f1f3'),
       ('b276cb54-ac52-4f8c-adb1-afce5ced67c4', 'ac47d822-ffc9-48b7-8492-4d49e921d4df'),
       ('b276cb54-ac52-4f8c-adb1-afce5ced67c4', 'e802e6d4-9559-4218-b17b-2c199907e867'),
       ('b276cb54-ac52-4f8c-adb1-afce5ced67c4', '95c74044-953d-4a4e-b226-b10b7553c525'),
       ('6360d1e7-ccd0-43d0-8bf5-d7bc807213d3', '07cab731-d981-4915-9444-cc997eec351f'),
       ('6360d1e7-ccd0-43d0-8bf5-d7bc807213d3', '1606f175-3502-4028-9501-6b591c00f1f3'),
       ('6360d1e7-ccd0-43d0-8bf5-d7bc807213d3', 'ac47d822-ffc9-48b7-8492-4d49e921d4df'),
       ('785bfeb5-b982-4c56-88ca-bbd7a49cee4c', '07cab731-d981-4915-9444-cc997eec351f'),
       ('785bfeb5-b982-4c56-88ca-bbd7a49cee4c', '1606f175-3502-4028-9501-6b591c00f1f3'),
       ('785bfeb5-b982-4c56-88ca-bbd7a49cee4c', 'ac47d822-ffc9-48b7-8492-4d49e921d4df');
