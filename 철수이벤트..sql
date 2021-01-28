USE [CW_Game]
GO



--100031	)	// 철수
--100032	)	// 미래
--100033	)	// 은하
--100034	)	// 루시
DECLARE	@UnitSubTye		INTEGER			= 100031 
DECLARE	@DateStart		SMALLDATETIME	= '2021-01-14 14:00:00'
DECLARE	@DateFinish		SMALLDATETIME	= '2021-01-28 09:59:59'


SELECT	RankUserMaxLevel.UnitUID --, RankUserMaxLevel.OwnerUID
FROM	
		(
			SELECT	UNT.OwnerUID, UNT.UID AS UnitUID , ROW_NUMBER() OVER(PARTITION BY UNT.OwnerUID ORDER BY UNT.UID ASC) AS UnitCreateRank
			FROM	[CW_Game].[dbo].[tbl_Unit] AS UNT WITH(NOLOCK),
					(
						SELECT	OwnerUID, MAX(Level) AS Max_level
						FROM	[CW_Game].[dbo].[tbl_Unit]  WITH(NOLOCK)
						WHERE	OwnerUID IN (
												SELECT	UNT2.OwnerUID
												FROM	[CW_Game].[dbo].[tbl_Unit] AS UNT2
												WHERE	UnitSubType	=	@UnitSubTye
												 AND	CreateTime	>=	@DateStart
												 AND	CreateTime	<	@DateFinish
												GROUP
												  BY	UNT2.OwnerUID
											)
						AND		Removed = 0
						GROUP 
						  BY	OwnerUID
					) AS LVL
			WHERE	UNT.OwnerUID	= LVL.OwnerUID
			 AND	UNT.Level		= LVL.Max_level
			 AND	Removed = 0
		) AS RankUserMaxLevel
WHERE	RankUserMaxLevel.UnitCreateRank = 1
ORDER
  BY	RankUserMaxLevel.UnitUID;
