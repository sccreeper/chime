import 'package:flutter/material.dart';

// Centered loading spinner,
class LoadingSpinner extends StatelessWidget {

  @override
  Widget build(BuildContext context) {
    return Center(
    child: 
      CircularProgressIndicator(
        color: Colors.yellow[800],
        value: null,
      )
    );
  }

}