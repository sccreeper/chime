class Library {

  final List<LibraryItem> albums;
  final List<LibraryItem> playlists;
  final List<LibraryItem> radios;

  Library({
    required this.albums,
    required this.playlists,
    required this.radios,
  });

  factory Library.fromJSON(Map<String, dynamic> json) => Library(
    albums: (json["albums"] as List<dynamic>).map((x) => LibraryItem.fromJSON(x)).toList(), 
    playlists: (json["playlists"] as List<dynamic>).map((x) => LibraryItem.fromJSON(x)).toList(), 
    radios: (json["radios"] as List<dynamic>).map((x) => LibraryItem.fromJSON(x)).toList()
  );

  Map<String, dynamic> toJson() => {
    "albums":albums.map((e) => e.toJson()),
    "playlists":playlists.map((e) => e.toJson()),
    "radios":playlists.map((e) => e.toJson())
  };

}

class LibraryItem {
  
  final String id;
  final String name;

  LibraryItem({
    required this.id,
    required this.name,
  });

  factory LibraryItem.fromJSON(Map<String, dynamic> json) => LibraryItem(
    id: json["id"], 
    name: json["name"]
  );

  Map<String, dynamic> toJson() => {
    "id":id,
    "name":name
  };

}

enum LibaryItemType {album, playlist, radio}

// Collections

class Collection {

  final String title;
  final String coverId;
  final bool isAlbum;
  final List<Track> tracks;
  final List<String> dates;
  final String description;
  final bool protected;

  Collection({
    required this.title,
    required this.coverId,
    required this.isAlbum,
    required this.tracks,
    required this.dates,
    required this.description,
    required this.protected
  });

  factory Collection.fromJSON(Map<String, dynamic> json) => Collection(
    title: json["title"], 
    coverId: json["cover"], 
    isAlbum: json["is_album"], 
    tracks: (json["tracks"] as List<dynamic>).map((e) => Track.fromJSON(e)).toList(), 
    dates: (json["dates"] as List<dynamic>).map((e) => e as String).toList(), 
    description: json["description"], 
    protected: json["protected"]
  );

  Map<String, dynamic> toJson() => {
    "title":title,
    "cover":coverId,
    "is_album":isAlbum,
    "tracks":tracks.map((e) => e.toJson()).toList(),
    "dates":dates,
    "description":description,
    "protected":protected,
  };

}

class Track {

  final String id;
  final String name;
  final String albumName;
  final int released;
  final String artist;
  final String albumId;
  final double duration;
  final String coverId;
  final int position;

  Track({
    required this.id,
    required this.name,
    required this.albumName,
    required this.released,
    required this.artist,
    required this.albumId,
    required this.duration,
    required this.coverId,
    required this.position,
  });

  factory Track.fromJSON(Map<String, dynamic> json) => Track(
    id: json["id"], 
    name: json["name"], 
    albumName: json["album_name"], 
    released: json["released"], 
    artist: json["artist"], 
    albumId: json["album_id"], 
    duration: json["duration"], 
    coverId: json["cover_id"],
    position: json["position"]
  );

  Map<String, dynamic> toJson() => {
    "id":id,
    "name":name,
    "album_name":name,
    "released":released,
    "artist":artist,
    "album_id":albumId,
    "duration":duration,
    "cover_id":coverId,
    "position":position
  };

  factory Track.fromMetadata(String id, TrackMetadata metadata) => Track(
    id: id, 
    name: metadata.title, 
    albumName: metadata.albumName, 
    released: metadata.released, 
    artist: metadata.artist, 
    albumId: metadata.albumId, 
    duration: metadata.duration, 
    coverId: metadata.coverId,
    position: metadata.position
  );

}

class TrackMetadata {

  final String title; //Title of the track
  final String albumName; //The name of the album this track belongs to
  final String albumId; //The ID of the album this track belongs to
  final String coverId; //The ID of the cover for the album this track belongs to
  final String artist; //The artist of this track
  final String originalFile; //The original file name for this track
  final String format; //The format of the track e.g. FLAC, MP3, WAV etc.
  final double duration; //The duration of the track in seconds
  final int released; //The year the track was released
  final int size; //The size of the track in bytes
  final int position;

  TrackMetadata({
    required this.title,
    required this.albumName,
    required this.albumId,
    required this.coverId,
    required this.artist,
    required this.originalFile,
    required this.format,
    required this.duration,
    required this.released,
    required this.size,
    required this.position,
  });

  factory TrackMetadata.fromJSON(Map<String, dynamic> json) => TrackMetadata(
    title: json["title"], 
    albumName: json["album_name"], 
    albumId: json["album_id"], 
    coverId: json["cover_id"], 
    artist: json["artist"], 
    originalFile: json["original_file"], 
    format: json["format"], 
    duration: json["duration"], 
    released: json["released"], 
    size: json["size"],
    position: json["position"]
  );

  Map<String, dynamic> toJson() => {
    "title":title,
    "album_name":albumName,
    "album_id":albumId,
    "cover_id":coverId,
    "artist":artist,
    "original_file":originalFile,
    "format":format,
    "duration":duration,
    "released":released,
    "size":size,
    "position":position
  };

}