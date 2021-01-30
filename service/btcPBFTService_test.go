package service

import (
	"reflect"
	"testing"
)

func Test_newBTCPBFTService(t *testing.T) {
	type args struct {
		ps *PBFTService
	}
	tests := []struct {
		name string
		args args
		want *btcPBFTService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newBTCPBFTService(tt.args.ps); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newBTCPBFTService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_btcPBFTService_UpdateSeq(t *testing.T) {
	type fields struct {
		PBFTService *PBFTService
	}
	type args struct {
		seq       int
		op        string
		auxiliary string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := &btcPBFTService{
				PBFTService: tt.fields.PBFTService,
			}
			ps.UpdateSeq(tt.args.seq, tt.args.op, tt.args.auxiliary)
		})
	}
}

func Test_btcPBFTService_CanOpSuccess(t *testing.T) {
	type fields struct {
		PBFTService *PBFTService
	}
	type args struct {
		op   string
		view int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := &btcPBFTService{
				PBFTService: tt.fields.PBFTService,
			}
			if err := ps.CanOpSuccess(tt.args.op, tt.args.view); (err != nil) != tt.wantErr {
				t.Errorf("btcPBFTService.CanOpSuccess() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_btcPBFTService_GetOpAuxiliary(t *testing.T) {
	type fields struct {
		PBFTService *PBFTService
	}
	type args struct {
		op   string
		view int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := &btcPBFTService{
				PBFTService: tt.fields.PBFTService,
			}
			got, err := ps.GetOpAuxiliary(tt.args.op, tt.args.view)
			if (err != nil) != tt.wantErr {
				t.Errorf("btcPBFTService.GetOpAuxiliary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("btcPBFTService.GetOpAuxiliary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_btcPBFTService_PrepareSeq(t *testing.T) {
	type fields struct {
		PBFTService *PBFTService
	}
	type args struct {
		view      int
		seq       int
		op        string
		auxiliary string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := &btcPBFTService{
				PBFTService: tt.fields.PBFTService,
			}
			if err := ps.PrepareSeq(tt.args.view, tt.args.seq, tt.args.op, tt.args.auxiliary); (err != nil) != tt.wantErr {
				t.Errorf("btcPBFTService.PrepareSeq() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_btcPBFTService_CommitSeq(t *testing.T) {
	type fields struct {
		PBFTService *PBFTService
	}
	type args struct {
		view      int
		seq       int
		op        string
		auxiliary string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := &btcPBFTService{
				PBFTService: tt.fields.PBFTService,
			}
			if err := ps.CommitSeq(tt.args.view, tt.args.seq, tt.args.op, tt.args.auxiliary); (err != nil) != tt.wantErr {
				t.Errorf("btcPBFTService.CommitSeq() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_btcPBFTService_newUTXO(t *testing.T) {
	type fields struct {
		PBFTService *PBFTService
	}
	type args struct {
		op string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantUtxos string
		wantErr   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := &btcPBFTService{
				PBFTService: tt.fields.PBFTService,
			}
			gotUtxos, err := ps.newUTXO(tt.args.op)
			if (err != nil) != tt.wantErr {
				t.Errorf("btcPBFTService.newUTXO() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUtxos != tt.wantUtxos {
				t.Errorf("btcPBFTService.newUTXO() = %v, want %v", gotUtxos, tt.wantUtxos)
			}
		})
	}
}
