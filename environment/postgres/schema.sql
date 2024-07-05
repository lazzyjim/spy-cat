CREATE TABLE spy_cats
(
    id                  serial primary key,
    name                varchar(50) NOT NULL,
    years_of_experience int         NOT NULL,
    breed               varchar(50) NOT NULL,
    breed_validation    boolean default false,
    salary              float       NOT NULL
);

CREATE TABLE missions
(
    id             serial primary key,
    cat_id         int references spy_cats (id),
    complete_state boolean default false
);

CREATE TABLE targets
(
    id             serial primary key,
    mission_id     int references missions (id) NOT NULL,
    name           varchar(50)                  NOT NULL,
    country        varchar(50)                  NOT NULL,
    notes          varchar(100)                 NOT NULL,
    complete_state boolean default false
);


