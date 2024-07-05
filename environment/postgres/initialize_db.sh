#!/bin/sh

echo "DROP DATABASE IF EXISTS ${POSTGRES_DB}; ->  psql -d postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgresql.loc/postgres"
echo "DROP DATABASE IF EXISTS ${POSTGRES_DB};" | psql -d postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgresql.loc/postgres
echo "CREATE DATABASE ${POSTGRES_DB}; ->  psql -d postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgresql.loc/postgres"
echo "CREATE DATABASE ${POSTGRES_DB};" | psql -d postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgresql.loc/postgres

psql -d postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgresql.loc/${POSTGRES_DB} -f /schema.sql