import 'package:app/api/api.dart';
import 'package:app/api/models/collections.dart';
import 'package:app/shared.dart';
import 'package:app/views/collectionview.dart';
import 'package:app/widgets/iconlabel.dart';
import 'package:app/widgets/loadingspinner.dart';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

class LibrayView extends StatefulWidget {

    const LibrayView({super.key});

    @override
    _LibaryViewState createState() => _LibaryViewState();

}


class _LibaryViewState extends State<LibrayView> {

  Widget _childWidget = LoadingSpinner();

  @override
  void initState() {
    _fetchLibaryData();
    super.initState();
  }

  void _fetchLibaryData() async {

    final Library lib = await ChimeAPI.getLibary();
    
    List<Widget> content = [];

    content.add(const IconLabel(icon: Icons.album, label: "Albums"));
    content.add(const Divider());

    for (var album in lib.albums) {      
      content.addAll(
        [LibaryItem(id: album.id, type: LibaryItemType.album, name: album.name,), const Divider()]
      );
    }

    content.add(const IconLabel(icon: Icons.playlist_play, label: "Playlists"));
    content.add(const Divider());

    for (var playlist in lib.playlists) {      
      content.addAll(
        [LibaryItem(id: playlist.id, type: LibaryItemType.playlist, name: playlist.name,), const Divider()]
      );
    }

    content.add(const IconLabel(icon: Icons.radio, label: "Radios"));
    content.add(const Divider());

    for (var radio in lib.radios) {      
      content.addAll(
        [LibaryItem(id: radio.id, type: LibaryItemType.radio, name: radio.name,), const Divider()]
      );
      
    }

    Provider.of<LibraryViewChangeNotifier>(context, listen: false).changeActiveWidget(ListView(
      children: content,
    ));

    setState(() {});
    
  }

  @override
  Widget build(BuildContext context) {
    return WillPopScope(
      onWillPop: () async {

        if (context.read<LibraryViewChangeNotifier>().activeWidget is ListView) {
          return false;
        } else {
          Provider.of<LibraryViewChangeNotifier>(context, listen: false).changeActiveWidget(LoadingSpinner());
          _fetchLibaryData();
          return false;
        }

      }, 
      child: Container(
        padding: EdgeInsets.all(16.0),
        child: context.watch<LibraryViewChangeNotifier>().activeWidget,
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
          Provider.of<LibraryViewChangeNotifier>(context, listen: false).changeActiveWidget(CollectionView(id: widget.id)); 
        }

      },
    );

  }

}