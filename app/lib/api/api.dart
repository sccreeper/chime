import 'dart:convert';
import 'dart:typed_data';

import 'package:app/api/models/collections.dart';
import 'package:app/shared.dart';
import 'package:app/api/endpoints.dart';
import 'package:app/shared.dart';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

class ChimeAPI {

  static Future<Library> getLibary() async {
    
    final req = await http.get(Uri.parse("${session.serverOrigin}${apiGetCollections}"), headers: {"Cookie": "session=${session.sessionBase64}"});
    return Library.fromJSON(jsonDecode(req.body));

  }

  static Future<Collection> getCollection(String id) async {
    
    final req = await http.get(Uri.parse("${session.serverOrigin}${apiGetCollection}/${id}"), headers: {"Cookie": "session=${session.sessionBase64}"});
    return Collection.fromJSON(jsonDecode(req.body));

  }

  static Future<Uint8List> getCover(String id) async {

    final req = await http.get(Uri.parse("${session.serverOrigin}${apiGetCover}/${id}"), headers: {"Cookie": "session=${session.sessionBase64}"});
    return req.bodyBytes;
  
  }

  static Future<TrackMetadata> getTrackMetadata(String id) async {

    final req = await http.get(Uri.parse("${session.serverOrigin}${apiTrackMetadata}/${id}"), headers: {"Cookie": "session=${session.sessionBase64}"});
    return TrackMetadata.fromJSON(jsonDecode(req.body));

  }

  static Future<List<String>> getTracks(int limit) async {

    final req = await http.post(
      Uri.parse("${session.serverOrigin}/${apiAllTracks}"),
      headers: {"Cookie":"session=${session.sessionBase64}"},
      body: jsonEncode(<String,int>{"limit":limit})
    );
    return jsonDecode(req.body) as List<String>;

  }

}