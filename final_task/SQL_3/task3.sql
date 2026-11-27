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
    'assignment_ids', COALESCE((
      SELECT json_agg(da.assignment_id)
      FROM dwarf_assignments da
      WHERE d.dwarf_id = da.dwarf_id
      AND (da.end_date IS NULL OR da.end_date >= CURRENT_DATE)
    ), '[]'::json),
    'squad_ids', COALESCE((
      SELECT json_agg(DISTINCT sm.squad_id)
      FROM squad_members sm
      WHERE d.dwarf_id = sm.dwarf_id
    ), '[]'::json),
    'equipment_ids', COALESCE((
      SELECT json_agg(DISTINCT de.equipment_id)
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
      SELECT json_agg(DISTINCT wc.dwarf_id)
      FROM workshop_craftsdwarves wc
      WHERE w.workshop_id = wc.workshop_id
    ), '[]'::json),
    'project_ids', COALESCE((
      SELECT json_agg(p.project_id)
      FROM projects p
      WHERE w.workshop_id = p.workshop_id
      AND p.status = 'active'
    ), '[]'::json),
    'input_material_ids', COALESCE((
      SELECT json_agg(DISTINCT wm.material_id)
      FROM workshop_materials wm
      WHERE w.workshop_id = wm.workshop_id
      AND wm.is_input = TRUE
    ), '[]'::json),
    'output_product_ids', COALESCE((
      SELECT json_agg(DISTINCT wp.product_id)
      FROM workshop_products wp
      WHERE w.workshop_id = wp.workshop_id
    ), '[]'::json)
  ) AS related_entities
FROM workshops w
;

SELECT
  s.squad_id,
  s.name,
  s.formation_type,
  s.leader_id,
  json_build_object(
    'member_ids', COALESCE((
      SELECT json_agg(DISTINCT sm.dwarf_id)
      FROM squad_members sm
      WHERE s.squad_id = sm.squad_id
    ), '[]'::json),
    'equipment_ids', COALESCE((
      SELECT json_agg(DISTINCT se.equipment_id)
      FROM squad_equipment se
      WHERE s.squad_id = se.squad_id
    ), '[]'::json),
    'operation_ids', COALESCE((
      SELECT json_agg(so.operation_id)
      FROM squad_operations so
      WHERE s.squad_id = so.squad_id
      AND so.start_date <= CURRENT_DATE
    ), '[]'::json),
    'training_schedule_ids', COALESCE((
      SELECT json_agg(st.schedule_id)
      FROM squad_training st
      WHERE s.squad_id = st.squad_id
    ), '[]'::json),
    'battle_report_ids', COALESCE((
      SELECT json_agg(sb.report_id)
      FROM squad_battles sb
      WHERE s.squad_id = sb.squad_id
    ), '[]'::json)
  ) AS related_entities
FROM military_squads s
;
