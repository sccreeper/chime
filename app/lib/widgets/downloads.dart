import 'package:app/api/downloads.dart';
import 'package:app/shared.dart';
import 'package:app/widgets/borderedchip.dart';
import 'package:flutter/material.dart';
import 'package:get_it/get_it.dart';

class CollectionDownloadButton extends StatefulWidget {
  final String id;

  const CollectionDownloadButton({super.key, required this.id});

  @override
  State<StatefulWidget> createState() => _CollectionDownloadButtonState();
}

class _CollectionDownloadButtonState extends State<CollectionDownloadButton> {

  bool isDownloaded = false;

  @override
  void initState() {
    GetIt.I<DownloadNotifier>().addListener(() {

      if (mounted) {
        setState(() {});
      }

    });

    (() async {

      isDownloaded = await dbMgr.collectionRecordExists(widget.id);
      setState(() {});

      log.fine(isDownloaded);
    
    })();

    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    
    bool isDownloading = GetIt.I<DownloadNotifier>().downloading && 
      widget.id == GetIt.I<DownloadNotifier>().downloadingId && 
      GetIt.I<DownloadNotifier>().downloadType == DownloadType.collection;

    if (isDownloading) {
      return BorderedChipButton(
        text: "Downloading ${GetIt.I<DownloadNotifier>().downloadingProgress}/${GetIt.I<DownloadNotifier>().downloadingTotal}", 
        icon: Icons.stop, 
        onTap: () => {}
      );      
    } else if(!isDownloaded) {

      return BorderedChipButton(
        text: "Download", 
        icon: Icons.download, 
        onTap: () => { DownloadManager.downloadCollection(widget.id) }
      );
    } else {
      return BorderedChipButton(
        text: "Delete download", 
        icon: Icons.delete, 
        onTap: () { DownloadManager.deleteCollection(widget.id); setState(() {});}
      );
    }
    
  }

}