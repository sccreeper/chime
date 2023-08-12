import 'package:app/shared.dart';
import 'package:app/views/dockedplayer.dart';
import 'package:app/views/libraryview.dart';
import 'package:app/views/searchview.dart';
import 'package:app/views/settingsview.dart';
import 'package:flutter/material.dart';
import 'package:get_it/get_it.dart';
import 'package:google_fonts/google_fonts.dart';

class MainScreen extends StatefulWidget {

    const MainScreen({super.key});

    @override
    _MainScreenState createState() => _MainScreenState();

}

class _MainScreenState extends State<MainScreen> {

  final List<String> _viewNames = ["Library", "Search", "Settings"];

  static const List<Widget> _widgetOptions = <Widget>[

     LibrayView(),
     SearchView(),
     SettingsView(),

  ];

  void _onItemTapped(int index) {
      GetIt.I<ScreenChangeNotifier>().currentIndex = index;
  }

  @override
  void initState() {
    
    GetIt.I<ScreenChangeNotifier>().addListener(() {
      setState(() {});
    });

    super.initState();
  }

  @override
  Widget build(BuildContext context) {

    log.fine("Switched to main screen, fetching data");

    return Scaffold(
      appBar: AppBar(
        title: Text(_viewNames[GetIt.I<ScreenChangeNotifier>().currentIndex]),
      ),
      body: Center (child: _widgetOptions.elementAt(GetIt.I<ScreenChangeNotifier>().currentIndex)),
      floatingActionButton: DockedPlayer(),
      floatingActionButtonLocation: FloatingActionButtonLocation.centerFloat,
      bottomNavigationBar: BottomNavigationBar(
        items: const [
          BottomNavigationBarItem(
            icon: Icon(Icons.library_music),
            label: "Library",
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.search),
            label: "Search",
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.settings),
            label: "Settings"
          )
        ],
        currentIndex: GetIt.I<ScreenChangeNotifier>().currentIndex,
        backgroundColor: Colors.grey[800],
        unselectedItemColor: Colors.white70,
        selectedItemColor: Colors.yellow[800],
        selectedLabelStyle: GoogleFonts.anuphan(),
        unselectedLabelStyle: GoogleFonts.anuphan(),
        onTap: _onItemTapped,
      ),
    );

  }

}

class ScreenChangeNotifier extends ChangeNotifier {

  int _currentIndex = 0;

  int get currentIndex => _currentIndex;

  set currentIndex(int val) {

    _currentIndex = val;
    notifyListeners();

  }

}