package duocli

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	duoapi "github.com/duosecurity/duo_api_golang"
)

// compareDuoAPI compares the fields of duoapi.DuoApi ignoring the various
// option fields.
func compareDuoAPI(t1 *duoapi.DuoApi, t2 *duoapi.DuoApi) bool {
	val1 := reflect.ValueOf(t1).Elem()
	val2 := reflect.ValueOf(t2).Elem()

	if val1.FieldByName("ikey").String() != val2.FieldByName("ikey").String() {
		return false
	}
	if val1.FieldByName("skey").String() != val2.FieldByName("skey").String() {
		return false
	}
	if val1.FieldByName("host").String() != val2.FieldByName("host").String() {
		return false
	}
	if val1.FieldByName("userAgent").String() != val2.FieldByName("userAgent").String() {
		return false
	}
	return true
}
func Test_loadConfig(t *testing.T) {

	buf1 := bytes.NewBuffer([]byte(`{"ikey": "DIXXXXXXXXXXXXXXXXXX","skey": "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX","api_host": "api-XXXXXXXX.duosecurity.com"}`))
	buf2 := bytes.NewBuffer([]byte(`{"skey": "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX","api_host": "api-XXXXXXXX.duosecurity.com"}`))
	buf3 := bytes.NewBuffer([]byte(`{"ikey": "DIXXXXXXXXXXXXXXXXXX","api_host": "api-XXXXXXXX.duosecurity.com"}`))
	buf4 := bytes.NewBuffer([]byte(`{"ikey": "DIXXXXXXXXXXXXXXXXXX","skey": "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"}`))

	type args struct {
		file io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    *duoapi.DuoApi
		wantErr bool
	}{
		{
			name: "Min config",
			args: args{file: buf1},
			want: duoapi.NewDuoApi("DIXXXXXXXXXXXXXXXXXX", "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX", "api-XXXXXXXX.duosecurity.com", "duocli"),
		},
		{
			name:    "Missing ikey config",
			args:    args{file: buf2},
			wantErr: true,
		},
		{
			name:    "Missing skey config",
			args:    args{file: buf3},
			wantErr: true,
		},
		{
			name:    "Missing host config",
			args:    args{file: buf4},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := loadConfig(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadConfig() error = %#v, wantErr %#v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !compareDuoAPI(got, tt.want) {
				t.Errorf("loadConfig() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
