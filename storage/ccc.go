package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ohsu-comp-bio/funnel/ccc/dts"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"os"
	"path/filepath"
	"strings"
)

// CCCProtocol defines the expected prefix of URL matching this storage system.
// e.g. "file:///path/to/file" matches the CCC storage system.
const CCCProtocol = "ccc://"

// CCCBackend provides access to a ccc-disk storage system.
type CCCBackend struct {
	conf  *config.CCCStorage
	local *LocalBackend
}

// NewCCCBackend returns a CCCBackend instance
func NewCCCBackend(conf config.CCCStorage) (*CCCBackend, error) {
	local, err := NewLocalBackend(config.LocalStorage{AllowedDirs: []string{"/cluster_share"}})
	if err != nil {
		return nil, err
	}
	b := &CCCBackend{
		conf:  &conf,
		local: local,
	}
	return b, nil
}

// Get copies a file from storage into the given hostPath.
func (ccc *CCCBackend) Get(ctx context.Context, url string, hostPath string, class tes.FileType) error {
	log.Info("Starting download", "url", url, "hostPath", hostPath)

	cccID := strings.TrimPrefix(url, CCCProtocol)
	cli, err := dts.NewClient(ccc.conf.DTSUrl)
	if err != nil {
		return err
	}
	record, err := cli.Get(cccID)
	if err != nil {
		return err
	}
	log.Debug("Resolved DTS Record", "cccID", cccID, "record", fmt.Sprintf("%+v", record))

	if ccc.conf.Strategy == "fetch_file" && record.HasSiteLocation(ccc.conf.Sites.Remote) && !record.HasSiteLocation(ccc.conf.Sites.Local) {
		path := record.SitePath(ccc.conf.Sites.Remote)
		if class == File {
			cli := NewSCPClient(ccc.conf.Sites.Remote)
			err = cli.SCPRemoteToLocal(path, path)
		} else if class == Directory {
			err = fmt.Errorf("SCP of directories not supported")
		} else {
			err = fmt.Errorf("Unknown file class: %s", class)
		}
		if err != nil {
			return fmt.Errorf("Failed to download from remote %s: %v", url, err)
		}
		record, err = ccc.updateDTSRecord(cccID, path, ccc.conf.Sites.Local)
		if err != nil {
			return fmt.Errorf("Failed update the DTS record: %v", err)
		}
	}

	if record.HasSiteLocation(ccc.conf.Sites.Local) {
		path := record.SitePath(ccc.conf.Sites.Local)
		err = ccc.local.Get(ctx, path, hostPath, class)
		if err != nil {
			return fmt.Errorf("Failed to download %s: %v", url, err)
		}
	} else {
		return fmt.Errorf("Failed to download %s: %v", url, err)
	}

	log.Info("Finished download", "url", url, "hostPath", hostPath)
	return nil
}

// Put copies a file from the hostPath into storage.
func (ccc *CCCBackend) Put(ctx context.Context, url string, hostPath string, class tes.FileType) error {
	log.Info("Starting upload", "url", url, "hostPath", hostPath)

	path := strings.TrimPrefix(url, CCCProtocol)
	cli, err := dts.NewClient(ccc.conf.DTSUrl)
	if err != nil {
		return err
	}
	record, err := cli.Get(path)
	if err == nil {
		return fmt.Errorf("CCCID %s conflicts with an existing record: %+v", path, record)
	}

	if ccc.conf.Strategy == "push_file" {
		if class == File {
			cli := NewSCPClient(ccc.conf.Sites.Remote)
			err = cli.SCPLocalToRemote(hostPath, path)
		} else if class == Directory {
			err = fmt.Errorf("SCP of directories not supported")
		} else {
			err = fmt.Errorf("Unknown file class: %s", class)
		}
		if err != nil {
			return fmt.Errorf("Failed to upload to remote %s: %v", url, err)
		}
	}

	err = ccc.local.Put(ctx, path, hostPath, class)
	if err != nil {
		return fmt.Errorf("Failed to upload %s: %v", url, err)
	}

	cerr := ccc.createDTSRecord(path)
	if cerr != nil {
		return fmt.Errorf("Failed to create DTS Record for output %s. %v", url, cerr)
	}

	log.Info("Finished upload", "url", url, "hostPath", hostPath)
	return nil
}

// Supports indicates whether this backend supports the given storage request.
// For the CCCBackend, the url must start with "ccc://"
func (ccc *CCCBackend) Supports(url string, hostPath string, class tes.FileType) bool {
	return strings.HasPrefix(url, CCCProtocol)
}

func (ccc *CCCBackend) createDTSRecord(path string) error {
	record, err := dts.GenerateRecord(path, ccc.conf.Sites.Local)
	if err != nil {
		return fmt.Errorf("Failed to generate DTS Record for output %s. %v", path, err)
	}

	if ccc.conf.Strategy == "push_file" {
		var l dts.Location
		l = record.Location[0]
		l.Site = ccc.conf.Sites.Remote
		record.Location = append(record.Location, l)
	}

	cli, err := dts.NewClient(ccc.conf.DTSUrl)
	if err != nil {
		return err
	}
	msg, err := json.Marshal(record)
	if err != nil {
		return err
	}
	err = cli.Post(msg)
	if err != nil {
		return err
	}
	log.Debug("Created DTS Record", "record", fmt.Sprintf("%+v", record))
	return err
}

func (ccc *CCCBackend) updateDTSRecord(cccID string, hostPath string, site string) (*dts.Record, error) {
	cli, err := dts.NewClient(ccc.conf.DTSUrl)
	if err != nil {
		return nil, err
	}
	record, err := cli.Get(cccID)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(hostPath)
	if err != nil {
		return nil, err
	}
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	name := fi.Name()
	size := fi.Size()
	if name != record.Name || size != record.Size {
		return nil, fmt.Errorf("Name: %s or Size: %d in update does not match record: %+v", name, size, record)
	}

	var l dts.Location
	l.Site = site
	l.Path = filepath.Dir(hostPath)
	l.User.Name = os.Getenv("USER")
	l.TimestampUpdated = fi.ModTime().Unix()
	record.Location = append(record.Location, l)

	msg, err := json.Marshal(record)
	if err != nil {
		return nil, err
	}
	err = cli.Put(msg)
	if err != nil {
		return nil, err
	}
	log.Debug("Updated DTS Record", "record", fmt.Sprintf("%+v", record))
	return record, err
}
