package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

var noErr = fmt.Sprintf("%v", nil)

func TestParseTegrastats(t *testing.T) {

	tests := []struct {
		s         string
		expect    Stats
		expectErr string
	}{
		{
			s: "RAM 1873/3830MB (lfb 144x4MB) CPU [8%@2016,26%@2034,24%@2035,7%@2026,7%@2034,6%@2034] EMC_FREQ 0%@1600 GR3D_FREQ 0%@114 APE 150 MTS fg 0% bg 1% PLL@39C MCPU@39C PMIC@50C Tboard@35C GPU@37C BCPU@39C thermal@38.5C Tdiode@37C VDD_SYS_GPU 77/77 VDD_SYS_SOC 539/539 VDD_IN 2930/3405 VDD_SYS_CPU 539/937 VDD_SYS_DDR 727/758",
			expect: Stats{
				EMC_FREQ: UsageStats{
					Load:      "0",
					Frequency: "1600",
				},
				GR3D_FREQ: UsageStats{
					Load:      "0",
					Frequency: "114",
				},
				VIC_FREQ: UsageStats{},
				APE:      "150",
				MTS_fg:   "0",
				MTS_bg:   "1",
				PLL:      "39",
				MCPU:     "39",
				PMIC:     "50",
				Tboard:   "35",
				GPU:      "37",
				BCPU:     "39",
				Thermal:  "38.5",
				Tdiode:   "37",
				VDD_SYS_GPU: PowerStats{
					Current: "77",
					Average: "77",
				},
				VDD_SYS_SOC: PowerStats{
					Current: "539",
					Average: "539",
				},
				VDD_IN: PowerStats{
					Current: "2930",
					Average: "3405",
				},
				VDD_SYS_CPU: PowerStats{
					Current: "539",
					Average: "937",
				},
				VDD_SYS_DDR: PowerStats{
					Current: "727",
					Average: "758",
				},
			},
			expectErr: noErr,
		},
		{
			s: "EMC_FREQ 12%",
			expect: Stats{
				EMC_FREQ: UsageStats{
					Load: "12",
				},
			},
			expectErr: noErr,
		},
		{
			s: "EMC_FREQ 1%@",
			expect: Stats{
				EMC_FREQ: UsageStats{
					Load: "1",
				},
			},
			expectErr: noErr,
		},
		{
			s: "EMC_FREQ 1%@1",
			expect: Stats{
				EMC_FREQ: UsageStats{
					Load:      "1",
					Frequency: "1",
				},
			},
			expectErr: noErr,
		},
	}

	for i, test := range tests {
		stats, err := ParseStats(test.s)

		if !strings.Contains(fmt.Sprintf(`%v`, err), test.expectErr) {
			t.Fatalf("test[%d] error %v, want %v", i, err, test.expectErr)
		}

		if reflect.DeepEqual(stats, test.expect) {
			t.Errorf("test[%d]: %v, want: %v", i, stats, test.expect)
		}
	}
}
