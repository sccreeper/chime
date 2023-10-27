import 'dart:convert';
import 'dart:io' as io;
import 'package:app/api/api.dart';
import 'package:app/api/downloads.dart';
import 'package:app/api/endpoints.dart';
import 'package:app/player.dart';
import 'package:app/shared.dart';
import 'package:app/mainscreen.dart';
import 'package:app/views/libraryview.dart';
import 'package:app/views/radioview.dart';
import 'package:connectivity_plus/connectivity_plus.dart';
import 'package:flutter/services.dart';
import 'package:get_it/get_it.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:http/http.dart' as http;

import 'package:app/login.dart';
import 'package:app/api/models/session.dart';
import 'package:flutter/material.dart';
import 'package:path_provider/path_provider.dart';
import 'package:logging/logging.dart';

void main() async {
  hierarchicalLoggingEnabled = true;

  log.level = Level.ALL;
  log.onRecord.listen((record) {

    print("${record.level.name}: ${record.time}: ${record.message}");

  });

  WidgetsFlutterBinding.ensureInitialized();

  session = UserSession.empty();
  dbMgr = DownloadDatabaseManager();
  Player.init();

  // Register change notifiers.
  GetIt.I.registerSingleton<LibraryViewChangeNotifier>(LibraryViewChangeNotifier());
  GetIt.I.registerSingleton<PlayerStatusNotifier>(PlayerStatusNotifier());
  GetIt.I.registerSingleton<RadioViewLoadedNotifier>(RadioViewLoadedNotifier());
  GetIt.I.registerSingleton<ScreenChangeNotifier>(ScreenChangeNotifier());
  GetIt.I.registerSingleton<ActiveMainViewNotifier>(ActiveMainViewNotifier());
  GetIt.I.registerSingleton<DownloadNotifier>(DownloadNotifier());

  // Load no_cover image

  docDirectory = (await getApplicationDocumentsDirectory()).path;

  runApp(const MaterialApp(home: MainApp()));
  
}

class MainApp extends StatefulWidget {
  const MainApp({super.key});

    @override
    MainAppState createState() => MainAppState();
}

class MainAppState extends State<MainApp> {

  late io.Directory appDocuments;

  // Widget _currentView = Scaffold(backgroundColor: Colors.grey[700],);

  @override
  void initState() {
    // Check config and switch view.
    
    GetIt.I<ActiveMainViewNotifier>().addListener(updateView);

    mainInit();
    super.initState();
  }

  void updateView() {

    if (mounted) {
      log.fine("Updating view");

      setState(() {});
    }
  }

  @override
  Widget build(BuildContext context) {

    final ThemeData baseTheme = ThemeData.light();

    return MaterialApp(
      theme: baseTheme.copyWith( 
        primaryTextTheme: GoogleFonts.anuphanTextTheme(),
        scaffoldBackgroundColor: const Color.fromARGB(255, 43, 43, 43),
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
        ),

        inputDecorationTheme: InputDecorationTheme(
          labelStyle: GoogleFonts.anuphan(color: Colors.white70),
          hintStyle: GoogleFonts.anuphan(color: Colors.white70),
          enabledBorder: const UnderlineInputBorder(
            borderSide: BorderSide(color: Colors.white54),
          ),
          focusedBorder: const UnderlineInputBorder(
            borderSide: BorderSide(color: Colors.orange),
          )
        ),

        elevatedButtonTheme: ElevatedButtonThemeData(
          style: ButtonStyle(
            foregroundColor: MaterialStateColor.resolveWith((states) => Colors.white),
            overlayColor: MaterialStateColor.resolveWith((states) => const Color(0xFFF57F17)),
            backgroundColor: MaterialStateColor.resolveWith((states) => const Color(0xFFF9A825)),
            textStyle: MaterialStateTextStyle.resolveWith((states) => GoogleFonts.anuphan(color: Colors.white))
          )
        ),

        textButtonTheme: TextButtonThemeData(
          style: ButtonStyle(
            textStyle: MaterialStateTextStyle.resolveWith((states) => GoogleFonts.anuphan(color: const Color(0xFFF9A825))),
            foregroundColor: MaterialStateColor.resolveWith((states) => const Color(0xFFF9A825)),
            overlayColor: MaterialStateColor.resolveWith((states) => const Color(0xFFF57F17))
          )
        ),

        dividerColor: Colors.white70,

        progressIndicatorTheme: ProgressIndicatorThemeData(
          color: Colors.amber,
          linearTrackColor: Colors.amber[50],
          linearMinHeight: 2.0
        ),

        sliderTheme: SliderThemeData(
          trackHeight: 2,
          activeTrackColor: Colors.amber,
          inactiveTrackColor: Colors.amber[50],
          thumbColor: Colors.yellow[800],
          overlayColor: Colors.grey.withOpacity(0.25),
          
        )

      ),
      home: GetIt.I<ActiveMainViewNotifier>().widget
    );
  }

  Future mainInit() async {

    // Check config
    
    appDocuments = await getApplicationDocumentsDirectory();

    log.fine("Checking if config file exists");

    if (io.File("${appDocuments.path}/config.json").existsSync()) {

      log.fine("Config file exists, checking formatting");
      
      // Check session

      var configFile = io.File("${appDocuments.path}/config.json").readAsStringSync();

      try {
        session = UserSession.fromJSON(jsonDecode(configFile));

        if (session.serverOrigin == "" || session.sessionBase64 == "" || session.sessionID == "" || session.username == "") {  
          log.warning("Session file has empty fields.");
          throw const FormatException("Session fields empty");
        }

      } on FormatException {
        log.warning("Error in config file formatting, changing to login screen");
        GetIt.I<ActiveMainViewNotifier>().widget = const LoginScreen();
        return;
      }

      // Check if connected to internet

      var connectivity = await (Connectivity().checkConnectivity());

      if (connectivity == ConnectivityResult.none) {
        connected = false;
      } else {

        var result = await io.InternetAddress.lookup(Uri.parse(session.serverOrigin).host);
        
        if (result.isNotEmpty && result[0].rawAddress.isNotEmpty) {
          connected = true;
        } else {
          
          // Check if server is "alive"

          var ping = await ChimeAPI.ping(session.serverOrigin);
          connected = ping.successful;

        }
      }

      if (connected) {
        // Auth if we are
      
        log.fine("Checking session with server and authenticating");

        // Try logging in with stored json.

        // ignore: unnecessary_brace_in_string_interps
        var resp = await http.get(Uri.parse("${session.serverOrigin}${apiAuthSessionExists}/${session.sessionID}"));
        var respJson = jsonDecode(resp.body);

        log.fine(resp.body);

        if (respJson["status"] == "exists") {
          log.fine("Session exists, continuing to main screen");
          
          GetIt.I<ActiveMainViewNotifier>().widget = const MainScreen();

        } else {
          log.fine("Session does not exist, continuing to login screen");
          
          GetIt.I<ActiveMainViewNotifier>().widget = const LoginScreen();

        } 
      
      } else {
        
        log.fine("Unable to establish session with server, reverting to offline mode");

        GetIt.I<ActiveMainViewNotifier>().widget = const MainScreen();
      
      }


    } else {

      log.fine("Config does not exist, creating config file");

      // Create directory and config file
      await io.Directory(appDocuments.path).create();
      await io.File("${appDocuments.path}/config.json").create();

      io.File('${appDocuments.path}/config.json').writeAsString(
        jsonEncode(UserSession.empty().toJson()));

      // Push login screen

      GetIt.I<ActiveMainViewNotifier>().widget = const LoginScreen();

    }

  }
}

class ActiveMainViewNotifier extends ChangeNotifier {

  Widget _widget = Scaffold(backgroundColor: Colors.grey[700],);
  Widget get widget => _widget;

  set widget(Widget val) {
    _widget = val;
    notifyListeners();
  }

}