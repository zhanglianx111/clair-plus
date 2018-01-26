package system

import (
	"sync"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/cpu"
	"container/list"
	"time"
	"github.com/astaxie/beego/logs"
	"github.com/zhanglianx111/clair-plus/models"
	"github.com/astaxie/beego"
)

type SystemInterface interface {
	GetAverageInterval() models.OS
}

var once sync.Once
var sysHandler *systemHandler
var storeOSNum int
var callOSCycle int

type systemHandler struct {
	CPUIntervalQueue list.List
	MemIntercalQueue list.List
}

func GetSystemHandler() SystemInterface {

	return sysHandler
}

func init() {

	once.Do(func() {
		sysHandler = &systemHandler{}
	})

	storeOSNum = beego.AppConfig.DefaultInt("storeOSNum", 5)
	callOSCycle = beego.AppConfig.DefaultInt("callOSCycle", 1)

	//周期性监控cpu与内存
	go func() {

		ticker := time.NewTicker(time.Second * time.Duration(callOSCycle))

		for range ticker.C {
			sysHandler.MemLastIntervalQueue()
			sysHandler.CPULastIntervalQueue()
		}
	}()
}

func (s *systemHandler) GetAverageInterval() models.OS {

	mem := sysHandler.GetMemAveragePercent()
	mem.Totle /= 1024
	mem.Available /= 1024

	cpu := sysHandler.GetCPUAveragePercent()

	return models.OS{
		Memary: mem,
		CPU: cpu,
	}
}

func (s *systemHandler) GetMemAveragePercent() models.Memary {

	var perSum float64
	var avalSum uint64
	var totle uint64

	memQueue := s.GetMemQueue()

	//求和
	for mem := memQueue.Front() ; mem != nil ; mem = mem.Next() {

		value, ok := mem.Value.(models.Memary)
		if !ok {
			logs.Error("err:", ok)
		}
		//logs.Debug("memary队列:", value.UsedPercent)

		perSum += value.UsedPercent
		avalSum += value.Available
		totle = value.Totle
	}

	return models.Memary {
		Totle: totle,
		UsedPercent: perSum / float64(memQueue.Len()),
		Available: avalSum / uint64(memQueue.Len()),
	}
}

func (s *systemHandler) GetCPUAveragePercent() models.CPU {

	var perSum float64
	var mhz float64
	var cores int

	cpuQueue := s.GetCPUQueue()

	//求和
	for cpu := cpuQueue.Front() ; cpu != nil ; cpu = cpu.Next() {

		value, ok := cpu.Value.(models.CPU)
		if !ok {
			logs.Error("err:", ok)
		}
		//logs.Debug("cpu队列:", value.UsedPercent)

		perSum += value.UsedPercent
		mhz = value.Mhz
		cores = value.Cores
	}

	return models.CPU {
		Cores: cores,
		Mhz: mhz,
		UsedPercent: perSum / float64(cpuQueue.Len()),
	}
}

func (s *systemHandler) GetMemQueue() list.List {

	return s.MemIntercalQueue
}

func (s *systemHandler) GetCPUQueue() list.List {

	return s.CPUIntervalQueue
}

func (s *systemHandler) MemLastIntervalQueue() {

	mem, err := sysHandler.GetCurrentMemary()
	if err != nil {
		logs.Error("获取系统内存失败:", err)
	}

	s.MemIntercalQueue.PushBack(mem)
	//logs.Debug("men:", mem, "入队")

	//如果队列的item，大于时间间隔，则队首出列
	if s.MemIntercalQueue.Len() > storeOSNum {
		obsoleteMem := s.MemIntercalQueue.Front()
		s.MemIntercalQueue.Remove(obsoleteMem)
		//logs.Debug("mem:", obsoleteMem, "出队")
	}
}

func (s *systemHandler) CPULastIntervalQueue() {

	cpu, err := sysHandler.GetCurrentCPU()
	if err != nil {
		logs.Error("获取系统CPU失败:", err)
	}

	s.CPUIntervalQueue.PushBack(cpu)
	//logs.Debug("cpu:", cpu, "入队")

	//如果队列的item，大于时间间隔，则队首出列
	if s.CPUIntervalQueue.Len() > storeOSNum {
		obsoleteCpu := s.CPUIntervalQueue.Front()
		s.CPUIntervalQueue.Remove(obsoleteCpu)
		//logs.Debug("cpu:", obsoleteCpu, "出队")
	}
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

	cCore, err := cpu.Counts(true)
	if err != nil {
		return
	}

	cPer, err := cpu.Percent(0, false)
	if err != nil {
		return
	}

	cpuInfo = models.CPU{
		Cores: cCore,
		Mhz: cInfo[0].Mhz,
		UsedPercent: cPer[0],
	}

	return
}