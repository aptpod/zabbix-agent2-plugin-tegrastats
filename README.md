# Tegrastats utility plugin for Zabbix Agent 2
This plugin provides a solution to monitor [Tegrastats utility](https://docs.nvidia.com/drive/drive-os-5.2.0.0L/drive-os/index.html#page/DRIVE_OS_Linux_SDK_Development_Guide/Utilities/util_tegrastats.html). 
The plugin can monitor memory usage and processor usage for Tegra-based devices with [Zabbix agent 2](https://www.zabbix.com/documentation/current/en/manual/concepts/agent2) using [Tegrastats utility](https://docs.nvidia.com/drive/drive-os-5.2.0.0L/drive-os/index.html#page/DRIVE_OS_Linux_SDK_Development_Guide/Utilities/util_tegrastats.html).

*Output is limited to use with [active checks](https://www.zabbix.com/documentation/current/en/manual/appendix/items/activepassive#active-checks) in the current version.*

## Requirements
* Zabbix agent 2
* Go >= 1.18
* git
* make

## Installation

### From Source

1. Clone the repository:
    ```sh
    git clone https://github.com/aptpod/zabbix-agent2-plugin-tegrastats.git
    cd zabbix-agent2-plugin-tegrastats
    ```

2. Build the plugin:
    ```sh
    make
    ```

3. Copy the plugin binary to the installation directory:
    ```sh
    sudo install -m 0755 target/**/zabbix-agent2-plugin-tegrastats /usr/sbin/zabbix-agent2-plugin/
    sudo install -m 0644 tegrastats.conf /etc/zabbix_agent2.d/plugins.d/
    ```

4. Restart the Zabbix agent to load the new plugin:
    ```sh
    sudo systemctl restart zabbix-agent2
    ```

## Configuration
Open Zabbix agent 2 tegrastats configuration file `zabbix_agent2.d/plugins.d/tegrastats.conf` and set the required parameters.

**Plugins.Tegrastats.IntervalInMilliSeconds=1000** — the interval between log prints in milliseconds.
*Default value:* 1000
*Limits:* 1-300000


## Supported keys
**tegrastats.emc.usage** - returns percentage of EMC memory bandwidth in use relative to the current running frequency.

**tegrastats.gpu.usage** - returns proportion of GPU activation time in a period. different GPCs have the same percentage.

**tegrastats.vic.usage** - returns VIC engine loading as a percentage of current VIC engine frequency.

**tegrastats.temp.pll** - returns PLL temperature in degrees celsius.

**tegrastats.temp.mcpu** - returns MCPU temperature in degrees celsius.

**tegrastats.temp.pmic** - returns PMIC temperature in degrees celsius.

**tegrastats.temp.tboard** - returns Tboard temperature in degrees celsius.

**tegrastats.temp.gpu** - returns GPU temperature in degrees celsius.

**tegrastats.temp.bcpu** - returns BCPU temperature in degrees celsius.

**tegrastats.temp.thermal** - returns Thermal temperature in degrees celsius.

**tegrastats.temp.tdiode** - returns Tdiode temperature in degrees celsius.

**tegrastats.power.vdd_sys_gpu[\<Mode\>]** - returns VDD_SYS_GPU power consumption in milliwatts.  
*Parameters:*
Mode — possible values: current (current power consumption, default), avg (average power consumption).  
*Note:*
This value can only be retrieved by running the zabbix-agent with root access.

**tegrastats.power.vdd_sys_soc[\<Mode\>]** - returns VDD_SYS_SOC power consumption in milliwatts.  
*Parameters:*
Mode — possible values: current (current power consumption, default), avg (average power consumption).  
*Note:*
This value can only be retrieved by running the zabbix-agent with root access.

**tegrastats.power.vdd_in[\<Mode\>]** - returns VDD_IN power consumption in milliwatts.  
*Parameters:*
Mode — possible values: current (current power consumption, default), avg (average power consumption).  
*Note:*
This value can only be retrieved by running the zabbix-agent with root access.

**tegrastats.power.vdd_sys_cpu[\<Mode\>]** - returns VDD_SYS_CPU power consumption in milliwatts.  
*Parameters:*
Mode — possible values: current (current power consumption, default), avg (average power consumption).  
*Note:*
This value can only be retrieved by running the zabbix-agent with root access.

**tegrastats.power.vdd_sys_ddr[\<Mode\>]** - returns VDD_SYS_DDR power consumption in milliwatts.  
*Parameters:*
Mode — possible values: current (current power consumption, default), avg (average power consumption).  
*Note:*
This value can only be retrieved by running the zabbix-agent with root access.
