package goroutines

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/evoliatis/buildup/load"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func GoLoad(server, INFLUXDB_URL, INFLUXDB_ORG, INFLUXDB_BUCKET, INFLUXDB_TOKEN string, tick int) {

	// Connexion
	client := influxdb2.NewClient(INFLUXDB_URL, INFLUXDB_TOKEN)
	defer client.Close()

	writeAPI := client.WriteAPIBlocking(INFLUXDB_ORG, INFLUXDB_BUCKET)

	ticker := time.NewTicker(time.Duration(tick) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// Récupération des valeurs
		reqUrl := "http://" + server + "/load"
		resp, err := http.Get(reqUrl)
		if err != nil {
			log.Println("Erreur de communication : " + reqUrl)
			continue
		}
		defer resp.Body.Close()

		var avg load.LoadInfo
		if err := json.NewDecoder(resp.Body).Decode(&avg); err != nil {
			log.Println("Impossible de récupérer : " + reqUrl)
			continue
		}

		// ---- point ----
		p := influxdb2.NewPoint(
			"load_average",
			map[string]string{ // tags (ex: host, env…)
				"host": server,
			},
			map[string]interface{}{ // Données à stocker
				"avg1":  avg.Avg.Load1,
				"avg5":  avg.Avg.Load5,
				"avg15": avg.Avg.Load15,
			},
			time.Now(),
		)

		// Ecriture
		if err := writeAPI.WritePoint(context.Background(), p); err != nil {
			log.Printf("write failed: %v\n", err)
		}

		log.Println("load average written to InfluxDB : ", reqUrl)
	}
}
