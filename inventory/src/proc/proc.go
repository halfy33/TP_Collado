package proc

import (
	"fmt"
	"strconv"

	"github.com/shirou/gopsutil/v4/process"
)

type Proc struct {
	Pid    int32   `json:"pid"`
	Name   string  `json:"name"`
	User   string  `json:"user"`
	CPU    float64 `json:"cpu"`
	Memory float32 `json:"memory"`
	Status string  `json:"status"`
}

// Chargement des Processus
func ReadProc(user string) (*[]Proc, error) {
	var myProcs []Proc
	procs, err := process.Processes()
	if err != nil {
		return nil, err
	}
	for _, proc := range procs {
		var myProc Proc
		myProc.User, _ = proc.Username()
		if user == "" || user == myProc.User {
			myProc.Pid = proc.Pid
			myProc.Name, _ = proc.Name()
			myProc.CPU, _ = proc.CPUPercent()
			myProc.Memory, _ = proc.MemoryPercent()
			status, _ := proc.Status()
			myProc.Status = status[0]
			myProcs = append(myProcs, myProc)
		}
	}
	return &myProcs, nil
}

// KillProc tue un processus par son PID
func KillProc(pid string) error {

	pid64, err := strconv.ParseInt(pid, 10, 32)
	if err != nil {
		return fmt.Errorf("PID invalide: %v", err)
	}
	pid32 := int32(pid64)

	proc, err := process.NewProcess(pid32)
	if err != nil {
		return fmt.Errorf("processus %d introuvable: %v", pid32, err)
	}

	exists, err := proc.IsRunning()
	if err != nil || !exists {
		return fmt.Errorf("processus %d n'est pas en cours d'exécution", pid32)
	}

	if err := proc.Kill(); err != nil {
		return fmt.Errorf("impossible de tuer le processus %d: %v", pid32, err)
	}

	return nil
}
