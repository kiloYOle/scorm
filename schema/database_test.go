package schema

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// os.Exit skips defer calls
	// so we need to call another function
	code, err := RunTestMain(m)
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(code)
}
