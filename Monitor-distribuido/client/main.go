package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

type Report struct {
	Agent  string  `json:"agent"`
	CPU    float64 `json:"cpu"`
	Memory float64 `json:"memory"`
}

// Future-style async executor
func Async[T any](fn func() T) <-chan T {
	ch := make(chan T)
	go func() {
		ch <- fn()
	}()
	return ch
}

func getCPUUsage() float64 {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil || len(percent) == 0 {
		return 0
	}
	return percent[0]
}

func getMemoryUsage() float64 {
	vm, err := mem.VirtualMemory()
	if err != nil {
		return 0
	}
	return vm.UsedPercent
}

func reportUsage(endpoint string, agent string) {
	for {
		cpuFuture := Async(getCPUUsage)
		memFuture := Async(getMemoryUsage)

		cpuUsage := <-cpuFuture
		memUsage := <-memFuture

		report := Report{
			Agent:  agent,
			CPU:    cpuUsage,
			Memory: memUsage,
		}

		body, _ := json.Marshal(report)
		_, err := http.Post(endpoint, "application/json", bytes.NewBuffer(body))
		if err != nil {
			fmt.Println("Error reportando:", err)
		}

		time.Sleep(2 * time.Second)
	}
}

func main() {
	endpoint := "http://server:8080/report"
	agent := os.Getenv("AGENT_NAME")
	if agent == "" {
		agent = "unknown"
	}
	fmt.Println("Iniciando cliente:", agent)
	reportUsage(endpoint, agent)
}
