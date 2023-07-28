import 'package:app/api/api.dart';
import 'package:app/api/models/collections.dart';
import 'package:app/shared.dart';
import 'package:app/views/libraryview.dart';
import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';

class MainScreen extends StatefulWidget {

    const MainScreen({super.key});

    @override
    _MainScreenState createState() => _MainScreenState();

}

class _MainScreenState extends State<MainScreen> {

  int _selectedIndex = 0;
  final List<String> _viewNames = ["Library", "Search", "Settings"];

  static const List<Widget> _widgetOptions = <Widget>[

     LibrayView(),
     Text("Index 1: Search"),
     Text("Index 2: Settings")

  ];

  void _onItemTapped(int index) {
    setState(() {
      _selectedIndex = index;
    });
  }

  @override
  Widget build(BuildContext context) {

    log.fine("Switched to main screen, fetching data");

    return Scaffold(
      appBar: AppBar(
        title: Text(_viewNames[_selectedIndex]),
      ),
      body: Center (child: _widgetOptions.elementAt(_selectedIndex)),
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
        currentIndex: _selectedIndex,
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