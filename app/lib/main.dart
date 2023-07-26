import 'dart:convert';
import 'dart:io' as io;
import 'package:app/api/endpoints.dart';
import 'package:app/backend/login.dart';
import 'package:app/backend/shared.dart';
import 'package:app/mainscreen.dart';
import 'package:http/http.dart' as http;

import 'package:app/login.dart';
import 'package:app/api/models/session.dart';
import 'package:flutter/material.dart';
import 'package:path_provider/path_provider.dart';
import 'package:logging/logging.dart';

void main() {
  hierarchicalLoggingEnabled = true;

  log.level = Level.ALL;
  log.onRecord.listen((record) {

    print("${record.level.name}: ${record.time}: ${record.message}");

  });

  session = UserSession.empty();

  runApp(const MaterialApp(
    home: MainApp(),
  ));
}

class MainApp extends StatefulWidget {
  const MainApp({super.key});

    @override
    _MainAppState createState() => _MainAppState();
}

class _MainAppState extends State<MainApp> {

  late io.Directory appDocuments;
  late UserSession userSession;


  
  @override
  void initState() {
    // Check config and switch view.
    _checkConfig();
    super.initState();
  }
  
  @override
  Widget build(BuildContext context) {

    return const MaterialApp(
      home: Scaffold(
        backgroundColor: Colors.black,
      ),
    );
  }

  Future _checkConfig() async {
    
    appDocuments = await getApplicationDocumentsDirectory();

    log.fine("Checking if config file exists");

    if (io.File("${appDocuments.path}/config.json").existsSync()) {

      log.fine("Config file exists, checking formatting");
      
      // Check session

      var config_file = io.File("${appDocuments.path}/config.json").readAsStringSync();

      try {
        session = new UserSession.fromJSON(jsonDecode(config_file));

        if (session.serverOrigin == "" || session.sessionBase64 == "" || session.sessionID == "" || session.username == "") {  
          log.warning("Session file has empty fields.");
          throw FormatException("Session fields empty");
        }

      } on FormatException {
        log.warning("Error in config file formatting, changing to login screen");
        Navigator.pushReplacement(context, MaterialPageRoute(builder: (context) => const LoginScreen()));
        return;
      }
      
      log.fine("Checking session with server and authenticating");

      // Try logging in with stored json.

      // ignore: unnecessary_brace_in_string_interps
      var resp = await http.get(Uri.parse("${session.serverOrigin}${apiAuthSessionExists}/${session.sessionID}"));
      var respJson = jsonDecode(resp.body);

      log.fine(resp.body);

      if (respJson["status"] == "exists") {
        log.fine("Session exists, continuing to main screen");
        Navigator.pushReplacement(context, MaterialPageRoute(builder: (context) => const MainScreen()));

      } else {
        log.fine("Session does not exist, continuing to login screen");
        Navigator.pushReplacement(context, MaterialPageRoute(builder: (context) => const LoginScreen()));
      }


    } else {

      log.fine("Config does not exist, creating config file");

      // Create directory and config file
      await io.Directory(appDocuments.path).create();
      await io.File("${appDocuments.path}/config.json").create();

      io.File('${appDocuments.path}/config.json').writeAsString(
        jsonEncode(UserSession.empty().toJson()));

      // Push login screen

      Navigator.pushReplacement(context, MaterialPageRoute(builder: (context) => const LoginScreen()));

    }

    

  }
}
