package upyun

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"sync"
	"testing"
)

var (
	ROOT       = MakeTmpPath()
	NOTIFY_URL = "http://124.160.114.202:18989/echo?key=gosdk"
)

var up = NewUpYun(&UpYunConfig{
	Bucket:   "prog-test",
	Operator: "myworker",
	Password: "tyghbnTYGHBN",
})

func MakeTmpPath() string {
	return "/go-sdk/123456789"
}

func Equal(t *testing.T, actual, expected interface{}) {
	if !reflect.DeepEqual(actual, expected) {
		_, file, line, _ := runtime.Caller(1)
		t.Logf("\033[31m%s:%d:\n\n\tnexp: %#v\n\n\tgot:  %#v\033[39m\n\n",
			filepath.Base(file), line, expected, actual)
		t.FailNow()
	}
}

func NotEqual(t *testing.T, actual, expected interface{}) {
	if reflect.DeepEqual(actual, expected) {
		_, file, line, _ := runtime.Caller(1)
		t.Logf("\033[31m%s:%d:\n\n\tnexp: %#v\n\n\tgot:  %#v\033[39m\n\n",
			filepath.Base(file), line, expected, actual)
		t.FailNow()
	}
}

func Nil(t *testing.T, object interface{}) {
	if !isNil(object) {
		_, file, line, _ := runtime.Caller(1)
		t.Logf("\033[31m%s:%d:\n\n\t   <nil> (expected)\n\n\t!= %#v (actual)\033[39m\n\n",
			filepath.Base(file), line, object)
		t.FailNow()
	}
}

func NotNil(t *testing.T, object interface{}) {
	if isNil(object) {
		_, file, line, _ := runtime.Caller(1)
		t.Logf("\033[31m%s:%d:\n\n\tExpected value not to be <nil>\033[39m\n\n",
			filepath.Base(file), line, object)
		t.FailNow()
	}
}

func isNil(object interface{}) bool {
	if object == nil {
		return true
	}

	value := reflect.ValueOf(object)
	kind := value.Kind()
	if kind >= reflect.Chan && kind <= reflect.Slice && value.IsNil() {
		return true
	}

	return false

}

func TestMain(m *testing.M) {
	clean := func() {
		objs := make(chan *FileInfo, 20)
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			for obj := range objs {
				up.Delete(&DeleteObjectConfig{
					Path: path.Join(ROOT, obj.Name),
				})
			}
			up.Delete(&DeleteObjectConfig{
				Path: ROOT,
			})
			wg.Done()
		}()

		up.List(&GetObjectsConfig{
			Path:         ROOT,
			ObjectsChan:  objs,
			MaxListLevel: -1,
		})
		wg.Wait()

		if _, err := up.GetInfo(ROOT); err == nil {
			fmt.Println("Not cleanup")
			os.Exit(-1)
		}
	}

	flag.Parse()
	code := m.Run()

	clean()
	os.Exit(code)
}
