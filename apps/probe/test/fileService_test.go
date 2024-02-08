package test

import (
	"apps/common/message/data"
	fileservice "apps/probe/service/fs"
	"reflect"
	"testing"
)

func TestFileModify(t *testing.T) {
	type args struct {
		fm data.FileModify
	}
	tests := []struct {
		name    string
		args    args
		want    data.FileModifyResponse
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "FileModify",
			args: args{
				fm: data.FileModify{
					Path:    "/Users/meichuankutou/Public/testrewrite.txt",
					Changes: []data.Change{},
				},
			},
			want:    data.FileModifyResponse{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fileservice.FileModify(tt.args.fm)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileModify() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FileModify() = %v, want %v", got, tt.want)
			}
		})
	}
}
