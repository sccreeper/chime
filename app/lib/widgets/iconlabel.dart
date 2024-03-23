import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';

class IconLabel extends StatelessWidget {

  final IconData icon;
  final String label;

  const IconLabel({
    super.key,
    required this.icon,
    required this.label,
  });

  @override
  Widget build(BuildContext context) {
    return Row(children: [
      Expanded(
        flex: 1,
        child: Icon(icon, color: Colors.yellow[800],)
      ),
      Expanded(
        flex: 9,
        child: Text(label, style: GoogleFonts.ibmPlexSans(color: Colors.yellow[800], fontSize:16, fontWeight: FontWeight.bold)),
      )
    ],
    );
  }

}