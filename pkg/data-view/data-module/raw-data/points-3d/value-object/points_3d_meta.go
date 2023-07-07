/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package valueobject

import (
	"fmt"
	"strconv"
	"strings"
)

type Points3DMeta struct {
	FileName string
	Sha256   string
	// points number
	Size  int64
	Xmin  float64
	Xmax  float64
	Ymin  float64
	Ymax  float64
	Zmin  float64
	Zmax  float64
	Rmean float64
	Gmean float64
	Bmean float64
	Rstd  float64
	Gstd  float64
	Bstd  float64
}

func (meta *Points3DMeta) ParseFromString(str string) (err error) {
	items := strings.Split(str, " ")
	meta.FileName = items[0]
	meta.Sha256 = items[1]
	meta.Size, err = strconv.ParseInt(items[2], 10, 64)
	if err != nil {
		return err
	}
	meta.Xmin, err = strconv.ParseFloat(items[3], 64)
	if err != nil {
		return err
	}
	meta.Xmax, err = strconv.ParseFloat(items[4], 64)
	if err != nil {
		return err
	}
	meta.Ymin, err = strconv.ParseFloat(items[5], 64)
	if err != nil {
		return err
	}
	meta.Ymax, err = strconv.ParseFloat(items[6], 64)
	if err != nil {
		return err
	}
	meta.Zmin, err = strconv.ParseFloat(items[7], 64)
	if err != nil {
		return err
	}
	meta.Zmax, err = strconv.ParseFloat(items[8], 64)
	if err != nil {
		return err
	}

	meta.Rmean, err = strconv.ParseFloat(items[9], 64)
	if err != nil {
		return err
	}
	meta.Gmean, err = strconv.ParseFloat(items[10], 64)
	if err != nil {
		return err
	}
	meta.Bmean, err = strconv.ParseFloat(items[11], 64)
	if err != nil {
		return err
	}
	meta.Rstd, err = strconv.ParseFloat(items[12], 64)
	if err != nil {
		return err
	}
	meta.Gstd, err = strconv.ParseFloat(items[13], 64)
	if err != nil {
		return err
	}
	meta.Bstd, err = strconv.ParseFloat(items[14], 64)
	if err != nil {
		return err
	}

	return
}

func (meta *Points3DMeta) String() string {
	return fmt.Sprintf("%s %s %d "+
		"%f %f "+
		"%f %f "+
		"%f %f "+
		"%f %f %f "+
		"%f %f %f",
		meta.FileName, meta.Sha256, meta.Size,
		meta.Xmin, meta.Xmax,
		meta.Ymin, meta.Ymax,
		meta.Zmin, meta.Zmax,
		meta.Rmean, meta.Gmean, meta.Bmean,
		meta.Rstd, meta.Gstd, meta.Bstd,
	)
}
