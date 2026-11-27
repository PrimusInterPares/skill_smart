SELECT
  d.dwarf_id,
  d.name,
  d.age,
  d.profession,
  json_build_object(
    'skill_ids', COALESCE((
      SELECT json_agg(ds.skill_id)
      FROM dwarf_skills ds
      WHERE d.dwarf_id = ds.dwarf_id
    ), '[]'::json),
    'resource_ids', COALESCE((
      SELECT json_agg(da.assignment_id)
      FROM dwarf_assignments da
      WHERE d.dwarf_id = da.dwarf_id
    ), '[]'::json),
    'workshop_ids', COALESCE((
      SELECT json_agg(sm.squad_id)
      FROM squad_members sm
      WHERE d.dwarf_id = sm.dwarf_id
    ), '[]'::json),
    'squad_ids', COALESCE((
      SELECT json_agg(de.equipment_id)
      FROM dwarf_equipment de
      WHERE d.dwarf_id = de.dwarf_id
    ), '[]'::json)
  ) AS related_entities
FROM dwarves d
;
