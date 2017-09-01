---
title: Funnel Basic Scheduler

menu:
  main:
    parent: guides
    weight: 20
---

# Funnel's Basic Scheduler

Funnel provides a basic scheduler to handle task assignment to Funnel nodes. A
`Node` is a Funnel process that runs on a compute resource and report's available 
resources and starts worker processes. 

# Funnel's Scaler

Funnel also provides a scaler interface to automatically provision new nodes. 
Currently we have scalers for Google Cloud VMs and Openstack VMs.
