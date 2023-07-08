/*
 * Created on Sat Jul 08 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package repo

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jacklv111/aifs/pkg/data-view/data-view/do"
	vb "github.com/jacklv111/aifs/pkg/data-view/data-view/value-object"
	"github.com/jacklv111/common-sdk/log"
	. "github.com/jacklv111/common-sdk/test"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_dataViewRepoImpl(t *testing.T) {
	log.ValidateAndApply(log.LogConfig)
	defer DbSetUpAndTearDown()()

	Convey("Create data view", t, func() {
		createDo := do.DataViewDo{ID: uuid.New(), Name: "name0", ViewType: "raw-data", Description: "desc", RawDataType: "image"}
		Convey("Create success", func() {
			Sqlmocker.ExpectBegin()
			Sqlmocker.ExpectExec(regexp.QuoteMeta("INSERT INTO `data_views`")).
				WithArgs().
				WillReturnResult(sqlmock.NewResult(1, 1))
			Sqlmocker.ExpectCommit()
			err := DataViewRepo.Create(createDo)
			So(err, ShouldEqual, nil)
		})

		Convey("Create failed, insert conflict", func() {
			Sqlmocker.ExpectBegin()
			Sqlmocker.ExpectExec(regexp.QuoteMeta("INSERT INTO `data_views`")).
				WithArgs().
				WillReturnError(fmt.Errorf("insert conflict"))
			Sqlmocker.ExpectRollback()
			err := DataViewRepo.Create(createDo)
			So(err, ShouldNotEqual, nil)
		})
	})

	Convey("Get details by id", t, func() {
		testId := uuid.New()
		Convey("Success", func() {
			dataViewRow := sqlmock.NewRows([]string{"id", "name", "view_type", "description", "create_at", "update_at", "delete_at"}).
				AddRow(testId.String(), "name0", "raw-data", "", 1672887665782, 1672887665782, 0)
			Sqlmocker.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `data_views` WHERE `data_views`.`id` = ? AND `data_views`.`delete_at` = ? ORDER BY `data_views`.`id` LIMIT 1")).
				WithArgs(testId.String(), 0).
				WillReturnRows(dataViewRow)

			res, _ := DataViewRepo.GetById(testId)

			So(res.ID, ShouldEqual, testId)
			So(res.Name, ShouldEqual, "name0")
			So(res.ViewType, ShouldEqual, "raw-data")
		})
	})

	Convey("Get data view item count", t, func() {
		testId := uuid.New()

		Sqlmocker.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `data_view_items` WHERE data_view_id = ?")).
			WithArgs(testId.String()).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(10))

		count, _ := DataViewRepo.GetDataViewItemCount(testId)
		So(count, ShouldEqual, 10)
	})

	Convey("Get list", t, func() {
		testId1 := uuid.New()
		testId2 := uuid.New()
		dataViewIdList := []string{testId1.String(), testId2.String()}
		options := vb.DataViewListQueryOptions{Offset: 10, Limit: 10, DataViewIdList: dataViewIdList}

		Convey("Success", func() {
			dataViewRow := sqlmock.NewRows([]string{"id", "name", "type", "description", "create_at", "update_at", "delete_at"}).
				AddRow(testId1.String(), "name0", "raw-data", "", 1672887665782, 1672887665782, 0).
				AddRow(testId2.String(), "name0", "raw-data", "", 1672887665782, 1672887665782, 0)

			Sqlmocker.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `data_views` WHERE id in (?,?) AND `data_views`.`delete_at` = ? LIMIT 10 OFFSET 10")).
				WithArgs(testId1.String(), testId2.String(), 0).
				WillReturnRows(dataViewRow)

			res, err := DataViewRepo.GetList(options)
			So(err, ShouldEqual, nil)
			So(len(res), ShouldEqual, 2)
			So(res[0].ID, ShouldEqual, testId1)
			So(res[1].ID, ShouldEqual, testId2)
		})
	})

	Convey("Hard delete data view", t, func() {
		testId := uuid.New()
		Convey("Success", func() {
			Sqlmocker.ExpectBegin()
			Sqlmocker.ExpectExec(regexp.QuoteMeta("DELETE FROM `data_views` WHERE `data_views`.`id` = ?")).
				WithArgs(testId.String()).
				WillReturnResult(sqlmock.NewResult(1, 1))
			Sqlmocker.ExpectExec(regexp.QuoteMeta("DELETE FROM `data_view_items` WHERE data_view_id = ?")).
				WithArgs(testId.String()).
				WillReturnResult(sqlmock.NewResult(1, 1))
			Sqlmocker.ExpectCommit()
			err := DataViewRepo.HardDelete(testId)
			So(err, ShouldEqual, nil)
		})
		Convey("Failed", func() {
			Sqlmocker.ExpectBegin()
			Sqlmocker.ExpectExec(regexp.QuoteMeta("DELETE FROM `data_views` WHERE `data_views`.`id` = ?")).
				WithArgs(testId.String()).
				WillReturnResult(sqlmock.NewResult(1, 1))
			Sqlmocker.ExpectExec(regexp.QuoteMeta("DELETE FROM `data_view_items` WHERE data_view_id = ?")).
				WithArgs(testId.String()).
				WillReturnError(fmt.Errorf("delete data view item error"))
			Sqlmocker.ExpectRollback()
			err := DataViewRepo.HardDelete(testId)
			So(err, ShouldNotEqual, nil)
		})
	})

	Convey("Soft delete data view", t, func() {
		testId := uuid.New()
		Convey("Success", func() {
			Sqlmocker.ExpectBegin()
			Sqlmocker.ExpectExec(regexp.QuoteMeta("UPDATE `data_views` SET `delete_at`=? WHERE `data_views`.`id` = ? AND `data_views`.`delete_at` = ?")).
				WithArgs(sqlmock.AnyArg(), testId.String(), 0).
				WillReturnResult(sqlmock.NewResult(1, 1))
			Sqlmocker.ExpectCommit()
			err := DataViewRepo.SoftDelete(testId)
			So(err, ShouldEqual, nil)
		})
	})

	Convey("Delete data view item", t, func() {
		dataViewId := uuid.New()
		dataItemIdList := []uuid.UUID{uuid.New(), uuid.New()}
		Convey("Success", func() {
			Sqlmocker.ExpectBegin()
			Sqlmocker.ExpectExec(regexp.QuoteMeta("DELETE FROM `data_view_items` WHERE data_view_id = ? AND data_item_id in (?,?)")).
				WithArgs(dataViewId.String(), dataItemIdList[0], dataItemIdList[1]).
				WillReturnResult(sqlmock.NewResult(1, 2))
			Sqlmocker.ExpectCommit()

			err := DataViewRepo.DeleteDataViewItem(dataViewId, dataItemIdList)
			So(err, ShouldEqual, nil)
		})

		Convey("Failed, some data item do not exist", func() {
			Sqlmocker.ExpectBegin()
			Sqlmocker.ExpectExec(regexp.QuoteMeta("DELETE FROM `data_view_items` WHERE data_view_id = ? AND data_item_id in (?,?)")).
				WithArgs(dataViewId.String(), dataItemIdList[0], dataItemIdList[1]).
				WillReturnResult(sqlmock.NewResult(1, 1))
			Sqlmocker.ExpectRollback()

			err := DataViewRepo.DeleteDataViewItem(dataViewId, dataItemIdList)
			So(err, ShouldNotEqual, nil)
		})
	})

	Convey("Create raw data view items", t, func() {
		dataViewId := uuid.New()
		createDoList := []uuid.UUID{uuid.New(), uuid.New()}
		Convey("Create success", func() {
			Sqlmocker.ExpectBegin()
			Sqlmocker.ExpectExec(regexp.QuoteMeta("INSERT INTO `data_view_items` (`data_view_id`,`data_item_id`) VALUES (?,?),(?,?)")).
				WithArgs(dataViewId, createDoList[0], dataViewId, createDoList[1]).
				WillReturnResult(sqlmock.NewResult(2, 2))
			Sqlmocker.ExpectCommit()
			err := DataViewRepo.CreateDataViewItemsIgnoreConflict(dataViewId, createDoList)
			So(err, ShouldEqual, nil)
		})

		Convey("Create failed, insert conflict", func() {
			Sqlmocker.ExpectBegin()
			Sqlmocker.ExpectExec(regexp.QuoteMeta("INSERT INTO `data_view_items` (`data_view_id`,`data_item_id`) VALUES (?,?),(?,?)")).
				WithArgs().
				WillReturnError(fmt.Errorf("insert conflict"))
			Sqlmocker.ExpectRollback()
			err := DataViewRepo.CreateDataViewItemsIgnoreConflict(dataViewId, createDoList)
			So(err, ShouldNotEqual, nil)
			So(err.Error(), ShouldEqual, "insert conflict")
		})
	})

	Convey("Create annotation data view items", t, func() {
		annoTempId := uuid.New()
		dataViewId := uuid.New()
		createDoList := []do.DataViewItemDo{
			{DataViewId: dataViewId, DataItemId: uuid.New()},
			{DataViewId: dataViewId, DataItemId: uuid.New()},
		}
		annoIdToBeDeleted := []string{uuid.NewString(), uuid.NewString()}
		rows := sqlmock.NewRows([]string{"id"}).AddRow(annoIdToBeDeleted[0]).AddRow(annoIdToBeDeleted[1])
		Convey("Create success", func() {
			deleteSql := "DELETE FROM `data_view_items` WHERE data_item_id in (?,?) AND data_view_id = ?"
			Sqlmocker.ExpectBegin()
			Sqlmocker.ExpectQuery(regexp.QuoteMeta("SELECT annotations.id FROM `data_view_items` left join annotations on data_view_items.data_item_id=annotations.id WHERE data_view_items.data_view_id = ? AND annotations.annotation_template_id = ? AND annotations.data_item_id in (SELECT annotations.data_item_id FROM `annotations` WHERE annotations.id in (?,?))")).
				WithArgs(dataViewId, annoTempId, createDoList[0].DataItemId, createDoList[1].DataItemId).
				WillReturnRows(rows)

			Sqlmocker.ExpectExec(regexp.QuoteMeta(deleteSql)).
				WithArgs(annoIdToBeDeleted[0], annoIdToBeDeleted[1], dataViewId).
				WillReturnResult(sqlmock.NewResult(2, 2))

			Sqlmocker.ExpectExec(regexp.QuoteMeta("INSERT INTO `data_view_items` (`data_view_id`,`data_item_id`) VALUES (?,?),(?,?)")).
				WithArgs(createDoList[0].DataViewId, createDoList[0].DataItemId, createDoList[1].DataViewId, createDoList[1].DataItemId).
				WillReturnResult(sqlmock.NewResult(2, 2))
			Sqlmocker.ExpectCommit()
			err := DataViewRepo.CreateAnnotationDataViewItems(createDoList, annoTempId)
			So(err, ShouldEqual, nil)
		})
	})

	Convey("Get all data view items", t, func() {
		dataViewId := uuid.New()
		Convey("Create success", func() {
			dataViewRow := sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()).AddRow(uuid.New())
			Sqlmocker.ExpectQuery(regexp.QuoteMeta("SELECT `data_item_id` FROM `data_view_items` WHERE data_view_id = ?")).
				WithArgs(dataViewId).
				WillReturnRows(dataViewRow)
			idList, err := DataViewRepo.GetAllDataViewItems(dataViewId)
			So(err, ShouldEqual, nil)
			So(len(idList), ShouldEqual, 2)
		})
	})

	Convey("all exists", t, func() {
		id1 := uuid.New()
		id2 := uuid.New()
		Convey("Success, all exists", func() {
			row := sqlmock.NewRows([]string{"id"}).AddRow(id1.String()).AddRow(id2.String())
			Sqlmocker.ExpectQuery(regexp.QuoteMeta("SELECT `id` FROM `data_views` WHERE id in (?,?) AND `data_views`.`delete_at` = ?")).
				WithArgs(id1, id2, 0).
				WillReturnRows(row)
			res, err := DataViewRepo.GetInvalidId([]uuid.UUID{id1, id2})
			So(err, ShouldEqual, nil)
			So(len(res), ShouldEqual, 0)
		})
		Convey("Success, not all exists", func() {
			row := sqlmock.NewRows([]string{"id"}).AddRow(id1.String())
			Sqlmocker.ExpectQuery(regexp.QuoteMeta("SELECT `id` FROM `data_views` WHERE id in (?,?) AND `data_views`.`delete_at` = ?")).
				WithArgs(id1, id2, 0).
				WillReturnRows(row)
			res, err := DataViewRepo.GetInvalidId([]uuid.UUID{id1, id2})
			So(err, ShouldEqual, nil)
			So(len(res), ShouldEqual, 1)
			So(res[0], ShouldEqual, id2)
		})
	})

	Convey("exists by id", t, func() {
		id1 := uuid.New()
		Convey("Success, exists", func() {
			row := sqlmock.NewRows([]string{"count"}).AddRow(1)
			Sqlmocker.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `data_views` WHERE id = ? AND `data_views`.`delete_at` = ?")).
				WithArgs(id1, 0).
				WillReturnRows(row)
			res, err := DataViewRepo.ExistsById(id1)
			So(err, ShouldEqual, nil)
			So(res, ShouldBeTrue)
		})
	})

	Convey("get invalid data items", t, func() {
		dvId := uuid.New()
		queryDataItemIdList := []uuid.UUID{uuid.New(), uuid.New()}
		invalidId := queryDataItemIdList[0]
		Convey("Success", func() {
			row := sqlmock.NewRows([]string{"data_item_id"}).AddRow(invalidId)
			Sqlmocker.ExpectQuery(regexp.QuoteMeta("SELECT `data_item_id` FROM `data_view_items` WHERE data_view_id = (?) and data_item_id in (?,?)")).
				WithArgs(dvId, queryDataItemIdList[0], queryDataItemIdList[1]).
				WillReturnRows(row)
			res, err := DataViewRepo.GetInvalidDataItems(dvId, queryDataItemIdList)
			So(err, ShouldEqual, nil)
			So(res, ShouldResemble, []uuid.UUID{queryDataItemIdList[1]})
		})
	})

	Convey("get invalid data items", t, func() {
		srcAnnoViewId := uuid.New()
		destAnnoViewId := uuid.New()
		rawDataViewId := uuid.New()
		annoDataItemId1 := uuid.New()
		annoDataItemId2 := uuid.New()

		Convey("Success", func() {
			Sqlmocker.ExpectBegin()
			Sqlmocker.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `data_view_items` WHERE data_view_id = ?")).
				WithArgs(destAnnoViewId).
				WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

			Sqlmocker.ExpectQuery(regexp.QuoteMeta("SELECT annotations.id as id FROM `data_view_items` left join annotations on data_view_items.data_item_id=annotations.id WHERE data_view_items.data_view_id = ? AND annotations.data_item_id in (SELECT data_item_id FROM `data_view_items` WHERE data_view_id = ?)")).
				WithArgs(srcAnnoViewId, rawDataViewId).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(annoDataItemId1).AddRow(annoDataItemId2))

			Sqlmocker.ExpectExec(regexp.QuoteMeta("INSERT INTO `data_view_items` (`data_view_id`,`data_item_id`) VALUES (?,?),(?,?)")).
				WithArgs(destAnnoViewId, annoDataItemId1, destAnnoViewId, annoDataItemId2).
				WillReturnResult(sqlmock.NewResult(2, 2))

			Sqlmocker.ExpectCommit()
			err := DataViewRepo.FilterAnnotationsByRawData(srcAnnoViewId, rawDataViewId, destAnnoViewId)
			So(err, ShouldEqual, nil)
		})
	})

	Convey("data view move to", t, func() {
		srcViewId := uuid.New()
		destViewId := uuid.New()

		Convey("Success", func() {
			Sqlmocker.ExpectBegin()
			Sqlmocker.ExpectExec(regexp.QuoteMeta("UPDATE `data_view_items` SET `data_view_id`=? WHERE data_view_id = ?")).
				WithArgs(destViewId, srcViewId).
				WillReturnResult(sqlmock.NewResult(2, 2))
			Sqlmocker.ExpectCommit()
			err := DataViewRepo.MoveTo(srcViewId, destViewId)
			So(err, ShouldEqual, nil)
		})
	})
}
