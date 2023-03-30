package websocket

import "testing"

func Test_simpleChat(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{name: "chat"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			simpleChat()
		})
	}
}

func Test_fileMonitor(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{name: "file"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileMonitor()
		})
	}
}

func Test_image(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{name: "image"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			image()
		})
	}
}

func Test_chatRoom(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{name: "chatroom"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chatRoom()
		})
	}
}
