/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package constant

const (
	BATCH_SIZE = 1000
)

// 模型文件类型
const (
	MODEL_FILE = "model-file"
)

func IsModelFile(dataType string) bool {
	return dataType == MODEL_FILE
}
