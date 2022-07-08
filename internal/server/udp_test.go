package server

import "testing"

func TestUDPServer_NewUDP(t *testing.T) {
	tests := []struct {
		name string
		s    *UDPServer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.NewUDP()
		})
	}
}
