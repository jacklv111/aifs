/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package constant

// artifact 文件类型
const (
	ARTIFACT_FILE = "artifact-file"
)

func IsArtifactFile(dataType string) bool {
	return dataType == ARTIFACT_FILE
}
