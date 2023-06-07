package server_test

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	server "github.com/thecoreman/problematic-api-server/server/go"
)

func TestErrorsApiService_ErrorsPercentGet(t *testing.T) {
	type args struct {
		ctx          context.Context
		errorPercent int32
	}
	tests := []struct {
		name     string
		s        *server.ErrorsApiService
		args     args
		wantCode int
		wantErr  bool
	}{
		{
			name: "error_percent is 0",
			args: args{
				ctx:          context.Background(),
				errorPercent: 0,
			},
			wantCode: http.StatusOK,
			wantErr:  false,
		},
		{
			name: "error_percent is 100",
			args: args{
				ctx:          context.Background(),
				errorPercent: 100,
			},
			wantCode: http.StatusInternalServerError,
			wantErr:  false,
		},
		{
			name: "error_percent is 101",
			args: args{
				ctx:          context.Background(),
				errorPercent: 101,
			},
			wantCode: http.StatusBadRequest,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := server.NewErrorsApiService()
			got, err := s.ErrorsPercentGet(tt.args.ctx, tt.args.errorPercent)
			if (err != nil) != tt.wantErr {
				t.Errorf("ErrorsApiService.ErrorsPercentGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Code, tt.wantCode) {
				t.Errorf("ErrorsApiService.ErrorsPercentGet() = %v, want %v", got.Code, tt.wantCode)
			}
		})
	}
}
