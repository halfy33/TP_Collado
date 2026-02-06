package goroutines

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/evoliatis/buildup/cpu"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func GoCPU(server, INFLUXDB_URL, INFLUXDB_ORG, INFLUXDB_BUCKET, INFLUXDB_TOKEN string, tick int) {

	// Connexion
	client := influxdb2.NewClient(INFLUXDB_URL, INFLUXDB_TOKEN)
	defer client.Close()

	writeAPI := client.WriteAPIBlocking(INFLUXDB_ORG, INFLUXDB_BUCKET)

	ticker := time.NewTicker(time.Duration(tick) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// Récupération des valeurs
		reqUrl := "http://" + server + "/cpu"
		resp, err := http.Get(reqUrl)
		if err != nil {
			log.Println("Erreur de communication : " + reqUrl)
			continue
		}
		defer resp.Body.Close()

		var cpu cpu.CPUInfo
		if err := json.NewDecoder(resp.Body).Decode(&cpu); err != nil {
			log.Println("Impossible de récupérer : " + reqUrl)
			continue
		}

		for _, c := range cpu.Cores {
			tags := map[string]string{ // tags (ex: host, env…)
				"cpu":    c.Times.CPU,     // ex "cpu0"
				"vendor": c.Info.VendorID, // utile pour multi-host
				"model":  c.Info.ModelName,
				"host":   server,
			}
			fields := map[string]interface{}{
				"cpu": c.Usage,

				// Tous les champs TimesStat (numériques)
				"user":       c.Times.User,
				"system":     c.Times.System,
				"idle":       c.Times.Idle,
				"nice":       c.Times.Nice,
				"iowait":     c.Times.Iowait,
				"irq":        c.Times.Irq,
				"softirq":    c.Times.Softirq,
				"steal":      c.Times.Steal,
				"guest":      c.Times.Guest,
				"guest_nice": c.Times.GuestNice,
			}

			// ---- point ----
			p := influxdb2.NewPoint(
				"cpu",
				tags,
				fields,
				time.Now(),
			)

			// Ecriture
			if err := writeAPI.WritePoint(context.Background(), p); err != nil {
				log.Printf("write failed: %v \n", err)
			}
		}

		log.Println("CPU written to InfluxDB :" + reqUrl)
	}
}
