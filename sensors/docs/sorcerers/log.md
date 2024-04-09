
# Log Sorcerer Architecture

This service captures dmesg logs and any other log files you specify. 
    This recording is done periodically and only records the last X 
    lines of each log file.


## Log Sorcerer Configuration Parameters

logsweepinterval: An integer representing how often the logs should be
    recorded

    - default: 5

loglinescaptured: An integer representing how many lines should be 
    saved from the end of each file

    - default: 25

logfiles: An array of files which should be read and saved

    - default:
      - "/var/log/messages"
      - "/var/log/syslog"
      - "/var/log/auth.log"
      - "/var/log/daemon.log"
      - "/var/log/kern.log"

 


