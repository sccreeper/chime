package main

import (
	"math"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

type config struct {
	Port             uint `toml:"port" json:"port"`
	CastProxyEnabled bool `toml:"cast_proxy" json:"cast_proxy"`
}

var random *rand.Rand
var db_path string = "/var/lib/chime/data.db"
var database *gorm.DB

func init() {
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func generate_id(table string) int64 {
	var id_found bool
	var id int64
	var count int64

	for !id_found {

		id = random.Int63n(math.MaxInt64)
		database.Table(table).Where("id = ?", id).Count(&count)

		if count == 0 {
			id_found = true
		}

	}

	return id

}
