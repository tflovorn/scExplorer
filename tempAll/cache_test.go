package tempAll

import (
	"errors"
	"os"
	"testing"
)

func TestSaveAndLoadCache(t *testing.T) {
	wd, _ := os.Getwd()
	cachePath := wd + "/deleteme.cache_test"
	data := make([]interface{}, 2)
	errs := make([]error, 2)
	env, err := envDefaultEnv()
	if err != nil {
		t.Fatal(err)
	}
	data[0] = env
	data[1] = env
	errs[0] = errors.New("cache test error")
	errs[1] = nil
	err = SaveEnvCache(cachePath, data, errs)
	if err != nil {
		t.Fatal(err)
	}
	loadedData, loadedErrs, err := LoadEnvCache(cachePath)
	if err != nil {
		t.Fatal(err)
	}
	if loadedErrs[0].Error() != "cache test error" || loadedErrs[1] != nil {
		t.Fatalf("incorrect error loaded")
	}
	loadedEnv := loadedData[0].(*Environment)
	if loadedEnv.X != env.X || loadedEnv.Alpha != env.Alpha {
		t.Fatalf("incorrect Environment loaded")
	}
}
