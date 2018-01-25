package system

import (
	"testing"
	"github.com/astaxie/beego"
)

func TestGet(t *testing.T) {

	sysHandler := &systemHandler{}
	cpu, err := sysHandler.GetCurrentCPU()
	if err !=nil {
		beego.Error(err)
		return
	}

	beego.Debug("cpu核数：", cpu.Cores)
	beego.Debug("cpu赫兹：", cpu.Mhz)
	beego.Debug("cpu使用率：", cpu.UsedPercent)

	memary, err := sysHandler.GetCurrentMemary()
	if err != nil {
		beego.Error(err)
	}

	beego.Debug("内存总量：", memary.Totle/1024)
	beego.Debug("内存使用率：", memary.UsedPercent)
	beego.Debug("内存剩余：", memary.Available)
}