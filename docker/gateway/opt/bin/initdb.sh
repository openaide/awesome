#!/bin/bash
set -xe

function create_su() {
	local user=$1

	echo "Creating user $user"
	psql -v ON_ERROR_STOP=1 --username "postgres" <<-EOSQL
		CREATE USER $user;
		CREATE DATABASE $user;
		ALTER ROLE $user WITH SUPERUSER;
EOSQL
}

function create_db() {
	local db=$1
	local user=$2

	echo "Creating database $db by $user"
	psql -v ON_ERROR_STOP=1 --username "$user" <<-EOSQL
		CREATE USER $db;
		CREATE DATABASE $db;
		GRANT ALL PRIVILEGES ON DATABASE $db TO $db;
		GRANT ALL PRIVILEGES ON DATABASE $db TO $user;
EOSQL
}

create_su admin
create_db litellm admin

# FIXME
# export DATABASE_URL="postgresql://admin:password@localhost:5432/litellm"
# python /app/db_scripts/create_views.py
##
