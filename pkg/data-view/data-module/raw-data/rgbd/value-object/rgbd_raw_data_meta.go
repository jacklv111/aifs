/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package valueobject

import (
	"strconv"
	"strings"
)

type RgbdRawDataMeta struct {
	FileName string
	Sha256   string
	// bytes
	ImageSize   int64
	ImageWidth  int64
	ImageHeight int64
	// bytes
	DepthSize   int64
	DepthWidth  int64
	DepthHeight int64
}

func (meta *RgbdRawDataMeta) ParseFromString(str string) (err error) {
	items := strings.Split(str, " ")
	meta.FileName = items[0]
	meta.ImageSize, err = strconv.ParseInt(items[1], 10, 64)
	if err != nil {
		return
	}
	meta.ImageWidth, err = strconv.ParseInt(items[2], 10, 32)
	if err != nil {
		return
	}
	meta.ImageHeight, err = strconv.ParseInt(items[3], 10, 32)
	if err != nil {
		return
	}
	meta.DepthSize, err = strconv.ParseInt(items[4], 10, 64)
	if err != nil {
		return
	}
	meta.DepthWidth, err = strconv.ParseInt(items[5], 10, 32)
	if err != nil {
		return
	}
	meta.DepthHeight, err = strconv.ParseInt(items[6], 10, 32)
	if err != nil {
		return
	}
	meta.Sha256 = items[7]
	return
}

func (meta RgbdRawDataMeta) String() string {
	return strings.Join([]string{
		meta.FileName,
		strconv.FormatInt(meta.ImageSize, 10),
		strconv.FormatInt(meta.ImageWidth, 10),
		strconv.FormatInt(meta.ImageHeight, 10),
		strconv.FormatInt(meta.DepthSize, 10),
		strconv.FormatInt(meta.DepthWidth, 10),
		strconv.FormatInt(meta.DepthHeight, 10),
		meta.Sha256,
	}, " ")
}
