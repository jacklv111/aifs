/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package do

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// key points 类型数据的定义

type KeyPointDefType []string

type KeyPointSkeletonType [][]int32

func (kptDef KeyPointDefType) IsEmpty() bool {
	if kptDef == nil || len([]string(kptDef)) == 0 {
		return true
	}
	return false
}

func (kptDef KeyPointDefType) Value() (driver.Value, error) {
	if kptDef.IsEmpty() {
		return nil, nil
	}
	bytes, err := json.Marshal(kptDef)
	if err != nil {
		return nil, err
	}
	return string(bytes), nil
}

func (kptDef *KeyPointDefType) Scan(value interface{}) error {
	if value == nil {
		kptDef = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal key point def type value %v", value)
	}

	err := json.Unmarshal(bytes, &kptDef)
	return err
}

func (kptSkeleton KeyPointSkeletonType) IsEmpty() bool {
	if kptSkeleton == nil || len([][]int32(kptSkeleton)) == 0 {
		return true
	}
	return false
}

func (kptSkeleton KeyPointSkeletonType) Value() (driver.Value, error) {
	if kptSkeleton.IsEmpty() {
		return nil, nil
	}
	bytes, err := json.Marshal(kptSkeleton)
	if err != nil {
		return nil, err
	}
	return string(bytes), nil
}

func (kptSkeleton *KeyPointSkeletonType) Scan(value interface{}) error {
	if value == nil {
		kptSkeleton = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal key point def type value %v", value)
	}

	err := json.Unmarshal(bytes, &kptSkeleton)
	return err
}
