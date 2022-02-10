package base14

import (
	"testing"
)

func TestBase14(t *testing.T) {
	teststr := "一个测试293大大的啊定位为恶的我284的我……#@%@%@"
	e := EncodeString(teststr)
	es, err := UTF16be2utf8(e)
	if err == nil {
		t.Log(string(es))
		if string(es) != "蜮嘎惢磦筢貊豔耹嫹桊涖犧蟦癎摖壥禦籋萷犸粹瘛榞梄螢圓因苧璡屨灇炀瞸瘊暍严帉戀㴃" {
			t.Fail()
		}
		d, err := UTF82utf16be(es)
		if string(d) == string(e) {
			if err == nil {
				ds := DecodeString(d)
				t.Log(ds)
				if ds != teststr {
					t.Fail()
				}
			} else {
				t.Fatal(err)
			}
		} else {
			t.Fatal(d)
		}
	} else {
		t.Fatal(err)
	}
}
