package api

import (
	"testing"
)

var (
	validAMQP  = "amqp://guest:guest@localhost/vhost"
	validRedis = "redis://localhost:6379/0"
)

func TestNewAPI(t *testing.T) {
	type args struct {
		ampqURI  string
		redisURI string
		host     string
		port     int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"valid connection", args{validAMQP, validRedis, "0.0.0.0", 80}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewAPI(tt.args.ampqURI, tt.args.redisURI, tt.args.host, tt.args.port)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAPI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
