package files

import (
	"bot/lib/errorH"
	"bot/storage"
	"encoding/gob"
	"errors"
	"math/rand"
	"os"
	"path/filepath"
)

type Storage struct {
	basePath string
}

const perm = 0774

var ErrorNoSavedPage = errors.New("no saved pages")

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s *Storage) Save(page *storage.Page) (err error) {
	defer func() {
		err = errorH.WrapIfErr("cant save file!", err)
	}()

	fPath := filepath.Join(s.basePath, page.UserName)

	if err := os.MkdirAll(fPath, perm); err != nil {
		return err
	}

	fName, err := fileName(page)
	if err != nil {
		return err
	}

	fPath = filepath.Join(fPath, fName)

	f, err := os.Create(fPath)
	defer func() { _ = f.Close() }()
	if err != nil {
		return err
	}

	if err := gob.NewEncoder(f).Encode(page); err != nil {
		return err
	}

	return nil
}

func (s *Storage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() { err = errorH.WrapIfErr("Cant pick random!", err) }()

	path := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, ErrorNoSavedPage
	}

	n := rand.Intn(len(files))
	f := files[n]

	return s.decodePage(filepath.Join(path, f.Name()))
}

func (s *Storage) Remove(p *storage.Page) error {
	fileName, err := fileName(p)
	if err != nil {
		return errorH.Wrap("Cant remove file", err)
	}

	if err := os.Remove(fileName); err != nil {
		return errorH.Wrap("Cant remove file", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	if err := os.Remove(path); err != nil {
		return errorH.Wrap("Cant remove file!", err)
	}

	return nil
}

func (s *Storage) IfExist(p *storage.Page) (bool, error) {
	fileName, err := fileName(p)
	if err != nil {
		return false, errorH.Wrap("Cant find file!", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		return false, errorH.Wrap("Cant find file!", err)
	}

	return true, nil
}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, errorH.Wrap("Can`t open to decode file!", err)
	}
	defer func() { _ = f.Close() }()

	var p storage.Page

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, errorH.Wrap("Cant decode file!", err)
	}

	return &p, err
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
