package phash

import (
	"bytes"
	"errors"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"math"
	"strings"
)

const (
	Width  = 32
	Height = 32
)

func GetImageSimilarity(srcPath, destPath string) (int, error) {
	var (
		srcImage  image.Image
		destImage image.Image
		err       error
	)
	if srcImage, err = decodeImage(srcPath); err != nil {
		return 0, err
	}
	if destImage, err = decodeImage(destPath); err != nil {
		return 0, err
	}
	return imageSimilarity(srcImage, destImage), nil
}

func decodeImage(imagePath string) (decodedImage image.Image, err error) {
	rawImage, err := ioutil.ReadFile(imagePath)
	if err != nil {
		return nil, errors.New("read file error")
	}
	if strings.HasSuffix(imagePath, "jpg") || strings.HasSuffix(imagePath, "jpeg") {
		decodedImage, err = jpeg.Decode(bytes.NewReader(rawImage))
	} else if strings.HasSuffix(imagePath, "png") {
		decodedImage, err = png.Decode(bytes.NewReader(rawImage))
	} else {
		return nil, errors.New("not supported format")
	}
	if err != nil {
		// Maybe the suffix and format of the image don't match. Try again
		if strings.HasSuffix(imagePath, "png") {
			decodedImage, err = jpeg.Decode(bytes.NewReader(rawImage))
		} else {
			decodedImage, err = png.Decode(bytes.NewReader(rawImage))
		}
		if err != nil {
			return nil, errors.New("not supported format")
		}
	}
	return decodedImage, nil
}

func imageSimilarity(srcImage, destImage image.Image) int {
	srcHash := pHash(srcImage)
	dstHash := pHash(destImage)
	return ((8*8 - hmDistance(srcHash, dstHash)) * 100) / (8 * 8)
}

func pHash(img image.Image) string {
	resizeImg := resize.Resize(Width, Height, img, resize.Lanczos3)
	grayImg := grayingImage(resizeImg)
	imgMatrix := grayMatrix(grayImg)
	var resultMatrix [Height][Width]float64
	dct(&resultMatrix, imgMatrix, Width, Height)
	var sum float64 = 0
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			sum += resultMatrix[i][j]
		}
	}
	avg := sum / (8 * 8)
	sb := strings.Builder{}
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if resultMatrix[i][j] < avg {
				sb.WriteByte('0')
			} else {
				sb.WriteByte('1')
			}
		}
	}
	return sb.String()
}

func grayingImage(img image.Image) image.Image {
	rgba := image.NewRGBA(img.Bounds())
	for i := 0; i < img.Bounds().Dy(); i++ {
		for j := 0; j < img.Bounds().Dx(); j++ {
			r, g, b, a := color.GrayModel.Convert(img.At(j, i)).RGBA()
			rgba.SetRGBA(j, i, color.RGBA{
				R: uint8(r >> 8),
				G: uint8(g >> 8),
				B: uint8(b >> 8),
				A: uint8(a >> 8),
			})
		}
	}
	return rgba
}

func grayMatrix(img image.Image) *[Height][Width]float64 {
	var matrix [Height][Width]float64
	for i := 0; i < Height; i++ {
		for j := 0; j < Width; j++ {
			_, g, _, _ := img.At(i, j).RGBA()
			matrix[i][j] = float64(g >> 8)
		}
	}
	return &matrix
}

func hmDistance(src, dst string) int {
	if len(src) != len(dst) {
		panic("abnormal string length")
	}
	distance := 0
	srcBytes := []byte(src)
	dstBytes := []byte(dst)
	for i, c := range srcBytes {
		if dstBytes[i] != c {
			distance++
		}
	}
	return distance
}

func dct(DCTMatrix, Matrix *[Height][Width]float64, M, N int) {
	var (
		i = 0
		j = 0
		u = 0
		v = 0
	)
	for u = 0; u < N; u++ {
		for v = 0; v < M; v++ {
			(*DCTMatrix)[u][v] = 0
			for i = 0; i < N; i++ {
				for j = 0; j < M; j++ {
					(*DCTMatrix)[u][v] += (*Matrix)[i][j] * math.Cos(math.Pi/(float64(N))*(float64(i)+1./2.)*float64(u)) * math.Cos(math.Pi/(float64(M))*(float64(j)+1./2.)*float64(v))
				}
			}
		}
	}
}
