/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package sas

import (
	"os"
	"testing"
	"time"

	"k8s.io/klog/v2/ktesting"
)

type fakeFileInfo struct {
	name string
}

func (fi *fakeFileInfo) Name() string {
	return fi.name
}

func (fi *fakeFileInfo) Size() int64 {
	return 0
}

func (fi *fakeFileInfo) Mode() os.FileMode {
	return 777
}

func (fi *fakeFileInfo) ModTime() time.Time {
	return time.Now()
}
func (fi *fakeFileInfo) IsDir() bool {
	return false
}

func (fi *fakeFileInfo) Sys() interface{} {
	return nil
}

type fakeIOHandler struct{}

func (handler *fakeIOHandler) ReadDir(dirname string) ([]os.FileInfo, error) {
	switch dirname {
	case "/dev/disk/by-path/":
		f := &fakeFileInfo{
			name: "pci-0000:41:00.0-fc-0x500a0981891b8dc5-lun-0",
		}
		return []os.FileInfo{f}, nil
	case "/sys/block/":
		f := &fakeFileInfo{
			name: "dm-1",
		}
		return []os.FileInfo{f}, nil
	case "/dev/disk/by-id/":
		f := &fakeFileInfo{
			name: "wwn-0x3600508b400105e210000900000490000",
		}
		return []os.FileInfo{f}, nil
	}
	return nil, nil
}

func (handler *fakeIOHandler) Lstat(name string) (os.FileInfo, error) {
	return nil, nil
}

func (handler *fakeIOHandler) EvalSymlinks(path string) (string, error) {
	return "/dev/sda", nil
}

func (handler *fakeIOHandler) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return nil
}

func TestSearchDisk(t *testing.T) {
	logger, _ := ktesting.NewTestContext(t)
	fakeConnector := Connector{
		VolumeName: "fakeVol",
		VolumeWWN:  "3600508b400105e210000900000490000",
	}

	error := discoverDevices(logger, &fakeConnector, &fakeIOHandler{})

	if error != nil {
		t.Errorf("no disk found")
	}
}

func TestInvalidWWN(t *testing.T) {
	logger, _ := ktesting.NewTestContext(t)
	testWwn := "INVALIDWWN"
	dm, devices := findDiskById(logger, testWwn, &fakeIOHandler{})

	if len(devices) > 0 || dm != "" {
		t.Error("Found a disk with WWN that does not Exist")
	}
}
