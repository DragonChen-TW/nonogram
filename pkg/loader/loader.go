package loader

type FileLoader interface {
	Load(path string) (GameData, error)
}
