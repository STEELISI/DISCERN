package main;

import (
    "fmt"
    "net"
    "context"
    "google.golang.org/grpc" 
    "github.com/influxdata/influxdb-client-go/v2"

    Log "FusionCore/log"


    config "FusionCore/config"
    s_influx "FusionCore/databases/influx"
    s_postgres "FusionCore/databases/postgres"

    b_validation "FusionBridge/validation"
    s_validation "FusionCore/validation"

    b_jupyter "FusionBridge/control/jupyter"
    s_jupyter "FusionCore/ingest/control/jupyter"

    b_id "FusionBridge/metadata/id"
    s_id "FusionCore/ingest/metadata/id"

    b_log "FusionBridge/log"
    s_log "FusionCore/ingest/log"

    b_ansible "FusionBridge/control/ansible"
    s_ansible "FusionCore/ingest/control/ansible"

    b_bash "FusionBridge/control/bash"
    s_bash "FusionCore/ingest/control/bash"

    b_file_meta "FusionBridge/metadata/file"
    s_file_meta "FusionCore/ingest/metadata/file"

    b_network_meta "FusionBridge/metadata/network"
    s_network_meta "FusionCore/ingest/metadata/network"

    b_os_meta "FusionBridge/metadata/os"
    s_os_meta "FusionCore/ingest/metadata/os"
)

func main() {

    config.LoadConfig()

    s_influx.StartDB();
    s_postgres.StartDB();

    Log.LogInfo("Server Booting");

    // Setting up all the Server listening stuff
    lis, err := net.Listen(config.ListenProto, config.ListenIP+":"+config.GRPC_Port);
    if err != nil {
        Log.FatalError(fmt.Sprintf("listen failed with error: %v", err))
    }
    var opts = []grpc.ServerOption{
        grpc.MaxRecvMsgSize(10 * 1024 * 1024),
    }
    grpcServer := grpc.NewServer(opts...)


    // Setting up the database connection (bound to services through
    //     the NewServer handler)
    client := influxdb2.NewClient(config.FullAddr(), config.DBtoken);
    _, err = client.Health(context.Background())
    if err != nil {
        Log.FatalError(fmt.Sprintf("Heath of client connection bad: %v", err))
    }

    s_postgres.Connect()


    // REGISTER ALL OF YOUR SERVICES HERE

    // Basic check to see that the server is generally working 
    b_validation.RegisterSendAndRecvServer(grpcServer, 
        s_validation.NewServer())
    // End points to recieve all the juypter related information
    b_jupyter.RegisterJupyterServer(grpcServer, 
        s_jupyter.NewServer(client))
    // Register the Bash server
    b_bash.RegisterBashServer(grpcServer, 
        s_bash.NewServer(client))
    // Register the Ansible server
    b_ansible.RegisterAnsibleServer(grpcServer, 
        s_ansible.NewServer(client))
    // Register the Log server
    b_log.RegisterLogServer(grpcServer, 
        s_log.NewServer(client))
    // Register the File Metadata Server
    b_file_meta.RegisterFileServer(grpcServer, 
        s_file_meta.NewServer(client))
    // Register the Network metadata server
    b_network_meta.RegisterNetworkServer(grpcServer, 
        s_network_meta.NewServer(client))
    // Register the OS metadata server
    b_os_meta.RegisterOSServer(grpcServer, 
        s_os_meta.NewServer(client))
    // Register the ID metadata server
    b_id.RegisterIDServer(grpcServer, 
        s_id.NewServer(client))

    Log.LogInfo("Finished Registering All Services");

    Log.LogInfo("Running Serving")
    grpcServer.Serve(lis)
}

