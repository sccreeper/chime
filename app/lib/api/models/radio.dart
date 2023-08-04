class RadioModel {

  final String name;
  final String url;
  final String description;
  final String coverId;

  RadioModel({
    required this.name,
    required this.url,
    required this.description,
    required this.coverId,
  });

  factory RadioModel.fromJson(Map<String, dynamic> json) => RadioModel(
    name: json["name"], 
    url: json["url"], 
    description: json["description"], 
    coverId: json["cover_id"]
  );

  Map<String, dynamic> toJson() => {
    "name":name,
    "url":url,
    "description":description,
    "cover_id":coverId
  };

}