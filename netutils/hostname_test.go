// Acceldata Inc. and its affiliates.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// 	Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
