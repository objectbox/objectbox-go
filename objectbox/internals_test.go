package objectbox

import (
	"runtime"
	"testing"
)

func TestLargeArraySupport(t *testing.T) {
	t.Log(internalLibVersion())

	if runtime.GOARCH == `386` || runtime.GOARCH == `arm` {
		if supportsResultArray {
			t.Errorf("Expected large array support to be disabled on a 32-bit system (%s) but its enabled "+
				"in the ObjectBox core library", runtime.GOARCH)
		}
	} else if !supportsResultArray {
		t.Errorf("Expected large array support to be enabled on a 64-bit system (%s) but its disabled "+
			"in the ObjectBox core library", runtime.GOARCH)
	}
}
