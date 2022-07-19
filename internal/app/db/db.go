package db

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
)

type Repository interface {
	GetImageOffProduct(int) ([]ImagesData, error)
}

type repository struct {
	db *sql.DB
	qb sq.StatementBuilderType
}

func NewDbConnectClient(sqlConnect string) (Repository, error) {
	bd, err := sql.Open("postgres", sqlConnect) //postgres
	if err != nil {
		return nil, err
	}
	return &repository{bd, sq.StatementBuilder.PlaceholderFormat(sq.Dollar)}, nil
}

func (r repository) GetImageOffProduct(itemId int) ([]ImagesData, error) {
	Images := make([]ImagesData, 0, 5)

	textQuery := `
	SELECT item_id, image_id, ii.src_url, ii.url, ii.created_at, ii.updated_at, ii.checksum
	FROM (SELECT item_id, unnest(image_ids) as image_id from item_images WHERE item_id = $1) AS image_item
	JOIN image as ii on image_item.image_id = ii.id
	`

	rows, err := r.db.Query(textQuery, itemId)
	defer rows.Close()

	for rows.Next() {
		ImageData := ImagesData{}
		var ItemId, ImageId int
		var SrcUrl, Url, CreatedAt, UpdatedAt, Checksum string

		errScan := rows.Scan(&ItemId, &ImageId, &SrcUrl, &Url, &CreatedAt, &UpdatedAt, &Checksum)
		if errScan != nil {
			return nil, errScan
		}
		ImageData.ItemId, ImageData.ImageId = ItemId, ImageId
		ImageData.SrcUrl, ImageData.Url = SrcUrl, Url
		ImageData.CreatedAt, ImageData.UpdatedAt = CreatedAt, UpdatedAt
		ImageData.Md5Sum = Checksum

		Images = append(Images, ImageData)
	}

	if err != nil {
		return nil, err
	}

	return Images, nil
}
