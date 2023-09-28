/* sessionProperties: {"join_distribution_type":"BROADCAST"} */
select
  *,
  count(1) over() group_count
from
(
    select
      group_col,
      funnel_flow_array(max_step, 6) _col1,
      funnel_flow_array_date(date_max_step_map, 6) _col2,
      grouping(group_col) group_bit
    from
      (
        select
          virtual_user_id,
          group_col,
          funnel_max_step(packed_times, 86400000, 6) max_step,
          funnel_max_step_date(
            packed_times,
            86400000,
            6,
            cast('2023-09-21' as timestamp),
            date_add('second', 86400000 / 1000, timestamp '2023-09-27')
          ) date_max_step_map
        from
          (
            select
              virtual_user_id,
              min_by(
                group_col,
                if(
                  bitwise_and("@funnel_index_bit_set", 1) = 1,
                  "@vpc_tz_#event_time"
                )
              ) group_col,
              funnel_packed_time_collect("@vpc_tz_#event_time", "@funnel_index_bit_set") packed_times
            from
              (
                select
                  *,
                  cast(null as varchar) group_col,
                  8 "@funnel_index_bit_set"
                from
                  (
                    select
                      *,
                      "#user_id" virtual_user_id
                    from
                      (
                        select
                          *,
                          "#event_time" "@vpc_tz_#event_time"
                        from
                          (
                            select
                              "#event_name",
                              "#user_id",
                              "#event_time",
                              "$part_date",
                              "$part_event"
                            from
                              v_event_2
                          )
                      )
                  )
                where
                  ((("$part_event" IN ('region_select'))))
                  and (
                    "$part_date" between '2023-09-21'
                    and '2023-09-28'
                  )
                UNION
                ALL
                select
                  *,
                  cast(null as varchar) group_col,
                  4 "@funnel_index_bit_set"
                from
                  (
                    select
                      *,
                      "#user_id" virtual_user_id
                    from
                      (
                        select
                          *,
                          "#event_time" "@vpc_tz_#event_time"
                        from
                          (
                            select
                              "#event_name",
                              "#user_id",
                              "#event_time",
                              "$part_date",
                              "$part_event"
                            from
                              v_event_2
                          )
                      )
                  )
                where
                  ((("$part_event" IN ('version_update'))))
                  and (
                    "$part_date" between '2023-09-21'
                    and '2023-09-28'
                  )
                UNION
                ALL
                select
                  *,
                  cast(null as varchar) group_col,
                  2 "@funnel_index_bit_set"
                from
                  (
                    select
                      *,
                      "#user_id" virtual_user_id
                    from
                      (
                        select
                          *,
                          "#event_time" "@vpc_tz_#event_time"
                        from
                          (
                            select
                              "#event_name",
                              "#user_id",
                              "#event_time",
                              "$part_date",
                              "$part_event"
                            from
                              v_event_2
                          )
                      )
                  )
                where
                  ((("$part_event" IN ('version_check'))))
                  and (
                    "$part_date" between '2023-09-21'
                    and '2023-09-28'
                  )
                UNION
                ALL
                select
                  *,
                  cast(null as varchar) group_col,
                  bitwise_or(
                    bitwise_or(
                      if(
                        "$part_date" between '2023-09-21'
                        and '2023-09-27',
                        1,
                        0
                      ),
                      16
                    ),
                    32
                  ) "@funnel_index_bit_set"
                from
                  (
                    select
                      *,
                      "#user_id" virtual_user_id
                    from
                      (
                        select
                          *,
                          "#event_time" "@vpc_tz_#event_time"
                        from
                          (
                            select
                              "#event_name",
                              "#user_id",
                              "#event_time",
                              "$part_date",
                              "$part_event"
                            from
                              v_event_2
                          )
                      )
                  )
                where
                  ((("$part_event" IN ('device_login'))))
                  and (
                    "$part_date" between '2023-09-21'
                    and '2023-09-28'
                  )
              ) ta_ev
            where
              virtual_user_id is not null
            group by
              virtual_user_id
            having
              bitwise_and(bitwise_or_agg("@funnel_index_bit_set"), 1) = 1
          ) a
      ) a
    group by
      GROUPING SETS((group_col), ())
    order by
      group_bit desc,
      _col1 desc
  ) a
limit
  1001


---


/* sessionProperties: {"join_distribution_type":"BROADCAST"} */
select
  *,
  count(1) over() group_count
from
(
    select
      group_col,
      funnel_flow_array(max_step, 6) _col1,
      funnel_flow_array_date(date_max_step_map, 6) _col2,
      grouping(group_col) group_bit
    from
      (
        select
          virtual_user_id,
          group_col,
          funnel_max_step(packed_times, 86400000, 6) max_step,
          funnel_max_step_date(
            packed_times,
            86400000,
            6,
            cast('2023-09-21' as timestamp),
            date_add('second', 86400000 / 1000, timestamp '2023-09-27')
          ) date_max_step_map
        from
          (
            select
              virtual_user_id,
              min_by(
                group_col,
                if(
                  bitwise_and("@funnel_index_bit_set", 1) = 1,
                  "@vpc_tz_#event_time"
                )
              ) group_col,
              funnel_packed_time_collect("@vpc_tz_#event_time", "@funnel_index_bit_set") packed_times
            from
              (
                select
                  *,
                  "#app_version" group_col,
                  8 "@funnel_index_bit_set"
                from
                  (
                    select
                      *,
                      "#user_id" virtual_user_id
                    from
                      (
                        select
                          *,
                          "#event_time" "@vpc_tz_#event_time"
                        from
                          (
                            select
                              "#event_name",
                              "#app_version",
                              "#user_id",
                              "#event_time",
                              "$part_date",
                              "$part_event"
                            from
                              v_event_2
                          )
                      )
                  )
                where
                  ((("$part_event" IN ('region_select'))))
                  and (
                    "$part_date" between '2023-09-21'
                    and '2023-09-28'
                  )
                UNION
                ALL
                select
                  *,
                  "#app_version" group_col,
                  4 "@funnel_index_bit_set"
                from
                  (
                    select
                      *,
                      "#user_id" virtual_user_id
                    from
                      (
                        select
                          *,
                          "#event_time" "@vpc_tz_#event_time"
                        from
                          (
                            select
                              "#event_name",
                              "#app_version",
                              "#user_id",
                              "#event_time",
                              "$part_date",
                              "$part_event"
                            from
                              v_event_2
                          )
                      )
                  )
                where
                  ((("$part_event" IN ('version_update'))))
                  and (
                    "$part_date" between '2023-09-21'
                    and '2023-09-28'
                  )
                UNION
                ALL
                select
                  *,
                  "#app_version" group_col,
                  2 "@funnel_index_bit_set"
                from
                  (
                    select
                      *,
                      "#user_id" virtual_user_id
                    from
                      (
                        select
                          *,
                          "#event_time" "@vpc_tz_#event_time"
                        from
                          (
                            select
                              "#event_name",
                              "#app_version",
                              "#user_id",
                              "#event_time",
                              "$part_date",
                              "$part_event"
                            from
                              v_event_2
                          )
                      )
                  )
                where
                  ((("$part_event" IN ('version_check'))))
                  and (
                    "$part_date" between '2023-09-21'
                    and '2023-09-28'
                  )
                UNION
                ALL
                select
                  *,
                  "#app_version" group_col,
                  bitwise_or(
                    bitwise_or(
                      if(
                        "$part_date" between '2023-09-21'
                        and '2023-09-27',
                        1,
                        0
                      ),
                      16
                    ),
                    32
                  ) "@funnel_index_bit_set"
                from
                  (
                    select
                      *,
                      "#user_id" virtual_user_id
                    from
                      (
                        select
                          *,
                          "#event_time" "@vpc_tz_#event_time"
                        from
                          (
                            select
                              "#event_name",
                              "#app_version",
                              "#user_id",
                              "#event_time",
                              "$part_date",
                              "$part_event"
                            from
                              v_event_2
                          )
                      )
                  )
                where
                  ((("$part_event" IN ('device_login'))))
                  and (
                    "$part_date" between '2023-09-21'
                    and '2023-09-28'
                  )
              ) ta_ev
            where
              virtual_user_id is not null
            group by
              virtual_user_id
            having
              bitwise_and(bitwise_or_agg("@funnel_index_bit_set"), 1) = 1
          ) a
      ) a
    group by
      GROUPING SETS((group_col), ())
    order by
      group_bit desc,
      _col1 desc
  ) a
limit
  1001


---


/* sessionProperties: {"join_distribution_type":"BROADCAST"} */
select
  *,
  count(1) over() group_count
from
(
    select
      group_col,
      funnel_flow_array(max_step, 2) _col1,
      funnel_flow_array_date(date_max_step_map, 2) _col2,
      grouping(group_col) group_bit
    from
      (
        select
          virtual_user_id,
          group_col,
          funnel_max_step(packed_times, 86400000, 2) max_step,
          funnel_max_step_date(
            packed_times,
            86400000,
            2,
            cast('2023-09-21' as timestamp),
            date_add('second', 86400000 / 1000, timestamp '2023-09-27')
          ) date_max_step_map
        from
          (
            select
              virtual_user_id,
              min_by(
                group_col,
                if(
                  bitwise_and("@funnel_index_bit_set", 1) = 1,
                  "@vpc_tz_#event_time"
                )
              ) group_col,
              funnel_packed_time_collect("@vpc_tz_#event_time", "@funnel_index_bit_set") packed_times
            from
              (
                select
                  *,
                  cast(null as varchar) group_col,
                  bitwise_or(
                    if(
                      "$part_date" between '2023-09-21'
                      and '2023-09-27',
                      1,
                      0
                    ),
                    2
                  ) "@funnel_index_bit_set"
                from
                  (
                    select
                      *,
                      "#user_id" virtual_user_id
                    from
                      (
                        select
                          *,
                          "#event_time" "@vpc_tz_#event_time"
                        from
                          (
                            select
                              "#event_name",
                              "#user_id",
                              "#event_time",
                              "$part_date",
                              "$part_event"
                            from
                              v_event_2
                          )
                      )
                  )
                where
                  ((("$part_event" IN ('device_login'))))
                  and (
                    "$part_date" between '2023-09-21'
                    and '2023-09-28'
                  )
              ) ta_ev
            where
              virtual_user_id is not null
            group by
              virtual_user_id
            having
              bitwise_and(bitwise_or_agg("@funnel_index_bit_set"), 1) = 1
          ) a
      ) a
    group by
      GROUPING SETS((group_col), ())
    order by
      group_bit desc,
      _col1 desc
  ) a
limit
  1001
