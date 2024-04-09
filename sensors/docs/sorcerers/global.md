# Global Configuration Options

To install the services, use DataSorcerers/build.sh. This will create
    a series of systemd services on your computer, as well as a series
    of binaries located at /usr/local/bin/discern-*

Use DataSorcerers/build-config.sh to enable or disable services and 
    any optional features (separate from the live config options present
    in /etc/discern/SorcererConfig.yaml)

All services have the format discern-<sorcerername>-sorcerer, so you
    can grep for these services quite easily

Every service uses the same config file for live parameter updates, 
    located at /etc/discern/SorcererConfig.yaml

An example config can be found at DataSorcerers/SorcererConfig.yaml

Every services logs its information to the same file, specified by the
    config parameter `logfile`. Each program adds a timestamp and its
    name to the log submission, so they are easy to parse. This file 
    defaults to the same file the FusionCore server writes to.


### All of these options need to match with the FusionCore server

grpc_port: An integer representing the port number of the Core server, 
    where all the services should send their data

    - default: 50051

grpc_ip: A string which holds the URL / IP for the Core server

    - default: localhost

grpc_proto: A string representing the protocol format used to send data to FusionCore

    - default: tcp

maxrecvmsgsize: An int used to set the MaxRecvMsgSize field in the gRPC connections

    - default : 10485760


## TLS Setings

tls: A boolean, used to determine if TLS is used in the connection
    
    - default: false

certfile: A string representing the location of the certification file 
    for the TLS connection

    - default: ""

keyfile: A string representing the location of the key file for the 
    TLS connection

    - default: ""


## Identification Settings

In order to uniquely identify all of your devices, each service will
    send a POST request to the config parameter IdApiUrl.
    - Default value for IdApiUrl: http://localhost:50052/api/id

The the response to this request should be a text/plain which contains
    whatever unique string you'd like to use as an identifier. This 
    string will be saved in the database to unique identify hosts in
    the data records. Any other attempt at implementation can be spoofed
    in some way.

To make life easier for other possible implementations, a body can be 
    passed into the request via the IdApiBody parameter.
    - Default value for IdApiBody: ""

An example implementation of the identification server can be found at 
    metadata/id/server. This server simply parses /proc/cmdline for 
    a parameter called "inframac" and uses that as a unique identifier

