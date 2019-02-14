# PatrIoT_router
This repository is used to build docker image of router.
Router is used in network-simulator to control network communication.
Router image has to be built before starting network generator.
Network generator uses image tag to create router containers.

## Parts
<ul>
<li>iptables-api</li>
<li>iproute2-api</li>
<li>iproute2-rest</li>
</ul>

### IPTables
Iptables api written in golang is used to communicate with linux ip tables.
Api is used to filter communication, NAT, ... <br>
This api is forked from: https://github.com/Oxalide/iptables-api


### IPRoute2
IPRoute2 api written in golang is used to control routing tables
using iproute2 commands. Api supports:
<ul>
<li>Default routes creationg and modification</li>
<li>Routes creation and modification</li>
<li>List interfaces</li>
<li>List routes</li>
</ul>

### IPRoute2 REST API
IProute2 REST is used by network simulator api to remotely
control routes on routers. API is written in golang, using http package.


## Build
```docker build ./```
