/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package repo

import "github.com/jacklv111/aifs/pkg/data-view/data-module/basic/repo"

type datasetZipRepoImpl struct {
	repo.BasicDataRepoImpl
}

var DatasetZipRepo datasetZipRepoInterface

func init() {
	DatasetZipRepo = &datasetZipRepoImpl{}
}
