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
  static List<Track> previousTrackQueue = [];
  static List<Track> viewingTracks = [];
  static Track? currentTrack;
  static int trackIndex = 0;
  static String currentCollectionId = "";

  static bool playingRandom = false;

  static void init() {

    audioPlayer.playerStateStream.listen((playerState) async {

      if (playerState.processingState == ProcessingState.completed) {
        nextTrack();
      }

    });

    audioPlayer.positionStream.listen((duration) {

      if (audioPlayer.duration != null) {
        GetIt.I<PlayerStatusNotifier>().setCurrentTime(duration.inSeconds.toDouble());
      } else {
        GetIt.I<PlayerStatusNotifier>().setCurrentTime(0.0);
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
    GetIt.I<PlayerStatusNotifier>().active = true;
    GetIt.I<PlayerStatusNotifier>().duration = audioPlayer.duration!.inSeconds.toDouble();

    audioPlayer.play();
  
  }

  static void playCollection(String id, String startingTrackId, int startingIndex, Track trackDetails) {

    trackQueue = viewingTracks;
    trackIndex = startingIndex;
    currentCollectionId = "";
    
    currentTrack = trackDetails;

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

  static void seek(double time) {
    audioPlayer.seek(Duration(seconds: time.toInt()));
  }

  static void nextTrack() async {
    if (GetIt.I<PlayerStatusNotifier>().shuffle && !playingRandom) {
        
        previousTrackQueue.add(currentTrack!);

        Track randomTrack = trackQueue[Random().nextInt(trackQueue.length)];
        currentTrack = randomTrack;
        playTrack(randomTrack);

      } else if (GetIt.I<PlayerStatusNotifier>().loop) {

        audioPlayer.seek(Duration.zero);
        audioPlayer.play();

      } else if(trackIndex+1 <= trackQueue.length) {

        trackIndex += 1;

        previousTrackQueue.add(currentTrack!);

        playTrackId(trackQueue[trackIndex].id);
        GetIt.I<PlayerStatusNotifier>().updateDetails(trackQueue[trackIndex]);
        currentTrack = trackQueue[trackIndex];

      } else if(playingRandom) {

        previousTrackQueue.add(currentTrack!);

        List<String> allTracks = await ChimeAPI.getTracks(0);
        String trackId = allTracks[Random().nextInt(allTracks.length)];
        TrackMetadata track = await ChimeAPI.getTrackMetadata(trackId);

        currentTrack = Track.fromMetadata(trackId, track);
        GetIt.I<PlayerStatusNotifier>().updateDetails(Track.fromMetadata(trackId, track));

        playTrackId(trackId);

      } else {
        currentCollectionId = "";
        trackQueue.clear();
        trackIndex = 0;
        playingRandom = true;
      }
  }

  static void previousTrack() async {

    if (previousTrackQueue.isEmpty) {
      
      audioPlayer.seek(Duration.zero);
      audioPlayer.play();

    } else {
    
      GetIt.I<PlayerStatusNotifier>().updateDetails(previousTrackQueue.last);
      playTrack(previousTrackQueue.last);
      previousTrackQueue.removeLast();

    }

  }

}

class PlayerStatusNotifier extends ChangeNotifier {

  String _currentTrackName = "No track playing";
  String _currentArtist = "No artist";
  bool _active = false;
  double _completion = 0.0;
  String _coverId = "0";
  bool _playing = false;
  bool _shuffle = false;
  bool _loop = false;
  bool _playingRadio = false;
  double _duration = 0.0;
  double _currentTime = 0.0;

  String get currentTrack => _currentTrackName;
  String get currentArtist => _currentArtist;
  double get completion {
    if (_duration == 0) {
      return 0.0;
    } else {
      return _currentTime / _duration;
    }
  }
  String get coverID => _coverId;
  bool get playing => _playing;
  bool get shuffle => _shuffle;
  bool get loop => _loop;
  bool get playingRadio => _playingRadio;
  double get duration => _duration;
  double get currentTime => _currentTime;
  bool get active => _active;


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

  void setCurrentTime(double currentTime) {
    _currentTime = currentTime;

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

  set duration(double val) {
    _duration = val;
    notifyListeners();
  }

  set active(bool val) {
    _active = val;
    notifyListeners();
  }

  set currentTime(double val) {
    _currentTime = val;
    Player.seek(val);
  }

}