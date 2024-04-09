#!/bin/bash

source ./build-config.sh

sudo apt update -y && sudo apt upgrade -y

# Install go if it doesn't exist
if ! command -v go &> /dev/null; then
    echo "Go is not installed. Initiating installation...";
    wget https://golang.org/dl/go1.21.3.linux-amd64.tar.gz;
    sudo tar -C /usr/local -xzf go1.21.3.linux-amd64.tar.gz
    sudo ln -s /usr/local/go/bin/go /usr/local/bin/go
else
    echo "Go is already installed."
fi


mkdir -p /etc/discern
mkdir -p /etc/discern/data
mkdir -p /etc/discern/data/os


if [ ! -f "/etc/discern/SocererConfig.yaml" ]; then
    cp ./SorcererConfig.yaml /etc/discern/SorcererConfig.yaml
fi


# Build validation
if [ "$run_validation" = "true" ]; then
    cd ./validation
    go mod tidy
    go build
    mv validation "$install_dir/discern-validation-sorcerer"
    cd ..
fi

# Build log
if [ "$run_logs" = "true" ]; then
    cd ./logs
    go mod tidy
    go build
    mv logs "$install_dir/discern-logs-sorcerer"
    cd ..
fi

# Build control/ansible
if [ "$run_control_ansible" = "true" ]; then
    cd ./control/ansible
    go mod tidy
    go build
    mv ansible "$install_dir/discern-ansible-sorcerer"
    cd ../..
fi

# Build control/bash
if [ "$run_control_bash" = "true" ]; then
    cd ./control/bash
    go mod tidy
    go build
    mv bash "$install_dir/discern-bash-sorcerer"
    ./install.sh
    cp ./start_ttylog.sh "$install_dir/discern-tty-log"
    cd ../..
fi

# Build control/jupyter
if [ "$run_control_jupyter" = "true" ]; then
    cd ./control/jupyter
    go mod tidy
    go build
    mv jupyter "$install_dir/discern-jupyter-sorcerer"
    cd ../..
fi

# Build metadata/file
if [ "$run_metadata_file" = "true" ]; then
    cd ./metadata/file
    go mod tidy
    go build
    mv file "$install_dir/discern-file-sorcerer"
    cd ../..
fi

# Build metadata/id
if [ "$run_metadata_id" = "true" ]; then
    cd ./metadata/id
    go mod tidy
    go build
    mv id "$install_dir/discern-id-sorcerer"
    cd ../..
fi

# Build metadata/id/server
if [ "$run_metadata_id_server" = "true" ]; then
    cd ./metadata/id/server
    go mod tidy
    go build
    mv server "$install_dir/discern-id-server-sorcerer"
    cd ../../..
fi

# Build metadata/network
if [ "$run_metadata_network" = "true" ]; then
    cd ./metadata/network
    go mod tidy
    go build
    mv network "$install_dir/discern-network-sorcerer"
    cd ../..
fi

# Build metadata/os
if [ "$run_metadata_network" = "true" ]; then
    cd ./metadata/os
    go mod tidy
    go build
    mv metaOS "$install_dir/discern-os-sorcerer"
    # Build the bpftrace sensors so they can be run
    mkdir -p /etc/discern/bpfsensors
    cp ./bpftrace/* /etc/discern/bpfsensors/
    cd ../..


    # Building the bpftrace sensors file
    echo "#!/bin/bash" > ./discern-bpftrace-sensors
    chmod +x ./discern-bpftrace-sensors


    # Iterate over files in the directory
    for file in /etc/discern/bpfsensors/*; do
        # Check if the file is a regular file
        if [[ -f "$file" ]]; then
            # Strip off the extension using parameter expansion
            filename=$(basename "$file")
            filename_no_ext="${filename%.*}"  # Remove extension
            # Find corresponding should run parameter
            var_name="run_${filename_no_ext}_syscall_capture"
            # Check if we should add this
            should_run=${!var_name}
            # Add if so
            if [ "$should_run" = "true" ]; then
                echo "sudo bpftrace -f json /etc/discern/bpfsensors/${filename} >> /tmp/discern/data/os/${filename_no_ext}-res.txt &" >> ./discern-bpftrace-sensors
            fi
        fi
    done

    mv ./discern-bpftrace-sensors "$install_dir/discern-bpftrace-sorcerer"

fi


# Copy the service
cp ./services/* /etc/systemd/system/
sudo systemctl daemon-reload

# sudo systemctl enable discern-.service
# sudo systemctl start discern-.service

echo "Loading the services. Some may fail, thats ok."
echo "Just make sure whatever you enabled doesn't fail"
for file in ./services/*; do
    # Check if the file is a regular file
    if [[ -f "$file" ]]; then
        # Strip off the extension using parameter expansion
        filename=$(basename "$file")
        sudo systemctl enable $filename
        sudo systemctl start $filename
    fi
done


