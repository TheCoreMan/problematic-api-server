package server_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/go-test/deep"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	server "github.com/thecoreman/problematic-api-server/server/go"
)

func TestCachingApiService_CacheableGet(t *testing.T) {
	viper.SetDefault("BOOKS_DIRECTORY", "../../books")
	type args struct {
		ctx         context.Context
		bookTitle   string
		lineNumber  int32
		withControl bool
	}
	tests := []struct {
		name    string
		args    args
		want    *server.ImplResponse
		wantErr bool
	}{
		{
			name: "Traversal attack attempt",
			args: args{
				ctx:         context.Background(),
				bookTitle:   "../server/go/api_caching_service.go",
				lineNumber:  1,
				withControl: false,
			},
			want: &server.ImplResponse{
				Code: 500,
			},
			wantErr: false,
		},
		{
			name: "Get a book and line that exist",
			args: args{
				ctx:         context.Background(),
				bookTitle:   "war-and-peace.txt",
				lineNumber:  43176,
				withControl: false,
			},
			want: &server.ImplResponse{
				Code: 200,
				Headers: map[string][]string{
					"": {},
				},
				Body: server.SuccessfulResponse{
					BookName:   "war-and-peace.txt",
					LineNumber: 43176,
					Text:       `"Vive l'Empereur! Vive le roi de Rome! Vive l'Empereur!" came`,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := server.NewCachingApiService(zerolog.Nop())
			got, err := s.CacheableGet(tt.args.ctx, tt.args.bookTitle, tt.args.lineNumber, tt.args.withControl)
			if (err != nil) != tt.wantErr {
				t.Errorf("CachingApiService.CacheableGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Code, tt.want.Code) {
				t.Errorf("CachingApiService.CacheableGet() = %v, want %v", got, tt.want)
			}
			if got.Code == 200 {
				// use deep.Equal to compare the body
				if diff := deep.Equal(got.Body, tt.want.Body); diff != nil {
					t.Error(diff)
				}
			}
		})
	}
}
