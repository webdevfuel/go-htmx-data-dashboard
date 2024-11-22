CREATE TABLE users (
    id serial PRIMARY KEY,
    name text NOT NULL,
    email text NOT NULL,
    status varchar(50) NOT NULL,
    created_at timestamp NOT NULL DEFAULT now(),
    activated_at timestamp
);

CREATE UNIQUE INDEX users_email_idx ON users (email);

CREATE TABLE metrics (
    id serial PRIMARY KEY,
    metric_date date NOT NULL,
    new_users int NOT NULL DEFAULT 0,
    new_activations int NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX metrics_metric_date_idx ON metrics (metric_date);

