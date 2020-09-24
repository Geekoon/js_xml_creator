SELECT o.id, c.id, c.parent_id, c.name, o.name, c.url, IFNULL(c.content, ''), IFNULL(o.barcode, ''), 
o.id_1c_offer, CAST(ob.value AS SIGNED), 
pid.code, pid.id_property_sex, pid.id_property_age, IFNULL(pid.structure, ''), pa.name, pik.name 

FROM tbl_offers AS o 
LEFT OUTER JOIN tbl_core AS c ON o.id_product_item = c.id 
LEFT OUTER JOIN tbl_offer_balance AS ob ON o.id = ob.id_offer 
LEFT OUTER JOIN tbl_product_item_detail AS pid ON c.id = pid.id_product_item 
LEFT OUTER JOIN tbl_product_articles AS pa ON pid.id_article = pa.id 
LEFT OUTER JOIN tbl_product_item_kind AS pik ON pik.id = pid.kind_id 
WHERE c.act=1 AND o.act=1 AND o.id_1c_offer != '00000000-0000-0000-0000-000000000000' AND ob.id_storage=2 AND ob.value != 0 #GROUP BY o.id #LIMIT 1000