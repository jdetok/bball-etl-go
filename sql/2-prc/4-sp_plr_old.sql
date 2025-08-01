/* 
insert players from intake.gm_player into lg.plr
left join on intake.player to insert players if they are not in intake.player
need to also consider joining in wplayer as well to capture the league id 
for players in w table as well
*/
create or replace procedure lg.sp_plr_old_load()
language plpgsql
as $$
begin
    insert into lg.plr (
        lg_id, player_id, plr_cde, player, last_first, from_year, to_year)
        select 
            b.lg_id, 
            a.player_id,
            c.playercode,
            a.player_name,
            c.display_last_comma_first,
            c.from_year,
            c.to_year
        from intake.gm_player a
        inner join lg.team b on b.team_id = a.team_id
        left join intake.player c on a.player_id = c.person_id
        group by b.lg_id, a.player_id, c.playercode, a.player_name,
                    c.display_last_comma_first, c.from_year, c.to_year
    on conflict (player_id) do nothing;
end; $$;

-- call lg.sp_plr_old_load();
