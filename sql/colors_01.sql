SELECT * FROM tbl_feature_values WHERE rgb IN (
SELECT parent_color FROM tbl_feature_values
WHERE id_feature=2 AND parent_color IS NOT NULL AND parent_color <> ''
GROUP BY parent_color)
AND id_feature=2 AND parent_color IS NOT NULL AND parent_color <> ''
AND rgb = parent_color
AND rgb <> '000000' AND rgb <> 'ffffff'