/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package repo

import (
	"github.com/jacklv111/aifs/pkg/data-view/data-module/basic/repo"
)

//go:generate mockgen -source=interface.go -destination=./mock/mock_interface.go -package=mock

type datasetZipRepoInterface interface {
	repo.BasicDataRepoInterface
}
