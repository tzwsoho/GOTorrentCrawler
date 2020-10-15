package upnp

import (
	"log"
	"net"
	"strings"

	"github.com/huin/goupnp/dcps/internetgateway1"
)

const dhtProtocol string = "udp"
const portName string = "TorrentCrawler"

func getLocalIP() string {
	if interfaces, err := net.Interfaces(); nil != err {
		log.Panicf("getLocalIP Interfaces err: %s", err.Error())
	} else {
		for _, inf := range interfaces {
			// https://blog.csdn.net/hzj_001/article/details/81587824
			if 0 == (inf.Flags&net.FlagUp) || ("本地连接" != inf.Name && "以太网" != inf.Name &&
				!strings.Contains(strings.ToLower(inf.Name), "wl") && // WLAN
				!strings.Contains(strings.ToLower(inf.Name), "ww") && // 4G
				!strings.Contains(strings.ToLower(inf.Name), "en") && // Ethernet
				!strings.Contains(strings.ToLower(inf.Name), "em") && // 板载网卡
				!strings.Contains(strings.ToLower(inf.Name), "eth")) {
				continue
			}

			if addrs, err := inf.Addrs(); nil != err {
				log.Panicf("getLocalIP Addrs err: %s", err.Error())
			} else {
				for _, addr := range addrs {
					ip := addr.(*net.IPNet)
					if nil != ip && !ip.IP.IsLoopback() && ip.IP.IsGlobalUnicast() && nil != ip.IP.To4() {
						return ip.IP.To4().String()
					}
				}
			}
		}
	}

	return ""
}

func PortMapping(port uint16) (intPort uint16) {
	intPort = port
	localIP := getLocalIP()
	if "" == localIP {
		log.Panicf("PortMapping getLocalIP empty!")
	}

	if clients, _, err := internetgateway1.NewWANIPConnection1Clients(); nil != err {
		log.Panicf("PortMapping NewWANIPConnection1Clients err: %s", err.Error())
	} else {
		for _, client := range clients {
			// 检查外网端口映射是否已存在
			var isCreated bool = false
			var extPort uint16 = 6881
			for {
				var isExists bool = false
				for i := uint16(0); ; i++ {
					if _, externalPort, protocol, internalPort, internalClient, enabled, _, _, err := client.GetGenericPortMappingEntry(i); nil != err {
						break
					} else {
						protocol = strings.ToLower(protocol)
						if extPort == externalPort && dhtProtocol == protocol && intPort == internalPort && localIP == internalClient && enabled {
							isCreated = true
							break
						}

						if extPort == externalPort && dhtProtocol == protocol {
							extPort++
							isExists = true
							break
						}

						if intPort == internalPort && dhtProtocol == protocol {
							intPort++
							isExists = true
							break
						}
					}
				}

				if isCreated {
					break
				}

				// 不存在映射，创建新的端口映射
				if !isExists {
					if err := client.AddPortMapping("", extPort, dhtProtocol, intPort, localIP, true, portName, 0); nil != err {
						log.Printf("PortMapping AddPortMapping err: %s\n", err.Error())
					}

					isCreated = true
					break
				}
			}

			if isCreated {
				break
			}
		}
	}

	return intPort
}
