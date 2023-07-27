import 'dart:typed_data';

import 'package:app/api/api.dart';
import 'package:app/api/models/collections.dart';
import 'package:app/shared.dart';
import 'package:app/widgets/loadingspinner.dart';
import 'package:flutter/material.dart';
import 'package:app/api/endpoints.dart';

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
    Uint8List image = await ChimeAPI.getCover(collection.coverId);

    setState(() {
      _childWidget = CollectionScaffold(
        collection: collection, 
        coverBytes: image
      );
    });

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

      padding: EdgeInsets.all(8.0),
      child: Column(children: [
          Image.memory(coverBytes, width: 300, height: 300,),
          Text(collection.title, style: Theme.of(context).textTheme.bodyLarge,)
      ]),

    );
    
  }


}