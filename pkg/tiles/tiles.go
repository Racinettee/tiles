package tiles

type TileSet struct {
	ImagePath             string
	TileWidth, TileHeight int
}
type TileLayer struct {
	Width, Height int
}
type LevelData struct {
	TileSets   []TileSet
	TileLayers []TileLayer
}
