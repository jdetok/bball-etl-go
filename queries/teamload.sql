-- QUERY RETURNS ALL CURRENT NBA AND WNBA TEAMS
select
    0,
    a.team_id,
    a.team_abbreviation,
    lower(a.team_abbreviation) || '_' || b.team_code,
    a.team_name,
    b.team_city,
    b.team_name
from intake.gm_team a
inner join intake.player b 
    on b.team_id = a.team_id
inner join (
    select 
        team_id as t_id,
        max(season_id) as s_id
    from intake.gm_team
    where left(cast(season_id as varchar(5)), 1) = '2'
    group by team_id
) c on c.t_id = a.team_id and c.s_id = a.season_id
-- player from year greater than current season year
where b.from_year >= right(cast(a.season_id as varchar(5)), 4)
and b.team_id > 0 -- no team_id = 0
group by a.season_id, a.team_id, a.team_abbreviation, 
    b.team_code, a.team_name, b.team_city, b.team_name
UNION
select
    1,
    a.team_id,
    a.team_abbreviation,
    lower(a.team_abbreviation) || '_' || b.team_code,
    a.team_name,
    b.team_city,
    b.team_name
from intake.gm_team a
inner join intake.wplayer b 
    on b.team_id = a.team_id
inner join (
    select 
        team_id as t_id,
        max(season_id) as s_id
    from intake.gm_team
    where left(cast(season_id as varchar(5)), 1) = '2'
    group by team_id
) c on c.t_id = a.team_id and c.s_id = a.season_id
-- player from year greater than current season year
where b.from_year >= right(cast(a.season_id as varchar(5)), 4)
and b.team_id > 0 -- no team_id = 0
group by a.season_id, a.team_id, a.team_abbreviation, 
    b.team_code, a.team_name, b.team_city, b.team_name
;
-- ALL seasons, including those no longer active
-- joined the max season subq first, which fixed most problems
-- player from_year before max szn, to_year after
-- unioned nba and wnba query together
-- NBA QUERY
select
    0,
    a.team_id,
    a.team_abbreviation,
    lower(a.team_abbreviation) || '_' || c.team_code,
    a.team_name,
    c.team_city,
    c.team_name
from intake.gm_team a
inner join (
    select 
        team_id as t_id,
        max(season_id) as s_id
    from intake.gm_team
    group by team_id
) b on b.t_id = a.team_id and b.s_id = a.season_id
inner join intake.player c
    on c.team_id = a.team_id
-- player from year greater than current season year
where c.from_year <= right(cast(a.season_id as varchar(5)), 4)
and c.to_year >= right(cast(a.season_id as varchar(5)), 4)
and c.team_id > 0 -- no team_id = 0
group by a.season_id, a.team_id, a.team_abbreviation, 
    c.team_code, a.team_name, c.team_city, c.team_name 
-- ===================================================================
union -- +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
-- WNBA QUERY
select
    1,
    a.team_id,
    a.team_abbreviation,
    lower(a.team_abbreviation) || '_' || c.team_code,
    a.team_name,
    c.team_city,
    c.team_name
from intake.gm_team a
inner join (
    select 
        team_id as t_id,
        max(season_id) as s_id
    from intake.gm_team
    group by team_id
) b on b.t_id = a.team_id and b.s_id = a.season_id
inner join intake.wplayer c
    on c.team_id = a.team_id
-- player from year greater than current season year
where c.from_year <= right(cast(a.season_id as varchar(5)), 4)
and c.to_year >= right(cast(a.season_id as varchar(5)), 4)
and c.team_id > 0 -- no team_id = 0
group by a.season_id, a.team_id, a.team_abbreviation, 
    c.team_code, a.team_name, c.team_city, c.team_name
order by team_id; 