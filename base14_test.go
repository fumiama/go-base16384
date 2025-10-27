package base14

import (
	"bytes"
	"encoding/hex"
	"errors"
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
	buf := make([]byte, 4096+1)
	_, err := rand.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	w := bytes.NewBuffer(make([]byte, 0, 4096+1))
	for i := 0; i <= 4096; i++ {
		e := NewEncoder(w)
		_, err = io.Copy(e, bytes.NewReader(buf[:i]))
		if err != nil {
			t.Fatal(err)
		}
		_ = e.Close()
		if !bytes.Equal(Encode(buf[:i]), w.Bytes()) {
			t.Log("expect", hex.EncodeToString(Encode(buf[:i])))
			t.Log("butgot", hex.EncodeToString(w.Bytes()))
			t.Fatal("unexpected at index", i)
		}
		w.Reset()
	}
}

func TestDecoder(t *testing.T) {
	buf := make([]byte, 4096+1)
	_, err := rand.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	w := bytes.NewBuffer(make([]byte, 0, 4096+1))
	for i := 0; i <= 4096; i++ {
		d := NewDecoder(bytes.NewReader(Encode(buf[:i])))
		_, err = io.Copy(w, d)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(buf[:i], w.Bytes()) {
			t.Fail()
		}
		w.Reset()
	}
}

//go:nosplit
func encodeToGeneric(b, encd []byte) (int, error) {
	outlen := len(b) / 7 * 8
	offset := len(b) % 7
	switch offset { //算上偏移标志字符占用的2字节
	case 0:
		break
	case 1:
		outlen += 4
	case 2, 3:
		outlen += 6
	case 4, 5:
		outlen += 8
	case 6:
		outlen += 10
	}
	if len(encd) < outlen {
		return 0, errors.New("encd too small")
	}
	encodeGeneric(offset, outlen, b, encd)
	return outlen, nil
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

func benchEncodeGeneric(b *testing.B, data []byte) {
	_, err := rand.Read(data)
	if err != nil {
		panic(err)
	}
	buf := make([]byte, EncodeLen(len(data)))
	b.SetBytes(int64(len(data)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = encodeToGeneric(data, buf)
	}
}

//go:nosplit
func decodeToGeneric(b []byte, decd []byte) (int, error) {
	outlen := len(b)
	offset := 0
	if b[len(b)-2] == '=' {
		offset = int(b[len(b)-1])
		switch offset { //算上偏移标志字符占用的2字节
		case 0:
			break
		case 1:
			outlen -= 4
		case 2, 3:
			outlen -= 6
		case 4, 5:
			outlen -= 8
		case 6:
			outlen -= 10
		}
	}
	outlen = outlen/8*7 + offset
	if len(decd) < outlen {
		return 0, errors.New("decd too small")
	}
	decodeGeneric(offset, outlen, b, decd)
	return outlen, nil
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

func benchDecodeGeneric(b *testing.B, data []byte) {
	_, err := rand.Read(data)
	if err != nil {
		panic(err)
	}
	buf := make([]byte, EncodeLen(len(data)))
	_, err = encodeToGeneric(data, buf)
	if err != nil {
		panic(err)
	}
	b.SetBytes(int64(len(buf)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = decodeToGeneric(buf, data)
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

func BenchmarkEncodeToGeneric(b *testing.B) {
	b.Run("16", func(b *testing.B) {
		data := make([]byte, 16)
		benchEncodeGeneric(b, data)
	})
	b.Run("256", func(b *testing.B) {
		data := make([]byte, 256)
		benchEncodeGeneric(b, data)
	})
	b.Run("4K", func(b *testing.B) {
		data := make([]byte, 1024*4)
		benchEncodeGeneric(b, data)
	})
	b.Run("32K", func(b *testing.B) {
		data := make([]byte, 1024*32)
		benchEncodeGeneric(b, data)
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

func BenchmarkDecodeToGeneric(b *testing.B) {
	b.Run("16", func(b *testing.B) {
		data := make([]byte, 16)
		benchDecodeGeneric(b, data)
	})
	b.Run("256", func(b *testing.B) {
		data := make([]byte, 256)
		benchDecodeGeneric(b, data)
	})
	b.Run("4K", func(b *testing.B) {
		data := make([]byte, 4096)
		benchDecodeGeneric(b, data)
	})
	b.Run("32K", func(b *testing.B) {
		data := make([]byte, 1024*32)
		benchDecodeGeneric(b, data)
	})
}

func benchEncoder(b *testing.B, cnt int64) {
	s := rand.New(rand.NewSource(0))
	buf := bytes.NewBuffer(make([]byte, 0, cnt))
	enc := NewEncoder(buf)
	b.SetBytes(cnt)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = io.CopyN(enc, s, cnt)
		buf.Reset()
	}
}

func benchDecoder(b *testing.B, cnt int64) {
	s := rand.New(rand.NewSource(0))
	buf := bytes.NewBuffer(make([]byte, 0, cnt))
	enc := NewEncoder(buf)
	_, err := io.CopyN(enc, s, cnt)
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
