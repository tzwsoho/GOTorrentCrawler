package mysqlmanager

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql" // 需要用到 MySQL 驱动
)

type InfoHashes struct {
	InfoHash    string    `json:"info_hash" def:"varchar(40) NOT NULL DEFAULT '' COMMENT '信息哈希'" unique:"info_hash"`
	InfoName    string    `json:"info_name" def:"varbinary(4096) NOT NULL DEFAULT '' COMMENT '名称'"`
	TotalLength uint64    `json:"total_length" def:"bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '所有文件总大小（字节）'"`
	TotalFiles  uint64    `json:"total_files" def:"bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '文件总数量'"`
	Files       []string  `json:"files" def:"longblob NOT NULL COMMENT '文件列表'"`
	Hot         uint64    `json:"hot" def:"bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '热度'"`
	CreatedAt   time.Time `json:"created_at" def:"datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '记录添加日期'"`
	UpdatedAt   time.Time `json:"updated_at" def:"datetime NOT NULL COMMENT '记录最后更新日期'"`
}

var DB *sql.DB

func InitDB(user, pwd, ip, db string, port int) {
	var err error
	DB, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, pwd, ip, port, db))
	if nil != err {
		panic(err)
	}

	DB.SetConnMaxLifetime(0)
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(10)

	time.Sleep(time.Second)

	createDBIfNotExists(db)
	createTableIfNotExists(db, InfoHashes{})
}

func CloseDB() {
	if nil != DB {
		DB.Close()
	}
}

func createDBIfNotExists(db string) {
	var database string
	row := DB.QueryRow("SHOW DATABASES WHERE `Database` = ?", db)
	if err := row.Scan(&database); nil != err {
		DB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4", db))
	}

	if _, err := DB.Exec(fmt.Sprintf("USE `%s`", db)); nil != err {
		log.Panicf("createDBIfNotExists Exec err: %s", err.Error())
	}
}

func createTableIfNotExists(db string, table interface{}) {
	var tableName, createSQL string
	var unique map[string][]string = make(map[string][]string)

	t := reflect.TypeOf(table)
	row := DB.QueryRow(fmt.Sprintf("SHOW TABLES FROM `%s` WHERE `Tables_in_%s` = ?", db, db), strings.ToLower(t.Name()))
	if err := row.Scan(&tableName); nil != err {
		createSQL = fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s` (\n", strings.ToLower(t.Name()))
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			createSQL += fmt.Sprintf("`%s` %s,\n", f.Tag.Get("json"), f.Tag.Get("def"))

			if v, ok := f.Tag.Lookup("unique"); ok && v != "" {
				if nil == unique[v] {
					unique[v] = make([]string, 0)
				}

				unique[v] = append(unique[v], v)
			}
		}

		for k, v := range unique {
			createSQL += fmt.Sprintf("UNIQUE KEY `%s` (`%s`) USING BTREE\n", k, strings.Join(v, "`, `"))
		}

		createSQL += ") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4"

		if _, err := DB.Exec(createSQL); nil != err {
			log.Panicf("createTableIfNotExists Exec %s err: %s", createSQL, err.Error())
		}
	}
}
