import 'package:app/api/api.dart';
import 'package:app/api/models/search.dart';
import 'package:app/mainscreen.dart';
import 'package:app/shared.dart';
import 'package:app/views/collectionview.dart';
import 'package:app/views/libraryview.dart';
import 'package:app/views/radioview.dart';
import 'package:app/widgets/borderedchip.dart';
import 'package:app/widgets/iconlabel.dart';
import 'package:flutter/material.dart';
import 'package:get_it/get_it.dart';
import 'package:google_fonts/google_fonts.dart';

class SearchView extends StatefulWidget {

  const SearchView({super.key});

  @override
  SearchViewState createState() => SearchViewState();

}

class SearchViewState extends State<SearchView> {

  String searchQuery = "";
  List<SearchTrack> tracks = [];
  List<SearchCollection> playlists = [];
  List<SearchCollection> albums = [];
  List<SearchRadio> radios = [];

  void search() {

    ChimeAPI.search(searchQuery).then((results) {
      tracks.clear();
      playlists.clear();
      albums.clear();
      radios.clear();

      log.fine("Recieved ${results.collections.length} collection(s)");
      log.fine("Recieved ${results.tracks.length} track(s)");
      log.fine("Recieved ${results.radios.length} radio(s)");

      tracks = results.tracks;

      for (var element in results.collections) {
        
        if (element.isAlbum) {
          albums.add(element);
        } else {
          playlists.add(element);
        }
      }

      radios = results.radios;

      setState(() {});
    });

  }

  @override
  Widget build(BuildContext context) {

    List<Widget> albumChildren = [];
    List<Widget> playlistChildren = [];
    List<Widget> radioChildren = [];

    for (var element in albums) {
      albumChildren.add(const SizedBox(width: 6,));
      albumChildren.add(LargeSearchResultScaffold(
        title: element.title, 
        id: element.id, 
        coverId: element.coverId, 
        type: SearchResultType.album
      ));
    }

    for (var element in playlists) {
      playlistChildren.add(const SizedBox(width: 6,));
      playlistChildren.add(LargeSearchResultScaffold(
        title: element.title, 
        id: element.id, 
        coverId: element.coverId, 
        type: SearchResultType.playlist
      ));
    }

    for (var element in radios) {
      radioChildren.add(const SizedBox(width: 6,));
      radioChildren.add(LargeSearchResultScaffold(
        title: element.name, 
        id: element.id, 
        coverId: element.coverId, 
        type: SearchResultType.radio
      ));
    }
    
    return Column(
      children: [
        Container(
          padding: const EdgeInsets.all(10.0),
          child: Row(
            children: [
              Expanded(
                flex: 9,
                child: TextField(
                  decoration: const InputDecoration(hintText: "Search"),
                  onChanged: (value) => {searchQuery = value},
                )
              ),
              Expanded(
                flex: 1,
                child: IconButton(
                  icon: const Icon(Icons.search_rounded),
                  color: Colors.white70,
                  onPressed: () => search(),
                )
              )
            ],
          ),
        ),
        Expanded(
          child: Scrollbar( 
            child: ListView(
              padding: const EdgeInsets.all(16.0),
              scrollDirection: Axis.vertical,
              shrinkWrap: true,
              children: [
                  const IconLabel(icon: Icons.album, label: "Albums"),
                  const Divider(),
                  SingleChildScrollView(
                    padding: const EdgeInsets.all(4.0),
                    scrollDirection: Axis.horizontal,
                    child: Row(children: albumChildren),
                  ),
                  const IconLabel(icon: Icons.playlist_play_rounded, label: "Playlists"),
                  const Divider(),
                  SingleChildScrollView(
                    padding: const EdgeInsets.all(4.0),
                    scrollDirection: Axis.horizontal,
                    child: Row(children: playlistChildren),
                  ),
                  const IconLabel(icon: Icons.radio_rounded, label: "Radios"),
                  const Divider(),
                  SingleChildScrollView(
                    padding: const EdgeInsets.all(4.0),
                    scrollDirection: Axis.horizontal,
                    child: Row(children: radioChildren),
                  ),
                  const IconLabel(icon: Icons.music_note_rounded, label: "Tracks"),
                  const Divider(),
                  ...tracks.map((e) => TrackSearchResultScaffold(track: e)).toList(),
                  const SizedBox(height: 64.0,)

                ],
            )
          )
        )
      ]
    );
  }

}

enum SearchResultType {
  album,
  playlist,
  radio,
  track,
}

class LargeSearchResultScaffold extends StatelessWidget {

  final String title;
  final String id;
  final String coverId;
  final SearchResultType type;

  const LargeSearchResultScaffold({
    super.key,
    required this.title,
    required this.id,
    required this.coverId,
    required this.type,
  });

  @override
  Widget build(BuildContext context) {
    
    return InkWell(
      onTap: () {

        GetIt.I<ScreenChangeNotifier>().currentIndex = 0;
        
        if (type == SearchResultType.album || type == SearchResultType.playlist) {
          GetIt.I<LibraryViewChangeNotifier>().changeActiveWidget(CollectionView(id: id));
        } else {
          GetIt.I<LibraryViewChangeNotifier>().changeActiveWidget(RadioView(id: id));
        }

        log.fine(GetIt.I<LibraryViewChangeNotifier>().activeWidget.toStringShort());

      },
      child: Container(
        padding: const EdgeInsets.all(4.0),
        width: 128,
        //color: Colors.grey[800],
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Center(
              child: Image(image: ChimeAPI.getCover(coverId),
                width: 100,  
                height: 100,
              )
            ),
            const SizedBox(height: 6,),
            Text(title, maxLines: 1, overflow: TextOverflow.ellipsis, style: GoogleFonts.ibmPlexSans(fontSize: 14),),
            const SizedBox(height: 4,),
            BorderedChip(text: type.name[0].toUpperCase() + type.name.substring(1))
          ],
        ),
      )
    );
    
  }

}

class TrackSearchResultScaffold extends StatelessWidget {

  final SearchTrack track;

  const TrackSearchResultScaffold({
    super.key, 
    required this.track
  });

  @override
  Widget build(BuildContext context) {
    
    return ListTile(
      title: SingleChildScrollView(
          scrollDirection: Axis.horizontal,
          child: Text(track.title, maxLines: 1, style: GoogleFonts.ibmPlexSans(color: Colors.white, fontWeight: FontWeight.w500))
      ),
      subtitle: Text("${track.artist} ● ${Util.convertDuration(track.duration)}"),
      dense: true,
      contentPadding: const EdgeInsets.all(0.0),
      onTap: () { 
        GetIt.I<LibraryViewChangeNotifier>().changeActiveWidget(CollectionView(id: track.albumId)); 
        GetIt.I<ScreenChangeNotifier>().currentIndex = 0;
      },
    );
    
  }





}