package base14

import (
	"bytes"
	"encoding/hex"
	"io"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBase14(t *testing.T) {
	assert.Equal(t, "蜮嘎惢磦筢貊豔耹嫹桊涖犧蟦癎摖壥禦籋萷犸粹瘛榞梄螢圓因苧璡屨灇炀瞸瘊暍严帉戀㴃", EncodeString("一个测试293大大的啊定位为恶的我284的我……#@%@%@"))
	assert.Equal(t, "婀㴁", EncodeString("1"))
	assert.Equal(t, "婌渀㴂", EncodeString("12"))
	assert.Equal(t, "婌焰㴃", EncodeString("123"))
	assert.Equal(t, "婌焳帀㴄", EncodeString("1234"))
	assert.Equal(t, "婌焳廔㴅", EncodeString("12345"))
	assert.Equal(t, "婌焳廔萀㴆", EncodeString("123456"))
	assert.Equal(t, "婌焳廔萷", EncodeString("1234567"))
	assert.Equal(t, "婌焳廔萷尀㴁", EncodeString("12345678"))
	buf := make([]byte, 4096)
	for i := 1; i < 4096; i++ {
		rand.Read(buf[:i])
		out := Decode(Encode(buf[:i]))
		if !assert.Equal(t, hex.EncodeToString(buf[:i]), hex.EncodeToString(out)) {
			t.Fatal()
		}
	}
}

func TestEncoder(t *testing.T) {
	buf := make([]byte, 1024*1024+1)
	_, err := rand.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	w := bytes.NewBuffer(make([]byte, 0, 1024*1024+1))
	for i := 0; i <= 1024*1024; i += rand.Intn(128) * 7 {
		e := NewEncoder(bytes.NewReader(buf[:i]))
		_, err = io.Copy(w, e)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(Encode(buf[:i]), w.Bytes()) {
			t.Fail()
		}
		w.Reset()
	}
}

func TestBufferedEncoder(t *testing.T) {
	buf := make([]byte, 1024*1024+1)
	_, err := rand.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	w := bytes.NewBuffer(make([]byte, 0, 1024*1024+1))
	for i := 0; i <= 1024*1024; i += rand.Intn(128) * 7 {
		e := NewBufferedEncoder(buf[:i])
		_, err = io.Copy(w, e)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(Encode(buf[:i]), w.Bytes()) {
			t.Fail()
		}
		w.Reset()
	}
}

func TestDecoder(t *testing.T) {
	buf := make([]byte, 1024*1024+1)
	_, err := rand.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	encd := Encode(buf)
	w := bytes.NewBuffer(make([]byte, 0, 1024*1024+1))
	for i := 0; i <= 1024*1024; i += rand.Intn(128) * 8 {
		d := NewDecoder(bytes.NewReader(encd[:i]))
		_, err = io.Copy(w, d)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(buf[:i/8*7], w.Bytes()) {
			t.Fail()
		}
		w.Reset()
	}
}

func TestBufferedDecoder(t *testing.T) {
	buf := make([]byte, 1024*1024+1)
	_, err := rand.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	encd := Encode(buf)
	w := bytes.NewBuffer(make([]byte, 0, 1024*1024+1))
	for i := 0; i <= 1024*1024; i += rand.Intn(128) * 8 {
		d := NewBufferedDecoder(encd[:i])
		_, err = io.Copy(w, d)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(buf[:i/8*7], w.Bytes()) {
			t.Fail()
		}
		w.Reset()
	}
}

func benchEncode(b *testing.B, data []byte) {
	_, err := rand.Read(data)
	if err != nil {
		panic(err)
	}
	buf := make([]byte, EncodeLen(len(data)))
	b.SetBytes(int64(len(data)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = EncodeTo(data, buf)
	}
}

func benchDecode(b *testing.B, data []byte) {
	_, err := rand.Read(data)
	if err != nil {
		panic(err)
	}
	buf := make([]byte, EncodeLen(len(data)))
	_, err = EncodeTo(data, buf)
	if err != nil {
		panic(err)
	}
	b.SetBytes(int64(len(buf)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = DecodeTo(buf, data)
	}
}

func BenchmarkEncodeTo(b *testing.B) {
	b.Run("16", func(b *testing.B) {
		data := make([]byte, 16)
		benchEncode(b, data)
	})
	b.Run("256", func(b *testing.B) {
		data := make([]byte, 256)
		benchEncode(b, data)
	})
	b.Run("4K", func(b *testing.B) {
		data := make([]byte, 1024*4)
		benchEncode(b, data)
	})
	b.Run("32K", func(b *testing.B) {
		data := make([]byte, 1024*32)
		benchEncode(b, data)
	})
}

func BenchmarkDecodeTo(b *testing.B) {
	b.Run("16", func(b *testing.B) {
		data := make([]byte, 16)
		benchDecode(b, data)
	})
	b.Run("256", func(b *testing.B) {
		data := make([]byte, 256)
		benchDecode(b, data)
	})
	b.Run("4K", func(b *testing.B) {
		data := make([]byte, 4096)
		benchDecode(b, data)
	})
	b.Run("32K", func(b *testing.B) {
		data := make([]byte, 1024*32)
		benchDecode(b, data)
	})
}

func benchEncoder(b *testing.B, cnt int64) {
	enc := NewEncoder(rand.New(rand.NewSource(0)))
	buf := bytes.NewBuffer(make([]byte, 0, cnt))
	b.SetBytes(cnt)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = io.CopyN(buf, enc, cnt)
		buf.Reset()
	}
}

func benchDecoder(b *testing.B, cnt int64) {
	enc := NewEncoder(rand.New(rand.NewSource(0)))
	buf := bytes.NewBuffer(make([]byte, 0, cnt))
	_, err := io.CopyN(buf, enc, cnt)
	if err != nil {
		panic(err)
	}
	buf2 := bytes.NewBuffer(make([]byte, 0, cnt))
	b.SetBytes(cnt)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = io.Copy(buf2, NewDecoder(bytes.NewReader(buf.Bytes())))
		buf2.Reset()
	}
}

func BenchmarkEncoder(b *testing.B) {
	b.Run("16", func(b *testing.B) {
		benchEncoder(b, 16)
	})
	b.Run("256", func(b *testing.B) {
		benchEncoder(b, 256)
	})
	b.Run("4K", func(b *testing.B) {
		benchEncoder(b, 1024*4)
	})
	b.Run("32K", func(b *testing.B) {
		benchEncoder(b, 1024*32)
	})
}

func BenchmarkDecoder(b *testing.B) {
	b.Run("16", func(b *testing.B) {
		benchDecoder(b, 16)
	})
	b.Run("256", func(b *testing.B) {
		benchDecoder(b, 256)
	})
	b.Run("4K", func(b *testing.B) {
		benchDecoder(b, 1024*4)
	})
	b.Run("32K", func(b *testing.B) {
		benchDecoder(b, 1024*32)
	})
}
