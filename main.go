package main

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"sync"

	"golang.zabbix.com/sdk/conf"
	"golang.zabbix.com/sdk/plugin"
	"golang.zabbix.com/sdk/plugin/container"
)

type Plugin struct {
	plugin.Base
	cfg PluginConfig

	ctx     context.Context
	cancel  context.CancelFunc
	muCtx   sync.Mutex
	ch      chan<- *Stats
	stats   *Stats
	muStats sync.RWMutex
}

type PluginConfig struct {
	System                 plugin.SystemOptions `conf:"optional"`
	IntervalInMilliSeconds int                  `conf:"optional,range=1:300000,default=1000"`
}

var impl Plugin

var (
	keyEmcUsage                  = "tegrastats.emc.usage"
	keyGpuUsage                  = "tegrastats.gpu.usage"
	keyVicUsage                  = "tegrastats.vic.usage"
	keyTemperaturePll            = "tegrastats.temp.pll"
	keyTemperatureMcpu           = "tegrastats.temp.mcpu"
	keyTemperaturePmic           = "tegrastats.temp.pmic"
	keyTemperatureTboard         = "tegrastats.temp.tboard"
	keyTemperatureGpu            = "tegrastats.temp.gpu"
	keyTemperatureBcpu           = "tegrastats.temp.bcpu"
	keyTemperatureThermal        = "tegrastats.temp.thermal"
	keyTemperatureTdiode         = "tegrastats.temp.tdiode"
	keyPowerConsumptionVddSysGpu = "tegrastats.power.vdd_sys_gpu"
	keyPowerConsumptionVddSysSoc = "tegrastats.power.vdd_sys_soc"
	keyPowerConsumptionVddIn     = "tegrastats.power.vdd_in"
	keyPowerConsumptionVddSysCpu = "tegrastats.power.vdd_sys_cpu"
	keyPowerConsumptionVddSysDdr = "tegrastats.power.vdd_sys_ddr"
)

// Start is called if this plugin is targeted by Active Checks Response
func (p *Plugin) Start() {
	p.muCtx.Lock()
	defer p.muCtx.Unlock()

	ctx, cancel := context.WithCancel(context.Background())

	ch := make(chan *Stats)
	logCh := make(chan string, 10)

	go func() {
		defer cancel()

		cmd := exec.CommandContext(p.ctx, "tegrastats", "--interval", fmt.Sprintf("%d", p.cfg.IntervalInMilliSeconds))
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			p.Errf("failed creating stdout pipe for 'tegrastats': %s", err)
			return
		}

		if err := cmd.Start(); err != nil {
			p.Errf("failed starting 'tegrastats': %s", err)
			return
		}
		defer cmd.Wait()

		scanner := bufio.NewScanner(stdout)
		for {
			select {
			case <-ctx.Done():
				close(ch)
				close(logCh)
				return
			default:
				if scanner.Scan() {
					line := scanner.Text()
					stats, err := ParseStats(line)
					if err != nil {
						p.Warningf("failed parsing 'tegrastats' output: %s", err)
						continue
					}
					p.SetStats(stats)
				} else if err := scanner.Err(); err != nil {
					p.Warningf("failed reading stdout from 'tegrastats': %s", err)
					return
				}
			}
		}
	}()

	p.ctx = ctx
	p.cancel = cancel
	p.ch = ch
}

// Stop is called if this plugin is not targeted by Active Checks Response
func (p *Plugin) Stop() {
	p.muCtx.Lock()
	defer p.muCtx.Unlock()
	p.cancel()
}

func (p *Plugin) Configure(globalOptions *plugin.GlobalOptions, privateOptions interface{}) {
	conf.Unmarshal(privateOptions, &p.cfg)
}

func (p *Plugin) Validate(privateOptions interface{}) error {
	if err := conf.Unmarshal(privateOptions, &p.cfg); err != nil {
		return fmt.Errorf("invalid configuration for tegrastats plugin: %w", err)
	}
	return nil
}

func (p *Plugin) Export(key string, params []string, ctx plugin.ContextProvider) (result interface{}, err error) {

	stats, ok := p.Stats()
	if !ok {
		return nil, nil
	}

	switch key {
	case keyEmcUsage:
		return stats.EMC_FREQ.Load, nil
	case keyGpuUsage:
		return stats.GR3D_FREQ.Load, nil
	case keyVicUsage:
		return stats.VIC_FREQ.Load, nil
	case keyTemperaturePll:
		return stats.PLL, nil
	case keyTemperatureMcpu:
		return stats.MCPU, nil
	case keyTemperaturePmic:
		return stats.PMIC, nil
	case keyTemperatureTboard:
		return stats.Tboard, nil
	case keyTemperatureGpu:
		return stats.GPU, nil
	case keyTemperatureBcpu:
		return stats.BCPU, nil
	case keyTemperatureThermal:
		return stats.Thermal, nil
	case keyTemperatureTdiode:
		return stats.Tdiode, nil
	case keyPowerConsumptionVddSysGpu:
		return getPowerStatsByMode(stats.VDD_SYS_GPU, params), nil
	case keyPowerConsumptionVddSysSoc:
		return getPowerStatsByMode(stats.VDD_SYS_SOC, params), nil
	case keyPowerConsumptionVddIn:
		return getPowerStatsByMode(stats.VDD_IN, params), nil
	case keyPowerConsumptionVddSysCpu:
		return getPowerStatsByMode(stats.VDD_SYS_CPU, params), nil
	case keyPowerConsumptionVddSysDdr:
		return getPowerStatsByMode(stats.VDD_SYS_DDR, params), nil
	default:
		return nil, plugin.UnsupportedMetricError
	}
}

func (p *Plugin) Stats() (Stats, bool) {
	p.muStats.RLock()
	defer p.muStats.RUnlock()

	if p.stats == nil {
		return Stats{}, false
	}
	return *p.stats, true
}

func (p *Plugin) SetStats(stats *Stats) {
	p.muStats.Lock()
	defer p.muStats.Unlock()
	p.stats = stats
}

func getPowerStatsByMode(stats PowerStats, param []string) string {
	current := true
	if len(param) > 0 && param[0] == "avg" {
		current = false
	}
	if current {
		return stats.Current
	}
	return stats.Average
}

func init() {
	plugin.RegisterMetrics(&impl, "Tegrastats",
		keyEmcUsage, "Percent of EMC memory bandwidth being used.",
		keyGpuUsage, "Percent of the GR3D that is being used.",
		keyVicUsage, "Percent of the VIC that is being used.",
		keyTemperaturePll, "Temperature of the PLL in degrees celsius.",
		keyTemperatureMcpu, "Temperature of the MCPU in degrees celsius.",
		keyTemperaturePmic, "Temperature of the PMIC in degrees celsius.",
		keyTemperatureTboard, "Temperature of the Tboard in degrees celsius.",
		keyTemperatureGpu, "Temperature of the GPU in degrees celsius.",
		keyTemperatureBcpu, "Temperature of the BCPU in degrees celsius.",
		keyTemperatureThermal, "Temperature of the thermal in degrees celsius.",
		keyTemperatureTdiode, "Temperature of the Tdiode in degrees celsius.",
		keyPowerConsumptionVddSysGpu, "Power consumption of VDD_SYS_GPU in milliwatts, usage: "+keyPowerConsumptionVddSysGpu+"[<current|avg>].",
		keyPowerConsumptionVddSysSoc, "Power consumption of VDD_SYS_SOC in milliwatts, usage: "+keyPowerConsumptionVddSysSoc+"[<current|avg>].",
		keyPowerConsumptionVddIn, "Power consumption of VDD_IN in milliwatts, usage: "+keyPowerConsumptionVddIn+"[<current|avg>].",
		keyPowerConsumptionVddSysCpu, "Power consumption of VDD_SYS_CPU in milliwatts, usage: "+keyPowerConsumptionVddSysCpu+"[<current|avg>].",
		keyPowerConsumptionVddSysDdr, "Power consumption of VDD_SYS_DDR in milliwatts, usage: "+keyPowerConsumptionVddSysDdr+"[<current|avg>].",
	)
}

func main() {
	h, err := container.NewHandler(impl.Name())
	if err != nil {
		panic(fmt.Sprintf("failed to create plugin handler %s", err.Error()))
	}
	impl.Logger = &h

	err = h.Execute()
	if err != nil {
		panic(fmt.Sprintf("failed to execute plugin handler %s", err.Error()))
	}
}
