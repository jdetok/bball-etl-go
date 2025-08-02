/*
insert new teams into lg.team from intake.gm_teama
creates new concatenated column team_cde in format lal_lakers, bos_boston etc
*/ 

create or replace procedure lg.sp_team_load()
language plpgsql
as $$
begin
    insert into lg.team (
        lg_id, team_id, team, team_cde, team_long, team_city, team_shrt)
        select
            0,
            a.team_id,
            a.team_abbreviation,
            lower(a.team_abbreviation) || '_' || b.team_code,
            a.team_name,
            b.team_city,
            b.team_name
        from intake.gm_team a
        inner join intake.player b on b.team_id = a.team_id
        inner join (
            select 
                team_id as t_id,
                max(season_id) as s_id
            from intake.gm_team
            group by team_id
        ) c on c.t_id = a.team_id and c.s_id = a.season_id
        group by a.season_id, a.team_id, a.team_abbreviation, 
            b.team_code, a.team_name, b.team_city, b.team_name
        order by a.season_id desc
    on conflict (team_id) do nothing;

    insert into lg.team (
        lg_id, team_id, team, team_cde, team_long, team_city, team_shrt)
        select
            1,
            a.team_id,
            a.team_abbreviation,
            lower(a.team_abbreviation) || '_' || b.team_code,
            a.team_name,
            b.team_city,
            b.team_name
        from intake.gm_team a
        inner join intake.wplayer b on b.team_id = a.team_id
        inner join (
            select 
                team_id as t_id,
                max(season_id) as s_id
            from intake.gm_team
            group by team_id
        ) c on c.t_id = a.team_id and c.s_id = a.season_id
        group by a.season_id, a.team_id, a.team_abbreviation, 
            b.team_code, a.team_name, b.team_city, b.team_name
        order by a.season_id desc
    on conflict (team_id) do nothing;
end; $$;

-- call lg.sp_team_load();
-- select * from lg.team;
--delete from lg.team;
