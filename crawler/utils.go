package crawler

import (
	"time"
)

type Callback func()

func Delay(s time.Duration, f Callback) {
	time.Sleep(s * time.Second)
	f()

}

func Repeat(s time.Duration, f Callback) {
	for {
		f()
		time.Sleep(s * time.Second)
	}
}
