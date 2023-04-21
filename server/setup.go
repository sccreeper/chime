package main

import (
	"encoding/base64"
	"encoding/binary"
	"log"
	"math"
	"os"

	"github.com/pelletier/go-toml/v2"
	"golang.org/x/crypto/scrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var pass_chars string = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890@!"

var default_config_toml config = config{
	Port: 80,
}

// Creates default db and configuration
func first_time_setup(admin_pass string, in_docker bool) {

	log.Println("Creating directory /var/lib/chime...")

	if !in_docker {
		err := os.Mkdir("/var/lib/chime", os.ModeDir)
		if err != nil {
			log.Println("There was an error with the setup process:")
			log.Println(err)
			os.Exit(1)
		}
	}

	// Create the database and tables

	log.Println("Creating DB...")

	os.Create(db_path)

	db, err := gorm.Open(sqlite.Open("file:/var/lib/chime/data.db"), &gorm.Config{})
	if err != nil {
		log.Println("There was an error with the setup process:")
		log.Println(err)
		os.Exit(1)
	}

	log.Println("Creating tables...")

	db.Table(table_users).AutoMigrate(&user_model{})
	db.Table(table_playlists).AutoMigrate(&playlist_model{})
	db.Table(table_tracks).AutoMigrate(&track_model{})
	db.Table(table_covers).AutoMigrate(&cover_model{})

	log.Println("Creating admin user account...")

	// Hash password

	var salt uint64 = random.Uint64()
	var salt_bytes []byte = make([]byte, 8)
	binary.LittleEndian.PutUint64(salt_bytes, salt)

	hash, err := scrypt.Key([]byte(admin_pass), salt_bytes, 1<<15, 8, 1, 64)
	if err != nil {
		log.Println("There was an error with the setup process:")
		log.Println(err)
		os.Exit(1)
	}

	log.Println("Initalizing DB values...")

	// Add initial values to DB.

	var id_found bool
	var admin_id int64
	var count int64

	for !id_found {

		admin_id = random.Int63n(math.MaxInt64)
		db.Table(table_users).Where("id = ?", admin_id).Count(&count)

		if count == 0 {
			id_found = true
		}

	}

	db.Table(table_users).Create(&user_model{
		ID:           admin_id,
		Username:     "admin",
		Password:     base64.StdEncoding.EncodeToString(hash),
		Salt:         salt_bytes,
		IsAdmin:      1,
		SettingsJSON: "",
		Favourites:   "",
	})

	db.Table(table_playlists).Create(&playlist_model{
		ID:          1,
		Name:        "Unsorted",
		IsAlbum:     1,
		Cover:       0,
		Tracks:      "",
		Dates:       "",
		Owner:       admin_id,
		Description: "Default album for tracks that don't have any metadata.",
	})

	// Create default config.toml

	log.Println("Creating default config...")

	f, err := os.Create("/var/lib/chime/config.toml")
	if err != nil {
		log.Println("There was an error with the setup process:")
		log.Println(err)
		os.Exit(1)
	}

	toml.NewEncoder(f).Encode(default_config_toml)
	f.Close()

	// Create data directories

	log.Println("Creating data directories...")

	os.Mkdir("/var/lib/chime/tracks", os.ModeDir)
	os.Mkdir("/var/lib/chime/covers", os.ModeDir)

	log.Println("Setup process successful!")

}

// Generates random password
func generate_password(length int) string {

	var pass string

	for i := 0; i < length; i++ {
		pass += string(pass_chars[random.Intn(len(pass_chars))])
	}

	return pass

}
