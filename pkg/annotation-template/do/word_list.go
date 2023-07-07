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

type WordListType []string

func (wordList WordListType) IsEmpty() bool {
	if wordList == nil || len([]string(wordList)) == 0 {
		return true
	}
	return false
}

func (wordList WordListType) Value() (driver.Value, error) {
	if wordList.IsEmpty() {
		return nil, nil
	}
	bytes, err := json.Marshal(wordList)
	if err != nil {
		return nil, err
	}
	return string(bytes), nil
}

func (wordList *WordListType) Scan(value interface{}) error {
	if value == nil {
		wordList = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal word list type value %v", value)
	}

	err := json.Unmarshal(bytes, &wordList)
	return err
}
