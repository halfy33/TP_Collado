package goroutines

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/evoliatis/buildup/proc"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func GoProc(server, INFLUXDB_URL, INFLUXDB_ORG, INFLUXDB_BUCKET, INFLUXDB_TOKEN string, tick int) {

	// Connexion
	client := influxdb2.NewClient(INFLUXDB_URL, INFLUXDB_TOKEN)
	defer client.Close()

	writeAPI := client.WriteAPIBlocking(INFLUXDB_ORG, INFLUXDB_BUCKET)

	ticker := time.NewTicker(time.Duration(tick) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// Récupération des valeurs
		reqUrl := "http://" + server + "/ps"
		resp, err := http.Get(reqUrl)
		if err != nil {
			log.Println("Erreur de communication : " + reqUrl)
			continue
		}
		defer resp.Body.Close()

		var procs []proc.Proc
		if err := json.NewDecoder(resp.Body).Decode(&procs); err != nil {
			log.Println("Impossible de récupérer : " + reqUrl)
			continue
		}

		for _, proc := range procs {

			tags := map[string]string{ // tags (ex: host, env…)
				"host": "linux",
				"name": proc.Name,
				"pid":  strconv.FormatInt(int64(proc.Pid), 10),
			}
			fields := map[string]interface{}{ // Données à stocker
				"CPU":    proc.CPU,
				"Memory": proc.Memory,
				"Status": proc.Status,
			}
			// ---- points ----
			p := influxdb2.NewPoint(
				"processes",
				tags,
				fields,
				time.Now(),
			)

			// Ecriture
			if err := writeAPI.WritePoint(context.Background(), p); err != nil {
				log.Printf("write failed: %v \n", err)
			}
		}

		log.Println("(", len(procs), ") procs written to InfluxDB : "+reqUrl)
	}
}
