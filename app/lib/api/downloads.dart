import 'dart:io';

import 'package:app/api/api.dart';
import 'package:app/api/endpoints.dart';
import 'package:app/api/models/collections.dart';
import 'package:flutter/material.dart';
import 'package:get_it/get_it.dart';
import 'package:http/http.dart' as http;
import 'package:app/shared.dart';
import 'package:path_provider/path_provider.dart';
import 'package:sqflite/sqflite.dart';

const String tableCollections = "collections";
const String tableTracks = "tracks";
const String downloadDirectory = "downloads";
const String coverDownloadDirectory = "$downloadDirectory/covers";
const String trackDownloadDirectory = "$downloadDirectory/tracks";

class DownloadManager {

  // Cancels the currently running download.
  static void cancelDownload() {

  }

  static void downloadCollection(String id) async {
    String basePath = (await getApplicationDocumentsDirectory()).path;

    log.fine("Downloading collection: $id");

    GetIt.I<DownloadNotifier>().downloading = true;
    GetIt.I<DownloadNotifier>().downloadingId = id;
    GetIt.I<DownloadNotifier>().downloadType = DownloadType.collection;

    // Get collection items.

    Collection collection = await ChimeAPI.getCollection(id);

    GetIt.I<DownloadNotifier>().downloadingTotal = collection.tracks.length;
    GetIt.I<DownloadNotifier>().downloadingProgress = 0;

    for (Track track in collection.tracks) {

      GetIt.I<DownloadNotifier>().downloadingProgress++;

      // Check if track and cover have already been downloaded to save us adding them to the db again.
      if (!File("$basePath/$coverDownloadDirectory/${track.coverId}").existsSync()) {
        log.fine("Downloading cover: ${track.coverId}");

        var coverResult = await http.get(Uri.parse(
            "${session.serverOrigin}$apiGetCover/${track.coverId}"), headers: Util.genAuthHeaders());

        File coverFile = File("$basePath/$coverDownloadDirectory/${track.coverId}");
        coverFile.createSync();
        await coverFile.writeAsBytes(coverResult.bodyBytes);
      }

      if (!Directory("$basePath/$trackDownloadDirectory/${track.id}").existsSync()) {
        log.fine("Downloading track: ${track.id}");

        var trackResult = await http.get(
            Uri.parse("${session.serverOrigin}$apiDownload/${track.id}"), headers: Util.genAuthHeaders());

        var metadata = await ChimeAPI.getTrackMetadata(track.id);

        File trackFile = File("$basePath/$trackDownloadDirectory/${track.id}");
        trackFile.createSync();
        await trackFile.writeAsBytes(trackResult.bodyBytes);

        // Add to database

        dbMgr.addTrackRecord(track, metadata.size, metadata.originalFile, metadata.format);
      }
    }

    // Finally add collection.
    dbMgr.addCollectionRecord(collection);

    GetIt.I<DownloadNotifier>().downloading = false;
  }

  static void deleteCollection(String id) async {

    // Figure out if tracks in collection are present in any other collections, if they are we can't delete them.

    Collection collection = await dbMgr.getCollectionRecord(id) as Collection; //can never be null, if it is something very bad has happened.

    for (var t in collection.tracks) {
      
      if (await dbMgr.countCollectionRecords(t.id) == 1) {
        
        File("${(await getApplicationDocumentsDirectory()).path}/$trackDownloadDirectory/${t.id}").deleteSync();
        dbMgr.deleteTrackRecord(t.id);

      } else {
        continue;
      }

    }

    dbMgr.deleteCollectionRecord(id);

  }

}

// Used for updating UI details and notifications.
class DownloadNotifier extends ChangeNotifier {
  bool _downloading = false;
  late DownloadType _downloadType;
  late String _downloadingId;
  int _downloadingProgress = 0;
  int _downloadingTotal = 0;

  set downloading(bool val) {
    _downloading = val;
    notifyListeners();
  }

  set downloadType(DownloadType val) {
    _downloadType = val;
    notifyListeners();
  }

  set downloadingId(String val) {
    _downloadingId = val;
    notifyListeners();
  }

  set downloadingProgress(int val) {
    _downloadingProgress = val;
    notifyListeners();
  }

  set downloadingTotal(int val) {
    _downloadingTotal = val;
    notifyListeners();
  }

  bool get downloading => _downloading;
  DownloadType get downloadType => _downloadType;
  String get downloadingId => _downloadingId;
  int get downloadingProgress => _downloadingProgress;
  int get downloadingTotal => _downloadingTotal;
}

// Enums for determining what is downloaded.
enum DownloadType { track, collection }

// Used if download is paused/resumed etc.
class DownloadStatus {
  bool isDownloading;
  DownloadType downloadType;
  String downloadingId;
  int downloadingProgress;
  int downloadingTotal;

  DownloadStatus({
    required this.isDownloading,
    required this.downloadType,
    required this.downloadingId,
    required this.downloadingProgress,
    required this.downloadingTotal,
  });

  factory DownloadStatus.fromJson(Map<String, dynamic> json) => DownloadStatus(
      isDownloading: json["is_downloading"],
      downloadType: json["download_type"],
      downloadingId: json["downloading_id"],
      downloadingProgress: json["downloading_progress"],
      downloadingTotal: json["downloading_total"]);

  Map<String, dynamic> toJson() => {
        "is_downloading": isDownloading,
        "downloading_type": downloadType,
        "downloading_id": downloadingId,
        "downloading_progress": downloadingProgress,
        "downloading_total": downloadingTotal
      };
}

// Abstract inserting/deleting records into database
class DownloadDatabaseManager {
  late Database db;

  // Init the database manager.
  DownloadDatabaseManager() {
    log.fine("Opening database...");

    getApplicationDocumentsDirectory().then((value) async {
      log.fine(value.path);

      if (File("${value.path}/downloads.db").existsSync()) {
        db = await openDatabase("${value.path}/downloads.db");
      } else {
        await File("${value.path}/downloads.db").create();
        db = await openDatabase("${value.path}/downloads.db");

        // Create tables

        await db.execute("""
          CREATE TABLE $tableCollections (
            id STRING PRIMARY KEY,
            name STRING,
            description STRING,
            cover_id STRING,
            is_album INTEGER,
            tracks STRING,
            dates STRING);
          """);

        await db.execute("""
          CREATE TABLE $tableTracks (
            id STRING PRIMARY KEY,
            name STRING,
            released INTEGER,
            artist STRING,
            album_id STRING,
            album_name STRING,
            duration REAL,
            cover STRING,
            original STRING,
            position INTEGER,
            size INTEGER,
            type STRING);
          """);
      }

      // Create directory for downloads

      if (!Directory("${value.path}/$downloadDirectory").existsSync()) {
        Directory("${value.path}/$downloadDirectory").createSync();
        Directory("${value.path}/$trackDownloadDirectory").createSync();
        Directory("${value.path}/$coverDownloadDirectory").createSync();
      }
    });
  }

  // Add a track to the database
  void addTrackRecord(Track track, int size, String original, String type) {
    db.execute(
        "INSERT INTO $tableTracks (id, name, released, artist, duration, album_id, cover, original, size, position, album_name, type) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)",
        [
          track.id,
          track.name,
          track.released,
          track.artist,
          track.duration,
          track.albumId,
          track.coverId,
          original,
          size,
          track.position,
          track.albumName,
          type,

        ]);
  }

  // Add a collection to the database
  void addCollectionRecord(Collection collection) {

    db.execute(
        "INSERT INTO $tableCollections (id, name, description, cover_id, is_album, tracks, dates) VALUES (?,?,?,?,?,?,?)",
        [
          collection.id,
          collection.title,
          collection.description,
          collection.coverId,
          collection.isAlbum ? 1 : 0,
          collection.tracks.map((e) => e.id).toList().join(","),
          collection.dates.join(",")
        ]);
  }

  // Remove a collection from the database (for when it is "un-downloaded")
  void deleteCollectionRecord(String id) async {
    await db.delete(tableCollections, where: "id = ?", whereArgs: [id]);
  }

  // Deletes a track from the database when it is "un-downloaded".
  void deleteTrackRecord(String id) async {
    await db.delete(tableTracks, where: "id = ?", whereArgs: [id]);
  }

  // Returns track data from the database, if the track is not found, returns null.
  Future<Track?> getTrackRecord(String id) async {
    List<Map<String, dynamic>> results =
        await db.query(tableTracks, where: "id = ?", whereArgs: [id]);

    if (results.isEmpty) {
      return null;
    } else {
      return Track(
          id: id,
          name: results[0]["name"],
          albumName: results[0]["album_name"],
          released: results[0]["released"],
          artist: results[0]["artist"],
          albumId: results[0]["album_name"],
          duration: results[0]["duration"],
          coverId: results[0]["cover"],
          position: results[0]["position"]);
    }
  }

  // Returns a collection from the database
  Future<Collection?> getCollectionRecord(String id) async {
    List<Map<String, dynamic>> results =
        await db.query(tableCollections, where: "id = ?", whereArgs: [id]);

    if (results.isEmpty) {
      log.fine("No collection with id $id in database");
      return null;
    } else {
      List<Track> tracks = [];

      // Query tracks.

      for (var t in (results[0]["tracks"] as String).split(",")) {
        Track? temp = await getTrackRecord(t);

        if (temp == null) {
          return null;
        } else {
          tracks.add(temp);
        }
      }

      return Collection(
          id: id,
          title: results[0]["name"],
          coverId: results[0]["cover_id"],
          isAlbum: results[0]["is_album"] == 1,
          tracks: tracks,
          dates: (results[0]["dates"] as int).toString().split(","),
          description: results[0]["description"],
          protected: false);
    }
  }

  Future<int> countCollectionRecords(String trackId) async {

    return (await db.query(tableCollections, where: "tracks LIKE ?", whereArgs: ["%$trackId%"])).length;
    
  }

  Future<bool> collectionRecordExists(String id) async {

    return (await getCollectionRecord(id)) != null;

  }
}
