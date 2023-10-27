class SearchResults {
  final List<SearchTrack> tracks;
  final List<SearchCollection> collections;
  final List<SearchRadio> radios;

  SearchResults({
    required this.tracks,
    required this.collections,
    required this.radios,
  });

  factory SearchResults.fromJson(Map<String, dynamic> json) => SearchResults(
    tracks: (json["tracks"] as List<dynamic>).map((e) => SearchTrack.fromJson(e)).toList(), 
    collections: (json["collections"] as List<dynamic>).map((e) => SearchCollection.fromJson(e)).toList(), 
    radios: (json["radios"] as List<dynamic>).map((e) => SearchRadio.fromJson(e)).toList()
  );
  
  Map<String, dynamic> toJson() => {
    "tracks":tracks.map((e) => e.toJson()).toList(),
    "collections":tracks.map((e) => e.toJson()).toList(),
    "radios":tracks.map((e) => e.toJson()).toList()
  };

}

class SearchTrack {

  final String id;
  final String albumId;
  final String artist;
  final String title;
  final double duration;
  final String coverId;

  SearchTrack({
    required this.id,
    required this.albumId,
    required this.artist,
    required this.title,
    required this.duration,
    required this.coverId,
  });

  factory SearchTrack.fromJson(Map<String, dynamic> json) => SearchTrack(
    id: json["id"], 
    albumId: json["album_id"], 
    artist: json["artist"], 
    title: json["title"], 
    duration: json["duration"], 
    coverId: json["cover"]
  );

  Map<String, dynamic> toJson() => {
    "id":id,
    "abum_id":albumId,
    "artist":artist,
    "title":title,
    "duration":duration,
    "cover":coverId,
  };

}

class SearchCollection {

  final String id;
  final String title;
  final String coverId;
  final bool isAlbum;

  SearchCollection({
    required this.id,
    required this.title,
    required this.coverId,
    required this.isAlbum,
  });

  factory SearchCollection.fromJson(Map<String, dynamic> json) => SearchCollection(
    id: json["id"], 
    title: json["title"], 
    coverId: json["cover"], 
    isAlbum: json["is_album"]
  );

  Map<String, dynamic> toJson() => {
    "id":id,
    "title":title,
    "cover":coverId,
    "is_album":isAlbum
  };

}

class SearchRadio {

  final String id;
  final String name;
  final String coverId;

  SearchRadio({
    required this.id,
    required this.name,
    required this.coverId
  });

  factory SearchRadio.fromJson(Map<String, dynamic> json) => SearchRadio(
    id: json["id"], 
    name: json["name"], 
    coverId: json["cover"]
  );

  Map<String, dynamic> toJson() => {
    "id":id,
    "name":name,
    "cover":coverId
  };

}