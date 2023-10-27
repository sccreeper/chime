import 'dart:convert';
import 'dart:io' as io;

import 'package:app/api/endpoints.dart';
import 'package:app/api/models/session.dart';
import 'package:app/shared.dart';
import 'package:app/mainscreen.dart';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:path_provider/path_provider.dart';

class LoginScreen extends StatefulWidget {
  const LoginScreen({super.key});

  @override
  LoginScreenState createState() => LoginScreenState();
}

class LoginScreenState extends State<LoginScreen> {
  String loginPassword = "";
  String loginUsername = "";
  String loginServerAddress = "";

  String errorText = "";

  @override
  Widget build(BuildContext context) {
    log.fine("Building login screen...");

    return Scaffold(
        appBar: AppBar(
          title: const Text("Chime"),
        ),
        body: Padding(
          padding: const EdgeInsets.all(16.0),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.start,
        children: [
        const Text("Chime"),
        Text(
          errorText,
          style: const TextStyle(color: Colors.red),
        ),
        TextField(
          onChanged: (value) => {loginUsername = value},
          decoration: const InputDecoration(hintText: "Username"),
        ),
        TextField(
            onChanged: (value) => {loginPassword = value},
            decoration: const InputDecoration(hintText: "Password"),
            obscureText: true),
        TextField(
            onChanged: (value) => {loginServerAddress = value},
            decoration: const InputDecoration(hintText: "Server Address")),
        ElevatedButton(
          style: Theme.of(context).elevatedButtonTheme.style,
          onPressed: () => {_login()},
          child: const Text("Login"),
        ),
      ]),
    ));
  }

  void _login() async {
    log.fine("Logging in...");

    final url = Uri.parse("$loginServerAddress$apiAuth");
    var request = http.MultipartRequest("POST", url);

    request.fields["u"] = loginUsername;
    request.fields["p"] = loginPassword;

    final response = await request.send();
      
    String jsonString = await response.stream.bytesToString();

    log.fine("Recieved JSON: $jsonString");
    Map<String,dynamic> responseJson = jsonDecode(jsonString);

    String sessionB64 =
        base64Encode(utf8.encode(jsonEncode(responseJson["session"])))
        .replaceAll(RegExp(r"/"), "-")
        .replaceAll(r"+", "_")
        .replaceAll(
            RegExp(
              r"=",
            ),
            ".");

    log.fine(responseJson["status"]);

    if (responseJson["status"].toString() == "correct") {
      session = UserSession(
        sessionID: responseJson["session"]["session_id"], 
        username: responseJson["user"]["username"], 
        sessionBase64: sessionB64, 
        serverOrigin: url.origin);

      // Cache session to disk.
      io.Directory appDocuments = await getApplicationDocumentsDirectory();
      io.File("${appDocuments.path}/config.json")
          .writeAsStringSync(jsonEncode(session.toJson()));

      log.fine("Writing session JSON: ${jsonEncode(session.toJson())}");

      connected = true;

      // Change login screen
      Navigator.pushReplacement(context, MaterialPageRoute(builder: (context) => const MainScreen()));
    } else {

      setState(() {
        errorText = "Incorrect password";
      });
      
    }
  }
}
