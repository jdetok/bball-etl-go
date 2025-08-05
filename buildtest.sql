-- ** meant to be called on EXISTING database
-- TEST DB BUILD WITHOUT STARTING UP A NEW DB
-- delete from all tables in all schemas except intake 
/*
this script assumes everything has already been built. the purpose is just to
simulate starting from right after the GO etl process finishes. 
the intake.[w]player tables, intake.gm_team, and intake.gm_player tables are
fully populated. no stored procedures have yet been run, so the tables in the 
lg, stats, and api schemas are all empty
*/

create or replace procedure sp_rebuild()
language plpgsql
as $$
begin
	raise notice 'deleting from api.plr_agg...';
	truncate api.plr_agg cascade;
	
	raise notice 'deleting from stats.pbox...';
	truncate stats.pbox cascade;
	
	raise notice 'deleting from stats.tbox...';
	truncate stats.tbox cascade;
	
	raise notice 'deleting from lg.plr...';
	truncate lg.plr cascade;
	
	raise notice 'deleting from lg.team...';
	delete from lg.team where team_id > 0;
	
	raise notice 'deleting from lg.szn...';
	delete from lg.szn where right(cast(szn_id as varchar(5)), 4) != '9999';

	-- load seasons
	raise notice 'inserting seasons...';
	call lg.sp_szn_load();
	raise notice 'seasons insert complete: %s', fn_cntstr('lg.szn');

	-- load all teams
	raise notice 'inserting all nba/wnba teams...';
	call lg.sp_team_all_load();
	raise notice 'team insert complete: %s', fn_cntstr('lg.team');

	-- load tbox table with team box stats
	raise notice 'inserting team box stats into stats.tbox...';
	call stats.sp_tbox();
	raise notice 'tbox insert complete: %s', fn_cntstr('stats.tbox');

	-- load all players
	raise notice 'inserting all nba/wnba players...';
	call lg.sp_plr_all_load();
	raise notice 'player insert complete: %s', fn_cntstr('lg.plr');

	/* INSERT A ROW INTO lg.plr FOR WNBA PLAYER ANGEL ROBINSON WITH PLAYER ID 
	202270 this player had the ID 202270 in 2014 and 202657 in all years after
	this was causing an error with loading stats.pbox table
	create a new record with identical data except player id
	MUST BE RUN AFTER THE lg.sp_player_all_load*/
	raise notice 'inserting 202270 copy of 202657, won''t work without...';
	insert into lg.plr 
		(lg_id, player_id, plr_cde, player, last_first, from_year, to_year)
	select
		1, 202270, 
		playercode, display_first_last, display_last_comma_first,
		from_year, to_year
	from intake.wplayer
	where person_id = 202657;
	raise notice 'player id 202270 insert complete: %s', fn_cntstr('lg.plr');

	-- load pbox table with player box scores after inserting player causing issue
	raise notice 'inserting player box stats into stats.pbox...';
	call stats.sp_pbox();
	raise notice 'pbox insert complete: %s', fn_cntstr('stats.pbox');

	-- load api.plr_agg table with pbox stats 
	raise notice 'inserting season/career stat aggregations into api.plr_agg...';
	call api.sp_plr_agg();
	raise notice 'player agg insert complete: %s', fn_cntstr('api.plr_agg');
end; $$;