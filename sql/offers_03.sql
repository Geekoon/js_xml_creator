SELECT o.id, c.id, c.parent_id, c.name, o.name, c.url, IFNULL(c.content, ''), IFNULL(o.barcode, ''), 
o.id_1c_offer, CAST(ob.value AS SIGNED), MAX(IF(pr.id_price=1, pr.value, NULL)), MAX(IF(pr.id_price=3, pr.value, NULL)), 
pid.code, pid.id_property_sex, pid.id_property_age, IFNULL(pid.structure, ''), pa.name, pik.name, 
MAX(CASE WHEN fv.id_feature=1 THEN fv.value END) AS size, 
MAX(CASE WHEN fv.id_feature=2 THEN fv.value END) AS color, 
IFNULL(MAX(CASE WHEN fv.id_feature=2 THEN fv.rgb END), '') AS rgb 
FROM tbl_offers AS o LEFT OUTER JOIN tbl_core AS c ON o.id_product_item = c.id 
LEFT OUTER JOIN tbl_offer_balance AS ob ON o.id = ob.id_offer 
LEFT OUTER JOIN tbl_offer_prices AS pr ON o.id = pr.id_offer 
LEFT OUTER JOIN tbl_product_item_detail AS pid ON c.id = pid.id_product_item 
LEFT OUTER JOIN tbl_product_articles AS pa ON pid.id_article = pa.id 
LEFT OUTER JOIN tbl_product_item_kind AS pik ON pik.id = pid.kind_id 
LEFT OUTER JOIN tbl_offer_features AS of ON o.id = of.id_offer 
LEFT OUTER JOIN tbl_feature_values AS fv ON of.id_feature_value = fv.id 
WHERE c.act=1 AND o.act=1 AND o.id_1c_offer != 0 AND ob.id_storage=2 AND ob.value != 0 AND pr.id_price != 2 GROUP BY o.id LIMIT 1000