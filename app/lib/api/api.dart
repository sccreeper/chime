import 'dart:async';
import 'dart:convert';
import 'dart:io';

import 'package:app/api/downloads.dart';
import 'package:app/api/models/collections.dart';
import 'package:app/api/models/misc.dart';
import 'package:app/api/models/radio.dart';
import 'package:app/api/models/search.dart';
import 'package:app/shared.dart';
import 'package:app/api/endpoints.dart';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

class ChimeAPI {

  static Future<Library> getLibrary() async {

    if (connected) {
      
      final req = await http.get(Uri.parse("${session.serverOrigin}$apiGetCollections"), headers: Util.genAuthHeaders());
      return Library.fromJSON(jsonDecode(req.body)); 
    
    } else {

      var results = await dbMgr.db.query("collections");

      Library lib = Library(
        albums: [], playlists: [], radios: []
      );

      for (var element in results) {
        
        if (element["is_album"] as int == 1) {
          
          lib.albums.add(LibraryItem(
              id: element["id"] as String, 
              name: element["name"] as String,
            )
          );

        } else {
          
            lib.playlists.add(LibraryItem(
              id: element["id"] as String, 
              name: element["name"] as String,
            )
          );

        }

      }

      return lib;

    }

  }

  static Future<Collection> getCollection(String id) async {

    if (connected) {
      
      final req = await http.get(Uri.parse("${session.serverOrigin}$apiGetCollection/$id"), headers: Util.genAuthHeaders());
      return Collection.fromJSON(jsonDecode(req.body), id); 
    
    } else {

      // Shouldn't ever be null.
      return (await dbMgr.getCollectionRecord(id))!;

    }

  }

  // static Future<Uint8List> getCover(String id) async {

  //   if (id == "0") {
      
  //     return ( await rootBundle.load("assets/no_cover.png") ).buffer.asUint8List();

  //   } else if (connected) {
        
  //       final req = await http.get(Uri.parse("${session.serverOrigin}$apiGetCover/$id"), headers: {"Cookie": "session=${session.sessionBase64}"});
  //       return req.bodyBytes; 
    
  //   } else {

  //     if (File("${(await getApplicationDocumentsDirectory()).path}/$coverDownloadDirectory/$id").existsSync()) {
    
  //       return await File("${(await getApplicationDocumentsDirectory()).path}/$coverDownloadDirectory/$id").readAsBytes();

  //     } else {
        
  //       return ( await rootBundle.load("assets/no_cover.png") ).buffer.asUint8List();
      
  //     }

  //   }
  
  // }

  static ImageProvider getCover(String id, {int width = 500, int height = 500}) {

    if (id == "0") {
      
      return Image.asset("assets/no_cover.png").image;

    } else if (connected) {
        
        return Image.network(
                  "${session.serverOrigin}$apiGetCover/$id?width=$width&height=$height",
                  headers: Util.genAuthHeaders(),
                ).image;
    
    } else {

      if (File("$docDirectory/$coverDownloadDirectory/$id").existsSync()) {
    
        return Image.file(File("$docDirectory/$coverDownloadDirectory/$id")).image;

      } else {
        
        return Image.asset("assets/no_cover.png").image;
      
      }

    }

  }

  static Future<TrackMetadata> getTrackMetadata(String id) async {

    if (connected) {
      
      final req = await http.get(Uri.parse("${session.serverOrigin}$apiTrackMetadata/$id"), headers: Util.genAuthHeaders());
      return TrackMetadata.fromJSON(jsonDecode(req.body)); 
    
    } else {

      return TrackMetadata.fromDatabaseMap(
        (await dbMgr.db.query("tracks", where: "id = ?", whereArgs: [id]))[0],
        id 
      );

    }

  }

  static Future<List<String>> getTracks(int limit) async {

    if (connected) {
      
      final req = await http.post(
        Uri.parse("${session.serverOrigin}$apiAllTracks"),
        headers: Util.genAuthHeaders(),
        body: jsonEncode(<String,int>{"limit":limit})
      );

      return jsonDecode(req.body) as List<String>; 

    } else {

      return (await dbMgr.db.query("tracks", columns: ["id"])).map((e) => e["id"] as String).toList();

    }

  }

  static Future<RadioModel> getRadio(String id) async {

    final req = await http.get(Uri.parse("${session.serverOrigin}$apiGetRadio/$id"), headers: Util.genAuthHeaders());
    return RadioModel.fromJson(jsonDecode(req.body));

  }

  static Future<SearchResults> search(String query) async {

    if (connected) {
      
      final req = await http.post(
        Uri.parse("${session.serverOrigin}$apiSearch"),
        headers: Util.genAuthHeaders(),
        body: jsonEncode(<String,String>{"query":query})
      );
      
      return SearchResults.fromJson(jsonDecode(req.body)); 
    
    } else {

      List<Map<String, Object?>> collectionResults = await dbMgr.db.query("collections", where: "name LIKE ?", whereArgs: ["%$query%"]);
      List<Map<String, Object?>> trackResults = await dbMgr.db.query("tracks", where: "name LIKE ?", whereArgs: ["%$query%"]);

      SearchResults results = SearchResults(
        tracks: [], collections: [], radios: []
      );

      for (var res in collectionResults) {
        
        results.collections.add(
          SearchCollection(
            id: res["id"] as String, 
            title: res["name"] as String, 
            coverId: res["cover_id"] as String, 
            isAlbum: (res["is_album"] as int) == 1
          )
        );

      }

      for (var res in trackResults) {
        
        results.tracks.add(
          SearchTrack(
            id: res["id"] as String, 
            albumId: res["album_id"] as String, 
            artist: res["artist"] as String, 
            title: res["name"] as String, 
            duration: res["duration"] as double, 
            coverId: res["cover"] as String
          )
        );

      }

      return results;
    
    }

  }

  static Future<PingResult> ping(String serverOrigin) async {

      final req = await http.get(Uri.parse("$serverOrigin$apiPing")).timeout(const Duration(seconds: 5));
      return PingResult.fromJson(jsonDecode(req.body));

  } 

}