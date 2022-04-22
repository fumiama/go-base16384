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
		assert.Equal(t, hex.EncodeToString(buf[:i]), hex.EncodeToString(out))
	}
}

func TestEncoder(t *testing.T) {
	buf := make([]byte, 42242141)
	_, err := rand.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	e := NewEncoder(bytes.NewReader(buf))
	w := bytes.NewBuffer(make([]byte, 0, 42242150))
	_, err = io.Copy(w, e)
	if err != nil {
		t.Fatal(err)
	}
	out := w.Bytes()
	assert.Equal(t, 48276736, w.Len())
	d := Decode(out)
	t.Log(len(out))
	assert.Equal(t, buf, d)
}

func TestBufferedEncoder(t *testing.T) {
	buf := make([]byte, 1024*1024+1)
	_, err := rand.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	e := NewBufferedEncoder(buf)
	w := bytes.NewBuffer(make([]byte, 0, 1024*1024+16))
	_, err = io.Copy(w, e)
	if err != nil {
		t.Fatal(err)
	}
	out := w.Bytes()
	t.Log(w.Len())
	d := Decode(out)
	if !bytes.Equal(buf, d) {
		t.Fail()
	}
}

func TestDecoder(t *testing.T) {
	buf := make([]byte, 1024*1024+1)
	_, err := rand.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	w := bytes.NewBuffer(make([]byte, 0, 1024*1024+1))
	d := NewDecoder(bytes.NewReader(Encode(buf)))
	_, err = io.Copy(w, d)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(w.Len())
	if !bytes.Equal(buf, w.Bytes()) {
		t.Fail()
	}
}

func TestBufferedDecoder(t *testing.T) {
	buf := make([]byte, 1024*1024+1)
	_, err := rand.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	w := bytes.NewBuffer(make([]byte, 0, 1024*1024+1))
	d := NewBufferedDecoder(Encode(buf))
	_, err = io.Copy(w, d)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(w.Len())
	if !bytes.Equal(buf, w.Bytes()) {
		t.Fail()
	}
}

func benchEncrypt(b *testing.B, data []byte) {
	_, err := rand.Read(data)
	if err != nil {
		panic(err)
	}
	buf := make([]byte, EncodeLen(len(data)))
	b.SetBytes(int64(len(data)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = EncodeTo(data, buf)
	}
}

func benchDecrypt(b *testing.B, data []byte) {
	_, err := rand.Read(data)
	if err != nil {
		panic(err)
	}
	buf := make([]byte, EncodeLen(len(data)))
	err = EncodeTo(data, buf)
	if err != nil {
		panic(err)
	}
	b.SetBytes(int64(len(buf)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = DecodeTo(buf, data)
	}
}

func BenchmarkEncodeTo(b *testing.B) {
	b.Run("16", func(b *testing.B) {
		data := make([]byte, 16)
		benchEncrypt(b, data)
	})
	b.Run("256", func(b *testing.B) {
		data := make([]byte, 256)
		benchEncrypt(b, data)
	})
	b.Run("4K", func(b *testing.B) {
		data := make([]byte, 1024*4)
		benchEncrypt(b, data)
	})
	b.Run("32K", func(b *testing.B) {
		data := make([]byte, 1024*32)
		benchEncrypt(b, data)
	})
}

func BenchmarkDecodeTo(b *testing.B) {
	b.Run("16", func(b *testing.B) {
		data := make([]byte, 16)
		benchDecrypt(b, data)
	})
	b.Run("256", func(b *testing.B) {
		data := make([]byte, 256)
		benchDecrypt(b, data)
	})
	b.Run("4K", func(b *testing.B) {
		data := make([]byte, 4096)
		benchDecrypt(b, data)
	})
	b.Run("32K", func(b *testing.B) {
		data := make([]byte, 1024*32)
		benchDecrypt(b, data)
	})
}
