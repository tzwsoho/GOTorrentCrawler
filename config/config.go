package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	Port            uint16   `json:"port"`
	AutoPortMapping bool     `json:"auto_port_mapping"`
	DHTNodes        []string `json:"dht_nodes"`

	BlacklistSize    int `json:"blacklist_size"`
	RequestQueueSize int `json:"request_queue_size"`
	WorkerQueueSize  int `json:"worker_queue_size"`

	MySQLIP   string `json:"mysql_ip"`
	MySQLPort int    `json:"mysql_port"`
	MySQLUser string `json:"mysql_user"`
	MySQLPwd  string `json:"mysql_pwd"`
	MySQLDB   string `json:"mysql_db"`
}

var Cfg Config

func init() {
	Cfg = Config{
		Port:            6881,
		AutoPortMapping: true,

		// https://github.com/ngosang/trackerslist/blob/master/trackers_best.txt
		DHTNodes: []string{
			"router.bittorrent.com:6881",
			"router.utorrent.com:6881",
			"dht.transmissionbt.com:6881",
			"router.magnets.im:6881",
			"9.rarbg.to:2710",
			"9.rarbg.me:2710",
			"tracker.leechers-paradise.org:6969",
			"tracker.cyberia.is:6969",
			"exodus.desync.com:6969",
			"explodie.org:6969",
			"tracker3.itzmx.com:6961",
			"tracker.tiny-vps.com:6969",
			"open.stealth.si:80",
			"tracker.ds.is:6969",
			"tracker.torrent.eu.org:451",
			"retracker.lanta-net.ru:2710",
			"tracker.moeking.me:6969",
			"ipv4.tracker.harry.lu:80",
			"cdn-2.gamecoast.org:6969",
			"cdn-1.gamecoast.org:6969",
		},

		BlacklistSize:    65536,
		RequestQueueSize: 1024,
		WorkerQueueSize:  1024,

		MySQLIP:   "127.0.0.1",
		MySQLPort: 3306,
		MySQLUser: "root",
		MySQLPwd:  "123456",
		MySQLDB:   "torrents",
	}
}

func InitConfig() {
	if cfgData, err1 := ioutil.ReadFile("config.json"); nil == err1 {
		var cfgJSON Config
		if err2 := json.Unmarshal(cfgData, &cfgJSON); nil == err2 {
			Cfg = cfgJSON
		}
	}

	log.Printf("Torrent Crawler Config:\n%+v\n", Cfg)
}
