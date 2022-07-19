package main

import (
	"crypto/md5"
	"encoding/csv"
	"fmt"
	db2 "github.com/Sergei3232/image_md5_check/internal/app/db"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

const fileNameCSV = "imageCSV.csv"

func initItems() []int {
	items := []int{
		111,
	}
	return items
}

func initMapImagesSku() map[int]int {
	ImagesSku := map[int]int{
		111: 111,
	}

	return ImagesSku
}

func main() {
	ErrorImages := make([]db2.ImagesData, 0, 100)

	ItemsIds := initItems()
	SkuMap := initMapImagesSku()

	db, err := db2.NewDbConnectClient("")
	if err != nil {
		log.Panicln(err)
	}

	for i := 0; i < len(ItemsIds); i++ {
		ls, err := db.GetImageOffProduct(ItemsIds[i])
		if err != nil {
			log.Panicln(err)
		}

		ImagesErr, err := ConsumerImages(ls)
		if err != nil {
			log.Panicln(err)
		}
		ErrorImages = append(ErrorImages, ImagesErr...)
	}

	for i := 0; i < len(ErrorImages); i++ {
		ErrorImages[i].Sku = SkuMap[ErrorImages[i].ItemId]
	}

	err = SaveCSV(ErrorImages, fileNameCSV)
	if err != nil {
		log.Panicln(err)
	}
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

func ComparisonFiles(PathSrc, FileUrl, md5Sum string) (bool, error) {
	resSrcUrl, err := http.Get(PathSrc)
	if err != nil {
		return false, err
	}

	resUrl, err := http.Get(FileUrl)
	if err != nil {
		return false, err
	}

	md5SrcUrl, err := GetMD5File(resSrcUrl.Body)
	if err != nil {
		return false, nil
	}

	_, err = GetMD5File(resUrl.Body)
	if err != nil {
		return false, nil
	}
	if md5Sum != md5SrcUrl {
		return true, nil
	}

	return false, nil
}

func ConsumerImages(images []db2.ImagesData) ([]db2.ImagesData, error) {
	errorImages := make([]db2.ImagesData, 0, 20)

	for i := 0; i < len(images); i++ {
		ok, err := ComparisonFiles(images[i].SrcUrl, images[i].Url, images[i].Md5Sum)
		if err != nil {
			return nil, err
		}
		if ok {
			errorImages = append(errorImages, images[i])
		}
	}

	return errorImages, nil
}

func SaveCSV(images []db2.ImagesData, fileName string) error {
	records := [][]string{
		{"sku", "src_url", "url", "created_at", "updated_at", "md5"},
	}

	for i := 0; i < len(images); i++ {
		csvData := []string{strconv.Itoa(int(images[i].Sku)), images[i].SrcUrl,
			images[i].Url, images[i].CreatedAt, images[i].UpdatedAt, images[i].Md5Sum}
		records = append(records, csvData)
	}

	file, errCreate := os.Create(fileName)
	if errCreate != nil {
		log.Panic(errCreate)
	}

	w := csv.NewWriter(file)

	for _, record := range records {
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	w.Flush()

	if err := w.Error(); err != nil {
		return err
	}
	return nil
}
