/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package repo

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	basicdo "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/do"
	"github.com/jacklv111/aifs/pkg/data-view/data-module/raw-data/image/do"
	. "github.com/jacklv111/common-sdk/test"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_imageRawDataRepoImpl(t *testing.T) {
	defer DbSetUpAndTearDown()()
	Convey("FindExistedBySha256", t, func() {
		var imageExtDoList []do.ImageExtDo
		imageExtDoList = append(imageExtDoList, do.ImageExtDo{ID: uuid.New(), Sha256: "testsha2561"})
		imageExtDoList = append(imageExtDoList, do.ImageExtDo{ID: uuid.New(), Sha256: "testsha2562"})

		Convey("Success", func() {
			rows := sqlmock.NewRows([]string{"id", "sha256"}).
				AddRow(imageExtDoList[0].ID, imageExtDoList[0].Sha256)

			Sqlmocker.ExpectQuery(regexp.QuoteMeta("SELECT `id`,`sha256` FROM `image_exts` WHERE sha256 in (?,?)")).
				WithArgs(imageExtDoList[0].Sha256, imageExtDoList[1].Sha256).WillReturnRows(rows)

			res, _ := ImageRawDataRepo.FindExistedByHash(do.GetHashList(imageExtDoList))

			So(len(res), ShouldEqual, 1)
			So(res[imageExtDoList[0].Sha256], ShouldEqual, imageExtDoList[0].ID)
		})
	})

	Convey("Create batch", t, func() {
		testId1 := uuid.New()
		testId2 := uuid.New()
		var imageExtDoList []do.ImageExtDo
		imageExtDoList = append(imageExtDoList, do.ImageExtDo{ID: testId1, Sha256: "testsha2561"})
		imageExtDoList = append(imageExtDoList, do.ImageExtDo{ID: testId2, Sha256: "testsha2562"})

		var dataItemDoList []basicdo.DataItemDo
		dataItemDoList = append(dataItemDoList, basicdo.DataItemDo{ID: testId1})
		dataItemDoList = append(dataItemDoList, basicdo.DataItemDo{ID: testId2})

		var imageScoreDoList []do.ImageScoreDo
		imageScoreDoList = append(imageScoreDoList, do.ImageScoreDo{ID: testId1})
		imageScoreDoList = append(imageScoreDoList, do.ImageScoreDo{ID: testId2})

		Convey("Success", func() {
			Sqlmocker.ExpectBegin()
			Sqlmocker.ExpectExec(regexp.QuoteMeta("INSERT INTO `data_items` (`id`,`name`,`type`,`create_at`) VALUES (?,?,?,?),(?,?,?,?) ON DUPLICATE KEY UPDATE `id`=`id`")).
				WithArgs().
				WillReturnResult(sqlmock.NewResult(2, 2))

			Sqlmocker.ExpectExec(regexp.QuoteMeta("INSERT INTO `image_exts` (`id`,`thumbnail`,`size`,`sha256`,`width`,`height`) VALUES (?,?,?,?,?,?),(?,?,?,?,?,?) ON DUPLICATE KEY UPDATE `id`=`id`")).
				WithArgs().
				WillReturnResult(sqlmock.NewResult(2, 2))

			Sqlmocker.ExpectExec(regexp.QuoteMeta("INSERT INTO `image_scores` (`id`,`light`,`dense`,`shelter`,`size`) VALUES (?,?,?,?,?),(?,?,?,?,?) ON DUPLICATE KEY UPDATE `id`=`id`")).
				WithArgs().
				WillReturnResult(sqlmock.NewResult(2, 2))

			Sqlmocker.ExpectCommit()

			err := ImageRawDataRepo.CreateBatch(dataItemDoList, imageExtDoList, imageScoreDoList)

			So(err, ShouldEqual, nil)
		})
	})

	Convey("Delete batch", t, func() {
		testId1 := uuid.New()
		testId2 := uuid.New()
		idList := []uuid.UUID{testId1, testId2}

		Convey("Success", func() {
			Sqlmocker.ExpectBegin()

			Sqlmocker.ExpectExec(regexp.QuoteMeta("DELETE FROM `data_items` WHERE id in (?,?)")).
				WithArgs(idList[0], idList[1]).
				WillReturnResult(sqlmock.NewResult(2, 2))

			Sqlmocker.ExpectExec(regexp.QuoteMeta("DELETE FROM `image_exts` WHERE id in (?,?)")).
				WithArgs(idList[0], idList[1]).
				WillReturnResult(sqlmock.NewResult(2, 2))

			Sqlmocker.ExpectExec(regexp.QuoteMeta("DELETE FROM `image_scores` WHERE id in (?,?)")).
				WithArgs(idList[0], idList[1]).
				WillReturnResult(sqlmock.NewResult(2, 2))

			Sqlmocker.ExpectCommit()

			err := ImageRawDataRepo.DeleteBatch(idList)

			So(err, ShouldEqual, nil)
		})
	})
}
