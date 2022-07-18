package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
)

type ImagesUrls struct {
	SrcUrl string
	Url    string
}

func initItems() []int {
	items := []int{
		313087867, 313087864,
		313087861, 313087866,
		313087868, 313087614,
		313087608, 313087033,
		313087032, 313087014,
		312797953, 312797950,
		312797945, 312675762,
		311917155, 66636515,
		313087607, 313087009,
		312671623, 66157959,
		66636933, 1771822,
	}
	return items
}

func main() {
	ItemsIds := initItems()
	fmt.Println(ItemsIds)
}

func GetMD5File(file io.Reader) (string, error) {

	hash := md5.New()
	_, err := io.Copy(hash, file)

	if err != nil {
		return "", err
	}

	h := fmt.Sprintf("%x", hash.Sum(nil))

	return h, nil
}

func ComparisonFiles(PathSrc, FileUrl string) (bool, error) {
	resSrcUrl, err := http.Get(PathSrc)
	if err != nil {
		return false, err
	}

	resUrl, err := http.Get(FileUrl)
	if err != nil {
		return false, err
	}

	if resSrcUrl != resUrl {
		return true, nil
	}

	return false, nil
}
