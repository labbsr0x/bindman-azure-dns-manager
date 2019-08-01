package manager

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/labbsr0x/bindman-azure-dns-manager/src/azure"
	hookTypes "github.com/labbsr0x/bindman-dns-webhook/src/types"
	"github.com/peterbourgon/diskv"
	"github.com/sirupsen/logrus"
)

type Builder struct {
	TTL          time.Duration
	RemovalDelay time.Duration
}

// Manager holds the information for managing a dns server
type Manager struct {
	*Builder
	DNSRecords *diskv.Diskv
	Door       *sync.RWMutex
	DNSUpdater azure.DNSUpdater
}

// New creates a new Manager instance
func (b *Builder) New(dnsupdater azure.DNSUpdater, basePath string) (*Manager, error) {
	if dnsupdater == nil {
		return nil, errors.New("not possible to start the Bindman Manager; Bindman Manager expects a valid non-nil DNSUpdater")
	}

	if strings.TrimSpace(basePath) == "" {
		return nil, errors.New("not possible to start the Bindman Manager; Bindman Manager expects a non-empty basePath")
	}

	result := &Manager{
		DNSRecords: diskv.New(diskv.Options{
			BasePath:     basePath,
			Transform:    func(s string) []string { return []string{} },
			CacheSizeMax: 1024 * 1024,
		}),
		Builder:    b,
		Door:       new(sync.RWMutex),
		DNSUpdater: dnsupdater,
	}
	return result, nil
}

// GetDNSRecords retrieves all the dns records being managed
func (m *Manager) GetDNSRecords() (records []hookTypes.DNSRecord, err error) {
	m.Door.RLock()
	defer m.Door.RUnlock()

	err = filepath.Walk(m.DNSRecords.BasePath, func(path string, info os.FileInfo, errr error) error {
		if strings.HasSuffix(path, Extension) {
			r, err := m.GetDNSRecord(m.getRecordNameAndType(info.Name()))
			if err != nil {
				return err
			}
			if r != nil {
				records = append(records, *r)
			}
		}
		return nil
	})
	return
}

// GetDNSRecord retrieves the dns record identified by name
func (m *Manager) HasDNSRecord(name, recordType string) bool {
	key := m.getRecordFileName(name, recordType)
	return m.DNSRecords.Has(key)
}

// GetDNSRecord retrieves the dns record identified by name
func (m *Manager) GetDNSRecord(name, recordType string) (record *hookTypes.DNSRecord, err error) {
	m.Door.RLock()
	defer m.Door.RUnlock()

	if !m.HasDNSRecord(name, recordType) {
		return nil, hookTypes.NotFoundError(fmt.Sprintf("No record found with name '%s' and type '%s'", name, recordType), nil)
	}

	var r []byte
	r, err = m.DNSRecords.Read(m.getRecordFileName(name, recordType))
	if err == nil {
		err = json.Unmarshal(r, &record)
	}
	return
}

// AddDNSRecord adds a new DNS record
func (m *Manager) AddDNSRecord(record hookTypes.DNSRecord) (err error) {
	err = m.DNSUpdater.AddRR(record, m.TTL)
	if err == nil {
		err = m.saveRecord(record)
	}
	return
}

// UpdateDNSRecord updates an existing dns record
func (m *Manager) UpdateDNSRecord(record hookTypes.DNSRecord) (err error) {
	err = m.DNSUpdater.UpdateRR(record, m.TTL)
	if err == nil {
		err = m.saveRecord(record)
	}
	return
}

// RemoveDNSRecord removes a DNS record
func (m *Manager) RemoveDNSRecord(name, recordType string) error {
	if !m.HasDNSRecord(name, recordType) {
		return hookTypes.NotFoundError(fmt.Sprintf("No record found with name '%s' and type '%s", name, recordType), nil)
	}
	go m.delayRemove(name, recordType)
	logrus.Infof("Record '%s' with type '%v' scheduled to be removed in %v seconds", name, recordType, m.RemovalDelay)
	return nil
}
