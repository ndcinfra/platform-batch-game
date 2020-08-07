package models

import "time"

// 유저의 게임 정보가 일 배치로 저장 된다.
// daily batch url: /v1/game/getGameDataDaily
/**
* see the file sql_game.sql. -- 1. get batch for game info
 */
type GameUnit struct {
	Id                            int64     `orm:"pk;auto" json:"id"` //
	GAccountUid                   int64     `json:"g_account_uid"`
	GAccountPublisherSn           int64     `json:"g_account_publisher_sn"`
	GAccountId                    string    `orm:"size(50);" json:"g_account_id"`
	GAccountInfo                  string    `json:"g_account_info"`
	GAccountNickName              string    `orm:"size(50);" json:"g_account_nick_name"`
	GAccountAccountLevel          int       `json:"g_account_account_level"`
	GUnitUid                      int64     `json:"g_unit_uid"`
	GUnitOwnerUid                 int64     `json:"g_unit_owner_uid"`
	GUnitUnitId                   string    `orm:"size(50)" json:"g_unit_unit_id"`
	GUnitUnitSubType              int       `json:"g_unit_unit_sub_type"`
	GUintLevel                    int8      `json:"g_uint_level"` //tinyint
	GUnitExp                      int       `json:"g_unit_exp"`
	GUnitGold                     int64     `json:"g_unit_gold"`
	GUnitPveSkillPoint            int16     `json:"g_unit_pve_skill_point"` //smallint
	GUnitPnaSkillPoint            int16     `json:"g_unit_pna_skill_point"`
	GUnitItemSkillPoint           int16     `json:"g_unit_item_skill_point"`
	GUnitClassUpgradeLevel        int8      `json:"g_unit_class_upgrade_level"`
	GUnitTutorialTask             int       `json:"g_unit_tutorial_task"`
	GUnitTitleId                  int       `json:"g_unit_title_id"`
	GUnitTownId                   int       `json:"g_unit_town_id"`
	GUnitPosX                     float32   `json:"g_unit_pos_x"`
	GUnitPosY                     float32   `json:"g_unit_pos_y"`
	GUnitPosZ                     float32   `json:"g_unit_pos_z"`
	GUnitTutorialCompleteLevel    int8      `json:"g_unit_tutorial_complete_level"`
	GUnitIsFirstSlot              int8      `json:"g_unit_is_first_slot"`
	GUnitQuickSlotUnlock          int8      `json:"g_unit_quick_slot_unlock"`
	GUnitRebirthCoin              int16     `json:"g_unit_rebirth_coin"`
	GUnitEntrancePoint            int16     `json:"g_unit_entrance_point"`
	GUnitCreateTime               time.Time `json:"g_unit_create_time"`
	GUnitLastUseTime              time.Time `json:"g_unit_last_use_time"`
	GUnitLastLeaveTime            time.Time `json:"g_unit_last_leave_time"`
	GUnitRemoved                  int8      `json:"g_unit_removed"`
	GUnitRemovedTime              time.Time `json:"g_unit_removed_time"`
	GUitUnionPoint                int       `json:"g_unit_union_point"`
	GUnitVisualTitleId            int       `json:"g_unit_visual_title_id"`
	GUnitAvatarMode               int8      `json:"g_unit_avatar_mode"`
	GUnitVisualFrameId            int       `json:"g_unit_visual_frame_id"`
	GUnitClosetSetSlotExpandCount int8      `json:"g_unit_closet_set_slot_expand_count"`
	GUnitArtifactEnergy           int16     `json:"g_unit_artifact_energy"`
	GUnitArtifactOnOff            bool      `json:"g_unit_artifact_on_off"`
	GUnitDamageFontId             int       `json:"g_unit_damage_font_id"`
	GUnitDungeonClearCount        int       `json:"g_unit_dungeon_clear_count"`
}
