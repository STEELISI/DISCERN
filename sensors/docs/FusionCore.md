
# Fusion Core Architecture

This architecture of this server is intentionally quite basic:

- A large server to create listeners for and store all the data 
    from the gRPC services
    - The specifications of these services can be found at ./sorcerers. 
        Each service is activated independently on a client node and 
        none are required to be run

- An influxDB instance running in a docker container to store all the 
    data from the gRPC services

This architecture is controlled entirely with a config file located at
    /etc/discern/CoreConfig.yaml

The server itself logs all information to /var/log/discern.log

The database logs information inside the docker container using the
    default influx settings

This program also features a systemd service which one can easily install
    using sudo ./build.sh
    - If you don't need the system service, then go build . works as well


### This following sections show every option you can set in the yaml file. 
### If not set, the default value is used

## gRPC Connection Options

- logfile
    - default: "/var/log/discern.log"
    - A string path of the log file you'd like the server to store 
        any information about its operation

- grpc_port
    - default: 50051
    - An integer value which will set the port the gRPC listeners will
        all listen on.

- listenip
    - default: localhost
    - A string value which is used to set the IP the gRPC server will
        listen on

- listenproto
    - default: tcp
    - A string value which is used to set the connection type gRPC 
        listen for. Technically an argument to net.Listen in golang


## Docker & inFlux Options

- internaldockerport
    - default: 8086
    - An integer which sets which port influxDB listens to INSIDE the 
        docker container

- externalport
    - default: 8086
    - An integer which sets the port number outside the docker container
        which the influx db port is forwarded to

- dburl
    - default: 'http://localhost'
    - A string which contained to URL where the DB is hosted. This allows
        for the database to be hosted separately from the Core server, 
        without breaking anything

- dbtoken
    - default: 'BIGElHSa291FOkrliGaBVc7ksnGgQ4vALbkfJzRuH02T2XB8qouH0H3IkYTJACE-XZ-QYV664CH5655LkbQDIQ::'
    - A string which contains the database access token the Core sever
        needs to create a new client. If creating a new database is 
        specified, the token will be automatically added to the database

- startnewdockerinstance
    - default: false
    - A boolean representing whether the previous docker instance should
        be killed and re-created on boot. Helpful for debugging

- bucket_name
    - default: DISCERN
    - A string used to set the bucket name in the influx instance

- org
    - default: ISI
    - A string used to set the org name in the influx instance

- username
    - default: default-user
    - A string used to set the default user name set up when creating
        the influx instance. This parameter has no effect if you do
        not used the built-in docker instance

- password
    - default: something
    - A string used to set the default password set for the default 
        user when creating the influx instance. This parameter has no 
        effect if you do not used the built-in docker instance

