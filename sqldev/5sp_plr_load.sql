/* 
load players from intake.player & wplayer tables into lg.plr
in db creation, this should be called before sp_plr_old_load
*/

create or replace procedure lg.sp_plr_load()
language plpgsql
as $$
begin
    insert into lg.plr (
        lg_id, player_id, plr_cde, player, last_first, from_year, to_year)
        select 
            0, 
            person_id,
            playercode,
            display_first_last,
            display_last_comma_first,
            from_year,
            to_year
        from intake.player
    on conflict (player_id) do nothing;

    insert into lg.plr (
        lg_id, player_id, plr_cde, player, last_first, from_year, to_year)
        select 
            1, 
            person_id,
            playercode,
            display_first_last,
            display_last_comma_first,
            from_year,
            to_year
        from intake.wplayer
    on conflict (player_id) do nothing;
end; $$;

-- call lg.sp_plr_load();
-- select * from lg.plr;