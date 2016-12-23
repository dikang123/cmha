## CMHA:Consul based MySQL High Availability  [![Release](https://img.shields.io/badge/release-v.1.1.5--beta-blue.svg)](https://github.com/upmio/cmha/releases/latest)

CMHA is a rock solid MySQL high availability cluster solution for mission critical use case like financial service and telco.

## Features of CMHA

1. Written in Golang, light weight and very easy to deploy as well as upgrade. 
2. No single point of failure design and strong data consistence, no data lose, no transaction mess.
3. Consul cluster maintains multiple HA groups for different applications.
4. Prevent network partition and MySQL nodes brain-split.
5. Stateless agent design with Run-As-Needed failure handler and health monitoring handler.
6. Interactive command line(CLI) console for DBA's troubleshooting and maintenance.
7. Clear and elegant web UI for status monitoring (embedded web service, no need for additional web server).
8. RESTFul interface for external automation ops tools/platform.
9. User-friendly and customizable deployment scripting
10. Upgrading on the fly 


## Getting Started

###1. Prerequisities

* 7 nodes for production (physical or VM)
* OS Supported : CentoS 6.x/7.x, RHEL 6.x/7.x
* Local YUM Repository
* Disable SELinux
* SSHD enabled
* NTP enabled (not necessary but suggested)
* Disable iptables
* Configure /etc/hosts on each node for 

###2. Role of cluster nodes

```
CS:   Consul server for cluster coordinator ( 1 node for testing purpose, 3 nodes for production)
CHAP: Application access node (1 node for testing purpose,2 nodes for production)
DB:   MySQL instance node (2 nodes for Master-Master configuration)
```
###3. Install from binary 

* download pre-built binary installer(with *auto-deployment* in the file name) from https://github.com/upmio/cmha/releases/latest
* extract the package to any Linux node as the deployment node (can be one of the cluster or additional node for temporary use only)
* deployment node require "expect" package installed
* install "openssl" package on CHAP nodes
* edit auto-deployment.ini in the extracted directory of the binary installer, follow the inline instruction
* execute deployment-check.sh script to verify the environment and necessary packages
* run the installer to auto deploy the cluster:
```
#sudo auto-deployment.sh cs
#sudo auto-deployment.sh ca
#sudo auto-deployment.sh db
#sudo auto-deployment.sh chap
```
* installation verification, use favorite brower, visit http://CS_IP:8500/ui to check the cluster status, see more details from the project wiki page (https://github.com/upmio/cmha/wiki)

## Contributing
1. Fork it!
2. Create your feature branch: git checkout -b my-new-feature
3. Commit your changes: git commit -am 'Add some feature'
4. Push to the branch: git push origin my-new-feature
5. Submit a pull request :D


## Get Support 
Please use our github issues(https://github.com/upmio/cmha/issues) for bugs, feature requests or questions !

## Licensing
CMHA is licensed under the Apache License, Version 2.0. See
[LICENSE](https://github.com/upmio/cmha/blob/master/LICENSE) for the full
license text.

