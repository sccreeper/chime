
import 'package:app/api/api.dart';
import 'package:app/player.dart';
import 'package:flutter/material.dart';
import 'package:get_it/get_it.dart';
import 'package:app/shared.dart';
import 'package:google_fonts/google_fonts.dart';

class FullPlayerView extends StatefulWidget {

  @override
  FullPlayerViewState createState() => FullPlayerViewState();

}

class FullPlayerViewState extends State<FullPlayerView> {

  @override
  void initState() {

    GetIt.I<PlayerStatusNotifier>().addListener(updateView);

    super.initState();
  }

  void updateView() async {
    if (mounted) {
      setState(() {}); 
    }
  }

  @override
  Widget build(BuildContext context) {
    
    return Dialog.fullscreen(

      backgroundColor: Colors.grey[800],

      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        //crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          Container(
            decoration: BoxDecoration(
              boxShadow: [
                BoxShadow(
                  color: Colors.black.withOpacity(0.25),
                  spreadRadius: 5,
                  blurRadius: 16,
                  offset: const Offset(0, 3)
                )
              ]
            ),
            padding: const EdgeInsets.all(8.0),
            child:  Image(image: ChimeAPI.getCover(GetIt.I<PlayerStatusNotifier>().coverID)),
          ),
          const SizedBox(height: 24,),
          Text(GetIt.I<PlayerStatusNotifier>().currentTrack, textAlign: TextAlign.center, style: GoogleFonts.ibmPlexSans(fontWeight: FontWeight.bold, fontSize: 20, ),),
          const SizedBox(height: 8,),
          Text(GetIt.I<PlayerStatusNotifier>().currentArtist, textAlign: TextAlign.center, style: GoogleFonts.ibmPlexSans(fontSize: 12),),
          const SizedBox(height: 16,),
          Row(
            children: [
              Expanded(child: Text(Util.convertDuration(GetIt.I<PlayerStatusNotifier>().currentTime), textAlign: TextAlign.center,), flex: 2,), 
              Expanded(flex: 6, child: Slider(
                value: GetIt.I<PlayerStatusNotifier>().currentTime,
                max: GetIt.I<PlayerStatusNotifier>().duration,
                onChangeStart: (val) => audioPlayer.pause(),
                onChangeEnd: (val) => audioPlayer.play(),
                onChanged: (double val) => GetIt.I<PlayerStatusNotifier>().currentTime = val
              )),
              Expanded(child: Text(Util.convertDuration(GetIt.I<PlayerStatusNotifier>().duration), textAlign: TextAlign.center,), flex: 2,), 
            ],
          ),
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              IconButton(
                icon: const Icon(Icons.shuffle_rounded),
                iconSize: 24.0,
                color: GetIt.I<PlayerStatusNotifier>().shuffle ? Colors.yellow[800] : Colors.white70,
                onPressed: () {
                  GetIt.I<PlayerStatusNotifier>().shuffle = !GetIt.I<PlayerStatusNotifier>().shuffle;
                },
              ),
              IconButton(
                onPressed: () => Player.previousTrack(), 
                icon: Icon(Icons.skip_previous_rounded, color: Colors.yellow[800],),
                iconSize: 48.0,
              ),
              IconButton(
                onPressed: () {
                  if (GetIt.I<PlayerStatusNotifier>().playing) {
                    audioPlayer.pause();
                  } else {
                    audioPlayer.play();
                  }
                },
                icon: Icon(GetIt.I<PlayerStatusNotifier>().playing ? Icons.pause_rounded : Icons.play_arrow_rounded, color: Colors.yellow[800],),
                iconSize: 72.0,
              ),
              IconButton(
                onPressed: () => Player.nextTrack(), 
                icon: Icon(Icons.skip_next_rounded, color: Colors.yellow[800]),
                iconSize: 48.0,
              ),
              IconButton(
                icon: const Icon(Icons.repeat_rounded),
                color: GetIt.I<PlayerStatusNotifier>().loop ? Colors.yellow[800] : Colors.white70,
                iconSize: 24.0,
                onPressed: () {
                  GetIt.I<PlayerStatusNotifier>().loop = !GetIt.I<PlayerStatusNotifier>().loop;
                },
              )
            ],
          ),
        ],
      ),

    );
    
  }

}