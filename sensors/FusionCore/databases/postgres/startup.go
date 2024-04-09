package postgres;

import (

    "fmt"
    "context"
    "os"
    "os/exec"
    "strings"
    "time"
    
    "github.com/jackc/pgx/v5/pgxpool"

    Log "FusionCore/log"
    "strconv"

    "github.com/docker/docker/errdefs"
    "github.com/docker/docker/client"

    config "FusionCore/config"

)

var Connection *pgxpool.Pool;
var err error;

func StartDB() {

    if !config.Settings.StartNewPostgresInstance { return; }

    old_dir, err := os.Getwd()
    if err != nil {
        out := fmt.Sprintf("Error changing working directories in postgres docker build: %v", err)
        Log.FatalError(out)
    }
    os.Chdir("/etc/discern/postgres")
    defer os.Chdir(old_dir)

    Log.LogInfo("Stopping old postgres container")
    // Stop the old containers
    args := strings.Fields("stop discernpsql")
    cmd := exec.Command("docker", args...)
    err = cmd.Run()
    // Remove it 
    args = strings.Fields("rm discernpsql")
    cmd = exec.Command("docker", args...)
    err = cmd.Run()

    // Remove old volume
    args = strings.Fields("volume rm discernpsql")
    cmd = exec.Command("docker", args...)
    err = cmd.Run()
    // Create new volume
    args = strings.Fields("volume create discernpsql")
    cmd = exec.Command("docker", args...)
    err = cmd.Run()


    /* 
    THIS WAS MOVED TO BUILD SCRIPT
    Log.LogInfo("Building new image")
    // Start up the image
    args = strings.Fields("buildx build -t discernpsql .")
    cmd = exec.Command("docker", args...)
    err = cmd.Start()
    if err != nil {
        Log.FatalError(fmt.Sprintf("Error in starting the docker container: ", err))
    }
    */


    Log.LogInfo("Starting up postgres docker container")
    // Start up the image
    str_port := strconv.FormatInt(config.Settings.ExternalPostgresPort, 10)
    args = strings.Fields("run -dit -p 5432:" + str_port + " -v discernpsql:" + 
        config.Settings.PostgresDataMount + " --name discernpsql discernpsql")
    cmd = exec.Command("docker", args...)
    err = cmd.Start()
    if err != nil {
        Log.FatalError(fmt.Sprintf("Error in starting the postgres docker container: ", err))
    }

    // Create a new client to figure out when the container starts
    cli, err := client.NewEnvClient()
    if err != nil {
        Log.FatalError(fmt.Sprintf("Error creating env client for postgres: %v", err))
    }

    // Define the container name or ID you want to check
    containerName := "discernpsql"
    containerRunning := false;

    // Busy wait for container to start
    for !containerRunning {
        // Use the Docker client to inspect the container
        containerInfo, err := cli.ContainerInspect(context.Background(), 
                                    containerName)

        time.Sleep(1 * time.Second)
        Log.LogInfo("Waiting on postgres container")
        if err != nil {
            // If the container is not found, an error is returned
            if _, ok := err.(errdefs.ErrNotFound); !ok {
                Log.FatalError(fmt.Sprintf("%v", err))
            }
        } else {
            // Check the container's status
            if containerInfo.State.Running {
                containerRunning = true;
            }
        } 
    }
}

func Connect() {
    // Kept as default
    connStr := fmt.Sprintf("postgresql://%v:%v@%v/Discern?sslmode=disable", 
        config.Settings.PostgresUser, config.Settings.PostgresPassword, 
        config.Settings.PostgresIP)

    Connection, err = pgxpool.New(context.Background(), connStr)
    if err != nil {
        Log.FatalError(fmt.Sprintf("Unable to connect to database: %v\n", err))
    }
    Log.LogInfo("Postgres Connection Established")
    /*
    Connection.Query(`
        DO $$ 
        BEGIN
            -- Attempt to insert the values
            INSERT INTO Logs
            VALUES ('John', 'Contents for John'); -- Replace with your desired name and contents
            
        EXCEPTION WHEN unique_name_contents THEN
            -- If a unique constraint violation occurs, retrieve the ID of the conflicting entry
            RAISE NOTICE 'Unique constraint violated. ID of conflicting entry: %',
                (SELECT id FROM unique_pairs WHERE name = 'John' AND contents = 'Contents for John');
        END $$;
        `)
    */
}

