import 'package:app/backend/shared.dart';
import 'package:flutter/material.dart';

class MainScreen extends StatefulWidget {

    const MainScreen({super.key});

    @override
    _MainScreenState createState() => _MainScreenState();

}

class _MainScreenState extends State<MainScreen> {

  @override
  Widget build(BuildContext context) {

    log.fine("Switched to main screen, fetching data");

    return const Scaffold(
      body: Center (child: Text("Logged in.")),
    );

  }

}