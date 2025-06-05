package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Report struct {
	Agent  string  `json:"agent"`
	CPU    float64 `json:"cpu"`
	Memory float64 `json:"memory"`
}

func main() {
	http.HandleFunc("/report", reportHandler)

	fmt.Println("Servidor escuchando en :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func reportHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var rep Report
	if err := json.NewDecoder(r.Body).Decode(&rep); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	log.Printf("Reporte recibido de %s -> CPU: %.2f%%, Memoria: %.2f%%\n", rep.Agent, rep.CPU, rep.Memory)
	w.WriteHeader(http.StatusOK)
}
