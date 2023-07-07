/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package repo

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jacklv111/aifs/pkg/annotation-template/do"
	vb "github.com/jacklv111/aifs/pkg/annotation-template/value-object"
	. "github.com/jacklv111/common-sdk/test"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_annotationTemplateRepoImpl(t *testing.T) {
	defer DbSetUpAndTearDown()()
	Convey("Create annotation template", t, func() {
		testId := uuid.New()
		var labels []do.LabelDo
		labels = append(labels, do.LabelDo{ID: uuid.New(), Name: "labelName1", Color: 10, AnnotationTemplateId: testId})
		labels = append(labels, do.LabelDo{ID: uuid.New(), Name: "labelName2", Color: 20, AnnotationTemplateId: testId})
		annoTempDo := do.AnnotationTemplateDo{Type: "testType", Name: "testName", Description: "testDesc", ID: testId}
		annoTempExtDo := do.AnnotationTemplateExtDo{}

		insLabelArgs := []driver.Value{
			sqlmock.AnyArg(), testId, labels[0].Name, "", labels[0].Color, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), nil, nil, "",
			sqlmock.AnyArg(), testId, labels[1].Name, "", labels[1].Color, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), nil, nil, "",
		}
		Convey("Create success", func() {
			Sqlmocker.ExpectBegin()
			Sqlmocker.ExpectExec(regexp.QuoteMeta("INSERT INTO `annotation_templates` (`id`,`name`,`type`,`description`,`create_at`,`update_at`,`delete_at`) VALUES (?,?,?,?,?,?,?)")).
				WithArgs(sqlmock.AnyArg(), annoTempDo.Name, annoTempDo.Type, annoTempDo.Description, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
				WillReturnResult(sqlmock.NewResult(1, 1))
			Sqlmocker.ExpectExec(regexp.QuoteMeta("INSERT INTO `labels` (`id`,`annotation_template_id`,`name`,`super_category_name`,`color`,`create_at`,`update_at`,`delete_at`,`key_point_def`,`key_point_skeleton`,`cover_image_url`) VALUES (?,?,?,?,?,?,?,?,?,?,?),(?,?,?,?,?,?,?,?,?,?,?)")).
				WithArgs(insLabelArgs...).
				WillReturnResult(sqlmock.NewResult(1, 1))
			Sqlmocker.ExpectCommit()

			err := AnnotationTemplateRepo.Create(annoTempDo, annoTempExtDo, labels)
			So(err, ShouldEqual, nil)
		})

		Convey("Create failed, insert conflict", func() {
			Sqlmocker.ExpectBegin()
			Sqlmocker.ExpectExec(regexp.QuoteMeta("INSERT INTO `annotation_templates`")).
				WithArgs(sqlmock.AnyArg(), annoTempDo.Name, annoTempDo.Type, annoTempDo.Description, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
				WillReturnResult(sqlmock.NewResult(1, 1))
			Sqlmocker.ExpectExec(regexp.QuoteMeta("INSERT INTO `labels` (`id`,`annotation_template_id`,`name`,`super_category_name`,`color`,`create_at`,`update_at`,`delete_at`,`key_point_def`,`key_point_skeleton`,`cover_image_url`) VALUES (?,?,?,?,?,?,?,?,?,?,?),(?,?,?,?,?,?,?,?,?,?,?)")).
				WithArgs(insLabelArgs...).
				WillReturnError(fmt.Errorf("conflict detected"))
			Sqlmocker.ExpectRollback()

			err := AnnotationTemplateRepo.Create(annoTempDo, annoTempExtDo, labels)
			So(err, ShouldNotEqual, nil)
		})
	})

	Convey("Get list", t, func() {
		options := vb.ListQueryOptions{Offset: 1, Limit: 10, AnnoTemplateIdList: []string{"1", "2"}}

		Convey("Success", func() {
			rows := sqlmock.NewRows([]string{"id", "name", "create_at", "view_type", "label_count", "delete_at"}).AddRow("uuid0", "name0", 1672887665782, "bounding-box", 10, 0).AddRow("uuid1", "name1", 1672887665782, "bounding-box", 5, 0)
			Sqlmocker.ExpectQuery(regexp.QuoteMeta("SELECT annotation_templates.*, count(labels.id) as label_count FROM `annotation_templates` left join labels on annotation_templates.id=labels.annotation_template_id WHERE annotation_templates.id in (?,?) AND `annotation_templates`.`delete_at` = ? GROUP BY `annotation_templates`.`id` LIMIT 10 OFFSET 1")).
				WithArgs(options.AnnoTemplateIdList[0], options.AnnoTemplateIdList[1], 0).
				WillReturnRows(rows)
			res, _ := AnnotationTemplateRepo.GetList(options)
			So(len(res), ShouldEqual, 2)

			So(res[0].LabelCount, ShouldEqual, 10)
			So(res[0].Id, ShouldEqual, "uuid0")
			So(res[0].Name, ShouldEqual, "name0")

			So(res[1].LabelCount, ShouldEqual, 5)
			So(res[1].Id, ShouldEqual, "uuid1")
			So(res[1].Name, ShouldEqual, "name1")
		})
	})

	Convey("Get details by id", t, func() {
		testId := uuid.New()
		Convey("Success", func() {
			annoTempRow := sqlmock.NewRows([]string{"id", "name", "type", "description", "create_at", "update_at", "delete_at"}).AddRow(testId.String(), "name0", "bounding-box", "", 1672887665782, 1672887665782, 0)
			wordList, _ := json.Marshal([]string{"a"})
			annoTempExtRow := sqlmock.NewRows([]string{"annotation_template_id", "word_list"}).AddRow(testId.String(), wordList)
			labelRow := sqlmock.NewRows([]string{"annotation_template_id", "name", "color"}).AddRow(testId.String(), "name1", 10).AddRow(testId.String(), "name2", 15)
			Sqlmocker.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `annotation_templates` WHERE `annotation_templates`.`id` = ? AND `annotation_templates`.`delete_at` = ? ORDER BY `annotation_templates`.`id` LIMIT 1")).WithArgs(testId.String(), 0).WillReturnRows(annoTempRow)
			Sqlmocker.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `annotation_template_exts` WHERE annotation_template_id = ? AND `annotation_template_exts`.`delete_at` = ? ORDER BY `annotation_template_exts`.`annotation_template_id` LIMIT 1")).WithArgs(testId.String(), 0).WillReturnRows(annoTempExtRow)
			Sqlmocker.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `labels` WHERE `labels`.`annotation_template_id` = ? AND `labels`.`delete_at` = ?")).WithArgs(testId.String(), 0).WillReturnRows(labelRow)
			annoTempDo, annoTempExtDo, labelDoList, _ := AnnotationTemplateRepo.GetById(testId)
			So(annoTempDo.ID, ShouldEqual, testId)
			So(annoTempDo.Name, ShouldEqual, "name0")
			So(annoTempDo.Type, ShouldEqual, "bounding-box")
			So(annoTempExtDo.WordList, ShouldResemble, do.WordListType([]string{"a"}))

			So(len(labelDoList), ShouldEqual, 2)
			So(labelDoList[0].Name, ShouldEqual, "name1")
			So(labelDoList[0].Color, ShouldEqual, 10)
			So(labelDoList[1].Name, ShouldEqual, "name2")
			So(labelDoList[1].Color, ShouldEqual, 15)
		})

		Convey("Failed, db error, get annotation template error", func() {
			Sqlmocker.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `annotation_templates` WHERE `annotation_templates`.`id` = ? AND `annotation_templates`.`delete_at` = ? ORDER BY `annotation_templates`.`id` LIMIT 1")).WithArgs(testId.String(), 0).WillReturnError(fmt.Errorf("db error"))
			_, _, _, err := AnnotationTemplateRepo.GetById(testId)
			So(err, ShouldNotEqual, nil)
			So(err.Error(), ShouldEqual, "db error")
		})

		Convey("Failed, no annotation template exists", func() {
			annoTempRow := sqlmock.NewRows([]string{"id", "name", "view_type", "description", "create_at", "update_at", "delete_at"})

			Sqlmocker.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `annotation_templates` WHERE `annotation_templates`.`id` = ? AND `annotation_templates`.`delete_at` = ? ORDER BY `annotation_templates`.`id` LIMIT 1")).WithArgs(testId.String(), 0).WillReturnRows(annoTempRow)
			_, _, _, err := AnnotationTemplateRepo.GetById(testId)
			So(err, ShouldNotEqual, nil)
			So(err.Error(), ShouldStartWith, "annotation template not found")
		})

		Convey("Failed, db error, get label error", func() {
			annoTempRow := sqlmock.NewRows([]string{"id", "name", "view_type", "description", "create_at", "update_at", "delete_at"}).
				AddRow(testId.String(), "name0", "bounding-box", "", 1672887665782, 1672887665782, 0)
			wordList, _ := json.Marshal([]string{"a"})
			annoTempExtRow := sqlmock.NewRows([]string{"annotation_template_id", "word_list"}).AddRow(testId.String(), wordList)
			Sqlmocker.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `annotation_templates` WHERE `annotation_templates`.`id` = ? AND `annotation_templates`.`delete_at` = ? ORDER BY `annotation_templates`.`id` LIMIT 1")).WithArgs(testId.String(), 0).WillReturnRows(annoTempRow)
			Sqlmocker.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `annotation_template_exts` WHERE annotation_template_id = ? AND `annotation_template_exts`.`delete_at` = ? ORDER BY `annotation_template_exts`.`annotation_template_id` LIMIT 1")).WithArgs(testId.String(), 0).WillReturnRows(annoTempExtRow)
			Sqlmocker.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `labels` WHERE `labels`.`annotation_template_id` = ? AND `labels`.`delete_at` = ?")).WithArgs(testId.String(), 0).WillReturnError(fmt.Errorf("db error"))
			_, _, _, err := AnnotationTemplateRepo.GetById(testId)
			So(err, ShouldNotEqual, nil)
			So(err.Error(), ShouldEqual, "db error")
		})
	})

	Convey("Delete by id", t, func() {
		testId := uuid.New()
		Convey("Success", func() {
			Sqlmocker.ExpectBegin()
			Sqlmocker.ExpectExec(regexp.QuoteMeta("UPDATE `annotation_templates` SET `delete_at`=? WHERE `annotation_templates`.`id` = ? AND `annotation_templates`.`delete_at` = ?")).
				WithArgs(sqlmock.AnyArg(), testId, 0).
				WillReturnResult(sqlmock.NewResult(1, 1))
			Sqlmocker.ExpectExec(regexp.QuoteMeta("UPDATE `annotation_template_exts` SET `delete_at`=? WHERE annotation_template_id = ? AND `annotation_template_exts`.`delete_at` = ?")).
				WithArgs(sqlmock.AnyArg(), testId, 0).
				WillReturnResult(sqlmock.NewResult(1, 1))
			Sqlmocker.ExpectExec(regexp.QuoteMeta("UPDATE `labels` SET `delete_at`=? WHERE annotation_template_id = ? AND `labels`.`delete_at` = ?")).
				WithArgs(sqlmock.AnyArg(), testId, 0).
				WillReturnResult(sqlmock.NewResult(1, 1))
			Sqlmocker.ExpectCommit()
			err := AnnotationTemplateRepo.Delete(testId)
			So(err, ShouldEqual, nil)
		})
		Convey("Failed, delete labels error", func() {
			Sqlmocker.ExpectBegin()
			Sqlmocker.ExpectExec(regexp.QuoteMeta("UPDATE `annotation_templates` SET `delete_at`=? WHERE `annotation_templates`.`id` = ? AND `annotation_templates`.`delete_at` = ?")).
				WithArgs(sqlmock.AnyArg(), testId, 0).
				WillReturnResult(sqlmock.NewResult(1, 1))
			Sqlmocker.ExpectExec(regexp.QuoteMeta("UPDATE `annotation_template_exts` SET `delete_at`=? WHERE annotation_template_id = ? AND `annotation_template_exts`.`delete_at` = ?")).
				WithArgs(sqlmock.AnyArg(), testId, 0).
				WillReturnResult(sqlmock.NewResult(1, 1))
			Sqlmocker.ExpectExec(regexp.QuoteMeta("UPDATE `labels` SET `delete_at`=? WHERE annotation_template_id = ? AND `labels`.`delete_at` = ?")).
				WithArgs(sqlmock.AnyArg(), testId, 0).
				WillReturnError(fmt.Errorf("db error"))
			Sqlmocker.ExpectRollback()
			err := AnnotationTemplateRepo.Delete(testId)
			So(err, ShouldNotEqual, nil)
		})
	})

	Convey("Update annotation template", t, func() {
		// fake src data
		testAnnoTempId := uuid.New()
		srcAnnoTempRow := sqlmock.NewRows([]string{"id", "name", "view_type", "description", "create_at", "update_at", "delete_at"}).
			AddRow(testAnnoTempId.String(), "name0", "bounding-box", "", 1672887665782, 1672887665782, 0)
		wordList, _ := json.Marshal([]string{})
		srcAnnoTempExtRow := sqlmock.NewRows([]string{"annotation_template_id", "word_list", "delete_at"}).
			AddRow(testAnnoTempId.String(), wordList, 0)
		srcLabel1Id := uuid.New()
		srcLabel2Id := uuid.New()
		srcLabelRow := sqlmock.NewRows([]string{"id", "annotation_template_id", "name", "color"}).
			AddRow(srcLabel1Id.String(), testAnnoTempId.String(), "name1", 10).
			AddRow(srcLabel2Id.String(), testAnnoTempId.String(), "name2", 15) // will delete it

		// fake dest data

		destAnnoTemp := do.AnnotationTemplateDo{ID: testAnnoTempId, Name: "name00", Type: "bounding-box", Description: "", CreateAt: 1672887665782, UpdateAt: 1672887665782, DeleteAt: 0}
		destAnnoTempExt := do.AnnotationTemplateExtDo{AnnotationTemplateId: testAnnoTempId, WordList: []string{}}
		var destLabels []do.LabelDo
		destLabel3Id := uuid.New()
		// update it
		destLabels = append(destLabels, do.LabelDo{ID: srcLabel1Id, AnnotationTemplateId: testAnnoTempId, Name: "name11", Color: 10, SuperCategoryName: ""})
		// insert it
		destLabels = append(destLabels, do.LabelDo{ID: destLabel3Id, AnnotationTemplateId: testAnnoTempId, Name: "name33", Color: 20})

		Convey("Success", func() {
			Sqlmocker.ExpectBegin()
			Sqlmocker.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `annotation_templates` WHERE `annotation_templates`.`id` = ? AND `annotation_templates`.`delete_at` = ? ORDER BY `annotation_templates`.`id` LIMIT 1")).
				WithArgs(testAnnoTempId.String(), 0).WillReturnRows(srcAnnoTempRow)
			Sqlmocker.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `annotation_template_exts` WHERE annotation_template_id = ? AND `annotation_template_exts`.`delete_at` = ? ORDER BY `annotation_template_exts`.`annotation_template_id` LIMIT 1")).
				WithArgs(testAnnoTempId.String(), 0).WillReturnRows(srcAnnoTempExtRow)
			Sqlmocker.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `labels` WHERE `labels`.`annotation_template_id` = ? AND `labels`.`delete_at` = ?")).WithArgs(testAnnoTempId.String(), 0).WillReturnRows(srcLabelRow)

			Sqlmocker.ExpectExec(regexp.QuoteMeta("UPDATE `annotation_templates` SET `name`=?,`type`=?,`update_at`=? WHERE `annotation_templates`.`delete_at` = ? AND `id` = ?")).
				WithArgs(destAnnoTemp.Name, destAnnoTemp.Type, sqlmock.AnyArg(), 0, testAnnoTempId.String()).
				WillReturnResult(sqlmock.NewResult(1, 1))

			Sqlmocker.ExpectExec(regexp.QuoteMeta("DELETE FROM `labels` WHERE `labels`.`id` = ?")).
				WithArgs(srcLabel2Id.String()).
				WillReturnResult(sqlmock.NewResult(1, 1))

			Sqlmocker.ExpectExec(regexp.QuoteMeta("UPDATE `labels` SET `name`=?,`color`=?,`update_at`=? WHERE `labels`.`delete_at` = ? AND `id` = ?")).
				WithArgs(destLabels[0].Name, destLabels[0].Color, sqlmock.AnyArg(), 0, srcLabel1Id.String()).
				WillReturnResult(sqlmock.NewResult(1, 1))

			Sqlmocker.ExpectExec(regexp.QuoteMeta("INSERT INTO `labels` (`id`,`annotation_template_id`,`name`,`super_category_name`,`color`,`create_at`,`update_at`,`delete_at`,`key_point_def`,`key_point_skeleton`,`cover_image_url`) VALUES (?,?,?,?,?,?,?,?,?,?,?)")).
				WithArgs(sqlmock.AnyArg(), testAnnoTempId.String(), destLabels[1].Name, destLabels[1].SuperCategoryName, destLabels[1].Color, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
				WillReturnResult(sqlmock.NewResult(1, 1))

			Sqlmocker.ExpectCommit()
			err := AnnotationTemplateRepo.Update(destAnnoTemp, destAnnoTempExt, destLabels)
			So(err, ShouldEqual, nil)
		})

		Convey("Failed, insert label conflict", func() {
			Sqlmocker.ExpectBegin()
			Sqlmocker.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `annotation_templates` WHERE `annotation_templates`.`id` = ? AND `annotation_templates`.`delete_at` = ? ORDER BY `annotation_templates`.`id` LIMIT 1")).WithArgs(testAnnoTempId.String(), 0).WillReturnRows(srcAnnoTempRow)
			Sqlmocker.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `labels` WHERE `labels`.`annotation_template_id` = ? AND `labels`.`delete_at` = ?")).WithArgs(testAnnoTempId.String(), 0).WillReturnRows(srcLabelRow)

			Sqlmocker.ExpectExec(regexp.QuoteMeta("UPDATE `annotation_templates` SET `name`=?,`type`=?,`update_at`=? WHERE `annotation_templates`.`delete_at` = ? AND `id` = ?")).
				WithArgs(destAnnoTemp.Name, destAnnoTemp.Type, sqlmock.AnyArg(), 0, testAnnoTempId.String()).
				WillReturnResult(sqlmock.NewResult(1, 1))

			Sqlmocker.ExpectExec(regexp.QuoteMeta("DELETE FROM `labels` WHERE `labels`.`id` = ?")).
				WithArgs(srcLabel2Id.String()).
				WillReturnResult(sqlmock.NewResult(1, 1))

			Sqlmocker.ExpectExec(regexp.QuoteMeta("UPDATE `labels` SET `name`=?,`color`=?,`update_at`=? WHERE `labels`.`delete_at` = ? AND `id` = ?")).
				WithArgs(testAnnoTempId.String(), destLabels[0].Name, destLabels[0].Color, sqlmock.AnyArg(), 0, srcLabel1Id.String()).
				WillReturnResult(sqlmock.NewResult(1, 1))

			Sqlmocker.ExpectExec(regexp.QuoteMeta("INSERT INTO `labels` (`id`,`annotation_template_id`,`name`,`color`,`create_at`,`delete_at`) VALUES (?,?,?,?,?,?)")).
				WithArgs(destLabel3Id.String(), testAnnoTempId.String(), destLabels[1].Name, destLabels[1].Color, sqlmock.AnyArg(), 0).
				WillReturnError(fmt.Errorf("conflict"))
			Sqlmocker.ExpectRollback()
			err := AnnotationTemplateRepo.Update(destAnnoTemp, destAnnoTempExt, destLabels)
			So(err, ShouldNotEqual, nil)
		})

		Convey("Failed, not found", func() {
			Sqlmocker.ExpectBegin()
			Sqlmocker.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `annotation_templates` WHERE `annotation_templates`.`id` = ? AND `annotation_templates`.`delete_at` = ? ORDER BY `annotation_templates`.`id` LIMIT 1")).WithArgs(testAnnoTempId.String(), 0).WillReturnError(fmt.Errorf("not found"))
			Sqlmocker.ExpectRollback()
			err := AnnotationTemplateRepo.Update(destAnnoTemp, destAnnoTempExt, destLabels)
			So(err, ShouldNotEqual, nil)
		})
	})
}
