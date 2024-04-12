## MergeTB Testbed Instrumentation
This guide helps to setup the environment for infrastructure dependent sensors responsible for data collection specific to attack scenarios.

### Source Code
The source code is placed at: https://gitlab.com/mergetb/tech/instrumentation
The `.deb` can be fetched from: https://gitlab.com/groups/mergetb/tech/-/packages

For our nodes (which are debian based), we can pull the images of latest `.deb` files and follow the following steps.

### Installation of DataSorcerers.deb on node

```bash
sudo apt-get update
sudo apt-get upgrade -y
sudo apt update
sudo apt upgrade -y
```

We install the `DataSorcerers.deb` by running the following command:

```bash
sudo dpkg -i DataSorcerers.deb
```

If the installation is encountering an error which states as follows:

```bash
Selecting previously unselected package datasorcerers.
(Reading database ... 60313 files and directories currently installed.)
Preparing to unpack DataSorcerers.deb ...
Hit:1 http://security.debian.org/debian-security bullseye-security InRelease
Hit:2 http://deb.debian.org/debian bullseye InRelease       
Hit:3 http://deb.debian.org/debian bullseye-updates InRelease
Hit:4 https://download.docker.com/linux/debian bullseye InRelease
Reading package lists... Done         
Building dependency tree... Done
Reading state information... Done
All packages are up to date.
Unpacking datasorcerers (1.0) ...
dpkg: dependency problems prevent configuration of datasorcerers:
 datasorcerers depends on strace; however:
  Package strace is not installed.
 datasorcerers depends on bpftrace; however:
  Package bpftrace is not installed.

dpkg: error processing package datasorcerers (--install):
 dependency problems - leaving unconfigured
Errors were encountered while processing:
 datasorcerers
```

Then, try to run the following commands:

```bash
sudo apt --fix-broken install
sudo dpkg -i DataSorcerers.deb
```

This should install `DataSorcerers.deb` and create the symbolic links.

### Installation of FusionCore.deb on node

We then install the `FusionCore.deb`

```bash
sudo dpkg -i FusionCore.deb
```

If the installation is encountering an error which states as follows:

```bash
Selecting previously unselected package fusioncore.
(Reading database ... 45450 files and directories currently installed.)
Preparing to unpack FusionCore.deb ...
Hit:1 http://security.debian.org/debian-security bullseye-security InRelease
Hit:2 http://deb.debian.org/debian bullseye InRelease
Hit:3 http://deb.debian.org/debian bullseye-updates InRelease
Reading package lists... Done
Unpacking fusioncore (1.0) ...
dpkg: dependency problems prevent configuration of fusioncore:
 fusioncore depends on docker-ce; however:
  Package docker-ce is not installed.
 fusioncore depends on docker-ce-cli; however:
  Package docker-ce-cli is not installed.
 fusioncore depends on containerd.io; however:
  Package containerd.io is not installed.
 fusioncore depends on docker-buildx-plugin; however:
  Package docker-buildx-plugin is not installed.
 fusioncore depends on docker-compose-plugin; however:
  Package docker-compose-plugin is not installed.

dpkg: error processing package fusioncore (--install):
 dependency problems - leaving unconfigured
Errors were encountered while processing:
 fusioncore
```

We need Docker for the FusionCore and hence, we run the following commands to avoid "Unmet Dependencies" error.

```bash
sudo apt install apt-transport-https ca-certificates curl software-properties-common
curl -fsSL https://download.docker.com/linux/debian/gpg | sudo apt-key add -
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/debian `lsb_release -cs` test"
sudo apt update
sudo apt install docker-ce
```

This should install `FusionCore.deb` and create the symbolic links.