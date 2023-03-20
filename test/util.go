package test

import (
	"bytes"
	"github.com/lunixbochs/struc"
	"testing"
)

func SizeOf(t *testing.T, i any) int {
	var buf bytes.Buffer
	if err := struc.Pack(&buf, i); err != nil {
		t.Fatal(err)
	}
	return buf.Len()
}

func SizeOfWithOpt(t *testing.T, i any, opt *struc.Options) int {
	var buf bytes.Buffer
	if err := struc.PackWithOptions(&buf, i, opt); err != nil {
		t.Fatal(err)
	}
	return buf.Len()
}
