package bytesequal

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io"
)

// MD5Hash создает хэш для ридера
func MD5Hash(src io.Reader, hashSize int64) (string, error) {
	hash := md5.New()
	if _, err := io.CopyN(hash, src, hashSize); err != nil && err != io.EOF {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// BytesAreEqual сравнивает два байт слайса и возвращает true если они равны
func BytesAreEqual(b1 []byte, b2 []byte) bool {
	return bytes.Compare(b1, b2) == 0
}
