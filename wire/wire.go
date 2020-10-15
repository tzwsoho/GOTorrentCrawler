package wire

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/shiyanhui/dht"
)

type FileInfo struct {
	Path   []interface{} `json:"path"`
	Length uint64        `json:"length"`
}

type MetadataInfo struct {
	InfoHash string     `json:"infohash"`
	Name     string     `json:"name"`
	Files    []FileInfo `json:"files,omitempty"`
	Length   uint64     `json:"length,omitempty"`
}

var W *dht.Wire

func InitWire(blacklistSize, requestQueueSize, workerQueueSize int, dbName string, db *sql.DB) {
	stmt, err := db.Prepare(fmt.Sprintf("INSERT `%s`.`infohashes`"+
		"(`info_hash`, `info_name`, `total_length`, `total_files`, `files`, `hot`) "+
		"VALUES (?, ?, ?, ?, ?, 1) "+
		"ON DUPLICATE KEY UPDATE "+
		"`updated_at` = ?, `hot` = `hot` + 1", dbName),
	)
	if nil != err {
		log.Panicf("InitWire Prepare err: %s\n", err.Error())
	}

	W = dht.NewWire(blacklistSize, requestQueueSize, workerQueueSize)
	go getMetadataInfo(stmt)
	go W.Run()
}

func getMetadataInfo(stmt *sql.Stmt) {
	for resp := range W.Response() {
		metadata, err := dht.Decode(resp.MetadataInfo)
		if nil != err {
			log.Printf("getMetadataInfo Decode err: %s\n", err.Error())
			continue
		}

		info := metadata.(map[string]interface{})
		infoName, ok := info["name"]
		if !ok {
			continue
		}

		h := hex.EncodeToString(resp.InfoHash)
		meta := MetadataInfo{
			InfoHash: h,
			Name:     infoName.(string),
		}
		log.Printf("%s:%d magnet:?xt=urn:btih:%s %s\n", resp.IP, resp.Port, h, infoName.(string))

		if length, ok := info["length"]; ok {
			meta.Length = uint64(length.(int))
		}

		if fls, ok := info["files"]; ok {
			files := fls.([]interface{})
			meta.Files = make([]FileInfo, len(files))

			for i, f := range files {
				fileInfo := f.(map[string]interface{})
				meta.Files[i] = FileInfo{
					Path:   fileInfo["path"].([]interface{}),
					Length: uint64(fileInfo["length"].(int)),
				}
			}
		}

		if filesJSON, err := json.Marshal(meta.Files); nil != err {
			log.Printf("getMetadataInfo Marshal err: %s\n", err.Error())
		} else {
			if _, err := stmt.Exec(meta.InfoHash, meta.Name, meta.Length, len(meta.Files), string(filesJSON), time.Now()); nil != err {
				log.Printf("getMetadataInfo Exec err: %s\n", err.Error())
			}
		}
	}
}
