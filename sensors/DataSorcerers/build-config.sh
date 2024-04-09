#!/bin/bash


install_dir="/usr/local/bin/"


# What packages would you like to run
run_validation=false

run_logs=true

run_control_ansible=true
run_control_bash=true
run_control_jupyter=true

run_metadata_file=true
run_metadata_id=true
run_metadata_id_server=true # run the custom server
run_metadata_network=true
run_metadata_os=true



# What kind of OS information would you like to record
#   Need to enable run_metadata_os for this to have an effect
run_close_syscall_capture=true
run_close_range_syscall_capture=true
run_execve_syscall_capture=true
run_execveat_syscall_capture=true
run_fork_syscall_capture=true
run_kill_syscall_capture=true
run_open_syscall_capture=true
run_openat_syscall_capture=true
run_openat2_syscall_capture=true
run_recvfrom_syscall_capture=true
run_recvmmsg_syscall_capture=true
run_recvmsg_syscall_capture=true
run_sendmmsg_syscall_capture=true
run_sendmsg_syscall_capture=true
run_sendto_syscall_capture=true
run_socket_syscall_capture=true
run_socketpair_syscall_capture=true
run_sysinfo_syscall_capture=true
run_tkill_syscall_capture=true
run_vfork_syscall_capture=true

