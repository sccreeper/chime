// Contains everything needed for managing the audio player.

import 'package:app/api/endpoints.dart';
import 'package:app/api/models/collections.dart';
import 'package:app/shared.dart';
import 'package:just_audio/just_audio.dart';

AudioPlayer? audioPlayer;

class Player {

  static void init() {

    audioPlayer!.playerStateStream.listen((playerState) async {

      if (playerState.processingState == ProcessingState.completed) {
        
        await audioPlayer!.setAudioSource(
          AudioSource.uri(
            Uri.parse("${session.serverOrigin}${apiStream}/317c37e910343dc"), //Hardcoded before playlist support was added.
            headers: {"Cookie":"session=${session.sessionBase64}"},
          )
        );

        audioPlayer!.play();

      }

    });

  }

  static void playTrack(Track track) async {

    await audioPlayer!.setAudioSource(
      AudioSource.uri(
        Uri.parse("${session.serverOrigin}${apiStream}/${track.id}"),
        headers: {"Cookie":"session=${session.sessionBase64}"},
      )
    );

    audioPlayer!.play();

  }

}