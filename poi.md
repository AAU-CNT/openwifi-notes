# Points of interrest

## Number of backoff waits

<https://github.com/open-sdr/openwifi-hw/blob/5705bb3a200270768235be0dc189d5aed0664e35/ip/xpu/src/csma_ca.v#L267>

Packages are sent when the statemachine is in `BACKOFF_RUN`.
`ch_idle_final` is on when CCA finds the channel idle.

This test could be implemented with a counter `BACKOFF_CH_BUSY`, which could be reset when `BACKOFF_RUN` is entered
or then a package is sent.
