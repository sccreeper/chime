import 'dart:ui';

import 'package:app/api/endpoints.dart';
import 'package:app/player.dart';
import 'package:app/shared.dart';
import 'package:app/views/fullplayerview.dart';
import 'package:flutter/material.dart';
import 'package:get_it/get_it.dart';
import 'package:google_fonts/google_fonts.dart';

class DockedPlayer extends StatefulWidget {

  @override
  _DockedPlayerState createState() => _DockedPlayerState();

}


class _DockedPlayerState extends State<DockedPlayer> {

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
          borderRadius: BorderRadius.all(Radius.circular(4.0)),
        ),
        child: SizedBox(
          height: 50,
          child: Stack(
            fit:  StackFit.expand,
            children: [
              GetIt.I<PlayerStatusNotifier>().coverID == "0" ? 
                Image.asset("assets/no_cover.png", fit: BoxFit.cover,) : 
                Image.network(
                  "${session.serverOrigin}${apiGetCover}/${GetIt.I<PlayerStatusNotifier>().coverID}",
                  headers: {"Cookie":"session=${session.sessionBase64}"},
                  fit: BoxFit.cover,
              ),
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
                                    child: GetIt.I<PlayerStatusNotifier>().coverID == "0" ? 
                                    Image.asset("assets/no_cover.png") : 
                                    Image.network(
                                      "${session.serverOrigin}${apiGetCover}/${GetIt.I<PlayerStatusNotifier>().coverID}",
                                      headers: {"Cookie":"session=${session.sessionBase64}"},
                                      )
                                  ),
                                  SizedBox(width: 5,),
                                  Expanded(
                                    flex: 5,
                                    child: Column(
                                      crossAxisAlignment: CrossAxisAlignment.start,
                                      children: [
                                        Text(
                                          GetIt.I<PlayerStatusNotifier>().currentTrack,
                                          style: GoogleFonts.anuphan(fontSize: 16),
                                          overflow: TextOverflow.ellipsis,
                                          maxLines: 1,
                                          softWrap: false,
                                        ),
                                        Text(
                                          GetIt.I<PlayerStatusNotifier>().currentArtist,
                                          style: GoogleFonts.anuphan(fontSize: 10),
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