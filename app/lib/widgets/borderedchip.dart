import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';

class BorderedChip extends StatelessWidget {

  final String text;

  const BorderedChip({super.key, required this.text});

  @override
  Widget build(BuildContext context) {
    
    return Container(
      padding: const EdgeInsets.all(3.0),
      decoration: BoxDecoration(
        border: Border.all(
          color: const Color(0xFFF9A825),
          width: 1.0,
          style: BorderStyle.solid,
        ),
        borderRadius: const BorderRadius.all(Radius.circular(14))
      ),
      child: Text(text, textAlign: TextAlign.center, style: GoogleFonts.anuphan(color: const Color(0xFFF9A825), fontSize: 14, fontWeight: FontWeight.w200),),
    );
    
  }

}

class BorderedChipButton extends StatelessWidget {
    
  final String text;
  final IconData icon;
  final Function() onTap;

  const BorderedChipButton({super.key, required this.text, required this.icon, required this.onTap});

  @override
  Widget build(BuildContext context) {
    
    return Container(
      padding: const EdgeInsets.all(3.0),
      decoration: BoxDecoration(
        border: Border.all(
          color: const Color(0xFFF9A825),
          width: 1.0,
          style: BorderStyle.solid,
        ),
        borderRadius: const BorderRadius.all(Radius.circular(14))
      ),
      
      child: InkWell(
        onTap: () async => {onTap()},
        child: Row(children: [
          Icon(icon, size: 14, color: const Color(0xFFF9A825),),
          Text(text, textAlign: TextAlign.center, style: GoogleFonts.anuphan(color: const Color(0xFFF9A825), fontSize: 14, fontWeight: FontWeight.w200),)
        ]),
      ),

    );
    
  }
}