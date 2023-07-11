/*
 * Created on Mon Jul 10 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package job

import (
	"context"
	"os"

	aifsclientgo "github.com/jacklv111/aifs-client-go"
	"github.com/jacklv111/aifs/pkg/job/dataset-zip/handler"
	"github.com/jacklv111/common-sdk/log"
	"github.com/jacklv111/common-sdk/utils"
)

func prepareData(handler handler.Handler, dataViewId string, workDir string) error {
	filePath, err := utils.DownloadZipCmd(dataViewId, workDir)
	if err != nil {
		return err
	}

	if err := handler.UpdateProgress(10.0, "download zip completed"); err != nil {
		return err
	}

	if err := utils.Decompress(filePath, workDir, false); err != nil {
		return err
	}

	if err := handler.UpdateProgress(20.0, "unzip completed"); err != nil {
		return err
	}
	err = os.Remove(filePath)

	if err != nil {
		return err
	}
	return nil
}

func UnzipDatasetZipView(client *aifsclientgo.APIClient, dataViewId string, workDir string) error {
	log.Infof("begin to decompression dataset zip view %s, work dir %s", dataViewId, workDir)
	details, _, err := client.DataViewApi.GetDataViewDetails(context.Background(), dataViewId).Execute()
	if err != nil {
		return err
	}

	// is completed
	if details.Progress != nil && *details.Progress+1e-8 >= 100.0 {
		return nil
	}

	handler, err := handler.BuildHandler(client, details, workDir)
	if err != nil {
		return err
	}

	err = prepareData(handler, dataViewId, workDir)
	if err != nil {
		return err
	}

	return handler.Exec()
}
