package util

import (
	"os"
	"testing"
	"text/template"
)

type UnitTest struct{
	Test1 string `ini:"test1" default:"test string"`
	Test2 int    `ini:"test2" default:"10"`
	Test3 int64  `ini:"test3" default:"13"`
}

func TestLoadConfig(t *testing.T) {
	file := "testConfig.ini"
	str1 := "Hello Unit Test!"
	str2 := 20
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		t.Fatal(err)
	}
	var to = &struct{
		Str1 string
		Str2 int
	}{
		str1,
		str2,
	}
	err = template.Must(template.New("").Parse(testScript)).Execute(f, to)
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	u := &UnitTest{}
	err = LoadConfig(file, u)
	if err != nil {
		t.Fatal(err)
	}

	if u.Test1 != "Hello Unit Test!" {
		t.Fatalf("u.Test1: %s != %s", u.Test1, str1)
	}
	if u.Test2 != str2 {
		t.Fatalf("u.Test1: %d != %d", u.Test2, str2)
	}
	if u.Test3 != 13 {
		t.Fatalf("u.Test3: %d != %d", u.Test3, 13)
	}

}




const testScript = `[UnitTest]
test1={{.Str1}}
test2={{.Str2}}
`