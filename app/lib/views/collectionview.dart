import 'dart:math';
import 'dart:typed_data';
import 'dart:ui';

import 'package:app/api/api.dart';
import 'package:app/api/downloads.dart';
import 'package:app/api/endpoints.dart';
import 'package:app/api/models/collections.dart';
import 'package:app/player.dart';
import 'package:app/shared.dart';
import 'package:app/widgets/borderedchip.dart';
import 'package:app/widgets/downloads.dart';
import 'package:app/widgets/loadingspinner.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:google_fonts/google_fonts.dart';

class CollectionView extends StatefulWidget {

    final String id;

    const CollectionView({super.key, required this.id});

    @override
    _CollectionViewState createState() => _CollectionViewState();

}

class _CollectionViewState extends State<CollectionView> {

  Widget _childWidget = LoadingSpinner();

  @override
  void initState() {
    
    _fetchCollection();

    super.initState();

  }

  void _fetchCollection() async {

    log.fine("Fetching details...");
    
    Collection collection = await ChimeAPI.getCollection(widget.id);

    Uint8List image;

    if (collection.coverId == "0") {
      image = ( await rootBundle.load("assets/no_cover.png") ).buffer.asUint8List();
    } else {
      image = await ChimeAPI.getCover(collection.coverId);
    }

    if (mounted) {
      setState(() {
        _childWidget = CollectionScaffold(
          collection: collection, 
          coverBytes: image
        );
      }); 
    }

  }
  
  @override
  Widget build(BuildContext context) {

    currentCollection = widget.id;
    
    return Container(
      child: _childWidget,
    );
    
  }

}

class CollectionScaffold extends StatelessWidget {

  final Collection collection;
  final Uint8List coverBytes;

  CollectionScaffold({super.key, required this.collection, required this.coverBytes});

  @override
  Widget build(BuildContext context) {
    
    // The cast on this line *is* required.
    List<int> contentsIndexes = Iterable<int>.generate(collection.tracks.length).toList();
    List<Widget> contents = contentsIndexes.map((i) => TrackScaffold(track: collection.tracks[i], index: i, collectionId: currentCollection) as Widget).toList();

    log.fine(contentsIndexes.length);

    Player.viewingTracks = collection.tracks;

    contents.add(
      const SizedBox(
        height: 48.0,
      )
    );
    
    return Container(

      child: Column(children: [
          
          // SizedBox with stack for blurred image background
          SizedBox(
            height: 200,
            child: Stack(
              fit: StackFit.expand,
              children: [
                Image.memory(coverBytes,fit: BoxFit.cover,),
                ClipRRect(
                  child: BackdropFilter(
                    filter: ImageFilter.blur(sigmaX: 10, sigmaY: 10),
                    child: Container(
                      alignment: Alignment.center,
                      color: Colors.grey.withOpacity(0.1),
                      child: Image.memory(coverBytes, width: 200, height: 200,),
                    ),  
                  ),
                )
              ],
            ),
          ),
          // Title and other details about collection
          const Divider(),
          Text(collection.title, textAlign: TextAlign.center, style: GoogleFonts.anuphan(color: Colors.white, fontSize: 24.0, fontWeight: FontWeight.bold)),
          const Divider(),
          Row(
            mainAxisAlignment: MainAxisAlignment.start,
            crossAxisAlignment: CrossAxisAlignment.center,
            children: [
              BorderedChip(text: collection.isAlbum ? "Album" : "Playlist"),
              const SizedBox(width: 5,),
              BorderedChip(text: "${collection.tracks.length} tracks"),
              const SizedBox(width: 5,),
              BorderedChip(text: Util.convertDurationVerbose(collection.tracks.fold(0.0, (prevValue, e) => prevValue + e.duration))),
              const SizedBox(width: 5,),
              CollectionDownloadButton(id: collection.id,)
            ],
          ),
          const Divider(),
          Expanded(
            child: Scrollbar(
              child: ListView(
                children: contents,
              )
            ),
          )
      ]),

    );
    
  }

}

class TrackScaffold extends StatelessWidget {

  final Track track;
  final int index;
  final String collectionId;

  const TrackScaffold({super.key, required this.track, required this.index, required this.collectionId});

  @override
  Widget build(BuildContext context) {
    
    return ListTile(
      title: SingleChildScrollView(
          scrollDirection: Axis.horizontal,
          child: Text(track.name, maxLines: 1, style: GoogleFonts.anuphan(color: Colors.white, fontWeight: FontWeight.w500))
      ),
      subtitle: Text("${track.artist} ● ${Util.convertDuration(track.duration)}"),
      dense: true,
      contentPadding: EdgeInsets.all(0.0),
      onTap: () {
          
        Player.playCollection(collectionId, track.id, index, track);

      },

      onLongPress: () {

        ChimeAPI.getTrackMetadata(track.id).then((TrackMetadata trackMetadata) {
          
          showDialog(
            context: context,
            builder: (BuildContext context) => AlertDialog(
              backgroundColor: Colors.grey[800],
              title: const Text("Track metadata"),
              contentTextStyle: GoogleFonts.anuphan(),
              titleTextStyle: GoogleFonts.anuphan(),
              content: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                mainAxisSize: MainAxisSize.min,
                children: [
                  Center(
                    child: (track.coverId == "0" ? 
                      Image.asset("assets/no_cover.png", width: 100, height: 100,) : 
                      Image.network(
                        "${session.serverOrigin}${apiGetCover}/${track.coverId}", 
                        headers: {"Cookie":"session=${session.sessionBase64}"}, 
                        width: 100,
                        height: 100,
                        )
                    ),
                  ),
                  Divider(),
                  Text(track.name),
                  Divider(),
                  RichText(
                    text: TextSpan(
                      children: [
                        const TextSpan(text: "Released: ", style: TextStyle(fontWeight: FontWeight.bold)),
                        TextSpan(text: trackMetadata.released.toString())
                      ]
                    ),
                  ),
                  RichText(
                    text: TextSpan(
                      children: [
                        const TextSpan(text: "Duration: ", style: TextStyle(fontWeight: FontWeight.bold)),
                        TextSpan(text: Util.convertDuration(trackMetadata.duration))
                      ]
                    ),
                  ),
                  RichText(
                    text: TextSpan(
                      children: [
                        const TextSpan(text: "Format: ", style: TextStyle(fontWeight: FontWeight.bold)),
                        TextSpan(text: trackMetadata.format)
                      ]
                    ),
                  ),
                  RichText(
                    text: TextSpan(
                      children: [
                        const TextSpan(text: "Original file: ", style: TextStyle(fontWeight: FontWeight.bold)),
                        TextSpan(text: trackMetadata.originalFile)
                      ]
                    ),
                  ),
                  RichText(
                    text: TextSpan(
                      children: [
                        const TextSpan(text: "File size: ", style: TextStyle(fontWeight: FontWeight.bold)),
                        TextSpan(text: "${(trackMetadata.size / pow(10, 6)).toStringAsFixed(2)} mb")
                      ]
                    ),
                  ),
                ],
              ),
              actions: [
                TextButton(
                  onPressed: () => Navigator.pop(context, "OK"), 
                  child: const Text("OK")
                )
              ],
            )
          );

        });

      },
    );
    
  }

}