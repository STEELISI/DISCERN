package main;

import (
    "fmt"
    "io/fs"
    "syscall"
    "os"
    "os/user"
    "regexp"
    "path/filepath"
    "time"
    "io/ioutil"
    "context"
    "google.golang.org/grpc" 
    "github.com/golang/protobuf/ptypes"

    noti "github.com/fsnotify/fsnotify"

    bridge "FusionBridge/metadata/file"

    "DataSorcerers/helpers"
)


var opts []grpc.DialOption;
var conn *grpc.ClientConn;
var err error; 
var client bridge.FileClient;

var watcher *noti.Watcher;

func main() {

    helpers.SetFileName("File")

    // Load config
    helpers.LoadConfig()
    if !helpers.Config.RunFile { return }


    helpers.LogInfo("Hello World from metadata.file logger")

    conn = helpers.CreateConnection()
    client = bridge.NewFileClient(conn); 


    // Set up watchers & callbacks
    watcher, err = noti.NewWatcher()
    if err != nil {
        out := fmt.Sprintf("Error adding watcher: %v", err)
        helpers.FatalError(out)
    }
   
    helpers.LogInfo("Recursively Adding listeners to folders:")
    for _, dir := range helpers.Config.StartingDirs {

        helpers.LogInfo(fmt.Sprintf("\t", dir))

        err = watcher.Add(dir)
        if err != nil {
            out := fmt.Sprintf("Error adding watcher: %v", err)
            helpers.FatalError(out)
        }
    
        err = filepath.Walk(dir, add_folders); 
        if  err != nil {
            out := fmt.Sprintf("Error in Walk function: %v", err);
            helpers.LogInfo(out)
        }
    }
   
    helpers.LogInfo(fmt.Sprintf("Listening for file updates"))
    for {
        select {
            case ev := <-watcher.Events:
                stamp := time.Now()
                log_event(ev, stamp)
            case err := <-watcher.Errors:
                out := fmt.Sprintln("Error in file watcher: %v", err)
                helpers.LogInfo(out)
        }
    }
}


func add_folders(path string, fd os.FileInfo, err error) error {

     if err != nil {
        if pathErr, ok := err.(*fs.PathError); ok {
            if os.IsPermission(pathErr.Err) {
                out := fmt.Sprintf("Permission denied on: %v. Passing", path) 
                helpers.LogInfo(out)
                return nil
            }
        }
        out := fmt.Sprintf("Error Crashed File Watch Adder: %v", err)
        helpers.LogInfo(out)
        return err
    }
    
    // Only need to add the listeners to the directories
    if !fd.IsDir() { return nil }
    
    for _, regex := range helpers.Config.BlackListDirs {
        if res, _ := regexp.MatchString(regex, path); res {
            // Prevent further descent if thats the case
            out := fmt.Sprintf("Skipping blacklisted: %v", path)
            helpers.LogInfo(out)
            return filepath.SkipDir
        }
    }

    err = watcher.Add(path)
    if err != nil {
        out := fmt.Sprintf("Error adding watcher: %v", err)
        helpers.FatalError(out)
        return err
    }
   
    return nil
}


func log_event(ev noti.Event, stamp time.Time) {
    timestamp, _ := ptypes.TimestampProto(time.Now())

    MSG := bridge.FsEvent{
        SubmissionNumber:0,
        TimeStamp:timestamp,
        Op: int32(ev.Op),
        Location: string(ev.Name),
        DevID: helpers.GetID(),
    }

    // Save file content updates
    if (ev.Op == noti.Write) { 
        data, err := ioutil.ReadFile(string(ev.Name)) 
        if err != nil { 
            out := fmt.Sprintf("Error reading file while trying to log FsEvent: %v", err)
            helpers.LogInfo(out)
        } else {
            MSG.Content = data
        }
    }

    // Save 
    if ev.Op == noti.Create || ev.Op == noti.Chmod {

        // Strip out permissions information
        fileInfo, err := os.Stat(string(ev.Name))
        if err != nil {
            out := fmt.Sprintf("Error getting file info for: %v; %v", string(ev.Name), err)
            helpers.LogInfo(out)
        } else {
            // Make sure we can get the fileInfo. This su
            perms := fileInfo.Mode().Perm()
            MSG.Permissions = uint32(perms)
       
            // Get the user owners
            owner := fileInfo.Sys().(*syscall.Stat_t).Uid // Owner UID
            ownerName, err := getUsernameFromUID(owner)
            if err != nil {
                out := fmt.Sprintf("Error getting owner name: %v", err)
                helpers.LogInfo(out)
            } else {
                MSG.Owner = ownerName
            }

            // Get the group owner
            group := fileInfo.Sys().(*syscall.Stat_t).Gid // Group GID
            groupName, err := getGroupNameFromGID(group)
            if err != nil {
                out := fmt.Sprintf("Error getting group name: %v", err)
                helpers.LogInfo(out)
            } else {
                MSG.Group = groupName
            }
        }
    }

    _, e := client.SaveFsEvent(context.Background(), &MSG)
    if e != nil {
        out := fmt.Sprintf("Error while saving FsEvent: %v", e)
        helpers.LogInfo(out)
    }
    out := fmt.Sprintf("File output: %v", MSG.Location)
    helpers.LogInfo(out)
}


func getUsernameFromUID(uid uint32) (string, error) {
    user, err := user.LookupId(fmt.Sprintf("%d", uid))
    if err != nil {
        return "", err
    }
    return user.Username, nil
}


func getGroupNameFromGID(gid uint32) (string, error) {
    group, err := user.LookupGroupId(fmt.Sprintf("%d", gid))
    if err != nil {
        return "", err
    }
    return group.Name, nil
}

