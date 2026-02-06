package main

import (
	"log"
	"net/http"
	"os"

	"github.com/evoliatis/buildup/goroutines"
	"github.com/joho/godotenv"
)

func main() {
	// Charge le .env
	if err := godotenv.Load(); err != nil {
		log.Fatalln("fichier .env non trouvé")
	}
	InToken := os.Getenv("INFLUXDB_TOKEN")
	InUrl := os.Getenv("INFLUXDB_URL")
	InOrg := os.Getenv("INFLUXDB_ORG")
	InBucket := os.Getenv("INFLUXDB_BUCKET")

	// Charge la liste des serveurs
	cfg, err := LoadConfig("servers.yaml")
	if err != nil {
		log.Fatalln("fichier servers.yaml non trouvé")
	}

	// Go routines
	for _, srv := range cfg.Servers {

		go goroutines.GoLoad(srv, InUrl, InOrg, InBucket, InToken, 10)
		go goroutines.GoProc(srv, InUrl, InOrg, InBucket, InToken, 10)
		go goroutines.GoCPU(srv, InUrl, InOrg, InBucket, InToken, 10)
		go goroutines.GoMem(srv, InUrl, InOrg, InBucket, InToken, 10)
		go goroutines.GoDisk(srv, InUrl, InOrg, InBucket, InToken, 10)
		go goroutines.GoNet(srv, InUrl, InOrg, InBucket, InToken, 10)
	}
	log.Println("listening on :8084")
	log.Fatal(http.ListenAndServe(":8084", router()))
}
