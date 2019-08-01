package manager

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	hookTypes "github.com/labbsr0x/bindman-dns-webhook/src/types"
	"github.com/sirupsen/logrus"
)

const (
	// Extension sets the extension of the files holding the records infos
	Extension = "bindman"
)

// delayRemove schedules the removal of a DNS Resource Record
// it cancels the operation when it identifies the name was read
func (m *Manager) delayRemove(name, recordType string) {
	if m.HasDNSRecord(name, recordType) {
		go m.removeRecord(name, recordType) // marks its removal intent
		ticker := time.NewTicker(m.RemovalDelay)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if _, err := m.DNSRecords.Read(m.getRecordFileName(name, recordType)); err == nil { // record has been read
					logrus.Infof("Cancelling delayed removal of '%s' '%s'", name, recordType)
					return
				}

				// only remove in case the record has not been read
				if err := m.DNSUpdater.RemoveRR(name, recordType); err != nil {
					logrus.Infof("Error occurred while trying to remove '%s' '%s': %s", name, recordType, err)
				} else {
					logrus.Infof("record name '%s' and type '%s' removed successfully", name, recordType)
				}
				return
			}
		}
	} else {
		logrus.Errorf("Service '%v' cannot be removed given it does not exist.", name)
	}
}

// saveRecord saves a record to the local storage
func (m *Manager) saveRecord(record hookTypes.DNSRecord) (err error) {
	var r []byte
	r, err = json.Marshal(record)
	if err == nil {
		m.Door.Lock()
		defer m.Door.Unlock()

		err = m.DNSRecords.Write(m.getRecordFileName(record.Name, record.Type), r)
	}
	return
}

// removeRecord removes the record
func (m *Manager) removeRecord(recordName, recordType string) {
	m.Door.Lock()
	defer m.Door.Unlock()
	// marks its removal
	recordFileName := m.getRecordFileName(recordName, recordType)
	if err := m.DNSRecords.Erase(recordFileName); err != nil {
		logrus.Errorf("error to erase record '%s': %s", recordFileName, err)
	}
}

// getRecordFileName return the name of the file holding the record information
func (m *Manager) getRecordFileName(recordName, recordType string) string {
	toReturn := fmt.Sprintf("%v.%v.%v", recordName, recordType, Extension)
	return toReturn
}

// getRecordName returns the name of a record from its fileName
func (m *Manager) getRecordNameAndType(fileName string) (string, string) {
	subName := strings.TrimSuffix(fileName, "."+Extension)
	i := strings.LastIndex(subName, ".")
	return subName[:i], subName[i+1:]
}
