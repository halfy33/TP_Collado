package goroutines

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	memory2 "github.com/evoliatis/buildup/memory"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func GoMem(server, INFLUXDB_URL, INFLUXDB_ORG, INFLUXDB_BUCKET, INFLUXDB_TOKEN string, tick int) {
	client := influxdb2.NewClient(INFLUXDB_URL, INFLUXDB_TOKEN)
	defer client.Close()

	writeAPI := client.WriteAPIBlocking(INFLUXDB_ORG, INFLUXDB_BUCKET)
	ticker := time.NewTicker(time.Duration(tick) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		reqUrl := "http://" + server + "/mem"
		resp, err := http.Get(reqUrl)
		if err != nil {
			log.Println("Erreur de communication : " + reqUrl)
			continue
		}
		defer resp.Body.Close()

		var memory memory2.MemInfo
		if err := json.NewDecoder(resp.Body).Decode(&memory); err != nil {
			log.Println("Impossible de récupérer : " + reqUrl)
			continue
		}

		// Tags (par ex. host, env…)
		tags := map[string]string{
			"host": server,
		}

		// Champs numériques
		fields := map[string]interface{}{
			"total":      memory.Virtual.Total,
			"available":  memory.Virtual.Available,
			"used":       memory.Virtual.Used,
			"free":       memory.Virtual.Free,
			"cached":     memory.Virtual.Cached,
			"buffers":    memory.Virtual.Buffers,
			"swap_total": memory.Virtual.SwapTotal,
			"swap_free":  memory.Virtual.SwapFree,
		}

		p := influxdb2.NewPoint(
			"memory",
			tags,
			fields,
			time.Now(),
		)

		if err := writeAPI.WritePoint(context.Background(), p); err != nil {
			log.Printf("write failed: %v \n", err)
		} else {
			log.Println("Memory written to InfluxDB: " + reqUrl)
		}
	}
}
