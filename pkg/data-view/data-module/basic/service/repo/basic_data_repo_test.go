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
	. "github.com/jacklv111/common-sdk/test"
	. "github.com/smartystreets/goconvey/convey"
)

func TestBasicDataRepoImpl(t *testing.T) {
	defer DbSetUpAndTearDown()()
	repo := BasicDataRepoImpl{}
	Convey("Get name test", t, func() {
		idList := []uuid.UUID{uuid.New(), uuid.New()}
		rows := sqlmock.NewRows([]string{"id", "name", "type"}).
			AddRow(idList[0], "name1", "image").
			AddRow(idList[1], "name2", "image")

		Convey("Success", func() {
			Sqlmocker.ExpectQuery(regexp.QuoteMeta("SELECT `id`,`name`,`type` FROM `data_items` WHERE id in (?,?)")).
				WithArgs(idList[0], idList[1]).
				WillReturnRows(rows)

			res, err := repo.GetNameAndType(idList)
			So(err, ShouldEqual, err)
			So(len(res), ShouldEqual, 2)
			So(res[0].Name, ShouldEqual, "name1")
			So(res[1].Name, ShouldEqual, "name2")
			So(res[0].Type, ShouldEqual, "image")
			So(res[1].Type, ShouldEqual, "image")
		})
	})
}
