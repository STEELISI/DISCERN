package main;

import "bufio"
import "fmt"
import "strings"
import "os"
import "path/filepath"
import "io/ioutil"
import "regexp"
import "io/fs"
import "time"
import "context"
import "google.golang.org/grpc" 
import "github.com/golang/protobuf/ptypes"


import bridge "FusionBridge/control/ansible"

import "DataSorcerers/helpers"

var opts []grpc.DialOption;
var conn *grpc.ClientConn;
var err error; 
var client bridge.AnsibleClient;


func main() {
    
    helpers.SetFileName("Ansible")
    helpers.LoadConfig()
    if !helpers.Config.RunAnsible { return }


    conn = helpers.CreateConnection()
    client = bridge.NewAnsibleClient(conn); 

    /* 
        We can find ansible config files in 4 places (in the following 
        locations:
            1) $ANSIBLE_CONFIG
            2) ansible.cfg (current directory)
            3) ~/.ansible.cfg
            4) /etc/ansible/ansible.cfg
        for #2 we are going to sweep /home 

        We need to can /home for all .yml files and then check if they
            contain - hosts:
        If the contain this we will save them as "Unverified" ansible
            playbooks. "Verified" playbooks are going to be from bash 
        Verify there aren't any playbooks saved in /etc/ansible
    */

    helpers.LogInfo("Welcome from ansible scraper")


    first_interval  := 0
    second_interval := 0

    first_func := func () {}
    second_func := func () {}
    
    // Man I'm sorry for this. I hate nested logic but this is minimal
    // I miss the ternary operator. This could be 4 wonderous lines
    if helpers.Config.SaveAnsiblePlaybooks && helpers.Config.SaveAnsibleConfigs {
        if helpers.Config.AnsibleConfigSweepInterval > helpers.Config.AnsiblePlaybookSweepInterval {
            first_interval  = helpers.Config.AnsiblePlaybookSweepInterval
            second_interval = helpers.Config.AnsiblePlaybookSweepInterval - helpers.Config.AnsibleConfigSweepInterval

            first_func = ScrapeConfigs
            second_func = ScrapePlaybooks
        } else {
            first_interval  = helpers.Config.AnsibleConfigSweepInterval
            second_interval = helpers.Config.AnsibleConfigSweepInterval - helpers.Config.AnsiblePlaybookSweepInterval 

            first_func = ScrapePlaybooks
            second_func = ScrapeConfigs
        }
    } else if helpers.Config.SaveAnsiblePlaybooks {
        first_interval  = helpers.Config.AnsiblePlaybookSweepInterval
        first_func = ScrapePlaybooks
    } else if helpers.Config.SaveAnsibleConfigs {
        first_interval  = helpers.Config.AnsibleConfigSweepInterval
        first_func = ScrapeConfigs
    } else {
        return
    }


    for {
        first_func()
        time.Sleep(time.Second * time.Duration(first_interval))
        second_func()
        time.Sleep(time.Second * time.Duration(second_interval))
    }
}

func ScrapeConfigs() {
    // Save #1
    helpers.LogInfo("Saving config from env")
    err := read_env()
    if err != nil { 
        out := fmt.Sprintf("Failed to read env: %v", err.Error())
        helpers.FatalError(out)
    }

    // Save #2
    helpers.LogInfo("Saving config and playbooks from home")
    err = sweep_home()
    if err != nil { 
        out := fmt.Sprintf("Failed to sweep home: %v", err.Error())
        helpers.FatalError(out)
    }
    

    // Save #3 & #4
    helpers.LogInfo("Saving config from ~ and /etc")    
    err = read_default_locations()
    if err != nil { 
        out := fmt.Sprintf("Failed to read default locations for configs: %v", err.Error())
        helpers.FatalError(out)
    }
}

func ScrapePlaybooks() {
    helpers.LogInfo("Saving playbooks from /etc/ansible")
    err = read_default_playbooks()
    if err != nil { 
        out := fmt.Sprintf("Failed to read default locations for playbooks: %v", err.Error())
        helpers.LogInfo(out)
    }
}

// For Config #1
func read_env() error {
    location := os.Getenv("ANSIBLE_CONFIG")
    if location == "" { return nil}
    // wg.Add(1)
    err := save_config(location)
    return err
}


// For Config #2 and .yml playbooks
func walkFn(path string, fd os.FileInfo, err error) error {
    if err != nil {
        helpers.FatalError(fmt.Sprintf("WalkFn: %v", err.Error()))
        return err
    }

    if fd.IsDir() { return nil }

    // Save the ansible config
    res, err := regexp.MatchString(`ansible.cfg$`, fd.Name())
    if err != nil {
        helpers.FatalError(fmt.Sprintf("Regex Fail: %v", err.Error()))
        return err
    }
    if res {
        // Sync so we can wait for sending in the main function
        // wg.Add(1)
        go save_config(path);
    }

    // Save the playbooks
    res, err = regexp.MatchString(`\.ya?ml$`, fd.Name())
    if err != nil {
        helpers.FatalError(fmt.Sprintf("Regex Fail: %v", err.Error()))
        return err
    }
    if res && verifyPlaybook(path) {
        // Sync so we can wait for sending in the main function
        // wg.Add(1)
        go save_yaml(path, false);
    }
    return nil
}

func verifyPlaybook(path string) bool {
    file, err := os.Open(path)
    if err != nil { 
        err_str := fmt.Sprintf("Error in opening file to verify playbook: %v", err.Error())
        helpers.LogInfo(err_str) 
        return false
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()

        if strings.Contains(line, "- hosts") {
            return true
        }
    }
    if err := scanner.Err(); err != nil {
        out := fmt.Sprintf("Error in scanner while verifying playbook: %v", err)
        helpers.LogInfo(out)
    }
    return false
}

func sweep_home() error {
    if err := filepath.Walk("/home", walkFn); err != nil {
        out := fmt.Sprintf("Error in Ansible Scraping Walk function: %v", err);
        helpers.LogInfo(out)
    }
    return nil
}


// For Config #3 & Config #4
func read_default_locations() error {
    err := save_config("~/ansible.cfg")
    if err != nil { 
        out := fmt.Sprintf("Error saving config from ~: %v", err)
        helpers.LogInfo(out)
    }
    err = save_config("/etc/ansible/ansible.cfg")
    if err != nil { 
        out := fmt.Sprintf("Error saving config from /etc: %v", err)
        helpers.LogInfo(out)
    }
    return nil
}


func save_playbooks(path string, fd os.FileInfo, err error) error {
     if err != nil {
        out := fmt.Sprintf("Error save_playbooks: %v", err.Error())
        helpers.FatalError(out)
        return err
    }

    if fd.IsDir() { return nil }

    res, err := regexp.MatchString(`\.ya?ml$`, fd.Name())

    if err != nil {
        out := fmt.Sprintf("Error Regex Failure save_playbooks: %v", err.Error())
        helpers.FatalError(out)
        return err
    }
    if res {
        // Sync so we can wait for sending in the main function
        // Any playbooks from this folder is definitely verified
            // Going to assume there is no need to check
        // wg.Add(1)
        go save_yaml(path, true);
    }
    return nil
}

// For scraping /etc for playbooks
func read_default_playbooks() error {
    if _, err := os.Stat("/etc/ansible"); os.IsNotExist(err) { 
        return nil
    }
    if err := filepath.Walk("/etc/ansible", save_playbooks); err != nil {
        out := fmt.Sprintf("Error in Jupyter Scraping Walk function: %v", err.Error());
        helpers.LogInfo(out)
        return err
    }
    return nil
}


// General helper functions
func save_config(path string) error {
    
    // defer wg.Done()

    if path == "" { return nil }
    data, err := ioutil.ReadFile(path) 
    if err != nil { 
        if _, ok := err.(*fs.PathError); !ok {
            return err 
        }
    }
    if string(data) == "" { return nil }

    timestamp, err := ptypes.TimestampProto(time.Now())
    MSG := bridge.ConfigDetails{
        SubmissionNumber: 0,
        TimeStamp: timestamp,
        Location: path,
        Content: string(data),
        DevID: helpers.GetID(),
    }

    _, err = client.SaveAnsibleConfig(context.Background(), &MSG);
    if err != nil {
        out := fmt.Sprintf("Error Saving Ansible Config: %v", err)
        helpers.LogInfo(out)
        return err
    }
    return nil
}


func save_yaml(path string, verified bool) error {

    ctx, cancel := context.WithTimeout(context.Background(), 
        10*time.Second)
    defer cancel()

    if path == "" { return nil }
    data, err := ioutil.ReadFile(path) 
    if err != nil { 
        if _, ok := err.(*fs.PathError); !ok {
            return err 
        }
    }
    if string(data) == "" { return nil }

    timestamp, err := ptypes.TimestampProto(time.Now())
    MSG := bridge.PlaybookDetails{
        SubmissionNumber: 0,
        TimeStamp: timestamp,
        Location: path,
        Content: string(data),
        Verified: verified,
        DevID: helpers.GetID(),
    }

    _, err = client.SaveAnsiblePlaybook(ctx, &MSG);
    if err != nil {
        out := fmt.Sprintf("Error Saving Ansible Playbook: %v", err.Error())
        helpers.LogInfo(out)
        return err
    }
    return nil
}

