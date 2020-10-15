package dht

import (
	"fmt"
	"log"
	"time"

	"github.com/shiyanhui/dht"
)

func Crawle(port uint16, dhtNodes []string, w *dht.Wire) {
	config := dht.NewCrawlConfig()
	config.Address = fmt.Sprintf(":%d", port)
	config.PrimeNodes = dhtNodes
	config.OnAnnouncePeer = func(infoHash, ip string, port int) {
		w.Request([]byte(infoHash), ip, port)
	}

	d := dht.New(config)
	go d.Run()
}

func TestWire(w *dht.Wire) {
	d := dht.New(nil)
	d.OnGetPeersResponse = func(infoHash string, peer *dht.Peer) {
		// fmt.Printf("%s:%d\n", peer.IP, peer.Port)
		w.Request([]byte(infoHash), peer.IP.String(), peer.Port)
	}

	go func() {
		for {
			// ubuntu-14.04.2-desktop-amd64.iso
			err := d.GetPeers("3D8B16242B56A3AAFB8DA7B5FC83EF993EBCF35B")
			if err != nil && err != dht.ErrNotReady {
				log.Fatal(err)
			}

			if err == dht.ErrNotReady {
				time.Sleep(time.Second * 1)
				continue
			}

			break
		}
	}()

	go d.Run()
}
