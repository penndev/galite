package lib_test

import (
	"fmt"
	"time"

	"github.com/penndev/galite/internal/lib"
)

func ExampleInitTTLMap() {
	ttlmap, _ := lib.InitTTLMap("")
	ttlmap.SetAny("p", "penndev", 5*time.Second)
	var s string
	ttlmap.GetAny("p", &s)
	fmt.Println(s)
	// output:
	// penndev
}
