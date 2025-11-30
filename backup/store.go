package backup

import (
	"crypto/sha256"
	"fmt"
	goio "io"
	"os"
	"strings"
	"time"

	"github.com/Siposattila/go-backup/io"
	"github.com/Siposattila/go-backup/log"
	"github.com/Siposattila/go-backup/serializer"
)

const CHECKSUM_STORE_FILENAME = "checksum_store.json"
const OLD_CHECKSUM_STORE_FILENAME = "old_checksum_store.json"

type store struct {
	LastBackup string  `json:"lastBackup"`
	Entries    []entry `json:"entries"`
}

type entry struct {
	Name     string `json:"name"`
	Checksum string `json:"checksum"`
}

func (s *store) checksum(name string) string {
	file, _ := os.Open(name)
	hash := sha256.New()
	if _, err := goio.Copy(hash, file); err != nil {
		log.GetLogger().Fatal("Failed to calculate checksum for file: ", err.Error())
	}

	return fmt.Sprintf("%x", hash.Sum(nil))
}

func getStore() (s *store) {
	s = &store{}
	if raw, err := io.ReadFile(".", CHECKSUM_STORE_FILENAME); err != nil {
		s.Entries = make([]entry, 0)
		log.GetLogger().Warning("No checksum store was found!")
	} else {
		if err := serializer.Json.Serialize(raw, s); err != nil {
			log.GetLogger().Warning("Checksum file was corrupted not able to serialize it!")
		}
		if err := os.Rename(CHECKSUM_STORE_FILENAME, OLD_CHECKSUM_STORE_FILENAME); err != nil {
			log.GetLogger().Error("Was not able to do rename for store: ", err.Error())
		}
	}

	return s
}

func (s *store) saveStore() {
	s.LastBackup = time.Now().Format("2006-01-02")
	buffer, err := serializer.Json.Deserialize(s)
	if err != nil {
		log.GetLogger().Fatal("Failed to desirialize store: ", err.Error())
	}

	io.WriteFile(".", CHECKSUM_STORE_FILENAME, buffer)
}

func (s *store) add(newEntry entry) {
	s.Entries = append(s.Entries, newEntry)
}

func (s *store) contains(name string) (bool, int) {
	name = s.trimName(name)
	for i, entry := range s.Entries {
		if s.trimName(entry.Name) == name {
			return true, i
		}
	}

	return false, -1
}

func (s *store) trimName(name string) string {
	return strings.Trim(name, "<>:/|?*';!@#$%^&()[]{}=+~`,.\t\n\r ")
}
