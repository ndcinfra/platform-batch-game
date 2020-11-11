package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/astaxie/beego/logs"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
	"github.com/ndcinfra/platform-batch-game/libs"
	"github.com/ndcinfra/platform-batch-game/models"
)

var DailyBatchSql = "SELECT " +
	"a.UID as g_account_uid, " +
	"a.PublisherSN as g_account_publisher_sn, " +
	"a.ID as g_account_id, " +
	"ISNULL(CONVERT(varchar(max), a.Info, 0),'') as g_account_info, " +
	"ISNULL(a.NickName,'') as g_account_nick_name, " +
	"a.AccountLevel as g_account_account_level, " +
	"b.UID as g_unit_uid, " +
	"b.OwnerUID as g_unit_owner_uid, " +
	"b.UnitID as g_unit_unit_id, " +
	"b.UnitSubType as g_unit_unit_sub_type, " +
	"b.[Level] as g_uint_level, " +
	"b.[EXP] as g_unit_exp, " +
	"b.Gold as g_unit_gold, " +
	"b.PVESkillPoint as g_unit_pve_skill_point, " +
	"b.PnaSkillPoint as g_unit_pna_skill_point, " +
	"b.ItemSkillpoint as g_unit_item_skill_point, " +
	"b.ClassUpgradeLevel as g_unit_class_upgrade_level, " +
	"b.TutorialTask as g_unit_tutorial_task, " +
	"b.TitleID as g_unit_title_id, " +
	"b.TownID as g_unit_town_id, " +
	"b.PosX as g_unit_pos_x, " +
	"b.PosY as g_unit_pos_y, " +
	"b.PosZ as g_unit_pos_z, " +
	"b.TutorialCompleteLevel as g_unit_tutorial_complete_level, " +
	"b.IsFirstSlot as g_unit_is_first_slot, " +
	"b.QuickSlotUnlock as g_unit_quick_slot_unlock, " +
	"b.RebirthCoin as g_unit_rebirth_coin, " +
	"b.EntrancePoint as g_unit_entrance_point, " +
	"b.CreateTime as g_unit_create_time, " +
	"b.LastUseTime as g_unit_last_use_time, " +
	"b.LastLeaveTime as g_unit_last_leave_time, " +
	"b.Removed as g_unit_removed, " +
	"ISNULL(b.RemovedTime, '1900-01-01 00:00:00') as g_unit_removed_time, " +
	"b.UnionPoint as g_unit_union_point, " +
	"b.VisualTitleID as g_unit_visual_title_id, " +
	"b.AvatarMode as g_unit_avatar_mode, " +
	"b.VisualFrameID as g_unit_visual_frame_id, " +
	"b.ClosetSetSlotExpandCount as g_unit_closet_set_slot_expand_count, " +
	"b.ArtifactEnergy as g_unit_artifact_energy, " +
	"b.ArtifactOnOff as g_unit_artifact_on_off, " +
	"b.DamageFontID as g_unit_damage_font_id, " +
	"b.DungeonClearCount as g_unit_dungeon_clear_count " +
	"FROM [CW_Account].[dbo].[tbl_Account] a, [CW_Game].[dbo].[tbl_unit] b " +
	"WHERE a.UID = b.OwnerUID " +
	"ORDER BY b.LastUseTime DESC LIMIT 1000000 "

var insertSql = " INSERT INTO public.game_unit(" +
	"g_account_uid, g_account_publisher_sn, g_account_id, g_account_info, g_account_nick_name, g_account_account_level, g_unit_uid, g_unit_owner_uid, g_unit_unit_id, g_unit_unit_sub_type, g_uint_level, g_unit_exp, g_unit_gold, g_unit_pve_skill_point, g_unit_pna_skill_point, g_unit_item_skill_point, g_unit_class_upgrade_level, g_unit_tutorial_task, g_unit_title_id, g_unit_town_id, g_unit_pos_x, g_unit_pos_y, g_unit_pos_z, g_unit_tutorial_complete_level, g_unit_is_first_slot, g_unit_quick_slot_unlock, g_unit_rebirth_coin, g_unit_entrance_point, g_unit_create_time, g_unit_last_use_time, g_unit_last_leave_time, g_unit_removed, g_unit_removed_time, g_uit_union_point, g_unit_visual_title_id, g_unit_avatar_mode, g_unit_visual_frame_id, g_unit_closet_set_slot_expand_count, g_unit_artifact_energy, g_unit_artifact_on_off, g_unit_damage_font_id, g_unit_dungeon_clear_count) " +
	" VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42); "

var db *sql.DB

// Daily Cron Job for getting game data
func GetGameDataDaily(conn *pgx.Conn) {
	start := time.Now()
	logs.Info("start GetGameDataDaily: ", start)

	connString := os.Getenv("DBHOST_GAME_MSSQL")

	var err error

	// Create connection pool
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		logs.Error("Error creating connection pool: ", err.Error())
		return
	}

	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		logs.Error(err.Error())
		return
	}
	//fmt.Printf("Connected!\n")
	logs.Info("Connected Game DB")

	tsql := fmt.Sprintf(DailyBatchSql)

	// Execute query
	rows, err := db.QueryContext(ctx, tsql)
	if err != nil {
		logs.Error("Error execute query: ", err)
		return
	}

	defer rows.Close()

	var count int
	batch := &pgx.Batch{}
	//var gameUnits []models.GameUnit

	// Iterate through the result set.
	for rows.Next() {
		var gameUnit models.GameUnit

		// Get values from row.
		err := rows.Scan(
			&gameUnit.GAccountUid,
			&gameUnit.GAccountPublisherSn,
			&gameUnit.GAccountId,
			&gameUnit.GAccountInfo,
			&gameUnit.GAccountNickName,
			&gameUnit.GAccountAccountLevel,
			&gameUnit.GUnitUid,
			&gameUnit.GUnitOwnerUid,
			&gameUnit.GUnitUnitId,
			&gameUnit.GUnitUnitSubType,
			&gameUnit.GUintLevel,
			&gameUnit.GUnitExp,
			&gameUnit.GUnitGold,
			&gameUnit.GUnitPveSkillPoint,
			&gameUnit.GUnitPnaSkillPoint,
			&gameUnit.GUnitItemSkillPoint,
			&gameUnit.GUnitClassUpgradeLevel,
			&gameUnit.GUnitTutorialTask,
			&gameUnit.GUnitTitleId,
			&gameUnit.GUnitTownId,
			&gameUnit.GUnitPosX,
			&gameUnit.GUnitPosY,
			&gameUnit.GUnitPosZ,
			&gameUnit.GUnitTutorialCompleteLevel,
			&gameUnit.GUnitIsFirstSlot,
			&gameUnit.GUnitQuickSlotUnlock,
			&gameUnit.GUnitRebirthCoin,
			&gameUnit.GUnitEntrancePoint,
			&gameUnit.GUnitCreateTime,
			&gameUnit.GUnitLastUseTime,
			&gameUnit.GUnitLastLeaveTime,
			&gameUnit.GUnitRemoved,
			&gameUnit.GUnitRemovedTime,
			&gameUnit.GUitUnionPoint,
			&gameUnit.GUnitVisualTitleId,
			&gameUnit.GUnitAvatarMode,
			&gameUnit.GUnitVisualFrameId,
			&gameUnit.GUnitClosetSetSlotExpandCount,
			&gameUnit.GUnitArtifactEnergy,
			&gameUnit.GUnitArtifactOnOff,
			&gameUnit.GUnitDamageFontId,
			&gameUnit.GUnitDungeonClearCount,
		)
		if err != nil {
			logs.Error("Error rows.Scan: ", err)
			return
		}

		//fmt.Printf("g_account_uid: %d \n", gameUnit.GAccountUid)
		//gameUnits = append(gameUnits, gameUnit)

		batch.Queue(
			insertSql,
			gameUnit.GAccountUid,
			gameUnit.GAccountPublisherSn,
			gameUnit.GAccountId,
			gameUnit.GAccountInfo,
			gameUnit.GAccountNickName,
			gameUnit.GAccountAccountLevel,
			gameUnit.GUnitUid,
			gameUnit.GUnitOwnerUid,
			gameUnit.GUnitUnitId,
			gameUnit.GUnitUnitSubType,
			gameUnit.GUintLevel,
			gameUnit.GUnitExp,
			gameUnit.GUnitGold,
			gameUnit.GUnitPveSkillPoint,
			gameUnit.GUnitPnaSkillPoint,
			gameUnit.GUnitItemSkillPoint,
			gameUnit.GUnitClassUpgradeLevel,
			gameUnit.GUnitTutorialTask,
			gameUnit.GUnitTitleId,
			gameUnit.GUnitTownId,
			gameUnit.GUnitPosX,
			gameUnit.GUnitPosY,
			gameUnit.GUnitPosZ,
			gameUnit.GUnitTutorialCompleteLevel,
			gameUnit.GUnitIsFirstSlot,
			gameUnit.GUnitQuickSlotUnlock,
			gameUnit.GUnitRebirthCoin,
			gameUnit.GUnitEntrancePoint,
			gameUnit.GUnitCreateTime,
			gameUnit.GUnitLastUseTime,
			gameUnit.GUnitLastLeaveTime,
			gameUnit.GUnitRemoved,
			gameUnit.GUnitRemovedTime,
			gameUnit.GUitUnionPoint,
			gameUnit.GUnitVisualTitleId,
			gameUnit.GUnitAvatarMode,
			gameUnit.GUnitVisualFrameId,
			gameUnit.GUnitClosetSetSlotExpandCount,
			gameUnit.GUnitArtifactEnergy,
			gameUnit.GUnitArtifactOnOff,
			gameUnit.GUnitDamageFontId,
			gameUnit.GUnitDungeonClearCount)

		count++
	}

	logs.Info("total count: %d \n", count)

	// TODO: delete game_unit table before insert
	//_, err = conn.Exec(context.Background(), "delete from public.game_unit")
	_, err = conn.Exec(context.Background(), "TRUNCATE public.game_unit RESTART IDENTITY;")

	if err != nil {
		logs.Error("delete error: ", err)
		return
	}

	// bulk inserts
	br := conn.SendBatch(context.Background(), batch)
	for i := 0; i < count; i++ {
		_, err := br.Exec()
		if err != nil {
			//logs.Error("insert error: ", i, err)
			return
		}

		//logs.Info("count: ", i, "result: ", ct.RowsAffected())
	}

	end := time.Now()
	elapsed := time.Since(start)
	logs.Info("finish: ", end, " , elapsed: ", elapsed)

	libs.SendEmail(strconv.Itoa(count), start.String(), end.String(), elapsed.String())

	return
}

func main() {
	fmt.Printf("Start Get Game Data !\n")
	err := godotenv.Load()
	if err != nil {
		logs.Error("Error loading .env file")
	}

	//logging
	logs.SetLogger(logs.AdapterFile, `{"filename":"./logs/project.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":7,"color":true}`)

	//
	DB_CON := os.Getenv("DBHOST")
	conn, err := pgx.Connect(context.Background(), DB_CON)
	if err != nil {
		logs.Error("Unable to connect to database: %v", err)
		fmt.Printf("Unable to connect to database: %v\n", err)
		os.Exit(1)
		//fmt.Printf("Unable to connect to database: %v\n", err)
	}

	defer conn.Close(context.Background())

	logs.Info("connection success")
	//fmt.Printf("connection success")

	GetGameDataDaily(conn)

	fmt.Printf("End Get Game Data !\n")
	os.Exit(0)

}
