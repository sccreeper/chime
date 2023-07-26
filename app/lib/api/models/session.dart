class UserSession {
  final String sessionID;
  final String username;
  final String sessionBase64;
  final String serverOrigin;

  UserSession({ 
    required this.sessionID, 
    required this.username, 
    required this.sessionBase64, 
    required this.serverOrigin});

  factory UserSession.fromJSON(Map<String, dynamic> json) => UserSession(
     sessionID: json["session_id"],
     username: json["username"],
     sessionBase64: json["session"],
     serverOrigin: json["server_origin"],
  );

  factory UserSession.empty() => UserSession(
    sessionID: "", 
    username: "", 
    sessionBase64: "", 
    serverOrigin: "");

  Map<String, dynamic> toJson() => {
    "username":username,
    "session_id":sessionID,
    "session":sessionBase64,
    "server_origin":serverOrigin,
  };

}