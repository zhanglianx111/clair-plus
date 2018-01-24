package system

import (
	"sync"
	"github.com/shirou/gopsutil/mem"
	"github.com/zhanglianx111/clair-plus/models"
	"github.com/shirou/gopsutil/cpu"
)

type SystemInterface interface {
	GetCurrentMemary() (memary models.Memary, err error)
	GetCurrentCPU() (cpuInfo models.CPU, err error)
}

var once sync.Once
var sysHandler *systemHandler

type systemHandler struct {
}

func GetSystemHandler() SystemInterface {

	once.Do(func() {
		sysHandler = &systemHandler{}
	})

	return sysHandler
}

func (s *systemHandler) GetCurrentMemary() (memary models.Memary, err error) {

	m, err := mem.VirtualMemory()
	if err != nil {
		return
	}

	memary = models.Memary{
		Totle: m.Total,
		Available: m.Available,
		UsedPercent: m.UsedPercent,
	}

	return
}

func (s *systemHandler) GetCurrentCPU() (cpuInfo models.CPU, err error) {

	cInfo, err := cpu.Info()
	if err != nil {
		return
	}

	cPer, err := cpu.Percent(0, false)
	if err != nil {
		return
	}

	cpuInfo = models.CPU{
		Cores: cInfo[0].Cores,
		Mhz: cInfo[0].Mhz,
		UsedPercent: cPer[0],
	}

	return
}