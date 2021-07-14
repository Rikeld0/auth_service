package config

import (
	"flag"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

var c = config{}

func init() {
	v := reflect.ValueOf(&c).Elem()
	t := v.Type()
	var st string
	for i := 0; i < t.NumField(); i++ {
		tagName := t.Field(i).Tag.Get("name")
		st = t.Field(i).Tag.Get("default")
		v.Field(i).SetString(st)
		flag.StringVar((*string)(unsafe.Pointer(v.Field(i).Addr().Pointer())),
			tagName,
			st,
			strings.Join([]string{
				tagName,
				st,
			}, "="))
	}
	fmt.Printf("config: %+v\n", c)
}

type config struct {
	DBHost    string `name:"PG_HOST" default:"127.0.0.1"`
	DBPort    string `name:"PG_PORT" default:"5432"`
	DBName    string `name:"PG_NAME" default:"test"`
	DBUser    string `name:"PG_USER" default:"postgres"`
	DBPass    string `name:"PG_PASS" default:"root"`
	RedisHost string `name:"REDIS_HOST" default:"127.0.0.1"`
	RedisPort string `name:"REDIS_PORT" default:"6379"`
	RedisPass string `name:"REDIS_PASS" default:""`
	RedisDB   string `name:"REDIS_DB" default:"0"`
	Host      string `name:"HOST" default:"127.0.0.1"`
	Port      string `name:"PORT" default:"8080"`
}

func ConnInfo() string {
	return `host= ` + c.DBHost + ` port = ` + c.DBPort + ` dbname = ` + c.DBName + ` user =` + c.DBUser + ` password = ` + c.DBPass + ` sslmode = disable`
}

func RedisAddr() string {
	return strings.Join([]string{c.RedisHost, c.RedisPort}, ":")
}

func RedisPass() string {
	return c.RedisPass
}

func RedisDB() int {
	db, _ := strconv.Atoi(c.RedisDB)
	return db
}

func HostServer() string {
	return strings.Join([]string{c.Host, c.Port}, ":")
}
