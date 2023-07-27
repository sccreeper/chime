import 'dart:convert';
import 'dart:io' as io;
import 'package:app/api/endpoints.dart';
import 'package:app/shared.dart';
import 'package:app/shared.dart';
import 'package:app/mainscreen.dart';
import 'package:app/views/libraryview.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:http/http.dart' as http;

import 'package:app/login.dart';
import 'package:app/api/models/session.dart';
import 'package:flutter/material.dart';
import 'package:path_provider/path_provider.dart';
import 'package:logging/logging.dart';
import 'package:provider/provider.dart';

void main() {
  hierarchicalLoggingEnabled = true;

  log.level = Level.ALL;
  log.onRecord.listen((record) {

    print("${record.level.name}: ${record.time}: ${record.message}");

  });

  session = UserSession.empty();

  runApp(ChangeNotifierProvider(create: (_) => LibraryViewChangeNotifier(), child: const MaterialApp(home: MainApp()),));
  
}

class MainApp extends StatefulWidget {
  const MainApp({super.key});

    @override
    _MainAppState createState() => _MainAppState();
}

class _MainAppState extends State<MainApp> {

  late io.Directory appDocuments;

  Widget _currentView = Scaffold(backgroundColor: Colors.grey[700],);

  @override
  void initState() {
    // Check config and switch view.
    _checkConfig();
    super.initState();
  }

  @override
  Widget build(BuildContext context) {

    final ThemeData baseTheme = ThemeData.light();

    return MaterialApp(
      theme: baseTheme.copyWith( 
        primaryTextTheme: GoogleFonts.anuphanTextTheme(),
        scaffoldBackgroundColor: Colors.grey[700],
        appBarTheme: AppBarTheme(
          backgroundColor: Colors.grey[800],
          titleTextStyle: GoogleFonts.anuphan(color: Colors.white, fontSize: 24, fontWeight: FontWeight.bold)
        ),
        textTheme: baseTheme.textTheme.copyWith(
          bodySmall: GoogleFonts.anuphan(color: Colors.white, fontSize: 14),
          bodyMedium: GoogleFonts.anuphan(color: Colors.white, fontSize: 16),
          bodyLarge: GoogleFonts.anuphan(color: Colors.white, fontSize: 18),
          headlineSmall: GoogleFonts.anuphan(color: Colors.white, fontSize: 24, fontWeight: FontWeight.bold),
          titleMedium: GoogleFonts.anuphan(color: Colors.white)
        ),
        primaryColor: Colors.yellow[800],
        colorScheme: ColorScheme.fromSwatch(
          accentColor: Colors.amber[600],
          primaryColorDark: Colors.yellow[800],
        ),

        inputDecorationTheme: InputDecorationTheme(
          labelStyle: GoogleFonts.anuphan(color: Colors.white70),
          hintStyle: GoogleFonts.anuphan(color: Colors.white70),
          enabledBorder: UnderlineInputBorder(
            borderSide: BorderSide(color: Colors.white54),
          ),
          focusedBorder: UnderlineInputBorder(
            borderSide: BorderSide(color: Colors.orange),
          )
        ),

        elevatedButtonTheme: ElevatedButtonThemeData(
          style: ButtonStyle(
            foregroundColor: MaterialStateColor.resolveWith((states) => Colors.white),
            overlayColor: MaterialStateColor.resolveWith((states) => Colors.deepOrange),
            backgroundColor: MaterialStateColor.resolveWith((states) => Colors.orange),
            textStyle: MaterialStateTextStyle.resolveWith((states) => GoogleFonts.anuphan(color: Colors.white))
          )
        )

      ),
      home: _currentView
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
        _currentView = LoginScreen();
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
        
        setState(() {
          _currentView = MainScreen();
        });

      } else {
        log.fine("Session does not exist, continuing to login screen");
        
        setState(() {
          _currentView = LoginScreen();          
        });

      }


    } else {

      log.fine("Config does not exist, creating config file");

      // Create directory and config file
      await io.Directory(appDocuments.path).create();
      await io.File("${appDocuments.path}/config.json").create();

      io.File('${appDocuments.path}/config.json').writeAsString(
        jsonEncode(UserSession.empty().toJson()));

      // Push login screen

      setState(() {
        _currentView = LoginScreen();
      });

    }

    

  }
}
