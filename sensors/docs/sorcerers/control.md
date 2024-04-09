
# Control Sorcerers Architecture

There are 3 Control Data Sorcerers which scrape data for ansible, bash, 
    and jupyter notebooks respectively 
    - Each of these services can be run independently and do not require
        the others to run

Configuration for each of these services are saved in /etc/discern/SorcererConfig.yaml

The Ansible and Jupyter sorcerer should be run on any control nodes

The Bash sorcerer should be run locally on any node whos bash interactions you
    want to record

## Ansible Implementation And Configuration

This service periodically scans the system for ansible config files and playbooks

    - Ansible config files can be found at $ANSIBLE_CONFIG, ~/.ansible.cfg,
        ansible.cfg (in current directory), or /etc/ansible/ansible.cfg

    - Ansible playbooks can be identified by searching for yaml files containing
        the line `- hosts:`

Sweeping for config files can be turned on and off using the option:
    
    saveansibleconfigs: true # default

Sweeping for playbooks can be controlled using:

    saveansibleplaybooks: true

The period of time between config sweeps is controlled with:

    ansibleconfigsweepinterval: 5 # default in seconds

The period of time between playbook sweeps is controlled with:

    ansibleplaybooksweepinterval: 5



## Bash Implementation And Configuration

This service runs strace to listen for prints to STDIO and keyboard 
    entries. These outputs are saved, with timetamps, files in 
    /var/log/ttylog/

This program is split into 2 separate services. The first listens for
    user IO and records the text and timestamps to a file. The second 
    reads this file and dumps it to the remote server. To install the 
    recorder individually, simply run DataSorcer/control/bash/install.sh.
    The service which dumps the data to the remote server is created 
    using go build .

The output files are then collected and periodically dumped to the Core
    server.

The interval to dump the bash information is controlled by:

    bashsweepinterval: 5 # default in seconds



## Jupyter Implementation And Configuration

Jupyter notebooks can be identified with the file extension .ipynb. The
    jupyter scraper simply scans directories with this file extension and 
    saves them to the database

Scraping every directory in the system is quite a large task, so the
    search space is reduced to recursively searching the directories 
    listed in:

    jupyterpaths:
      - "/home"

And the frequency of this scraping is controlled with:

    jupytersweepinterval: 5 # default in seconds






