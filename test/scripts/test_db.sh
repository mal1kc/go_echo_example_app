#!/bin/env sh

docker_cli="podman"

test_postgres_user="postgres"
test_postgres_pass="postgres"
test_postgres_dbname="productapp"
test_max_conn="10"
test__max_conn_idle_time="30s"

$docker_cli run --name postgres-test --rm -e POSTGRES_USER="$test_postgres_user" -e POSTGRES_PASSWORD="test_postgres_pass" -p 6432:5432 -d docker.io/postgres:latest

echo "postgressql test db starting ..."
sleep 3

$docker_cli exec -it postgres-test psql -U postgres -d postgres -c "CREATE DATABASE $test_postgres_dbname"
sleep 3
echo "db productapp created"

$docker_cli exec -it postgres-test psql -U postgres -d productapp -c "
create table if not exists products
(
  id bigserial not null primary key,
  name varchar(255) not null,
  price double precision not null,
  discount double precision,
  store varchar(255) not null
);
"

sleep 3

echo "Table products created"
ConnString="host=localhost port=6432 user=$test_postgres_user password=$test_postgres_pass dbname=$test_postgres_dbname sslmode=disable statement_cache_mode=describe pool_max_conns=$test_max_conn pool_max_conn_idle_time=$test__max_conn_idle_time"
export APROD_APP_PQSQL_CONN_STR="$ConnString"
echo "env var set to db created container's data or use this connection string: APROD_APP_PQSQL_CONN_STR=$ConnString"
