package goroutines

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/evoliatis/buildup/disk"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func GoDisk(server, INFLUXDB_URL, INFLUXDB_ORG, INFLUXDB_BUCKET, INFLUXDB_TOKEN string, tick int) {
	client := influxdb2.NewClient(INFLUXDB_URL, INFLUXDB_TOKEN)
	defer client.Close()

	writeAPI := client.WriteAPIBlocking(INFLUXDB_ORG, INFLUXDB_BUCKET)
	ticker := time.NewTicker(time.Duration(tick) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		reqUrl := "http://" + server + "/disk"
		resp, err := http.Get(reqUrl)
		if err != nil {
			log.Println("Erreur de communication : " + reqUrl)
			continue
		}
		defer resp.Body.Close()

		var disks []disk.DiskFS
		if err := json.NewDecoder(resp.Body).Decode(&disks); err != nil {
			log.Println("Impossible de récupérer : " + reqUrl)
			continue
		}

		for _, d := range disks {
			if d.Usage == nil {
				log.Println("Skipping nil disk info for mountpoint")
				continue
			}

			tags := map[string]string{
				"host":       server,
				"mountpoint": d.Partition.Mountpoint,
			}

			fields := map[string]interface{}{
				"total":       d.Usage.Total,
				"used":        d.Usage.Used,
				"free":        d.Usage.Free,
				"usedPercent": d.Usage.UsedPercent,
			}

			p := influxdb2.NewPoint(
				"disk",
				tags,
				fields,
				time.Now(),
			)

			if err := writeAPI.WritePoint(context.Background(), p); err != nil {
				log.Printf("write failed: %v \n", err)
			} else {
				log.Println("Disk written to InfluxDB: " + d.Partition.Mountpoint)
			}
		}

	}
}
