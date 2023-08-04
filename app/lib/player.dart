// Contains everything needed for managing the audio player.

import 'dart:math';

import 'package:app/api/api.dart';
import 'package:app/api/endpoints.dart';
import 'package:app/api/models/collections.dart';
import 'package:app/api/models/radio.dart';
import 'package:app/shared.dart';
import 'package:flutter/foundation.dart';
import 'package:get_it/get_it.dart';
import 'package:just_audio/just_audio.dart';

AudioPlayer audioPlayer = AudioPlayer();

class Player {

  static List<Track> trackQueue = [];
  static List<Track> viewingTracks = [];
  static Track? currentTrack;
  static int trackIndex = 0;
  static String currentCollectionId = "";

  static bool playingRandom = false;

  static void init() {

    audioPlayer.playerStateStream.listen((playerState) async {

      if (playerState.processingState == ProcessingState.completed) {
        
        if (GetIt.I<PlayerStatusNotifier>().shuffle && !playingRandom) {

          playTrack(trackQueue[Random().nextInt(trackQueue.length)]);

        } else if (GetIt.I<PlayerStatusNotifier>().loop) {

          audioPlayer.seek(Duration.zero);
          audioPlayer.play();

        } else if(!(trackIndex+1 >= trackQueue.length)) {

          trackIndex += 1;

          playTrackId(trackQueue[trackIndex].id);
          GetIt.I<PlayerStatusNotifier>().updateDetails(trackQueue[trackIndex]);

        } else if(playingRandom) {

          List<String> allTracks = await ChimeAPI.getTracks(0);
          String trackId = allTracks[Random().nextInt(allTracks.length)];
          TrackMetadata track = await ChimeAPI.getTrackMetadata(trackId);

          GetIt.I<PlayerStatusNotifier>().updateDetails(Track.fromMetadata(trackId, track));

          playTrackId(trackId);

        } else {
          currentCollectionId = "";
          trackQueue.clear();
          trackIndex = 0;
          playingRandom = true;

        }

      }

    });

    audioPlayer.positionStream.listen((duration) {

      if (audioPlayer.duration != null) {
        GetIt.I<PlayerStatusNotifier>().setCompletion( duration.inMilliseconds / audioPlayer.duration!.inMilliseconds);
      } else {
        GetIt.I<PlayerStatusNotifier>().setCompletion(0.0);
      }

    });

    audioPlayer.playingStream.listen((playing) {
      GetIt.I<PlayerStatusNotifier>().playPause(playing);
    });

  }

  static void playTrack(Track track) async {

    playTrackId(track.id);
    GetIt.I<PlayerStatusNotifier>().updateDetails(track);

  }

  static void playTrackId(String trackId) async {
    
    await audioPlayer.setAudioSource(
      AudioSource.uri(
        Uri.parse("${session.serverOrigin}${apiStream}/${trackId}"),
        headers: {"Cookie":"session=${session.sessionBase64}"},
      )
    );

    GetIt.I<PlayerStatusNotifier>().playingRadio = false;

    audioPlayer.play();
  
  }

  static void playCollection(String id, String startingTrackId, int startingIndex, Track trackDetails) {

    trackQueue = viewingTracks;
    trackIndex = startingIndex;
    currentCollectionId = "";

    playTrack(trackDetails);

  }

  static void playRadio(RadioModel radio) async {

    await audioPlayer.setAudioSource(
      AudioSource.uri(Uri.parse(radio.url))
    );

    GetIt.I<PlayerStatusNotifier>().playingRadio = true;
    GetIt.I<PlayerStatusNotifier>().updateDetails(radio);
    
    audioPlayer.play();

  }

}

class PlayerStatusNotifier extends ChangeNotifier {

  String _currentTrackName = "No track playing";
  String _currentArtist = "No artist";
  double _completion = 0.0;
  String _coverId = "0";
  bool _playing = false;
  bool _shuffle = false;
  bool _loop = false;
  bool _playingRadio = false;

  String get currentTrack => _currentTrackName;
  String get currentArtist => _currentArtist;
  double get completion => _completion;
  String get coverID => _coverId;
  bool get playing => _playing;
  bool get shuffle => _shuffle;
  bool get loop => _loop;
  bool get playingRadio => _playingRadio;

  void updateDetails<T>(T item) {

    if (item is Track) {
      _currentTrackName = item.name;
      _currentArtist = item.artist;
      _coverId = item.coverId; 

      notifyListeners();
    } else if (item is RadioModel) {
      _currentTrackName = item.name;
      _currentArtist = "Playing internet radio";
      _coverId = item.coverId;

      notifyListeners();
    } else {
      throw ArgumentError('Unsupported $T type');
    }

  }

  void setCompletion(double completion) {
    _completion = completion;

    notifyListeners();

  }

  void playPause(bool playing) {

    _playing = playing;
    notifyListeners();

  }

  set shuffle(bool val) {
    _shuffle = val;
    notifyListeners();
  }

  set loop(bool val) {
    _loop = val;
    notifyListeners();
  }

  set playingRadio(bool val) {
    _playingRadio = val;
    notifyListeners();
  }

}