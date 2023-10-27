import "dart:io";
import "dart:typed_data";

import "package:app/api/downloads.dart";
import "package:app/api/models/session.dart";
import "package:logging/logging.dart";

var log = Logger("Chime");
late UserSession session;
late DownloadDatabaseManager dbMgr;
late String docDirectory;
bool connected = false;

String currentCollection = "";

class Util {

  static String convertDuration(double duration) {

    int minutes = (duration / 60).floor();
    int seconds = (duration - (minutes * 60)).floor();

    return "${minutes.toString().padLeft(2, "0")}:${seconds.toString().padLeft(2, "0")}";

  }

  static String convertDurationVerbose(double duration) {

    int hours = (duration / 3600).floor();
    int minutes = ((duration % 3600) / 60).floor();

    if (hours == 0) {
      return "${minutes.toString()} min";
    } else {
      return "${hours.toString()} hr ${minutes.toString()} min";
    }

  }

  static int directorySize(String path) {

    int size = 0;
    var dir = Directory(path);

    dir.listSync(recursive: true, followLinks: false).forEach((item) {

      if (item is File) {
       size += item.lengthSync(); 
      }

    });

    return size;

  }

}
