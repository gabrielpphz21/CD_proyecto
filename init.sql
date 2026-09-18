CREATE table tenants(
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(128),
    tax_id VARCHAR(32)
);

CREATE table satelites(
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    tenant INTEGER REFERENCES tenants(id)
);

CREATE table rockets(
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    satelite INTEGER REFERENCES satelites(id)
);

CREATE table take_offs(
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    rocket INTEGER REFERENCES rockets(id),
    date TIMESTAMP,
    runway INTEGER
);

CREATE TABLE destinations(
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    planet VARCHAR(32),
    state VARCHAR(32)
);

CREATE TABLE energy_transactions(
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    volume FLOAT,
    origin INTEGER REFERENCES satelites(id),
    destination INTEGER REFERENCES destinations(id),
    state VARCHAR(128),
    date TIMESTAMP
);

CREATE TABLE hyper_satelite(
    id INTEGER REFERENCES satelites(id),
    energy INTEGER,
    parts JSONB,
    state VARCHAR(128),
    x double precision ,
    y double precision,
    z double precision,
    date TIMESTAMP
);

CREATE TABLE hyper_rocket(
    id INTEGER REFERENCES rockets(id),
    fuel JSONB,
    parts JSONB,
    state VARCHAR(128),
    x double precision ,
    y double precision,
    z double precision,
    date TIMESTAMP
);

