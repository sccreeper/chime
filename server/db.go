package main

import "gorm.io/gorm"

const (
	table_users     string = "users"
	table_playlists string = "playlists"
	table_tracks    string = "tracks"
	table_covers    string = "covers"
	table_radio     string = "radio"
)

// Contains models for use with GORM.

type user_model struct {
	gorm.Model
	ID           int64 `gorm:"primaryKey"`
	Username     string
	Password     string
	Salt         []byte
	IsAdmin      int64  //1 = Admin, 0 = Not admin
	SettingsJSON string //JSON object containing settings
	Favourites   string //Comma seprated list of favourite tracks hex ids
}

type playlist_model struct {
	gorm.Model
	ID          int64  `gorm:"primaryKey"`
	Name        string //Name of this collection
	Description string
	Cover       int64
	IsAlbum     int64  //If the collection is an album or not
	Tracks      string //Comma seperated list of hexadecimal track IDs.
	Dates       string //Comma seperated list of dates added (as unix time hex)
	Owner       int64
}

type track_model struct {
	gorm.Model
	ID       int64   `gorm:"primaryKey"`
	Name     string  // Name of track
	Released int64   //Year released
	Artist   string  //Arist
	AlbumID  int64   //The album this track belongs to
	Duration float64 //The duration of the track in seconds
	Cover    int64
	Owner    int64
	Original string //The original file name, used for streaming and downloading
	Size     int64  //The size of the original file in bytes
	Position int64  //Position of a track in album.
}

type radio_model struct {
	gorm.Model
	ID          int64 `gorm:"primaryKey"`
	Name        string
	CoverID     int64
	URL         string
	Description string
	Owner       int64
}

type cover_model struct {
	gorm.Model
	ID      int64 `gorm:"primaryKey"`
	AlbumID int64 //Album this cover belongs to, if custom will be 0
	Owner   int64
}

// https://yourbasic.org/golang/how-to-sort-in-go/
type by_position []track_model

func (t by_position) Len() int {
	return len(t)
}

func (t by_position) Less(i int, j int) bool {
	return t[i].Position < t[j].Position
}

func (t by_position) Swap(i int, j int) {
	t[i], t[j] = t[j], t[i]
}
