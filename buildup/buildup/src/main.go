package main

import (
	"log"
	"net/http"
	"os"

	"github.com/evoliatis/buildup/goroutines"
	"github.com/joho/godotenv"
)

func main() {
	// .env 파일을 불러오세요
	if err := godotenv.Load(); err != nil {
		log.Fatalln("fichier .env non trouvé")
	}
	InToken := os.Getenv("INFLUXDB_TOKEN")
	InUrl := os.Getenv("INFLUXDB_URL")
	InOrg := os.Getenv("INFLUXDB_ORG")
	InBucket := os.Getenv("INFLUXDB_BUCKET")

	// 서버 목록을 불러오세요
	cfg, err := LoadConfig("servers.yaml")
	if err != nil {
		log.Fatalln("fichier servers.yaml non trouvé")
	}

	// Go routines
	for _, srv := range cfg.Servers {

		go goroutines.GoLoad(srv, InUrl, InOrg, InBucket, InToken, 10)
		go goroutines.GoProc(srv, InUrl, InOrg, InBucket, InToken, 10)
		go goroutines.GoCPU(srv, InUrl, InOrg, InBucket, InToken, 10)
		// 다른 고루틴
	}
	log.Println("listening on :80")
	log.Fatal(http.ListenAndServe(":80", router()))
}
