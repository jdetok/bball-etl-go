package main

var nightlyTmGameLog = GetReq{
	Host:     HOST,
	Headers:  HDRS,
	Endpoint: "/stats/leaguegamelog",
	Params: []Pair{
		{"LeagueID", "10"},
		{"Season", "2025-26"},
		{"SeasonType", "Regular+Season"},
		{"Counter", "0"},
		{"Sorter", "DATE"},
		{"Direction", "DESC"},
		{"DateFrom", "07/24/2025"},
		{"DateTo", "07/24/2025"},
		{"PlayerOrTeam", "T"},
	},
}

var nightlyPlGameLog = GetReq{
	Host:     HOST,
	Headers:  HDRS,
	Endpoint: "/stats/leaguegamelog",
	Params: []Pair{
		{"LeagueID", "10"},
		{"Season", "2025-26"},
		{"SeasonType", "Regular+Season"},
		{"Counter", "0"},
		{"Sorter", "DATE"},
		{"Direction", "DESC"},
		{"DateFrom", "07/24/2025"},
		{"DateTo", "07/24/2025"},
		{"PlayerOrTeam", "P"},
	},
}

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
