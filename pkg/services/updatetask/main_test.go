package updatetask

import (
	"fmt"
	"os"
	"testing"
)

const minimumCodeCoverage = 1

func TestMain(m *testing.M) {
	rc := m.Run()
	if rc != 0 || testing.CoverMode() == "" {
		return
	}
	if c := testing.Coverage(); c < minimumCodeCoverage {
		fmt.Printf("Tests passed but coverage %v is less than minimum code coverage %v", c, minimumCodeCoverage)
		rc = -1
	}
	os.Exit(rc)
}
