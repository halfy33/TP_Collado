package goroutines

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	_ "github.com/shirou/gopsutil/v4/net"
)

type NetInfo struct {
	Name        string `json:"name"`
	BytesSent   uint64 `json:"bytesSent"`
	BytesRecv   uint64 `json:"bytesRecv"`
	PacketsSent uint64 `json:"packetsSent"`
	PacketsRecv uint64 `json:"packetsRecv"`
	Errin       uint64 `json:"errin"`
	Errout      uint64 `json:"errout"`
	Dropin      uint64 `json:"dropin"`
	Dropout     uint64 `json:"dropout"`
}

func GoNet(server, INFLUXDB_URL, INFLUXDB_ORG, INFLUXDB_BUCKET, INFLUXDB_TOKEN string, tick int) {
	client := influxdb2.NewClient(INFLUXDB_URL, INFLUXDB_TOKEN)
	defer client.Close()

	writeAPI := client.WriteAPIBlocking(INFLUXDB_ORG, INFLUXDB_BUCKET)
	ticker := time.NewTicker(time.Duration(tick) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		reqUrl := "http://" + server + "/net"
		resp, err := http.Get(reqUrl)
		if err != nil {
			log.Println("Erreur de communication : " + reqUrl)
			continue
		}
		defer resp.Body.Close()

		var nets []NetInfo
		if err := json.NewDecoder(resp.Body).Decode(&nets); err != nil {
			log.Println("Impossible de récupérer : " + reqUrl)
			continue
		}

		for _, n := range nets {
			tags := map[string]string{
				"host":      server,
				"interface": n.Name,
			}

			fields := map[string]interface{}{
				"bytesSent":   n.BytesSent,
				"bytesRecv":   n.BytesRecv,
				"packetsSent": n.PacketsSent,
				"packetsRecv": n.PacketsRecv,
				"errin":       n.Errin,
				"errout":      n.Errout,
				"dropin":      n.Dropin,
				"dropout":     n.Dropout,
			}

			p := influxdb2.NewPoint(
				"net",
				tags,
				fields,
				time.Now(),
			)

			if err := writeAPI.WritePoint(context.Background(), p); err != nil {
				log.Printf("write failed: %v \n", err)
			} else {
				log.Println("Net written to InfluxDB: " + n.Name)
			}
		}
	}
}
