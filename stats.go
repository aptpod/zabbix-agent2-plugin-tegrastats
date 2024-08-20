package main

import (
	"fmt"
	"regexp"
)

type Stats struct {
	EMC_FREQ    UsageStats
	GR3D_FREQ   UsageStats
	VIC_FREQ    UsageStats
	APE         string
	MTS_fg      string
	MTS_bg      string
	PLL         string
	MCPU        string
	PMIC        string
	Tboard      string
	GPU         string
	BCPU        string
	Thermal     string
	Tdiode      string
	VDD_SYS_GPU PowerStats
	VDD_SYS_SOC PowerStats
	VDD_IN      PowerStats
	VDD_SYS_CPU PowerStats
	VDD_SYS_DDR PowerStats
}

type UsageStats struct {
	Load      string
	Frequency string
}

type PowerStats struct {
	Current string
	Average string
}

var (
	reEMC_FREQ    = regexp.MustCompile(`EMC_FREQ (\d+)%(@\d+)?`)
	reGR3D_FREQ   = regexp.MustCompile(`GR3D_FREQ (\d+)%(@\d+)?`)
	reVIC_FREQ    = regexp.MustCompile(`VIC_FREQ (\d+)%(@\d+)?`)
	reAPE         = regexp.MustCompile(`APE (\d+)`)
	reMTS         = regexp.MustCompile(`MTS fg (\d+%) bg (\d+%)`)
	rePLL         = regexp.MustCompile(`PLL@(\d+\.?\d+?)C`)
	reMCPU        = regexp.MustCompile(`MCPU@(\d+\.?\d+?)C`)
	rePMIC        = regexp.MustCompile(`PMIC@(\d+\.?\d+?)C`)
	reTboard      = regexp.MustCompile(`Tboard@(\d+\.?\d+?)C`)
	reGPU         = regexp.MustCompile(`GPU@(\d+\.?\d+?)C`)
	reBCPU        = regexp.MustCompile(`BCPU@(\d+\.?\d+?)C`)
	reThermal     = regexp.MustCompile(`thermal@(\d+\.?\d+?)C`)
	reTdiode      = regexp.MustCompile(`Tdiode@(\d+\.?\d+?)C`)
	reVDD_SYS_GPU = regexp.MustCompile(`VDD_SYS_GPU (\d+)/(\d+)`)
	reVDD_SYS_SOC = regexp.MustCompile(`VDD_SYS_SOC (\d+)/(\d+)`)
	reVDD_IN      = regexp.MustCompile(`VDD_IN (\d+)/(\d+)`)
	reVDD_SYS_CPU = regexp.MustCompile(`VDD_SYS_CPU (\d+)/(\d+)`)
	reVDD_SYS_DDR = regexp.MustCompile(`VDD_SYS_DDR (\d+)/(\d+)`)
)

func ParseStats(output string) (*Stats, error) {
	stats := &Stats{}
	var err error

	if match := reEMC_FREQ.FindStringSubmatch(output); match != nil {
		stats.EMC_FREQ, err = parseUsageStats(match)
		if err != nil {
			return nil, err
		}
	}
	if match := reGR3D_FREQ.FindStringSubmatch(output); match != nil {
		stats.GR3D_FREQ, err = parseUsageStats(match)
		if err != nil {
			return nil, err
		}
	}
	if match := reVIC_FREQ.FindStringSubmatch(output); match != nil {
		stats.VIC_FREQ, err = parseUsageStats(match)
		if err != nil {
			return nil, err
		}
	}
	if match := reAPE.FindStringSubmatch(output); match != nil {
		stats.APE = match[1]
	}
	if match := reMTS.FindStringSubmatch(output); match != nil {
		stats.MTS_fg = match[1]
		stats.MTS_bg = match[2]
	}
	if match := rePLL.FindStringSubmatch(output); match != nil {
		stats.PLL = match[1]
	}
	if match := reMCPU.FindStringSubmatch(output); match != nil {
		stats.MCPU = match[1]
	}
	if match := rePMIC.FindStringSubmatch(output); match != nil {
		stats.PMIC = match[1]
	}
	if match := reTboard.FindStringSubmatch(output); match != nil {
		stats.Tboard = match[1]
	}
	if match := reGPU.FindStringSubmatch(output); match != nil {
		stats.GPU = match[1]
	}
	if match := reBCPU.FindStringSubmatch(output); match != nil {
		stats.BCPU = match[1]
	}
	if match := reThermal.FindStringSubmatch(output); match != nil {
		stats.Thermal = match[1]
	}
	if match := reTdiode.FindStringSubmatch(output); match != nil {
		stats.Tdiode = match[1]
	}
	if match := reVDD_SYS_GPU.FindStringSubmatch(output); match != nil {
		stats.VDD_SYS_GPU.Current = match[1]
		stats.VDD_SYS_GPU.Average = match[2]
	}
	if match := reVDD_SYS_SOC.FindStringSubmatch(output); match != nil {
		stats.VDD_SYS_SOC.Current = match[1]
		stats.VDD_SYS_SOC.Average = match[2]
	}
	if match := reVDD_IN.FindStringSubmatch(output); match != nil {
		stats.VDD_IN.Current = match[1]
		stats.VDD_IN.Average = match[2]
	}
	if match := reVDD_SYS_CPU.FindStringSubmatch(output); match != nil {
		stats.VDD_SYS_CPU.Current = match[1]
		stats.VDD_SYS_CPU.Average = match[2]
	}
	if match := reVDD_SYS_DDR.FindStringSubmatch(output); match != nil {
		stats.VDD_SYS_DDR.Current = match[1]
		stats.VDD_SYS_DDR.Average = match[2]
	}

	return stats, nil
}

func parseUsageStats(match []string) (UsageStats, error) {
	switch len(match) {
	case 2:
		return UsageStats{Load: match[1]}, nil
	case 3:
		return UsageStats{Load: match[1], Frequency: match[2]}, nil
	default:
		panic(fmt.Sprintf("unexpected number of matches: %v", match))
	}
}
