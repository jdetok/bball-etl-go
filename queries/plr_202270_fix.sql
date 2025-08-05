/* 
INSERT A ROW INTO lg.plr FOR WNBA PLAYER ANGEL ROBINSON WITH PLAYER ID 202270
this player had the ID 202270 in 2014 and 202657 in all years after
this was causing an error with loading stats.pbox table
create a new record with identical data except player id
MUST BE RUN AFTER THE lg.sp_player_all_load
*/

insert into lg.plr (lg_id, player_id, plr_cde, player, last_first, from_year, to_year)
select
	1, 202270, 
	playercode, display_first_last, display_last_comma_first,
	from_year, to_year
from intake.wplayer
where person_id = 202657;