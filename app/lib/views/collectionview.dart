import 'dart:typed_data';
import 'dart:ui';

import 'package:app/api/api.dart';
import 'package:app/api/models/collections.dart';
import 'package:app/shared.dart';
import 'package:app/widgets/borderedchip.dart';
import 'package:app/widgets/loadingspinner.dart';
import 'package:flutter/material.dart';
import 'package:app/api/endpoints.dart';
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
              BorderedChip(text: Util.convertDurationVerbose(collection.tracks.fold(0.0, (prevValue, e) => prevValue + e.duration)))
            ],
          ),
          const Divider(),
          Expanded(
            child: Scrollbar(
              child: ListView(
                children: collection.tracks.map((e) => TrackScaffold(track: e)).toList(),
              )
            ),
          )
      ]),

    );
    
  }

}

class TrackScaffold extends StatelessWidget {

  final Track track;

  const TrackScaffold({super.key, required this.track});

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
    );


    // return Column(
    //   crossAxisAlignment: CrossAxisAlignment.start,
    //   children: [
    //     SingleChildScrollView(
    //       scrollDirection: Axis.horizontal,
    //       child: Text(track.name, maxLines: 1, style: GoogleFonts.anuphan(color: Colors.white, fontWeight: FontWeight.w500))
    //     ),
    //     Row(
    //       children: [            
    //         Text(track.artist, style: Theme.of(context).textTheme.bodySmall,),
    //         const SizedBox(width: 2.5,),
    //         Text("●", style: Theme.of(context).textTheme.bodySmall,),
    //         const SizedBox(width: 2.5,),
    //         Text(Util.convertDuration(track.duration), style: Theme.of(context).textTheme.bodySmall)
    //       ],
    //     )
    //   ],
    // );
    
  }


}