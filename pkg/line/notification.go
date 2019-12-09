package line

import (
	"net/url"
	"strconv"
)

type Notification struct {
	Message              string
	StickerPackageID     int
	StickerID            int
	NotificationDisabled bool
}

func (n Notification) URLValues() string {
	values := url.Values{}

	values.Set("message", n.Message)
	values.Set("notificationDisabled", strconv.FormatBool(n.NotificationDisabled))

	if n.StickerPackageID > 0 && n.StickerID > 0 {
		values.Set("stickerPackageId", strconv.Itoa(n.StickerPackageID))
		values.Set("stickerId", strconv.Itoa(n.StickerID))
	}

	return values.Encode()
}
