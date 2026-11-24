SELECT d.*, s.*
FROM dwarves d
LEFT JOIN squads s
ON d.squad_id = s.squad_id
WHERE d.squad_id IS NOT NULL;

SELECT d.*
FROM dwarves d
WHERE d.squad_id IS NULL
AND d.profession = 'miner';

SELECT t.*
FROM tasks t
WHERE t.priority = (
  SELECT MIN(priority)
  FROM tasks
)
AND status = 'pending';

SELECT d.dwarf_id, COUNT(i.item_id) AS items_count
FROM dwarves d
INNER JOIN items i
ON i.owner_id = d.dwarf_id
GROUP BY d.dwarf_id;

SELECT s.squad_id, COUNT(d.dwarf_id) AS dwarves_count
FROM squads s
LEFT JOIN dwarves d
ON d.squad_id = s.squad_id
GROUP BY s.squad_id;

SELECT d.profession
FROM dwarves d
INNER JOIN tasks t
ON t.assigned_to = d.dwarf_id
WHERE t.status IN ('pending', 'in_progress')
GROUP BY d.profession
HAVING COUNT(t.task_id) = (
  SELECT MAX(tasks_count)
  FROM (
    SELECT COUNT(t1.task_id) AS tasks_count
    FROM tasks t1
    INNER JOIN dwarves d1 
    ON t1.assigned_to = d1.dwarf_id
    WHERE t1.status IN ('pending', 'in_progress')
    GROUP BY d1.profession
  ) AS subq
);

SELECT i.type, ROUND(AVG(d.age)::numeric, 0) AS avg_dwarf_age
FROM items i
INNER JOIN dwarves d
ON d.dwarf_id = i.owner_id
GROUP BY i.type;

SELECT d.*
FROM dwarves d
WHERE d.age > (SELECT AVG(age) FROM dwarves)
AND d.dwarf_id NOT IN (
  SELECT owner_id FROM items WHERE owner_id IS NOT NULL
);
