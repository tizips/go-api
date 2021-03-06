package stay

import (
	"encoding/json"
	"fmt"
	"github.com/bsm/redislock"
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"saas/app/constant"
	"saas/app/model"
	"saas/app/request/admin/dormitory/stay"
	stayResponse "saas/app/response/admin/dormitory/stay"
	"saas/kernel/config"
	"saas/kernel/data"
	"saas/kernel/response"
	"time"
)

func ToPeopleByPaginate(ctx *gin.Context) {

	var request stay.ToPeopleByPaginate
	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	tx := data.Database.Where(fmt.Sprintf("%s.`status`=?", model.TableDorPeople), request.Status)

	if request.Floor > 0 {
		tx = tx.Where("`floor_id`=?", request.Floor)
	} else if request.Building > 0 {
		tx = tx.Where("`building_id`=?", request.Building)
	}

	if request.IsTemp > 0 {
		tx = tx.Where(fmt.Sprintf("%s.`is_temp`=?", model.TableDorPeople), request.IsTemp)
	}

	if request.Keyword != "" {

		condition := data.Database.Select("1")

		if request.Type == "mobile" {
			condition = condition.
				Table(model.TableMemMember).
				Where(fmt.Sprintf("%s.`member_id`=%s.`id`", model.TableDorPeople, model.TableMemMember)).
				Where(fmt.Sprintf("%s.`mobile`=?", model.TableMemMember), request.Keyword)
		} else if request.Type == "room" {
			condition = condition.
				Table(model.TableDorRoom).
				Where(fmt.Sprintf("%s.`room_id`=%s.`id`", model.TableDorPeople, model.TableDorRoom)).
				Where(fmt.Sprintf("%s.`name`=?", model.TableDorRoom), request.Keyword)
		} else {
			condition = condition.
				Table(model.TableMemMember).
				Where(fmt.Sprintf("%s.`member_id`=%s.`id`", model.TableDorPeople, model.TableMemMember)).
				Where(fmt.Sprintf("%s.`name`=?", model.TableMemMember), request.Keyword)
		}

		tx = tx.Where("exists (?)", condition)
	}

	tc := tx

	responses := response.Paginate{
		Total: 0,
		Page:  request.GetPage(),
		Size:  request.GetSize(),
		Data:  make([]any, 0),
	}

	tc.Model(&model.DorPeople{}).Count(&responses.Total)

	if responses.Total > 0 {

		var peoples []model.DorPeople

		tx.
			Preload("Member", func(db *gorm.DB) *gorm.DB { return db.Unscoped() }).
			Preload("Staff", func(db *gorm.DB) *gorm.DB { return db.Unscoped() }).
			Preload("Certification", func(db *gorm.DB) *gorm.DB { return db.Unscoped() }).
			Preload("Category", func(db *gorm.DB) *gorm.DB { return db.Unscoped() }).
			Preload("Building", func(db *gorm.DB) *gorm.DB { return db.Unscoped() }).
			Preload("Floor", func(db *gorm.DB) *gorm.DB { return db.Unscoped() }).
			Preload("Room", func(db *gorm.DB) *gorm.DB { return db.Unscoped() }).
			Preload("Bed", func(db *gorm.DB) *gorm.DB { return db.Unscoped() }).
			Order(fmt.Sprintf("%s.`id` desc", model.TableDorPeople)).
			Offset(request.GetOffset()).
			Limit(request.GetLimit()).
			Find(&peoples)

		for _, item := range peoples {

			results := stayResponse.ToPeopleByPaginate{
				Id:        item.Id,
				Category:  item.Category.Name,
				Building:  item.Building.Name,
				Floor:     item.Floor.Name,
				Room:      item.Room.Name,
				Bed:       item.Bed.Name,
				Name:      item.Member.Name,
				Mobile:    item.Member.Mobile,
				IsTemp:    item.Category.IsTemp,
				Start:     item.Start.ToDateString(),
				Remark:    item.Remark,
				CreatedAt: item.CreatedAt.ToDateTimeString(),
			}
			if item.End != nil {
				results.End = item.End.ToDateTimeString()
			}
			if item.Staff != nil && item.Staff.Id > 0 {
				results.Staff = item.Staff.Status
				results.Titles = item.Staff.Title
			}
			if item.Certification != nil && item.Certification.Id > 0 {
				certification := stayResponse.ToPeopleByPaginateOfCertification{
					No: item.Certification.No,
				}
				results.Certification = certification
			}
			responses.Data = append(responses.Data, results)
		}
	}

	response.SuccessByPaginate(ctx, responses)
}

func DoPeopleByCreate(ctx *gin.Context) {

	var request stay.DoPeopleByCreate
	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	var bed model.DorBed
	data.Database.Preload(clause.Associations).Where("`is_enable`=?", constant.IsEnableYes).Find(&bed, request.Bed)
	if bed.Id <= 0 {
		response.NotFound(ctx, "???????????????")
		return
	}

	var category model.DorStayCategory
	data.Database.Where("is_enable", constant.IsEnableYes).Find(&category, request.Category)
	if category.Id <= 0 {
		response.NotFound(ctx, "???????????????")
		return
	}

	var count int64 = 0
	data.Database.Model(model.DorPeople{}).Joins(fmt.Sprintf("left join `%s` on `%s`.`member_id`=`%s`.`id`", model.TableMemMember, model.TableDorPeople, model.TableMemMember)).Where(fmt.Sprintf("`%s`.`mobile`=? and `%s`.`status`=?", model.TableMemMember, model.TableDorPeople), request.Mobile, model.DorPeopleStatusLive).Count(&count)
	if count > 0 {
		response.Fail(ctx, "???????????????????????????????????????????????????")
		return
	}

	var member model.MemMember
	data.Database.Where("`mobile`=?", request.Mobile).Find(&member)
	if member.Id == "" {

		lock, err := redislock.New(data.Redis).Obtain(ctx, "lock:member:"+request.Mobile, time.Second*30, nil)
		if err != nil {
			response.Fail(ctx, "????????????")
			return
		}

		defer lock.Release(ctx)

		node, err := snowflake.NewNode(config.Values.Server.Node)
		if err != nil {
			response.Fail(ctx, "????????????")
			return
		}

		member = model.MemMember{
			Id:       node.Generate().String(),
			Mobile:   request.Mobile,
			Name:     request.Name,
			Nickname: request.Name,
			IsEnable: constant.IsEnableYes,
		}

		if t := data.Database.Create(&member); t.RowsAffected <= 0 {
			response.Fail(ctx, "????????????")
			return
		}
	}

	tx := data.Database.Begin()

	people := model.DorPeople{
		CategoryId: category.Id,
		BuildingId: bed.BuildingId,
		FloorId:    bed.FloorId,
		RoomId:     bed.RoomId,
		BedId:      bed.Id,
		TypeId:     bed.TypeId,
		MemberId:   member.Id,
		Start:      carbon.Date{Carbon: carbon.Parse(request.Start)},
		Status:     model.DorPeopleStatusLive,
		IsTemp:     category.IsTemp,
		Remark:     request.Remark,
	}

	if request.End != "" {
		people.End = &carbon.Date{Carbon: carbon.Parse(request.End)}
	}

	if t := tx.Create(&people); t.RowsAffected <= 0 {
		tx.Rollback()
		response.Fail(ctx, "????????????")
		return
	}

	var masterId int = 0
	masterName := member.Name

	var master model.DorPeople
	tx.Preload("Master.Member").Where("`bed_id`=? and `master_id`<>? and `status`=?", people.BedId, 0, model.DorPeopleStatusLive).Find(&master)
	if master.MasterId > 0 {
		masterId = master.MasterId
		masterName = master.Member.Name
	} else {
		masterId = people.Id
	}

	people.MasterId = masterId

	if t := tx.Model(&model.DorPeople{}).Where("`id`=?", people.Id).UpdateColumn("`master_id`=?", masterId); t.RowsAffected <= 0 {
		tx.Rollback()
		response.Fail(ctx, "????????????")
		return
	}

	var details = make(map[string]any, 11)

	details["category"] = category.Name
	details["building"] = bed.Building.Name
	details["floor"] = bed.Floor.Name
	details["room"] = bed.Room.Name
	details["bed"] = bed.Name
	details["type"] = bed.Type.Name
	details["name"] = member.Name
	details["mobile"] = member.Mobile
	details["master"] = masterName
	details["is_temp"] = category.IsTemp
	details["start"] = people.Start
	details["end"] = people.End

	str, _ := json.Marshal(details)

	log := model.DorPeopleLog{
		PeopleId: people.Id,
		MemberId: member.Id,
		Status:   model.DorPeopleLogStatusLive,
		Detail:   string(str),
		Remark:   people.Remark,
	}

	if t := tx.Create(&log); t.RowsAffected <= 0 {
		tx.Rollback()
		response.Fail(ctx, "????????????")
		return
	}

	tx.Commit()

	response.Success(ctx)
}
