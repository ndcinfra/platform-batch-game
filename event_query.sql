-- 이관 데이터 캐릭터 조회
SELECT a.*
FROM [CW_Account].[dbo].[tbl_Account] a, [CW_Game].[dbo].[tbl_unit] b 
WHERE a.UID = b.OwnerUID 
and a.id = '1608642026971677494'


--미래 사전 예약, max level
use [CW_Game]
SELECT	RankUserMaxLevel.UnitUID
FROM	
		(
      SELECT	UNT.UID AS UnitUID , ROW_NUMBER() OVER(PARTITION BY UNT.OwnerUID ORDER BY UNT.UID ASC) AS UnitCreateRank
      FROM	[CW_Game].[dbo].[tbl_Unit] AS UNT WITH(NOLOCK),
      (
        SELECT	OwnerUID, MAX(Level) AS Max_level
        FROM	[CW_Game].[dbo].[tbl_Unit]  WITH(NOLOCK)
        WHERE	OwnerUID IN (
                SELECT b.OwnerUID --, count(*)
                FROM  [CW_Game].[dbo].[tbl_Unit] b 
                where UnitSubType = 100032
                and  CreateTime >= '2020-12-10 00:00:00'
                and CreateTime < '2020-12-14 00:00:00'
                group by  b.OwnerUID
          )
        AND	Removed = 0
        GROUP BY	OwnerUID
        ) AS LVL
      WHERE	UNT.OwnerUID	= LVL.OwnerUID
      AND	UNT.Level		= LVL.Max_level
      AND	Removed = 0
) AS RankUserMaxLevel
WHERE	RankUserMaxLevel.UnitCreateRank = 1
order by RankUserMaxLevel.UnitUID;

--미래 굿즈
SELECT b.OwnerUID , a.UID, a.ID
FROM  [CW_Account].[dbo].[tbl_Account] a,[CW_Game].[dbo].[tbl_Unit] b 
WHERE a.UID = b.OwnerUID
and UnitSubType = 100032
and  CreateTime >= '2020-12-03 00:00:00'
and CreateTime < '2020-12-23 10:00:00'



-- 피로도 10이상 던전 10회 이상 클리어한 계정 UID와 플레이 수
USE [CW_GameLog_Archive]
GO


DECLARE	@StartDate	SMALLDATETIME = '2020-12-20 09:00';
DECLARE	@EndtDate	SMALLDATETIME = '2021-01-07 03:59';

--******************************************************
--******************************************************
-- tbl_actionlog_YYMMDD ���̺� ���� �����ؼ� �ڷ� ���.
--******************************************************
--******************************************************



--******************************************************
-- �� �Ѿ �� ���. over month 
--******************************************************
SELECT	UserUID, SUM(DungeonCount) AS TotalDungeonClearCount
FROM	(
			SELECT	UserUID, COUNT(1) AS DungeonCount FROM	tbl_actionlog_202012 WITH(NOLOCK) WHERE	ActionID IN (641) AND Value9 >= 10 AND @StartDate <= [Date] AND [Date] <= @EndtDate GROUP BY UserUID
			UNION ALL
			SELECT	UserUID, COUNT(1) AS DungeonCount FROM	tbl_actionlog_202101 WITH(NOLOCK) WHERE	ActionID IN (641) AND Value9 >= 10 AND @StartDate <= [Date] AND [Date] <= @EndtDate GROUP BY UserUID
		) AS DungeonLog
GROUP
  BY	UserUID
HAVING	SUM(DungeonCount) >= 10




--******************************************************
-- ���� ���̸� ���. only month
--******************************************************
SELECT	UserUID, COUNT(1) AS TotalDungeonClearCount FROM tbl_actionlog_202012 WITH(NOLOCK) WHERE ActionID IN (641) AND Value9 >= 10 AND @StartDate <= [Date] AND [Date] <= @EndtDate GROUP BY UserUID HAVING COUNT(1) >= 10;



--던젼 이벤트
USE [CW_GameLog_Archive]
GO
DECLARE @StartDate SMALLDATETIME = '2021-01-14 12:00:00';
DECLARE @EndtDate SMALLDATETIME = '2021-01-21 03:59:59';

SELECT UserUID, COUNT(1) AS TotalDungeonClearCount
FROM tbl_actionlog_202101 WITH(NOLOCK)
WHERE ActionID IN (641) AND Value9 >=10
AND @StartDate <= [Date] AND [Date] <= @EndtDate
GROUP BY UserUID HAVING COUNT(1) >= 10;


-- 던젼 이벤프 틀랫폼 u_id 추출
--1
SELECT g_account_id, gameuid
FROM public.z_temp_dungeon_event a, game_unit b
where a.gameuid = b.g_account_uid
group by g_account_id, gameuid

	--2
insert into public.z_temp_dungeon_event
select g_account_id::bigint, 0
from game_unit
where g_account_uid in(
	select gameuid
	FROM public.z_temp_dungeon_event_from_game
)
and g_account_id not in (
	SELECT u_id::character varying FROM public.payment_transaction 
where item_id = 1004
group by u_id )
group by g_account_id