package main

import (
	"testing"
)

//func Test_checkFormat(t *testing.T) {
//	type args struct {
//		str *string
//	}
//	tests := []struct {
//		name    string
//		args    args
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if err := checkFormat(tt.args.str); (err != nil) != tt.wantErr {
//				t.Errorf("checkFormat() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}

func BenchmarkFixedArray(b *testing.B) {
	const size = 1000000000
	var fixedArray = [size]bool{}
	println(len(fixedArray))
	println(cap(fixedArray))

	for i := 0; i < b.N; i++ {
		fixedArray[i] = true
	}
	println(cap(fixedArray))
}
