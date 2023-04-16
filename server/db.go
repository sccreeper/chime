package main

import "gorm.io/gorm"

const (
	table_users     string = "users"
	table_playlists string = "playlists"
	table_tracks    string = "tracks"
	table_covers    string = "covers"
)

// Contains models for use with GORM.

type user_model struct {
	gorm.Model
	ID           int64 `gorm:"primaryKey"`
	Username     string
	Password     string
	Salt         []byte
	IsAdmin      int64
	SettingsJSON string
	Favourites   string
}

type playlist_model struct {
	gorm.Model
	ID      int64 `gorm:"primaryKey"`
	Name    string
	Cover   string
	IsAlbum int64
	Tracks  string
	Dates   string
	Owner   int64
}

type track_model struct {
	gorm.Model
	ID       int64 `gorm:"primaryKey"`
	Name     string
	Released int64
	Artist   string
	AlbumID  int64
	Duration int64
	Cover    int64
	Owner    int64
}

type cover_model struct {
	gorm.Model
	ID      int64 `gorm:"primaryKey"`
	AlbumID int64
	Path    string
	Owner   int64
}
