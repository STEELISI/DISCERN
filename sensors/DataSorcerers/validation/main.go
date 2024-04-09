package main;

import (
    "fmt"
    "context"
    "time"
    "DataSorcerers/helpers"

    FB "FusionBridge/validation"
)


func main() {

    helpers.LoadConfig()


    conn := helpers.CreateConnection()
    client := FB.NewSendAndRecvClient(conn); 

    helpers.SetFileName("Validation")

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    to_send := FB.HelloWorld{ Tmp : 12345 }
    msg, err := client.HelloWorldBasic(ctx, &to_send)
    if err != nil {
        out := fmt.Sprintf("Error running validation: %v\n", err)
        helpers.LogInfo(out)
        return
    }
    out := fmt.Sprintf("Client got response: %v\n", msg.Tmp);
    helpers.LogInfo(out)
}
