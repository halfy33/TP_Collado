package main

import (
	"encoding/json"
	"net/http"

	"github.com/evoliatis/buildup/cpu"
	"github.com/evoliatis/buildup/disk"
	"github.com/evoliatis/buildup/load"
	"github.com/evoliatis/buildup/memory"
	"github.com/evoliatis/buildup/netcard"
)

func ReadLoadHandler(w http.ResponseWriter, r *http.Request) {
	data, err := load.ReadLoad()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func ReadDiskhandler(w http.ResponseWriter, r *http.Request) {
	data, err := disk.ReadDisk()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func ReadCpuHandler(w http.ResponseWriter, r *http.Request) {
	data, err := cpu.ReadCPU()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func ReadMemoryHandler(w http.ResponseWriter, r *http.Request) {
	data, err := memory.ReadMemory()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func ReadNetcardHandler(w http.ResponseWriter, r *http.Request) {
	data, err := netcard.ReadNetwork("toto")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
