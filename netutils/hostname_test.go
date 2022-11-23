package netutils

import (
	"os"
	"testing"
)

func TestGetHostName(t *testing.T) {
	type args struct {
		hostNameCommand string
		cmdTimeout      int
	}
	_ = os.Setenv("PULSE_HOSTNAME", "hdp1000.trl.iti.acceldata.dev")

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Env test",
			args: args{
				hostNameCommand: "ENV",
				cmdTimeout:      0,
			},
			want:    "hdp1000.trl.iti.acceldata.dev",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetHostName(tt.args.hostNameCommand, tt.args.cmdTimeout)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetHostName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetHostName() got = %v, want %v", got, tt.want)
			}
		})
	}
}
