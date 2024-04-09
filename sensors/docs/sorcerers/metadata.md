
# Metadata Sorcerers Architecture

Each of these services are implemented independently and do not require
    the other to be run



## File Metadata Sorcerer

This service uses fsnotify (an inotify wrapper in golang) to trigger 
    data recording events on file creation, deletion, renaming, writing, 
    and permissions changing.

The files which this programs listen to can be changed using the 
    config option: startingdirs. These directories are recursively 
    iterated down and listeners are added to all files and directories,

Certain folders can generate lots of data, so a blacklist is also 
    included to prevent the program from listening to certain files 
    and folders. This can be specified with: blacklistdirs.

blacklistdirs: # These are regexes
  - '/home/\..*/\.config'
  - '/home/\..*/\..*'
startingdirs: # These are regexes
  - "/tmp"
  - "/home"


### When & What Data is Recorded: 

When a file is writen to, the binary data of the updated file is saved 
    to the database. 

When a file is created or has permission changed, the owner and group
    is saved to the database.

When a file has its name changed, the new name is recorded to the 
    database

When a file is deleted only a timestamp is recorded


### Note

You may need to increase the following parameters, depending on how 
    many files you watch:

    fs.inotify.max_user_watches=124983
    fs.inotify.max_user_instances=128



## ID Metadata Sorcerer

This service scrapes 2 kinds of information, users on the machine and
    interface information
    - User information is scraped from /home

The duration for scraping user information is controlled by:
    
    scrapforusers: 30000 # default in seconds

The duration for scraping interface inforamtion is controlled by:

    scrapeinterfaceinfo: 300 # default in seconds


### Interface information

The information for each interface is:

    - MAC address
    - Name
    - An array of hostname and IP pairs

This should allow for correlation between the packets data and the 
    device which sent the data, even in the event of changing IPs



## Network Metadata Sorcerer

This service sets a listener on every interface accessible via 
    gopacket/pcap and sends detailed information on the packet:

    For ARP Packets:
        - Source Hardware Address
        - Source Protocol Address
        - Destination Hardware Address
        - Destination Protocol Address
        - ARP Operation Code
        - ARP Protocol

    For Ethernet Packets:
        - Source MAC Address
        - Destination MAC Address
        - Payload length

    For IP Packets:
        - Source IP Address
        - Destination IP Address
        - V4 / V6 boolean

    For TCP Packets:
        - Source Port
        - Destination Port

    For UDP Packets:
        - Srouce Port
        - Destination Port

    For DNS Packets:
        - DNS Questions
        - DNS Resource Records

    TLS and ICMP application layer protocols are also saved

    The device the data was sent on

These packets are created quite freqently and can quickly busy the 
    network. To get around this, you can adjust the parameter 
    networkslicelength. This parameter sets how many (packets, timestamp)
    pairs are saved before the data is sent to the Core Server

    networkslicelength: 25 # default in number of packets



## OS Metadata Sorcerer

The OS metadata is separated into two sections. The first section is a 
    series of bpftrace programs which log data, as jsons, into files 
    in /tmp/discern/data/os/. The second part is a scrapper, written in 
    go, which parses the json and sends the data, periodically, to the
    Core data server. Each bpftrace program has a different json output,
    so we cannot generalize the client to allow for an arbitrary bpftrce
    program. Thankfully, implementing new bpftrace prgrams is rather 
    trivial.

To kill the data recording program, you will need to kill the series of
    bpftrace programs which spawn when the service is started. The PIDs
    can easily be found with : `sudo ps -e | grep bpftrace`

The data directory is /tmp/discern/data/os/ by default but can be 
    specified with osdatadir

The interval between dumps is specified by osinfodumpinterval and is
    set to 3001 by default

A following is a list of all the systemcalls being recorded and the
    data recorded along side the system call:

    - close: # of times called by a give (pid, uid, gid) tuple over a 
        60 second interval

    - close_range: # of times called by a give (pid, uid, gid) pair 
        over a 60 second interval

    - execve: timestamp, pid, uid, gid, and args0-4. Because kprobs 
        don't allow for looping, reading all the arguments is hard. 4
        seemed like a sweet spot to get all the necessary information

    - execveat: timestamp, pid, uid, gid, and args0-4. Because kprobs 
        don't allow for looping, reading all the arguments is hard. 4
        seemed like a sweet spot to get all the necessary information

    - fork: Counts the number of times a given (pid, tid, uid, gid) 
        tuple calls fork

    - kill: Records the signal and program PID a given (pid, uid, gid)
        called kill on

    - open: Conuts the number of times a given (pid, uid, gid) opens
        a particular file 

    - openat: Conuts the number of times a given (pid, uid, gid) opens
        a particular file 

    - recvfrom: Counts the number of times a given (pid, uid gid)
        recieves from a given file descriptor

    - recvmmsg: Counts the number of times a given (pid, uid gid)
        recieves from a given file descriptor

    - recvmsg: Counts the number of times a given (pid, uid gid)
        recieves from a given file descriptor

    - sendmmsg: Counts the number of times a given (pid, uid gid)
        recieves from a given file descriptor and the length

    - sendmsg: Counts the number of times a given (pid, uid gid)
        recieves from a given file descriptor and the length

    - sendto: Counts the number of times a given (pid, uid gid)
        recieves from a given file descriptor and the length

    - socket: Records how many sockets with a given family, type, 
        and protocol a (pid, uid, gid) creates

    - socketpair: Records how many sockets with a given family, type, 
        and protocol a (pid, uid, gid) creates

    - socketpair: Records how many sockets with a given family, type, 
        and protocol a (pid, uid, gid) creates

    - sysinfo: Records when a given (pid, uid, gid) calls to ask for 
        system information

    - tkill: Counts how many times a given (pid, uid, gid) calls tkill
        as well as the signal and argument pid that it gets called on

    - vfork: Counts how many time a given (pid, uid, gid) calls vfork

