package shared

import (
	"encoding/json"
	"fmt"
	"path"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type TilemapJSON struct {
	Layers   []*TilemapLayerJSON `json:"layers"`
	Tilesets []*Tileset          `json:"tilesets"`
}

type TilemapLayerJSON struct {
	Data   []int `json:"data"`
	Width  int   `json:"width"`
	Height int   `json:"height"`
}

type Tileset struct {
	Firstgid int    `json:"firstgid"`
	Source   string `json:"source"`
	Data     struct {
		Columns   int    `json:"columns"`
		ImagePath string `json:"image"`
	}
}

func NewTilemapJSON(filepath string) (*TilemapJSON, error) {
	data, err := AssetsFS.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var tilemapJSON TilemapJSON
	err = json.Unmarshal(data, &tilemapJSON)
	if err != nil {
		return nil, err
	}

	for _, tileset := range tilemapJSON.Tilesets {
		data, err := AssetsFS.ReadFile(fmt.Sprintf("assets/maps/%v", tileset.Source))
		var tilesetData Tileset
		err = json.Unmarshal(data, &tilesetData.Data)
		if err != nil {
			return nil, err
		}
		tileset.Data = tilesetData.Data
		fmt.Printf("%+v\n", tileset)
	}

	return &tilemapJSON, nil
}

func GetTileImgIndex(id int, tilemapJSON *TilemapJSON) int {
	for i := range tilemapJSON.Tilesets {
		if id < tilemapJSON.Tilesets[i].Firstgid {
			return i - 1
		}
	}

	return len(tilemapJSON.Tilesets) - 1
}

func NewTileImgList(tilemapJSON *TilemapJSON) ([]*ebiten.Image, error) {
	imgMap := make([]*ebiten.Image, len(tilemapJSON.Tilesets))
	for i := range tilemapJSON.Tilesets {
		fmt.Printf("image path: %v,\n columns: %v\n", path.Clean("/"+tilemapJSON.Tilesets[i].Data.ImagePath), tilemapJSON.Tilesets[i].Data.Columns)
		img, _, err := ebitenutil.NewImageFromFileSystem(AssetsFS, (fmt.Sprintf("assets%v", path.Clean("/"+tilemapJSON.Tilesets[i].Data.ImagePath))))
		if err != nil {
			return nil, err
		}
		imgMap[i] = img
	}

	return imgMap, nil
}
