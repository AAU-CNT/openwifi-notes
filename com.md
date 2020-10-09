# Communication

Different ways to extract the measured data.

## Xilinx ILA

Seems to be a function in Xilinx vivado.
Recommended by openwifi as a method to extract state of internal state machines.

<https://github.com/open-sdr/openwifi/tree/master/doc>

## Using AXI without DMA

Here one could use one of the free interrupts to interrupt the the cpu when new data is available.
The CPU could then extract the info using AXI from either the xpu or a custom info module.

## Using DMA on top of AXI

This is probably better if large amounts of data should be moved.
However seems to be much more complicated.

## Using PMOD ports

Not use the linux kernel at all and send data directly to the PMOD ports from the logic unit.
This requires from hardware to receive it on the other end.
