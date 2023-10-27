class PingResult {

  final String serverId;
  final String message;
  final bool successful;

  PingResult({
    required this.serverId,
    required this.message,
    required this.successful,
  });

  factory PingResult.fromJson(Map<String, dynamic> json) => PingResult(
    serverId: json["server_id"],
    message: json["message"],
    successful: true,
  );

  Map<String, dynamic> toJson() => {
    "server_id":serverId,
    "message":message,
    "successful":successful,
  };

}