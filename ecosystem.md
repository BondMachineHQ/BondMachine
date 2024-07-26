
***The BondMachine ecosystem*** is formed by two parts. The first concerns the BMs devices building with their internal behaviour and architecture, the seconds the BMs interconnecting, both among themselves and with the external world.

### **BondMachines building**

A BondMachine is a set of computing units and non computing objects shared among them, all packed within a single hardware device.

#### *Connecting Processors*

The *Connecting Processor* (CP) is the computing core of the BondMachine. The name Connecting Processor describes the capability of the processor core to be configured in such a way to be connected to other processors and to Shared Objects. CPs are as simple as possible, specialized to do a single task and optimized for doing it. In fact, the CPs inside the BondMachine architecture can have different number of registers, number of input/output register, different instruction sets (i.e. opcodes) with respect to the other ones.

#### *Shared Objects*

These are objects that are shared between the CPs (SO). Several kind of objects can be implemented to increase the processing capability and the functionality of the BMs improving the high-speed synchronization and communication between tasks running on separate CPs and all other non computing capabilities.
Three kinds of objects have been currently implemented: Channels, Shared Memories and Barriers.

### **BondMachines interconnecting**

BMs may be connected together to form clusters or can interact with the external world.

#### *EtherBond protocol*

BMs may comunicate using a native protocol called EtherBond. It is a protocol over ethernet (ethertype 0x8888) whose purpose is to replicate the electronic behavior of BMs registers and to extend it over the device boundaries. In other words clusters of BMs may be created and their behaviour is driven be the same rules of the BM devices. The main objective is to handle devices and cluster the same way.

#### *EtherBond on Linux*

The EtherBond protocol has been ported to Linux. With this port BMs can comunicate with standard PC software.

#### *Board drivers*

BMs may be implemented in hardware using FPGA, several drivers, specific to different FPGA model have been developed in order to use displays, leds, switch et al.
