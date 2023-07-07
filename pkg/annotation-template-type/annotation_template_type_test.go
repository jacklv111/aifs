/*
 * Created on Fri Jul 07 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package annotationtemplatetype

import (
	"reflect"
	"testing"
)

func TestGetAnnotationTemplateList(t *testing.T) {
	type args struct {
		offset int
		limit  int
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "Success, want 1 get 1", args: args{0, 1}, want: []string{CATEGORY}},
		{name: "Success, want 100 get 1", args: args{0, 100}, want: []string{CATEGORY, COCO_TYPE, OCR, POINTS_3D, RGBD, SEGMENTATION_MASKS}},
		{name: "Success, get nothing", args: args{100, 100}, want: []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAnnotationTemplateList(tt.args.offset, tt.args.limit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAnnotationTemplateList() = %v, want %v", got, tt.want)
			}
		})
	}
}
