package utils

import (
	"fmt"
	"time"
)

func TimeFormat(t time.Time) string {
	d := time.Since(t).Round(time.Second)
	day := int(d.Hours() / 24)
	hour := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	second := int(d.Seconds()) % 60

	if d >= 24*time.Hour {
		return fmt.Sprintf("%dd %dh %dm %ds", day, hour, minutes, second)
	} else if d >= time.Hour {
		return fmt.Sprintf("%dh %dm %ds", hour, minutes, second)
	} else if d >= time.Minute {
		return fmt.Sprintf("%dm %ds", minutes, second)
	} else {
		return fmt.Sprintf("%ds", second)
	}
}
