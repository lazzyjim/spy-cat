CREATE TABLE spy_cats
(
    id                  serial primary key,
    name                varchar(50) unique NOT NULL,
    years_of_experience int                NOT NULL,
    breed               varchar(50)        NOT NULL,
    breed_validation    boolean default false,
    salary              float              NOT NULL
);

CREATE TABLE missions
(
    id             serial primary key,
    name           varchar(50) unique NOT NULL,
    cat_id         int,
    complete_state boolean default false
);

CREATE TABLE targets
(
    id             serial primary key,
    mission_id     int          NOT NULL,
    name           varchar(50)  NOT NULL,
    country        varchar(50)  NOT NULL,
    notes          varchar(100) NOT NULL,
    complete_state boolean default false
);


