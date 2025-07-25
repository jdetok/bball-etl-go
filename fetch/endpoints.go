package main

var leagueGameLog = GetReq{
	Host:     HOST,
	Headers:  HDRS,
	Endpoint: "/stats/leaguegamelog",
	Params: []Pair{
		{"LeagueID", "00"},
		{"Season", "2024-25"},
		{"SeasonType", "Regular+Season"},
		{"Counter", "0"},
		{"PlayerOrTeam", "T"},
		{"Sorter", "DATE"},
		{"DateFrom", ""},
		{"DateTo", ""},
		{"Direction", "DESC"},
	},
}

var commonPlayerInfo = GetReq{
	Host:     HOST,
	Headers:  HDRS,
	Endpoint: "/stats/commonplayerinfo",
	Params:   []Pair{{"LeagueID", "10"}, {"PlayerID", "2544"}},
}

var playerAwards = GetReq{
	Host:     HOST,
	Headers:  HDRS,
	Endpoint: "/stats/playerawards",
	Params:   []Pair{{"PlayerID", "2544"}},
}

var commonAllPlayers = GetReq{
	Host:     HOST,
	Headers:  HDRS,
	Endpoint: "/stats/commonallplayers",
	Params: []Pair{
		{"LeagueID", "00"},
		{"IsOnlyCurrentSeason", "1"},
		{"Season", "2024-25"},
	},
}

var leagueStandings = GetReq{
	Host:     HOST,
	Headers:  HDRS,
	Endpoint: "/stats/leaguestandings",
	Params: []Pair{
		{"LeagueID", "00"},
		{"Season", "2024-25"},
		{"SeasonType", "Regular+Season"},
	},
}
