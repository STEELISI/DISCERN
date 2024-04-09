package influx;

import (

    "fmt"
    "context"
    "os/exec"
    "strings"
    "time"
    "bytes"
    Log "FusionCore/log"

    "github.com/docker/docker/errdefs"
    "github.com/docker/docker/client"

    config "FusionCore/config"

)

func StartDB() {

    if !config.StartNewDockerInstance { return; }

    Log.LogInfo("Stopping old influx container")
    // Stop the old containers
    args := strings.Fields("kill influx")
    cmd := exec.Command("docker", args...)
    err := cmd.Run()
    // Remove it 
    args = strings.Fields("rm influx")
    cmd = exec.Command("docker", args...)
    err = cmd.Run()

    Log.LogInfo("Starting up influx docker container")
    // Start up the image
    args = strings.Fields("run --name influx -p " + config.InternalDockerPort +
            ":" + config.ExternalPort + " influxdb")
    cmd = exec.Command("docker", args...)
    err = cmd.Start()
    if err != nil {
        Log.FatalError(fmt.Sprintf("Error in starting the influx docker container: ", err))
    }

    // Create a new client to figure out when the container starts
    cli, err := client.NewEnvClient()
    if err != nil {
        Log.FatalError(fmt.Sprintf("%v", err))
    }

    // Define the container name or ID you want to check
    containerName := "influx"
    containerRunning := false;

    // Busy wait for container to start
    for !containerRunning {
        // Use the Docker client to inspect the container
        containerInfo, err := cli.ContainerInspect(context.Background(), 
                                    containerName)

        time.Sleep(1 * time.Second)
        Log.LogInfo("Waiting on influxDB container")
        if err != nil {
            // If the container is not found, an error is returned
            if _, ok := err.(errdefs.ErrNotFound); !ok {
                Log.FatalError(fmt.Sprintf("Error creating env client for influxDB %v", err))
            }
        } else {
            // Check the container's status
            if containerInfo.State.Running {
                containerRunning = true;
            }
        } 
    }


    args = strings.Fields(fmt.Sprintf(`exec influx influx setup 
            --bucket %v 
            --org %v 
            --password %v 
            --username %v 
            --force 
            `, config.BUCKET_NAME, config.ORG, config.PASSWORD, 
                config.USERNAME))
    cmd = exec.Command("docker", args...)
    err = cmd.Run()


    args = strings.Fields(fmt.Sprintf(`exec influx influx auth list`))
    cmd = exec.Command("docker", args...)
    var outb, errb bytes.Buffer
    cmd.Stdout = &outb
    cmd.Stderr = &errb
    err = cmd.Run()
    if err != nil {
        Log.FatalError(fmt.Sprintf("Error from influxdb run: %v", err))
    }
    // This is going to bit rot so quickly but they have a garbage 
        // interface for this
    el := strings.Split(outb.String(), "\n")[1] 
    el2 := strings.Split(el, " ")[1]
    config.DBtoken = strings.Split(el2, "\t")[1]
    
}

