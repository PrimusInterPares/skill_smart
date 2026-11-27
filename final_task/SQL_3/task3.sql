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

SELECT
  w.workshop_id,
  w.name,
  w.type,
  w.quality,
  json_build_object(
    'craftsdwarf_ids', COALESCE((
      SELECT json_agg(wc.dwarf_id)
      FROM workshop_craftsdwarves wc
      WHERE w.workshop_id = wc.workshop_id
    ), '[]'::json),
    'project_ids', COALESCE((
      SELECT json_agg(p.project_id)
      FROM projects p
      WHERE w.workshop_id = p.workshop_id
    ), '[]'::json),
    'input_material_ids', COALESCE((
      SELECT json_agg(wm.material_id)
      FROM workshop_materials wm
      WHERE w.workshop_id = wm.workshop_id
      AND wm.is_input = TRUE
    ), '[]'::json),
    'output_product_ids', COALESCE((
      SELECT json_agg(wp.product_id)
      FROM workshop_products wp
      INNER JOIN products p ON wp.product_id = p.product_id
      INNER JOIN workshop_materials wm ON p.material_id = wm.material_id
      WHERE w.workshop_id = wp.workshop_id
      AND wm.is_input = FALSE
    ), '[]'::json)
  )
FROM workshops w
;
