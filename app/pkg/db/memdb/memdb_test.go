package memdb

import (
	"reflect"
	"sync"
	"testing"
)

var db *DB

func TestMain(m *testing.M) {
	db = NewDB(10)

	m.Run()
}

func TestNewURLStorage(t *testing.T) {
	type args struct {
		urlsize uint
	}
	tests := []struct {
		name string
		args args
		want *URLStorage
	}{
		{
			name: "Test 1",
			args: args{urlsize: 10},
			want: &URLStorage{urlsize: 10, longToShort: map[string]string{}, shortToLong: map[string]string{}, mu: &sync.Mutex{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewURLStorage(tt.args.urlsize); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewURLStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestURLStorage_AddURL(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		s    *URLStorage
		args args
	}{
		{
			name: "Test 1",
			s:    NewURLStorage(10),
			args: args{"ozon.ru"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.AddURL(tt.args.url); got != tt.s.longToShort[tt.args.url] {
				t.Errorf("URLStorage.AddURL() = %v, want %v", got, tt.s.longToShort[tt.args.url])
			}
		})
	}
}

func TestDB_MakeShort(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		db      *DB
		args    args
		wantErr bool
	}{
		{
			name:    "Add URL",
			db:      db,
			args:    args{url: "http://example.com"},
			wantErr: false,
		},
		{
			name:    "Exising URL",
			db:      db,
			args:    args{url: "http://example.com"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.db.MakeShort(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("DB.MakeShort() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != db.storage.longToShort[tt.args.url] {
				t.Errorf("DB.MakeShort() = %v, want %v", got, db.storage.longToShort[tt.args.url])
			}
		})
	}
}
