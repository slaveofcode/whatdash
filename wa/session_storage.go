package wa

import (
	"encoding/gob"
	"log"
	"os"
	"path/filepath"

	whatsapp "github.com/slaveofcode/go-whatsapp"
)

type WASession struct {
	number  string
	session whatsapp.Session
}

type WASessions []WASession

type SessionStorage struct{}

func (s *SessionStorage) storePath() string {
	return os.TempDir() + "/whatdash/"
}

func (s *SessionStorage) FetchAll(storedSessions *WASessions) {
	s.prepareDir()

	var files []string
	err := filepath.Walk(s.storePath(), func(path string, info os.FileInfo, err error) error {
		fileName := info.Name()
		ext := filepath.Ext(fileName)
		if ext == ".gob" {
			cleanName := fileName[0 : len(fileName)-len(ext)]
			files = append(files, cleanName)
		}
		return nil
	})

	if err != nil {
		log.Fatalln("Cannot read from dir", s.storePath())
	}

	for _, number := range files {
		sess, err := s.Get(number)
		if err == nil {
			*storedSessions = append(*storedSessions, WASession{
				number:  number,
				session: sess,
			})
		}
	}

}

func (s *SessionStorage) prepareDir() {
	err := os.MkdirAll(s.storePath(), os.ModePerm)
	if err != nil {
		log.Println("Failed creating directory storage:", s.storePath())
	}
}

func (s *SessionStorage) Save(number string, session whatsapp.Session) error {
	s.prepareDir()

	file, err := os.Create(s.storePath() + number + ".gob")
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(session)
	if err != nil {
		return err
	}
	return nil
}

func (s *SessionStorage) Get(number string) (whatsapp.Session, error) {
	s.prepareDir()

	session := whatsapp.Session{}
	file, err := os.Open(s.storePath() + number + ".gob")
	if err != nil {
		return session, err
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&session)
	if err != nil {
		return session, err
	}
	return session, nil
}

func (s *SessionStorage) Destroy(number string) error {
	err := os.Remove(filepath.Join(s.storePath(), number))
	return err
}

func (s *SessionStorage) Reset() error {
	dir := s.storePath()
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}
