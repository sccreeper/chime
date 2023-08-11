import 'package:app/api/api.dart';
import 'package:app/api/models/collections.dart';
import 'package:app/shared.dart';
import 'package:app/views/collectionview.dart';
import 'package:app/views/radioview.dart';
import 'package:app/widgets/iconlabel.dart';
import 'package:app/widgets/loadingspinner.dart';
import 'package:flutter/material.dart';
import 'package:get_it/get_it.dart';

class LibrayView extends StatefulWidget {

    const LibrayView({super.key});

    @override
    _LibaryViewState createState() => _LibaryViewState();

}


class _LibaryViewState extends State<LibrayView> {

  final _service = GetIt.I<LibraryViewChangeNotifier>();

  @override
  void initState() {
    _service.addListener(updateView);
    _fetchLibaryData();
    super.initState();
  }

  void updateView() {
    if (mounted) {
      setState(() {}); 
    }
  }

  void _fetchLibaryData() async {

    // To prevent switching back if collection view is already present.
    if (GetIt.I<LibraryViewChangeNotifier>().activeWidget.runtimeType == CollectionView) {
      return;
    }

    final Library lib = await ChimeAPI.getLibary();
    
    List<Widget> content = [];

    content.add(const IconLabel(icon: Icons.album, label: "Albums"));
    content.add(const Divider());

    for (var album in lib.albums) {      
      content.addAll(
        [LibaryItem(id: album.id, type: LibaryItemType.album, name: album.name,), const Divider()]
      );
    }

    content.add(const IconLabel(icon: Icons.playlist_play_rounded, label: "Playlists"));
    content.add(const Divider());

    for (var playlist in lib.playlists) {      
      content.addAll(
        [LibaryItem(id: playlist.id, type: LibaryItemType.playlist, name: playlist.name,), const Divider()]
      );
    }

    content.add(const IconLabel(icon: Icons.radio_rounded, label: "Radios"));
    content.add(const Divider());

    for (var radio in lib.radios) {      
      content.addAll(
        [LibaryItem(id: radio.id, type: LibaryItemType.radio, name: radio.name,), const Divider()]
      );
      
    }

    GetIt.I<LibraryViewChangeNotifier>().changeActiveWidget(ListView(
      children: content,
    ));

    setState(() {});
    
  }

  @override
  Widget build(BuildContext context) {
    return WillPopScope(
      onWillPop: () async {

        if (GetIt.I<LibraryViewChangeNotifier>().activeWidget is ListView) {
          return false;
        } else {
          GetIt.I<LibraryViewChangeNotifier>().changeActiveWidget(LoadingSpinner());
          _fetchLibaryData();
          return false;
        }

      }, 
      child: Container(
        padding: EdgeInsets.all(16.0),
        child: GetIt.I<LibraryViewChangeNotifier>().activeWidget,
      )
    );
  }

}

class LibraryViewChangeNotifier extends ChangeNotifier {
  
  Widget activeWidget = LoadingSpinner();

  void changeActiveWidget(Widget newWidget) {
    activeWidget = newWidget;
    notifyListeners();
  }

}

class LibaryItem extends StatefulWidget {

    final String name;
    final String id;
    final LibaryItemType type;

    LibaryItem({super.key,required this.id, required this.type, required this.name});

    @override
    _LibraryItemState createState() => _LibraryItemState();

}

class _LibraryItemState extends State<LibaryItem> {
  

  @override
  Widget build(BuildContext context) {

    return InkWell(
      child: Text(widget.name, style: Theme.of(context).textTheme.bodyMedium,),
      onTap: () {
        log.fine("Opening ${widget.type.name} ${widget.id}");

        if (widget.type == LibaryItemType.album || widget.type == LibaryItemType.playlist) {
          
          GetIt.I<LibraryViewChangeNotifier>().changeActiveWidget(CollectionView(id: widget.id));

        } else {

          GetIt.I<LibraryViewChangeNotifier>().changeActiveWidget(RadioView(id: widget.id));
          
        }

      },
    );

  }

}