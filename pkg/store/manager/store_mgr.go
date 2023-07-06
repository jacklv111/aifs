/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package manager

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/google/uuid"
	"github.com/jacklv111/aifs/pkg/store/constant"
	"github.com/jacklv111/aifs/pkg/store/do"
	"github.com/jacklv111/aifs/pkg/store/repo"
	vb "github.com/jacklv111/aifs/pkg/store/value-object"
	"github.com/jacklv111/common-sdk/log"
	"github.com/jacklv111/common-sdk/s3"
)

type StoreMgrImpl struct {
}

var StoreMgr StoreMgrInterface

func init() {
	StoreMgr = &StoreMgrImpl{}
}

func (mgr *StoreMgrImpl) Upload(params vb.UploadParams) error {
	log.Infof("start to upload data, data type: %s, len data items %d", params.DataType, len(params.DataItemList))

	if params.IsEmpty() {
		return nil
	}

	var uKeyList []do.LocationUkey
	for _, data := range params.DataItemList {
		uKeyList = append(uKeyList, do.LocationUkey{DataItemId: data.DataItemId, Name: data.Name, Environment: constant.S3})
	}

	existedSet, err := repo.LocationRepo.FindExistedUkey(uKeyList)
	if err != nil {
		return err
	}

	var locationDoList []do.LocationDo
	pathRoot := uuid.New().String() + time.Now().Format("2006-01-02-15:04:05")
	var readerMappers []s3.ReaderMapper

	for _, dataItem := range params.DataItemList {
		// 如果已经存在，则不再上传
		if existedSet.Contains(do.LocationUkey{DataItemId: dataItem.DataItemId, Name: dataItem.Name, Environment: constant.S3}) {
			continue
		}
		existedSet.Add(do.LocationUkey{DataItemId: dataItem.DataItemId, Name: dataItem.Name, Environment: constant.S3})

		locationDo := do.LocationDo{
			Name:        dataItem.Name,
			DataItemId:  dataItem.DataItemId,
			BucketName:  s3.S3Config.Bucket,
			ObjectKey:   filepath.Join(params.DataType, pathRoot, dataItem.GetUniqueName()),
			Environment: constant.S3,
		}
		locationDoList = append(locationDoList, locationDo)
		readerMappers = append(readerMappers, s3.ReaderMapper{Key: locationDo.ObjectKey, Reader: dataItem.Reader})
	}

	batchIterator := s3.NewBatchUploadIterator(s3.S3Config.Bucket, readerMappers)
	err = s3.S3Uploader.UploadWithIterator(aws.BackgroundContext(), batchIterator)
	if err != nil {
		return err
	}

	err = repo.LocationRepo.Create(locationDoList)

	if err != nil {
		log.Info("create location error %s, try to rollback", err)
		batchDeleteIter := s3.NewBatchDeleteIterator(s3.S3Config.Bucket, s3.GetReaderMapperKeyList(readerMappers))
		s3err := s3.S3Deleter.Delete(aws.BackgroundContext(), batchDeleteIter)
		if s3err != nil {
			err = fmt.Errorf("delete data because create location err %s, but delete data err %s", err, s3err)
			return err
		}
		return err
	}

	log.Info("upload success")
	return nil
}

func (mgr *StoreMgrImpl) Download(params vb.DownloadParams) error {
	if params.IsEmpty() {
		return nil
	}

	locationMap, err := repo.LocationRepo.FindByIdList(params.GetIdList())
	if err != nil {
		return err
	}

	var writerMappers []s3.WriterMapper

	for _, dataItem := range params.DataItemList {
		ukey := do.LocationUkey{DataItemId: dataItem.DataItemId, Name: dataItem.Name, Environment: constant.S3}
		locationDo, ok := locationMap[ukey]
		if !ok {
			continue
		}
		writerMappers = append(writerMappers, s3.WriterMapper{Writer: dataItem.WriterAt, Key: locationDo.ObjectKey})
	}

	batchIterator := s3.NewBatchDownloadIterator(s3.S3Config.Bucket, writerMappers)
	err = s3.S3Downloader.DownloadWithIterator(aws.BackgroundContext(), batchIterator)
	if err != nil {
		return err
	}

	return nil
}

func (mgr *StoreMgrImpl) Delete(params vb.DeleteParams) error {
	if params.IsEmpty() {
		return nil
	}

	locationMap, err := repo.LocationRepo.FindByIdList(params.GetIdList())
	if err != nil {
		return err
	}

	s3DeleteKeyList := make([]string, 0)
	locationList := make([]do.LocationUkey, 0)
	for uk, data := range locationMap {
		locationList = append(locationList, uk)

		switch uk.Environment {
		case constant.S3:
			s3DeleteKeyList = append(s3DeleteKeyList, data.ObjectKey)
		}
	}

	batchIterator := s3.NewBatchDeleteIterator(s3.S3Config.Bucket, s3DeleteKeyList)
	err = s3.S3Deleter.Delete(aws.BackgroundContext(), batchIterator)
	if err != nil {
		return err
	}

	err = repo.LocationRepo.DeleteByUk(locationList)
	if err != nil {
		return err
	}

	return nil
}

func (mgr *StoreMgrImpl) GetByIdList(idList []uuid.UUID) (map[uuid.UUID]vb.LocationResult, error) {
	locationMap, err := repo.LocationRepo.FindByIdList(idList)
	if err != nil {
		return nil, err
	}
	res := make(map[uuid.UUID]vb.LocationResult, 0)
	for key, data := range locationMap {
		res[key.DataItemId] = vb.LocationResult{ID: key.DataItemId, ObjectKey: data.ObjectKey}
	}
	return res, nil
}

func (mgr *StoreMgrImpl) GetUrlListByIdList(idList []uuid.UUID) (map[uuid.UUID]string, error) {
	locationMap, err := repo.LocationRepo.FindByIdList(idList)
	if err != nil {
		return nil, err
	}
	res := make(map[uuid.UUID]string, 0)
	for key, data := range locationMap {
		req, _ := s3.GetObjectRequest(data.ObjectKey)
		url, err := req.Presign(5 * time.Minute)
		if err != nil {
			return nil, err
		}
		res[key.DataItemId] = url
	}
	return res, nil
}
