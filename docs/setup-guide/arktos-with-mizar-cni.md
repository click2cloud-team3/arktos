# Arktos to Enforce the Multi-tenant with Mizar Network Feature

This document captures the steps applied to an Arktos cluster lab enabling the unique multi-tenant network feature. The machines in this lab used are AWS EC2 t2-large (2 CPUs, 8GB mem), Ubuntu 18.04 LTS.

The steps might change with the progress of development.

If you would like to try with Flannel cni plugin, please ensure to read [multi-node setup guide](multi-node-dev-cluster.md).

1. Prepare lab machines. Particularly, check for local DNS cache (applicable to 127.0.0.53 name server; as it would cause coreDNS to crash with DNS loopback lookup):
```bash
sudo rm -f /etc/resolv.conf
sudo ln -s /run/systemd/resolve/resolv.conf /etc/resolv.conf
```
  Note: This step is required only when deploying arktos with mizar cni on On-Premise machine. 

2. Start Arktos cluster
```bash
CNIPLUGIN=mizar ./hack/arktos-up.sh
```

3. Leave the "arktos-up.sh" terminal and opened a another terminal to the master node. Run the following command to confirm that the first network, "default", in system tenant, has already been created. Its state is empty at this moment.
```bash
./cluster/kubectl.sh get net
NAME      TYPE    VPC                      PHASE   DNS
default   mizar   system-default-network    
```

Now, the default network of system tenant should be Ready.
```bash
./cluster/kubectl.sh get net
NAME      TYPE   VPC                       PHASE   DNS
default   mizar  system-default-network    Ready   10.0.0.207
```

From now on, you should be able to play with multi-tenant and the network features.