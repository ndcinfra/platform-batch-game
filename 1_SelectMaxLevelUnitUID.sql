USE [CW_Game]
GO



SELECT	*
FROM	
		(
			SELECT	UNT.OwnerUID AS UserUID, UNT.UID AS UnitUID, UNT.CreateTime, UNT.[Level], ROW_NUMBER() OVER(PARTITION BY UNT.OwnerUID ORDER BY UNT.UID ASC) AS UnitCreateRank
			FROM	tbl_Unit AS UNT WITH(NOLOCK),
					(
						SELECT	OwnerUID, MAX(Level) AS Max_level
						FROM	tbl_Unit WITH(NOLOCK)
						WHERE	OwnerUID IN (
												SELECT UserUID FROM tbl_Work_MaxLevelUnitSelect WITH(NOLOCK)
											)
						 AND	Removed = 0
						GROUP
						  BY	OwnerUID
					) AS LVL
			WHERE	UNT.OwnerUID	= LVL.OwnerUID
			 AND	UNT.Level		= LVL.Max_level
			 AND	Removed = 0
		) AS RankUserMaxLevel
WHERE	RankUserMaxLevel.UnitCreateRank = 1
;





/*
select * from (
select OwnerUID, UID, Level, CreateTime, RANK() over(PARTITION BY OwnerUID ORDER BY Level DESC, CreateTime ASC) as ranking from tbl_Unit with(nolock) where UnitSubType = 100032 and Removed = 0 
group by OwnerUID, UID, Level, CreateTime) as t1
where t1.ranking = 1



DECLARE @UserUID bigint = 0
select OwnerUID, UID, Level, CreateTime, RANK() over(PARTITION BY OwnerUID ORDER BY Level DESC, CreateTime ASC) as ranking from tbl_Unit with(nolock)
where OwnerUID = @UserUID and UnitSubType = 100032 and Removed = 0
*/





