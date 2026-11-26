SELECT squad_id, name
FROM squads
WHERE leader_id IS NULL;

SELECT d.dwarf_id, d.name
FROM dwarves d
WHERE d.age > 150
AND d.profession = 'Warrior';

SELECT DISTINCT d.dwarf_id, d.name
FROM dwarves d
INNER JOIN items i ON i.owner_id = d.dwarf_id
WHERE i.type = 'weapon';

SELECT d.dwarf_id, d.name, t.status, COUNT(t.task_id)
FROM dwarves d
LEFT JOIN tasks t ON d.dwarf_id = t.assigned_to
GROUP BY t.status, d.dwarf_id
ORDER BY d.name, t.status;

SELECT t.task_id, t.description
FROM tasks t
INNER JOIN dwarves d ON d.dwarf_id = t.assigned_to
INNER JOIN squads s ON d.squad_id = s.squad_id
WHERE s.name = 'Guardians';

SELECT d.dwarf_id, d.name, rs.relationship, r.dwarf_id, r.name
FROM dwarves d
LEFT JOIN relationships rs ON d.dwarf_id = rs.dwarf_id
LEFT JOIN dwarves r ON rs.related_to = r.dwarf_id;
