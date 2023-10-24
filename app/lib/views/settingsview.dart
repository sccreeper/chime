import 'dart:math';

import 'package:app/api/models/session.dart';
import 'package:app/shared.dart';
import 'package:app/widgets/iconlabel.dart';
import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'dart:io' as io;

import 'package:path_provider/path_provider.dart';
import 'package:restart_app/restart_app.dart';

class SettingsView extends StatefulWidget {

  const SettingsView({super.key});

  @override
  SettingsViewState createState() => SettingsViewState();

}

class SettingsViewState extends State<SettingsView> {

  String _documentsDirectory = "";
  String _downloadSize = "";

  @override
  void initState() {
    
    () async {

      _documentsDirectory = (await getApplicationDocumentsDirectory()).path;

      int size = Util.directorySize(_documentsDirectory);
      _downloadSize = "${(size / pow(10, 6)).toStringAsFixed(2)} mb";

      setState(() {});

    }();

  

    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    
    return Padding(
      padding: const EdgeInsets.all(8.0),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const IconLabel(icon: Icons.account_circle, label: "Account"),
          const Divider(),
          RichText(
            text: TextSpan(
              children: [
                const TextSpan(text: "Logged in as: ", style: TextStyle(fontWeight: FontWeight.bold)),
                TextSpan(text: session.username)
              ]
            )
          ),
          RichText(
            text: TextSpan(
              children: [
                const TextSpan(text: "Connected to Chime on: ", style: TextStyle(fontWeight: FontWeight.bold)),
                TextSpan(text: Uri.parse(session.serverOrigin).host)
              ]
          ), textAlign: TextAlign.left,),
          TextButton(
            onPressed: () => {

              showDialog(context: context, builder: (BuildContext context) => AlertDialog(
                backgroundColor: Colors.grey[800],
                title: const Text("Logout"),
                titleTextStyle: GoogleFonts.anuphan(color: Colors.white, fontSize: 14.0),
                contentTextStyle: GoogleFonts.anuphan(),
                content: const Text("Are you sure you want to log out? This will delete all your downloaded content from your device."),
                actions: [
                  TextButton(
                    onPressed: () => Navigator.pop(context, "Cancel"),
                    child: const Text("Cancel")
                  ),
                  TextButton(
                    onPressed: () async {

                      // Delete files.

                      getApplicationDocumentsDirectory().then((appDir) {
                        
                        io.File("${appDir.path}/config.json").delete().then((e) {

                          session = UserSession.empty();
                          Navigator.pop(context, "Logout");
                          
                          Restart.restartApp();

                        });

                      });

                    }, 
                    child: const Text("Logout")
                  )
                ],
              ))

            }, 
            child: const Text("Logout")
          ),
          const Divider(),
          const IconLabel(icon: Icons.download_rounded, label: "Downloads"),
          const Divider(),
          RichText(
            text: TextSpan(
              children: [
                const TextSpan(text: "Downloads stored in: ", style: TextStyle(fontWeight: FontWeight.bold)),
                TextSpan(text: _documentsDirectory)
              ]
            )
          ),
          RichText(
            text: TextSpan(
              children: [
                const TextSpan(text: "Total download size: ", style: TextStyle(fontWeight: FontWeight.bold)),
                TextSpan(text: _downloadSize)
              ]
            )
          ),
          TextButton(onPressed: () => {}, child: const Text("Clear downloads"))

        ],
      )
    );
    
  }

}