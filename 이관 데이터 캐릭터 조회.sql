-- 이관 데이터 캐릭터 조회
SELECT a.*
FROM [CW_Account].[dbo].[tbl_Account] a, [CW_Game].[dbo].[tbl_unit] b 
WHERE a.UID = b.OwnerUID 
and a.id = '1611316142132565292'
