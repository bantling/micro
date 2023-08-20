# Postgres Docker container

# BUILD ARGS:
# ARG_NAME   (myapp) = the database to create
# ARG_LOCALE (en_US) = the locale inherited by all tables
# ARG_USER           = the user to connect as
# ARG_PASSWORD       = the password to connect with

FROM postgres
ARG ARG_NAME   = myapp
ARG ARG_LOCALE = en_US
ARG ARG_USER
ARG ARG_PASSWORD

# Initialise the database with improved security in hba.conf:
# type = host
# database = all
# user = all
# address = 127.0.0.1/32
# method = scram-sha-256
RUN su postgres -c 'initdb --auth=ident --encoding=UTF8 --locale=en_US --username=postgres --no-instructions'

RUN su postgres -c 'pg_ctl start --wait'

# Create the database
RUN su postgres -c 'createdb --username=postgres $ARG_DATABASE_NAME'
