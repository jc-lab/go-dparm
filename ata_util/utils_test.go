package ata_util

import (
	"github.com/jc-lab/go-dparm/ata"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsDma(t *testing.T) {
	type args struct {
		op ata.OpCode
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, IsDma(tt.args.op), "IsDma(%v)", tt.args.op)
		})
	}
}

func TestIsNeedsLba48(t *testing.T) {
	type args struct {
		op    ata.OpCode
		lba   uint64
		nsect uint
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, IsNeedsLba48(tt.args.op, tt.args.lba, tt.args.nsect), "IsNeedsLba48(%v, %v, %v)", tt.args.op, tt.args.lba, tt.args.nsect)
		})
	}
}

func TestTfInit(t *testing.T) {
	type args struct {
		tf    ata.Tf
		op    ata.OpCode
		lba   uint64
		nsect uint
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			TfInit(tt.args.tf, tt.args.op, tt.args.lba, tt.args.nsect)
		})
	}
}
