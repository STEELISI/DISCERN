FROM postgres

# Default user is postgres

ENV POSTGRES_PASSWORD=*Something*
# ENV POSTGRES_PASSWORD_FILE=./password

# Five the db an actual name
ENV POSTGRES_DB=Discern

# Copy all .sql files you want to initialize into this folder so it will 
# create all the tables on startup
# THE FILES ARE RUN IN NAME SORTED ORDER. BIG DEAL
RUN mkdir -p /docker-entrypoint-initdb.d

# Need to actually import this from the build order stuff and do it that way. 
# Works preliminarily
RUN mkdir -p /schema

COPY ./schema/ /docker-entrypoint-initdb.d

RUN echo hello world

# To plainl run this, will exit
# sudo docker run --name some-postgres -e POSTGRES_PASSWORD=mysecretpassword -d postgres testing:v1
# sudo docker run -p 5432:5432 testing:v1 -d postgres
# To build this container 
# sudo docker buildx build -t testing:v1 .

