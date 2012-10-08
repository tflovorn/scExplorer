package tempAll

import (
	"errors"
	"os"
	"testing"
)

func TestSaveAndLoadCache(t *testing.T) {
	wd, _ := os.Getwd()
	cachePath := wd + "/deleteme.cache_test"
	data := make([]*Environment, 1)
	errs := make([]error, 1)
	env, err := envDefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	data[0] = env
	errs[0] = errors.New("cache test error")
	err = SaveEnvCache(cachePath, data, errs)
	if err != nil {
		t.Fatal(err)
	}
	loadedData, loadedErrs, err := LoadEnvCache(cachePath)
	if err != nil {
		t.Fatal(err)
	}
	if loadedErrs[0].Error() != "cache test error" {
		t.Fatalf("incorrect error loaded")
	}
	loadedEnv := loadedData[0].(*Environment)
	if loadedEnv.X != env.X || loadedEnv.Alpha != env.Alpha {
		t.Fatalf("incorrect Environment loaded")
	}
}
