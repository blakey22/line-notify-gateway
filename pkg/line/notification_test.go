package line

import (
	"testing"
)

func TestNotification_URLValues(t *testing.T) {
	tests := []struct {
		name         string
		notification Notification
		want         string
	}{
		{
			name: "only message",
			notification: Notification{
				Message: "my_message",
			},
			want: "message=my_message&notificationDisabled=false",
		},
		{
			name: "message and sticker",
			notification: Notification{
				Message: "my_message",
				StickerPackageID: 1,
				StickerID: 1,
			},
			want: "message=my_message&notificationDisabled=false&stickerId=1&stickerPackageId=1",
		},
		{
			name: "message and invalid sticker package",
			notification: Notification{
				Message: "my_message",
				StickerPackageID: 0,
				StickerID: 1,
			},
			want: "message=my_message&notificationDisabled=false",
		},
		{
			name: "message and invalid sticker",
			notification: Notification{
				Message: "my_message",
				StickerPackageID: 1,
				StickerID: 0,
			},
			want: "message=my_message&notificationDisabled=false",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := tt.notification.URLValues(); got != tt.want {
				t.Errorf("Notification.URLValues() = %v, want %v", got, tt.want)
			}
		})
	}
}
