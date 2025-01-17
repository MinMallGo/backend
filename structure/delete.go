package structure

import "time"

type Delete struct {
	Status     int32
	DeleteTime time.Time
}
