package main

import "time"

var YESTERDAY string = time.Now().Add(-24 * time.Hour).Format("01/02/2006")

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
		{"PlayerOrTeam", "T"},
		{"DateFrom", YESTERDAY},
		{"DateTo", YESTERDAY},
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
		{"PlayerOrTeam", "P"},
		{"DateFrom", YESTERDAY},
		{"DateTo", YESTERDAY},
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
