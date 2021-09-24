package main

type Settings struct {
	//TODO
	RootURL			string
	SessionURL		string
	ProfileURL		string
	BoardsURL		string

	ServerAddress	string
}

func InitSettings() (settings Settings) {
	//TODO
	settings = Settings{
		RootURL: "/api",
		SessionURL: "/sessions",
		ProfileURL: "/profile",
		BoardsURL: "/boards",
		ServerAddress: "0.0.0.0:8080",
	}
	return
}