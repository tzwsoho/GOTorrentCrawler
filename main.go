package main

import (
	"TorrentCrawler/config"
	"TorrentCrawler/dht"
	"TorrentCrawler/mysqlmanager"
	"TorrentCrawler/upnp"
	"TorrentCrawler/wire"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func initSignalHandler() {
	signalChan := make(chan os.Signal, 2)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		select {
		case <-signalChan:
			mysqlmanager.CloseDB()
			log.Println("GoodBye!")
			os.Exit(0)
		}
	}()
}

func main() {
	// 使 Ctrl + C 可以关闭程序
	initSignalHandler()

	// 读取配置
	config.InitConfig()

	// 初始化数据库
	mysqlmanager.InitDB(config.Cfg.MySQLUser, config.Cfg.MySQLPwd, config.Cfg.MySQLIP, config.Cfg.MySQLDB, config.Cfg.MySQLPort)
	defer mysqlmanager.CloseDB()

	// 在路由器上创建端口映射
	var port uint16 = config.Cfg.Port
	if config.Cfg.AutoPortMapping {
		port = upnp.PortMapping(port)
	}

	// 初始化获取 torrent 文件元信息
	wire.InitWire(config.Cfg.BlacklistSize, config.Cfg.RequestQueueSize, config.Cfg.WorkerQueueSize, config.Cfg.MySQLDB, mysqlmanager.DB)

	// 开始爬取 DHT 网络
	dht.Crawle(port, config.Cfg.DHTNodes, wire.W)

	// dht.TestWire(wire.W)

	select {}
}
