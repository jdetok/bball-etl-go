/*
insert current palyers into lg.plr_crnt. player has to exist in lg.plr before they can be 
inserted into plr_crnt. 
*/

create or replace procedure lg.sp_plr_crnt_load()
language plpgsql
as $$
begin
    insert into lg.plr_crnt (
        lg_id, team_id, player_id, plr_cde)
        select 
            0, 
            team_id,
            person_id,
            playercode
        from intake.player
    on conflict (player_id, team_id) do update set team_id = excluded.team_id;

    insert into lg.plr_crnt (
        lg_id, team_id, player_id, plr_cde)
        select 
            1, 
            team_id,
            person_id,
            playercode
        from intake.wplayer
    on conflict (player_id, team_id) do update set team_id = excluded.team_id;
end; $$;

-- call lg.sp_plr_crnt_load();
-- select * from lg.plr_crnt;
-- select * from lg.plr where player_id = 1628304;
-- select * from intake.wplayer where person_id = 1628304; 
