package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/evoliatis/inventory/cpu"
	"github.com/evoliatis/inventory/disk"
	"github.com/evoliatis/inventory/load"
	"github.com/evoliatis/inventory/memory"
	"github.com/evoliatis/inventory/netcard"
	"github.com/evoliatis/inventory/proc"
)

func HealthHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}

// Liste des processus
func PSHandler(w http.ResponseWriter, r *http.Request) {
	myProcs, err := proc.ReadProc("")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(myProcs)
}

// Liste des processus par user
func PSUserHandler(w http.ResponseWriter, r *http.Request) {
	user := r.PathValue("user")
	myProcs, err := proc.ReadProc(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(myProcs)
}

func KillHandler(w http.ResponseWriter, r *http.Request) {
	pid := r.PathValue("pid")
	if pid == "" {
		http.Error(w, "PID manquant", http.StatusBadRequest)
		return
	}

	err := proc.KillProc(pid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"message": "Processus %s tué avec succès"}`, pid)
}

func NetHandler(w http.ResponseWriter, r *http.Request) {
	out, err := netcard.ReadNetwork("")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(out)
}

func NetNameHandler(w http.ResponseWriter, r *http.Request) {
	card := r.PathValue("card")
	out, err := netcard.ReadNetwork(card)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(out)
}

func MemHandler(w http.ResponseWriter, r *http.Request) {
	out, err := memory.ReadMemory()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(out)
}

func DiskHandler(w http.ResponseWriter, r *http.Request) {
	out, err := disk.ReadDisk()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(out)
}

func LoadHandler(w http.ResponseWriter, r *http.Request) {
	out, err := load.ReadLoad()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(out)
}

func CPUHandler(w http.ResponseWriter, r *http.Request) {
	out, err := cpu.ReadCPU()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(out)
}
