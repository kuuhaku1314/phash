package phash

import (
	"fmt"
	"testing"
)

func TestGetImageSimilarity(t *testing.T) {
	src := "test.jpg"
	dst := "test_gray.jpg"
	similarity, err := GetImageSimilarity(src, dst)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(similarity)
	dst = "test_resize_300x500.jpg"
	similarity, err = GetImageSimilarity(src, dst)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(similarity)
}
