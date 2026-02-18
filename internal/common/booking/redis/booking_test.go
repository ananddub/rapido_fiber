package booking_common

import (
	"fmt"
	"testing"
)

type Temp struct {
	x int
	y int
}

func split(sum int) (w *Temp) {
	w = &Temp{
		x: sum / 2,
		y: sum / 2,
	}
	return
}

func TestMain(t *testing.T) {
	fmt.Println(split(17))
}
