import 'package:app/api/api.dart';
import 'package:app/api/endpoints.dart';
import 'package:app/api/models/radio.dart';
import 'package:app/player.dart';
import 'package:app/shared.dart';
import 'package:app/widgets/loadingspinner.dart';
import 'package:flutter/material.dart';
import 'package:get_it/get_it.dart';
import 'package:google_fonts/google_fonts.dart';

class RadioView extends StatefulWidget {

  final String id;

  const RadioView({super.key, required this.id});

  @override
  RadioViewState createState() => RadioViewState();

}

class RadioViewState extends State<RadioView> {

  @override
  void initState() {
    
    GetIt.I<RadioViewLoadedNotifier>().addListener(() {

      if (mounted) {
        setState(() {});
      }

    });

    _getInfo();

    super.initState();
    
  }

  void _getInfo() async {

    RadioModel radio = await ChimeAPI.getRadio(widget.id);

    GetIt.I<RadioViewLoadedNotifier>().setCurrentRadioWidget(
      RadioScaffold(radio: radio)
    );

  }
  
  @override
  Widget build(BuildContext context) {
    return GetIt.I<RadioViewLoadedNotifier>().currentRadioWidget;
  }

}

class RadioScaffold extends StatelessWidget {
  
  final RadioModel radio;

  const RadioScaffold({super.key, required this.radio});
  
  @override
  Widget build(BuildContext context) {
    
    return Column(

        mainAxisAlignment: MainAxisAlignment.center,
        crossAxisAlignment: CrossAxisAlignment.center,

        children: [

          Image(image: ChimeAPI.getCover(radio.coverId), width: 256, height: 256,),

          const SizedBox(height: 32,),

          Text(radio.name, textAlign: TextAlign.center, style: GoogleFonts.anuphan(color: Colors.white, fontSize: 24.0, fontWeight: FontWeight.bold)),

          const SizedBox(height: 16,),

          ElevatedButton.icon(
            onPressed: () => Player.playRadio(radio), 
            icon: const Icon(Icons.play_arrow_rounded), 
            label: const Text("Play")
          )
        ],
    );

  }
}

class RadioViewLoadedNotifier extends ChangeNotifier {

  Widget currentRadioWidget = LoadingSpinner();

  void setCurrentRadioWidget(Widget val) {
    currentRadioWidget = val;
    notifyListeners();
  }

}