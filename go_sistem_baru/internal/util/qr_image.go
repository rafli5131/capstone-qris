package util

import (
	"bytes"
	"errors"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
)

// DecodeQRPayloadFromImage membaca gambar dan mengembalikan string payload QR pertama yang ditemukan.
func DecodeQRPayloadFromImage(r io.Reader) (string, error) {
	if r == nil {
		return "", errors.New("image reader is nil")
	}

	// 1. Baca data dari reader
	data, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}

	// 2. Decode menjadi image.Image
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return "", err
	}

	// 3. Konversi image.Image ke BinaryBitmap (Dibutuhkan oleh ZXing)
	// NewBinaryBitmapFromImage secara otomatis menangani LuminanceSource dan Binarizer
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return "", errors.New("failed to create binary bitmap: " + err.Error())
	}

	// 4. Inisialisasi QRCodeReader
	qrReader := qrcode.NewQRCodeReader()

	// 5. Decode QR Code
	// Kita gunakan nil untuk hints, atau bisa ditambahkan jika ingin spesifik
	result, err := qrReader.Decode(bmp, nil)
	if err != nil {
		// Jika tidak ada QR yang terdeteksi, ZXing mengembalikan error NotFoundException
		return "", errors.New("no qr code detected or failed to decode")
	}

	payload := result.GetText()
	if payload == "" {
		return "", errors.New("qr code payload is empty")
	}

	return payload, nil
}
