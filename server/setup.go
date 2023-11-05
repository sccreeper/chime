package main

import (
	_ "embed"
	"encoding/base64"
	"encoding/binary"
	"log"
	"math"
	"os"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"golang.org/x/crypto/scrypt"
	"golang.org/x/exp/slices"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var pass_chars string = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890@!"

var default_config_toml config = config{
	Port:             80,
	CastProxyEnabled: false,
}

// Words for creating unique user friendly server ID.

//go:embed words.txt
var words string

var words_list []string
var words_blacklist []string = []string{"fuck", "shit", "cunt", "cum", "cumshot", "cumshots", "cock", "cocks", "penis", "dick", "dicks", "dildo"}

func init() {

	words_list = strings.Split(words, "\n")

}

const server_id_length int = 4

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
	db.Table(table_radio).AutoMigrate(&radio_model{})

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

	var admin_id int64 = random.Int63n(math.MaxInt64)
	var unsorted_id int64 = random.Int63n(math.MaxInt64)

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
		ID:          unsorted_id,
		Name:        "Unsorted",
		IsAlbum:     1,
		Cover:       0,
		Tracks:      "",
		Dates:       "",
		Owner:       admin_id,
		Description: "Default album for tracks that don't have any metadata.",
		Protected:   1,
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
	os.Mkdir("/var/lib/chime/backups", os.ModeDir)
	os.Mkdir("/var/lib/chime/cache", os.ModeDir)

	// Create unique server ID

	uuid_file, err := os.Create("/var/lib/chime/server_id.txt")
	if err != nil {
		log.Println("There was an error with the setup process:")
		log.Println(err)
		os.Exit(1)
	}

	var server_id []string

	for i := 0; i < server_id_length; i++ {

		word := ""

		for {

			word = words_list[random.Intn(len(words_list))]
			if slices.Contains(words_blacklist, word) {
				continue
			} else {
				break
			}

		}

		server_id = append(server_id, word)

	}

	uuid_file.Write([]byte(strings.Join(server_id, "-")))
	uuid_file.Close()

	log.Printf("Server ID: %s\n", strings.Join(server_id, "-"))

	log.Println("Setup process successful!")

}
