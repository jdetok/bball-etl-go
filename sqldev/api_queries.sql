-- replacement for mariadb Teams query
select
	b.lg, a.team_id, a.team, a.team_long
from lg.team a
inner join lg.league b on a.lg_id = b.lg_id
where a.lg_id < 2; -- captures only NBA & WNBA

-- replacement for mariadb RSeasons query
select szn_id, szn_desc, wszn_desc
from lg.szn
where sznt_id in (2, 4);

-- api stats query testing (replace mariadb api_player_stats)
-- the addition of szn_long and wszn_long matches the Player query exactly

/* EDIT: TOOK TEAM COLS OUT OF GROUP BY, QUERY MAX
THIS HANDLES PLAYERS WHO WERE TRADED MID YEAR
*/ 
select a.player_id, max(a.team_id), d.lg, a.szn_id, 'tot' as "stype", b.player, 
	max(c.team), max(c.team_long), count(distinct a.game_id) as "gp", e.szn_desc, 
	sum(a.pts) as "points", sum(a.ast) as "assists", 
	sum(a.reb) as "rebounds", sum(a.stl) as "steals", sum(a.blk) as "blocks", 
	sum(a.fgm) as "fgm", sum(a.fga) as "fga",
	coalesce(
		cast(round(avg(a.fgp) * 100, 2) as varchar(10)) || '%', '0%')
	as "fgp",
	sum(a.f3m) as "f3m", sum(a.f3a) as "f3a",
	coalesce(
		cast(round(avg(a.f3p) * 100, 2) as varchar(10)) || '%', '0%')
	as "f3p",
	sum(a.ftm) as "ftm", sum(a.fta) as "fta",
	coalesce(
		cast(round(avg(a.ftp) * 100, 2) as varchar(10)) || '%', '0%')
	as "ftp"
from stats.pbox a
inner join lg.plr b on b.player_id = a.player_id
inner join lg.team c on c.team_id = a.team_id
inner join lg.league d on d.lg_id = b.lg_id
inner join lg.szn e on e.szn_id = a.szn_id
where b.lg_id < 2 --and e.sznt_id = 2
group by a.player_id, d.lg, a.szn_id, b.player, e.szn_desc, e.wszn_desc
order by a.szn_id desc;


-- AVG
select a.player_id, max(a.team_id), d.lg, a.szn_id, 'avg' as "stype", b.player, 
	max(c.team), max(c.team_long), count(distinct a.game_id) as "gp", e.szn_desc, 
	e.wszn_desc, sum(a.pts) as "points", sum(a.ast) as "assists", 
	sum(a.reb) as "rebounds", sum(a.stl) as "steals", sum(a.blk) as "blocks", 
	round(avg(a.fgm), 2) as "fgm", round(avg(a.fga), 2) as "fga",
	coalesce(
		cast(round(avg(a.fgp) * 100, 2) as varchar(10)) || '%', '0%')
	as "fgp",
	round(avg(a.f3m), 2) as "f3m", round(avg(a.f3a), 2) as "f3a",
	coalesce(
		cast(round(avg(a.f3p) * 100, 2) as varchar(10)) || '%', '0%')
	as "f3p",
	round(avg(a.ftm), 2) as "ftm", round(avg(a.fta), 2) as "fta",
	coalesce(
		cast(round(avg(a.ftp) * 100, 2) as varchar(10)) || '%', '0%')
	as "ftp"
from stats.pbox a
inner join lg.plr b on b.player_id = a.player_id
inner join lg.team c on c.team_id = a.team_id
inner join lg.league d on d.lg_id = b.lg_id
inner join lg.szn e on e.szn_id = a.szn_id
where b.lg_id < 2 
group by a.player_id, d.lg, a.szn_id, b.player, e.szn_desc, e.wszn_desc
order by a.szn_id desc;