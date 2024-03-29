import 'dart:ui';

import 'package:app/api/api.dart';
import 'package:app/player.dart';
import 'package:app/views/fullplayerview.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:get_it/get_it.dart';
import 'package:google_fonts/google_fonts.dart';

class DockedPlayer extends StatefulWidget {

  @override
  DockedPlayerState createState() => DockedPlayerState();

}


class DockedPlayerState extends State<DockedPlayer> {

  @override
  void initState() {

      super.initState();
      GetIt.I<PlayerStatusNotifier>().addListener(updatePlayerDetails);

  }

  void updatePlayerDetails() {
    setState(() {});
  }

  @override
  Widget build(BuildContext context) {
    
    return GestureDetector(
      onTap: () {
        if (GetIt.I<PlayerStatusNotifier>().active && !GetIt.I<PlayerStatusNotifier>().playingRadio) {
          showDialog(context: context, builder: (BuildContext context) => FullPlayerView());
        }
      },
      child: Container(
        width: MediaQuery.of(context).size.width-30, 
        height: 50, 
        decoration: BoxDecoration(
          color: Colors.grey[800],
          boxShadow: [
            BoxShadow(
              color: Colors.black.withOpacity(0.25),
              spreadRadius: 4,
              blurRadius: 7,
              offset: const Offset(0, 3)
            )
          ],
          borderRadius: const BorderRadius.all(Radius.circular(4.0)),
        ),
        child: SizedBox(
          height: 50,
          child: Stack(
            fit:  StackFit.expand,
            children: [
              Image(image: ChimeAPI.getCover(GetIt.I<PlayerStatusNotifier>().coverID, width: 300, height: 300), fit: BoxFit.cover,),
              ClipRRect(
                child: BackdropFilter(
                  filter: ImageFilter.blur(sigmaX: 10, sigmaY: 10),
                  child: Container(
                    alignment: Alignment.center,
                    color: Colors.grey.withOpacity(0.1),
                    child: Stack(
                      children: [
                        Container(
                          padding: const EdgeInsets.fromLTRB(4.0, 1.0, 4.0, 0.0),
                          child: Column(
                            children: [
                              Row(
                                mainAxisAlignment: MainAxisAlignment.center,
                                children: [
                                  Expanded(
                                    flex: 1,
                                    child: Image(image: ChimeAPI.getCover(GetIt.I<PlayerStatusNotifier>().coverID, width: 100, height: 100),)
                                  ),
                                  const SizedBox(width: 5,),
                                  Expanded(
                                    flex: 5,
                                    child: Column(
                                      crossAxisAlignment: CrossAxisAlignment.start,
                                      children: [
                                        Text(
                                          GetIt.I<PlayerStatusNotifier>().currentTrack,
                                          style: GoogleFonts.ibmPlexSans(fontSize: 16),
                                          overflow: TextOverflow.ellipsis,
                                          maxLines: 1,
                                          softWrap: false,
                                        ),
                                        Text(
                                          GetIt.I<PlayerStatusNotifier>().currentArtist,
                                          style: GoogleFonts.ibmPlexSans(fontSize: 10),
                                        )
                                      ],
                                    ),
                                  ),
                                  Expanded(
                                    flex: 4,
                                    child: Row(
                                      children: [
                                        IconButton(
                                          icon: GetIt.I<PlayerStatusNotifier>().playing ? const Icon(Icons.pause_rounded) : const Icon(Icons.play_arrow_rounded), 
                                          iconSize: 24, 
                                          color: Colors.white70, 
                                          onPressed: () {
                                            if (GetIt.I<PlayerStatusNotifier>().playing) {
                                              audioPlayer.pause();
                                            } else {
                                              audioPlayer.play();
                                            }
                                          },
                                        ),
                                        IconButton(
                                          icon: const Icon(Icons.shuffle_rounded),
                                          color: GetIt.I<PlayerStatusNotifier>().shuffle ? Colors.yellow[800] : Colors.white70,
                                          onPressed: () {
                                            GetIt.I<PlayerStatusNotifier>().shuffle = !GetIt.I<PlayerStatusNotifier>().shuffle;
                                          },
                                        ),
                                        IconButton(
                                          icon: const Icon(Icons.repeat_rounded),
                                          color: GetIt.I<PlayerStatusNotifier>().loop ? Colors.yellow[800] : Colors.white70,
                                          onPressed: () {
                                            GetIt.I<PlayerStatusNotifier>().loop = !GetIt.I<PlayerStatusNotifier>().loop;
                                          },
                                        )
                                      ],
                                    )
                                  )
                                ],
                              ),
                            ],
                          )
                        ),
                        Positioned(
                          bottom: 0.0,
                          left: 0.0,
                          right: 0.0,
                          child: LinearProgressIndicator(
                            value: GetIt.I<PlayerStatusNotifier>().completion,
                          )
                        ),
                      ],
                    ),
                  ),  
                ),
              )
            ],
          ),
        )
      )
    );
    
  }

}