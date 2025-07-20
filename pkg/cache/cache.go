package cache

import "time"

type Interface interface {
	GetAny(k string, s any) error
	SetAny(k string, s any, exp time.Duration) error
	Delete(k string) error
}
