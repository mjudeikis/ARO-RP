// Code generated for package deploy by go-bindata DO NOT EDIT. (@generated)
// sources:
// deploy/staticresources/aro.openshift.io_clusters.yaml
// deploy/staticresources/master/deployment.yaml
// deploy/staticresources/master/rolebinding.yaml
// deploy/staticresources/master/service.yaml
// deploy/staticresources/master/serviceaccount.yaml
// deploy/staticresources/namespace.yaml
// deploy/staticresources/worker/deployment.yaml
// deploy/staticresources/worker/role.yaml
// deploy/staticresources/worker/rolebinding.yaml
// deploy/staticresources/worker/serviceaccount.yaml
package deploy

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _aroOpenshiftIo_clustersYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb4\x58\xcd\x6e\xdb\x48\x12\xbe\xeb\x29\x0a\xde\x43\x0e\x6b\xc9\x09\x16\x0b\xec\xea\x66\xd8\x93\xac\x30\x93\x8c\x11\x07\x99\xc3\x78\x0e\xa5\x66\x89\xac\x75\xb3\x9b\x53\xdd\x94\xa3\x2c\xf6\xdd\x17\xd5\x24\x25\x4a\x22\x65\x2b\xc6\xf2\x22\xa8\x7f\xbe\xfa\xed\xaf\xab\x7a\x32\x9d\x4e\x27\x58\xf1\x57\x92\xc0\xde\xcd\x01\x2b\xa6\x6f\x91\x9c\xfe\x0b\xb3\xc7\x7f\x84\x19\xfb\xab\xf5\xbb\xc9\x23\xbb\x6c\x0e\x37\x75\x88\xbe\xfc\x4c\xc1\xd7\x62\xe8\x96\x56\xec\x38\xb2\x77\x93\x92\x22\x66\x18\x71\x3e\x01\x40\xe7\x7c\x44\x1d\x0e\xfa\x17\xc0\x78\x17\xc5\x5b\x4b\x32\xcd\xc9\xcd\x1e\xeb\x25\x2d\x6b\xb6\x19\x49\x02\xef\x44\xaf\xdf\xce\xfe\x3e\x7b\x3b\x01\x30\x42\x69\xfb\x17\x2e\x29\x44\x2c\xab\x39\xb8\xda\xda\x09\x80\xc3\x92\xe6\x60\x6c\x1d\x22\x49\x98\xa1\xf8\x99\xaf\xc8\x85\x82\x57\x71\xc6\x7e\x12\x2a\x32\x2a\x33\x17\x5f\x57\x73\x38\x9a\x6f\x10\x5a\xb5\x5a\x93\x1a\xb0\x34\x62\x39\xc4\x9f\xfb\xa3\xbf\x70\x88\x69\xa6\xb2\xb5\xa0\xdd\x89\x4e\x83\x81\x5d\x5e\x5b\x94\xed\xf0\x04\x20\x18\x5f\x51\x1f\xb5\x35\x2f\xc9\x9c\xb6\x06\xac\xdf\xa1\xad\x0a\x7c\xd7\xa0\x98\x82\x4a\x6c\x54\x02\x50\x75\xaf\xef\x16\x5f\xff\x76\xbf\x37\x0c\x90\x51\x30\xc2\x55\x4c\xae\x6a\xe1\x81\x03\xc4\x82\xa0\x59\x0b\x2b\x2f\xe9\x6f\xa7\x24\x5c\xdf\x2d\xb6\xfb\x2b\xf1\x15\x49\xe4\xce\xfa\xe6\xeb\x85\xbe\x37\x7a\x20\xed\x8d\x2a\xd4\xac\x82\x4c\x63\x4e\x8d\xd8\xd6\x34\xca\x5a\x1b\xc0\xaf\x20\x16\x1c\x40\xa8\x12\x0a\xe4\x9a\x2c\xd8\x03\x06\x5d\x84\x0e\xfc\xf2\xdf\x64\xe2\x0c\xee\x49\x14\x06\x42\xe1\x6b\x9b\x69\xaa\xac\x49\x22\x08\x19\x9f\x3b\xfe\xbe\xc5\x0e\x10\x7d\x12\x6a\x31\x52\x1b\x94\xdd\xc7\x2e\x92\x38\xb4\xb0\x46\x5b\xd3\x25\xa0\xcb\xa0\xc4\x0d\x08\xa9\x14\xa8\x5d\x0f\x2f\x2d\x09\x33\xf8\xe8\x85\x80\xdd\xca\xcf\xa1\x88\xb1\x0a\xf3\xab\xab\x9c\x63\x97\xf2\xc6\x97\x65\xed\x38\x6e\xae\x52\xf6\xf2\xb2\x8e\x5e\xc2\x55\x46\x6b\xb2\x57\x81\xf3\x29\x8a\x29\x38\x92\x89\xb5\xd0\x15\x56\x3c\x4d\xaa\xbb\x94\xf6\xb3\x32\xfb\x8b\xb4\x87\x24\xbc\xd9\xd3\x35\x6e\x34\x3d\x42\x14\x76\x79\x6f\x22\xe5\xe2\x89\x08\x68\x56\x6a\xb4\xb1\xdd\xda\x58\xb1\x73\xb4\x0e\xa9\x77\x3e\xff\x74\xff\x05\x3a\xd1\x29\x18\x87\xde\x4f\x7e\xdf\x6d\x0c\xbb\x10\xa8\xc3\xd8\xad\x48\x9a\x20\xae\xc4\x97\x09\x93\x5c\x56\x79\x76\xb1\xcd\x2d\x26\x77\xe8\xfe\x50\x2f\x4b\x8e\x1a\xf7\x3f\x6b\x0a\x51\x63\x35\x83\x9b\xc4\x03\xb0\x24\xa8\xab\x0c\x23\x65\x33\x58\x38\xb8\xc1\x92\xec\x0d\x06\xfa\xbf\x07\x40\x3d\x1d\xa6\xea\xd8\x97\x85\xa0\x4f\x61\x87\x8b\x1b\xaf\xf5\x26\x3a\xa2\x19\x89\x57\x7b\x3e\xef\x2b\x32\x7b\x27\x26\xa3\xc0\xa2\x39\x1d\x31\x92\x9e\x84\x3e\xfb\x74\xdf\xf0\x49\xd5\x0f\x8d\xdc\xfa\x12\xd9\x1d\x4e\x8c\x1a\x05\xcd\x19\x5f\xb8\xb8\xb8\x3b\x6f\x53\xcf\xbb\x83\x0c\xb1\xdb\xaf\x87\x2f\x3f\xb0\x01\x00\xbf\xff\xe4\xd6\x2c\xde\x95\xe4\xe2\x59\xa2\xb3\xf3\x4d\x5c\x11\xaa\xa2\x47\x0e\x3b\x08\xcb\xfb\x76\xd9\x5e\x5c\xae\x3f\xff\xaa\xac\x2b\x18\xbd\x74\x40\x90\x2b\xcb\x1c\x81\x8d\x47\x46\x3f\xe5\x18\x67\xd8\xd2\xb5\x25\x89\xbf\xd1\xb2\xf0\xfe\x71\x68\x61\x67\xca\xd2\x7b\x4b\x78\xc8\x8f\x7b\x50\xb7\x9f\xee\x3f\x62\xf8\xf3\x95\x28\x1f\xc8\xd1\x1a\x7f\xf1\x79\xce\x2e\x7f\x25\xd6\x47\xef\x38\x7a\x8d\xc1\x8d\x77\x2b\x7e\x2d\xdc\xa7\xfb\x0f\x83\xce\x3c\x07\xc2\x67\x74\x2b\xc8\x8e\xe4\x95\x48\x77\xb5\xb5\xf7\x64\x84\x06\x12\xf6\x2c\xa0\xcf\xbe\x8e\xf4\x9e\xbf\xbd\x12\xe6\x37\x2f\x8f\x28\xbe\x76\x59\xb8\xd9\xd6\x50\x3f\x82\x39\xc2\x62\xfa\xe5\xa7\x73\xe3\x74\xca\x9b\x94\x02\xa3\xfc\x90\x00\x30\xea\xdd\x3c\x87\x37\xbf\xbf\x9d\xfe\xf3\x8f\xbf\xce\x9a\x9f\x37\x27\xac\x18\x3c\xe2\xfa\x95\xdb\xdc\xfb\x70\x73\x7f\x92\x5e\xf4\x23\x57\x97\xc3\x33\x53\xb8\x65\xcc\x9d\x0f\x91\x4d\xb8\x13\x9f\x8d\xac\xfa\x72\x5c\x69\xbc\x40\xcf\x13\xce\x66\xb7\x12\x5c\x64\x67\x71\x1b\xbb\x5c\x28\x84\x33\xf9\xbb\xa9\x88\x28\xde\x14\x64\x1e\x87\x92\xe6\x74\x60\x6b\xb1\x23\xc7\x92\x23\x95\x23\x53\xcf\xc6\xaf\x5b\x80\x22\xb8\x39\xc7\x6f\xd6\x9b\x54\x4a\x9e\xe5\x82\xae\x0c\x1a\xf2\xf7\xde\xb5\xd0\xf5\x33\x8b\xdb\xae\xa0\xbe\xfe\xae\x97\xc0\x0e\xa0\xa9\x6c\xa9\x57\xe7\xbf\x58\x8b\xb5\xa3\x78\x56\xc4\xc7\x2a\x8e\x88\xb1\x0e\x2f\xa8\x39\xd2\xba\xbd\xaa\xc3\x2f\x83\x96\x78\x3f\x5c\x76\x18\xef\x32\xee\xf5\x73\xe3\x2a\x6c\x17\xb6\xb5\x2a\xc5\x24\xad\x1b\x06\x76\x21\xa2\x33\x14\x66\x47\x40\xa3\x79\xb5\x27\xe1\x62\x87\xb5\x2b\x61\x9b\x7e\x42\x6d\x4c\x49\xb2\xd7\x61\xbc\x39\xbe\xc5\x3b\x6f\xd2\xac\xaf\x30\x0a\xe9\xae\x6d\xf3\x0b\x25\x99\x02\x1d\x87\x32\x9d\x25\x97\x51\xa6\x2d\x88\x96\xb3\x81\x86\x09\xe3\xa9\x20\xd7\x96\x79\x11\xd9\x86\xad\x22\x3b\xd5\x54\x8a\x56\xc5\x08\x95\xb0\x17\x86\x47\xe7\x9f\x1c\x78\x81\x27\xed\x7f\x06\x61\xd3\xfa\xaa\xb2\x1b\x95\x8f\xd6\xee\xbc\x98\x04\x40\xce\x6b\x72\xa0\x1d\xc2\x0c\x1e\x5c\xdf\xa6\xa6\xa9\x1a\x04\x5d\x12\x60\xd6\xda\x44\xdf\x2a\xcb\x86\xa3\xdd\x34\xfd\xd7\xa6\x97\x0b\x10\x0b\x8c\x6a\xb2\x84\xd4\x55\x19\x5f\x56\xde\xa9\xd7\x07\x61\x4d\x72\xe3\xd2\xd7\x11\x04\x63\x91\x7a\x09\x74\xa9\x31\x60\x69\x9a\x14\x1f\x68\x0f\x3f\xf9\x34\xf5\x1d\x32\xe2\xd7\xd4\x89\xf8\x84\xd6\xf3\x65\x98\xc1\xaf\xce\x50\x9b\xe9\xd9\x65\xf2\x7c\x49\xe8\x54\x4c\x72\xcc\xd6\x13\x23\xaa\x3a\x68\x1b\x14\x0d\x74\x4e\x19\xa0\x2c\x39\x0a\x0a\xdb\x0d\x4c\x81\x75\xce\xf8\x92\x02\x54\x28\xb1\xe3\x80\xeb\xbb\x45\x6a\x30\x07\x41\x0b\x6c\x8e\x5c\xc0\x92\x60\x89\xe6\xf1\x09\x25\x0b\xd3\xe4\xba\x95\x97\xe6\x9f\xfa\x10\x23\x2f\xd9\x72\x4c\x2e\x37\x24\x4e\x83\x39\x08\x89\x6e\xd3\x1a\x7f\xa0\xc5\xec\x62\x60\xfd\x69\x5a\x07\xb0\x18\xe2\x17\x41\x17\xb8\x7b\x61\x19\xe3\xf2\x95\x97\x12\xe3\x1c\xb4\x77\x9b\x46\x2e\xe9\x47\x39\xbf\xa4\x10\x30\x1f\x95\xf3\xec\x7e\x21\x0c\x63\xd5\xc5\x18\x01\x7d\x4e\x7b\x94\x85\x0e\x0e\x2f\x82\x77\x34\x7d\xf2\x92\x5d\xee\x7a\xd1\x11\x68\x38\x78\xc8\xd8\xde\x02\x18\x29\xf7\xb2\xd1\xff\x06\xeb\x40\xdb\x89\x5a\x84\x5c\x6c\xb9\xfa\x98\xe3\xba\x6f\x11\x07\x34\x53\x5a\x01\x76\x29\x1f\x58\x31\xeb\x58\xd5\xf1\x12\x42\x6d\x0a\xc0\x90\xf4\xb6\xec\xc6\x95\x7d\xac\x97\x64\xa2\x85\x5c\x59\xb7\xdd\xac\x79\xc7\x0e\x42\x5d\x96\x28\xfc\x3d\x1d\x0d\xd3\xa8\xd9\xf2\x47\x32\x60\x54\xd7\x67\x83\x33\x74\x2d\x9d\xb1\x3d\x2d\x78\x49\x64\x77\xc4\xff\x65\x53\x51\x77\x4f\xeb\xf6\xad\xf3\xb7\x37\xc3\xd8\xe1\xd4\x4f\x37\x6e\x2a\x36\x68\xed\x46\x29\xa2\x4b\x81\x0c\x34\x27\x94\x88\x43\xe1\x25\x42\x55\x48\x7a\xa4\xe8\x13\xea\x28\x68\x7a\x6a\xe8\x9e\xb0\xd8\x65\xac\x19\xd2\xde\xb6\xdc\x5c\x09\x0f\x17\xb8\x74\x7a\xa2\xec\x34\x4a\x4d\x0f\x17\x50\x79\x8b\xc2\x71\x33\x9e\x26\xef\xbd\x00\x7d\xc3\xb2\xb2\x74\x09\x7c\x68\x65\x27\x27\x34\xf7\x0e\x2a\x20\x9b\x4d\x93\x59\x6b\xb4\x9c\x5d\x8e\x2b\x9c\x34\xe2\x00\x69\xdd\xc3\x05\x18\x0c\xc9\xa9\x95\xf8\x25\x2e\xf5\xaa\x29\xf4\xa2\x92\xf2\x12\x82\xdf\x17\x3c\x0a\xda\xda\xaf\x7c\x8a\xd6\x52\x06\x0f\x17\x0b\xd7\x0a\x18\xe4\x2a\x78\x3e\x43\x9a\x8b\x83\x06\xea\x27\xad\xcf\x9b\xe4\x1b\x9c\x52\xdc\x81\x89\x13\x35\xe6\xa9\xe2\xb4\x7b\x1f\x78\xe6\x25\x64\xa4\x0c\xcd\xfe\x85\xf1\x67\xda\x84\xbb\x86\x4b\x8e\x77\x8f\xd6\x3e\x2f\xe8\x33\x8e\xd5\x1d\xb4\xf1\x68\xb0\x29\x09\xe7\xa0\xd9\xd8\x0c\x44\x2f\x4a\xd3\xbd\x91\x7a\xb9\x7d\xc5\xec\xb4\x6b\xcf\x3b\xfc\xe7\xbf\x93\xdd\xd1\x47\x63\xa8\x8a\x94\x7d\x3a\x7c\x5c\xbf\x68\xc2\xde\xbd\x9e\xa7\xbf\xbd\x6a\x12\x7e\xff\x63\xd2\x08\xa6\xec\x6b\xf7\x4e\xae\x83\xff\x0b\x00\x00\xff\xff\x4c\x08\x0b\x20\x97\x18\x00\x00")

func aroOpenshiftIo_clustersYamlBytes() ([]byte, error) {
	return bindataRead(
		_aroOpenshiftIo_clustersYaml,
		"aro.openshift.io_clusters.yaml",
	)
}

func aroOpenshiftIo_clustersYaml() (*asset, error) {
	bytes, err := aroOpenshiftIo_clustersYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "aro.openshift.io_clusters.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _masterDeploymentYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x84\x52\xc1\x6e\xdb\x30\x0c\xbd\xfb\x2b\x88\xde\xdd\xa4\xb7\x42\xb7\x62\x0d\x7a\x19\x82\x62\x59\x77\x67\x64\x26\x16\x22\x8b\x02\x49\x07\x75\xbf\x7e\x10\x92\x28\xce\x0a\x64\x3a\x19\x7c\x8f\xef\x3d\xd2\xc4\x1c\xfe\x90\x68\xe0\xe4\x00\x73\xd6\xc5\xf1\xa9\x39\x84\xd4\x39\x78\xa5\x1c\x79\x1a\x28\x59\x33\x90\x61\x87\x86\xae\x01\x88\xb8\xa5\xa8\xe5\x0b\x4a\x83\x03\x14\x6e\x39\x93\xa0\xb1\xb4\x03\xaa\x91\x34\x00\x09\x07\xba\x87\x69\x46\x4f\x0e\x38\x53\xd2\x3e\xec\xac\xc5\xaf\x51\xa8\x92\x1b\xcd\xe4\x8b\x89\x50\x8e\xc1\xa3\x3a\x78\x6a\x00\x94\x22\x79\x63\x39\xd9\x0f\x68\xbe\xff\x39\xcb\x73\x37\x91\x9a\xa0\xd1\x7e\x3a\x51\x85\x63\x0c\x69\xff\x91\x3b\x34\xba\x74\x0f\xf8\xb9\x19\x65\x4f\x27\xb3\x73\xe5\x23\xe1\x11\x43\xc4\x6d\x24\x07\xcb\x06\xc0\x68\xc8\xb1\x76\xcd\x77\x53\x5e\xbc\xc9\x73\x37\x11\xc0\x65\xca\xf2\x3c\x27\xc3\x90\x48\x6a\x73\x0b\x9e\x87\x01\x53\x77\x55\x6b\x8b\xd4\x55\x5b\xf6\x3a\xc7\xea\xf6\xae\xa5\x99\x59\x79\x61\xc0\x32\xde\xdb\x6a\xbd\xfa\xf5\xf2\x7b\xf5\x5a\x81\xef\xff\xab\x42\x99\xc5\x6e\x6c\x6a\xd2\x77\x16\x73\xf0\xbc\x7c\x5e\x56\xf4\xa2\xd4\x9b\xe5\x5a\x8c\xe1\x48\x89\x54\xdf\x85\xb7\xe4\x66\xdc\xc2\x7a\x23\x9b\x97\x00\x32\x5a\xef\x60\xd1\x13\x46\xeb\xbf\x16\x42\xd8\x4d\xb7\x84\x7f\x6d\x13\x77\xb4\xb9\x39\x8d\x4b\xb5\x15\x8e\xf4\x78\x18\xb7\x24\x89\x8c\xf4\x31\xf0\xe2\xb4\x12\x07\x0f\x0f\x67\xaa\x92\x1c\x83\xa7\x17\xef\x79\x4c\xb6\xbe\x73\xb9\xdf\xd9\xf7\x98\x59\x02\x4b\xb0\xe9\x47\x44\xd5\x93\xac\x4e\x6a\x34\xb4\x3e\x8e\x85\xd7\x7a\x09\x16\x3c\xc6\x73\x83\x71\x2c\x3a\x81\xd3\xec\x06\x0e\x34\xb9\xff\xcc\x52\x47\xbe\xe4\x70\xb0\xfa\x0c\x6a\x5a\x01\xda\xed\xc8\x9b\x83\x35\x6f\x7c\x4f\xdd\x18\xa9\xf9\x1b\x00\x00\xff\xff\x57\x5c\x5d\xa2\xfa\x03\x00\x00")

func masterDeploymentYamlBytes() ([]byte, error) {
	return bindataRead(
		_masterDeploymentYaml,
		"master/deployment.yaml",
	)
}

func masterDeploymentYaml() (*asset, error) {
	bytes, err := masterDeploymentYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "master/deployment.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _masterRolebindingYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\x8e\xb1\x4e\x03\x31\x0c\x40\x77\x7f\x85\x7f\x20\x87\xd8\x50\x36\x60\x60\x2f\x12\xbb\x9b\xb8\xd4\xf4\x62\x47\x8e\xd3\xa1\x5f\x8f\xaa\xa2\x5b\x90\x6e\xb5\xdf\xf3\x33\x75\xf9\x62\x1f\x62\x9a\xd1\x8f\x54\x16\x9a\x71\x36\x97\x1b\x85\x98\x2e\x97\x97\xb1\x88\x3d\x5d\x9f\xe1\x22\x5a\x33\xbe\xaf\x73\x04\xfb\xc1\x56\x7e\x13\xad\xa2\xdf\xd0\x38\xa8\x52\x50\x06\x44\xa5\xc6\x19\xc9\x2d\x59\x67\xa7\x30\x4f\x8d\xee\x02\xb8\xad\x7c\xe0\xd3\x1d\xa2\x2e\x1f\x6e\xb3\xef\x04\x01\xf1\x5f\x6f\x3b\x5f\x1e\xb3\x44\xb5\x89\xc2\x98\xc7\x1f\x2e\x31\x32\xa4\x3f\xe7\x93\xfd\x2a\x85\x5f\x4b\xb1\xa9\xb1\xfb\xd5\x63\x37\x3a\x15\xce\x68\x9d\x75\x9c\xe5\x14\x89\x6e\xd3\x79\x83\xe1\x37\x00\x00\xff\xff\x4f\x98\xa4\x7c\x24\x01\x00\x00")

func masterRolebindingYamlBytes() ([]byte, error) {
	return bindataRead(
		_masterRolebindingYaml,
		"master/rolebinding.yaml",
	)
}

func masterRolebindingYaml() (*asset, error) {
	bytes, err := masterRolebindingYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "master/rolebinding.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _masterServiceYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x74\x8d\x41\xca\xc2\x40\x0c\x46\xf7\x73\x8a\x5c\x60\xa0\xff\xae\xcc\x29\x7e\x10\xdc\x87\xe9\xa7\x1d\xb4\x93\x90\xc4\x2e\x3c\xbd\xd4\x16\x5d\xb9\x0b\xef\x7b\xbc\xb0\xb6\x33\xcc\x9b\xf4\x42\xeb\x5f\xba\xb5\x3e\x15\x3a\xc1\xd6\x56\x91\x16\x04\x4f\x1c\x5c\x12\x51\xe7\x05\x85\xd8\x24\x8b\xc2\x38\xc4\xf2\xc2\x1e\xb0\x63\x73\xe5\x8a\x42\xa2\xe8\x3e\xb7\x4b\x64\x7e\x3e\x0c\x1f\x39\xb9\xa2\x6e\x1d\xc7\x1d\x35\xc4\xb6\x9b\x88\x55\x7f\x45\x55\x2c\x7c\xb7\xf2\xf1\x7d\x8e\xd0\x37\xd8\xd7\x42\xe3\x30\x0e\x07\x08\xb6\x2b\xe2\xff\x8b\x5f\x01\x00\x00\xff\xff\x10\x70\xf6\x36\xda\x00\x00\x00")

func masterServiceYamlBytes() ([]byte, error) {
	return bindataRead(
		_masterServiceYaml,
		"master/service.yaml",
	)
}

func masterServiceYaml() (*asset, error) {
	bytes, err := masterServiceYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "master/service.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _masterServiceaccountYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x3c\xca\x31\x8e\x02\x31\x0c\x05\xd0\x3e\xa7\xf0\x05\x52\x6c\xeb\x6e\xcf\x80\x44\xff\x95\xf9\x08\x0b\xc5\x8e\x1c\xcf\x14\x9c\x9e\x06\x51\xbf\x87\x65\x77\xe6\xb6\x70\x95\xeb\xaf\xbd\xcc\x0f\x95\x1b\xf3\xb2\xc1\xff\x31\xe2\xf4\x6a\x93\x85\x03\x05\x6d\x22\x8e\x49\x15\x64\xf4\x58\x4c\x54\x64\x9f\xd8\xc5\xfc\xda\x5e\x18\x54\x89\x45\xdf\x4f\x7b\x54\xc7\xfb\x4c\xfe\x72\xfb\x04\x00\x00\xff\xff\xe4\xf5\x04\x25\x70\x00\x00\x00")

func masterServiceaccountYamlBytes() ([]byte, error) {
	return bindataRead(
		_masterServiceaccountYaml,
		"master/serviceaccount.yaml",
	)
}

func masterServiceaccountYaml() (*asset, error) {
	bytes, err := masterServiceaccountYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "master/serviceaccount.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _namespaceYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x44\xca\xb1\x0d\x02\x31\x0c\x05\xd0\x3e\x53\x58\xd7\x07\x44\x9b\x21\x28\xe9\xbf\x2e\x1f\x61\x41\xec\x28\x36\x14\x4c\x8f\xa8\xae\x7f\x98\x7a\xe3\x0a\x75\x6b\xf2\xb9\x94\xa7\x5a\x6f\x72\xc5\x60\x4c\xec\x2c\x83\x89\x8e\x44\x2b\x22\x86\xc1\x26\x3e\x69\xf1\xd0\x7b\x56\x7c\xdf\x8b\xd5\x27\x17\xd2\x57\x11\x81\x99\x27\x52\xdd\xe2\xef\xe5\xb0\x27\xf5\xb3\x79\x67\x0d\xbe\xb8\xa7\xaf\x26\xdb\x56\x7e\x01\x00\x00\xff\xff\xc1\xaf\xa6\x4c\x7c\x00\x00\x00")

func namespaceYamlBytes() ([]byte, error) {
	return bindataRead(
		_namespaceYaml,
		"namespace.yaml",
	)
}

func namespaceYaml() (*asset, error) {
	bytes, err := namespaceYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "namespace.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _workerDeploymentYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\x52\xcb\x6e\xdb\x40\x0c\xbc\xeb\x2b\x88\xdc\x15\x27\xb7\x60\x6f\x41\x63\xe4\x52\x04\x45\xd3\xf4\x4e\xaf\xa6\xd6\xc2\xfb\x02\x49\xbb\x55\xbe\xbe\x10\x64\xcb\x32\x02\x88\x27\x61\x38\x9c\x19\x2e\xc5\x35\xfc\x86\x68\x28\xd9\x11\xd7\xaa\x9b\xd3\x63\x73\x08\xb9\x73\xf4\x82\x1a\xcb\x90\x90\xad\x49\x30\xee\xd8\xd8\x35\x44\x91\x77\x88\x3a\x7e\xd1\x38\xe0\x88\xa5\xb4\xa5\x42\xd8\x8a\xb4\x7f\x8b\x1c\x20\x0d\x51\xe6\x84\xb5\x9e\x56\xf6\x70\x54\x2a\xb2\xf6\xe1\x8f\xb5\xfc\x79\x14\xcc\xe4\x46\x2b\xfc\x68\x22\xa8\x31\x78\x56\x47\x8f\x0d\x91\x22\xc2\x5b\x91\xc9\x3e\xb1\xf9\xfe\xfb\x22\xcf\x6a\x22\x35\x61\xc3\x7e\x98\xa8\x52\x62\x0c\x79\xff\x51\x3b\x36\x5c\xa6\x13\xff\x7b\x3f\xca\x1e\x93\xd9\x19\xf9\xc8\x7c\xe2\x10\x79\x17\xe1\xe8\xa1\x21\x32\xa4\x1a\xe7\xa9\xe5\xdb\x8c\x15\x6f\xf2\xac\x26\x22\xba\x6c\x39\x96\x2f\xd9\x38\x64\xc8\x3c\xdc\x92\x2f\x29\x71\xee\xae\x6a\xed\x28\x75\xd5\x96\xbd\x2e\x7b\xf3\xeb\x5d\xa1\x85\xd9\x58\x21\xf1\xb8\xde\xeb\xf6\x6d\xfb\xf3\xf9\xd7\xf6\x65\x6e\x7c\xbd\xd7\xdc\x8a\xe1\x84\x0c\xd5\x1f\x52\x76\xb8\xda\x11\xf5\x66\xf5\x15\xb6\x84\x88\x2a\x5b\xef\x68\xd3\x83\xa3\xf5\x9f\x1b\x01\x77\xc3\x2d\xa1\x88\x39\x7a\x7a\x78\x7a\x38\xc3\xb9\x74\x78\xbf\x39\xec\x05\x6d\xa5\x44\xdc\x1f\x8e\x3b\x48\x86\x41\xef\x43\xd9\x4c\x0b\x39\xba\xbb\x3b\x53\x15\x72\x0a\x1e\xcf\xde\x97\x63\xb6\xb7\x95\xff\xee\x2b\x7b\x8d\x59\x25\x14\x09\x36\x7c\x8b\xac\x3a\xc9\xea\xa0\x86\xd4\xfa\x78\x54\x83\xb4\x5e\x82\x05\xcf\xb1\xf9\x1f\x00\x00\xff\xff\x4f\x57\x4a\x02\x45\x03\x00\x00")

func workerDeploymentYamlBytes() ([]byte, error) {
	return bindataRead(
		_workerDeploymentYaml,
		"worker/deployment.yaml",
	)
}

func workerDeploymentYaml() (*asset, error) {
	bytes, err := workerDeploymentYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "worker/deployment.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _workerRoleYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xa4\x8e\xb1\x6e\x2c\x31\x08\x45\x7b\xbe\x82\x1f\xb0\x57\xaf\x7b\x72\x9b\x22\x7d\x14\xa5\x67\x3d\x24\x83\xc6\x63\x2c\xc0\xbb\x52\xbe\x3e\x9a\xd9\x6d\x53\xa5\xe2\x0a\x1d\x0e\x17\x52\x4a\x40\x43\x3e\xd8\x5c\xb4\x17\xb4\x2b\xd5\x4c\x33\x56\x35\xf9\xa6\x10\xed\x79\xfb\xef\x59\xf4\x72\xfb\x07\x9b\xf4\xa5\xe0\x4b\x9b\x1e\x6c\x6f\xda\x18\x76\x0e\x5a\x28\xa8\x00\x62\x35\x3e\x0f\xde\x65\x67\x0f\xda\x47\xc1\x3e\x5b\x03\xc4\x4e\x3b\x17\x24\xd3\xa4\x83\x8d\x42\x2d\xdd\xd5\x36\x36\xb0\xd9\xd8\x0b\x24\xa4\x21\xaf\xa6\x73\xf8\x61\x4a\x07\x9b\x75\x70\xf7\x55\x3e\x23\x8b\x02\xa2\xb1\xeb\xb4\xca\x4f\xa2\x3e\x5a\x38\x20\xde\xd8\xae\xcf\xed\x17\xc7\x39\x9b\xf8\x23\xdc\x29\xea\xfa\x17\xff\xc5\x83\x62\xfe\xf2\x66\x9c\xf6\x23\xcd\xb1\x50\x30\xfc\x04\x00\x00\xff\xff\x30\x78\x19\x41\x50\x01\x00\x00")

func workerRoleYamlBytes() ([]byte, error) {
	return bindataRead(
		_workerRoleYaml,
		"worker/role.yaml",
	)
}

func workerRoleYaml() (*asset, error) {
	bytes, err := workerRoleYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "worker/role.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _workerRolebindingYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\x8d\x31\x6e\xc3\x30\x0c\x45\x77\x9d\x82\x17\x90\x8b\x6e\x85\xb6\xb6\x43\x77\x17\xe8\x4e\xcb\x74\xcd\xda\x26\x05\x8a\x72\x01\x9f\x3e\x08\x12\x64\x09\xe0\xf9\xbf\xf7\x1f\x16\xfe\x21\xab\xac\x92\xc0\x06\xcc\x1d\x36\x9f\xd5\xf8\x40\x67\x95\x6e\x79\xab\x1d\xeb\xcb\xfe\x1a\x16\x96\x31\xc1\xe7\xda\xaa\x93\xf5\xba\xd2\x07\xcb\xc8\xf2\x1b\x36\x72\x1c\xd1\x31\x05\x00\xc1\x8d\x12\xa0\x69\xd4\x42\x86\xae\x16\xff\xd5\x16\xb2\x60\xba\x52\x4f\xd3\x15\xc2\xc2\x5f\xa6\xad\x9c\x04\x03\xc0\x53\xef\xf4\xbe\xb6\xe1\x8f\xb2\xd7\x14\xe2\xdd\xfc\x26\xdb\x39\xd3\x7b\xce\xda\xc4\x4f\xe5\xdb\x56\x0b\x66\x4a\xa0\x85\xa4\xce\x3c\x79\xc4\xa3\x19\x3d\xe0\x70\x09\x00\x00\xff\xff\x73\xce\x57\x9b\x2a\x01\x00\x00")

func workerRolebindingYamlBytes() ([]byte, error) {
	return bindataRead(
		_workerRolebindingYaml,
		"worker/rolebinding.yaml",
	)
}

func workerRolebindingYaml() (*asset, error) {
	bytes, err := workerRolebindingYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "worker/rolebinding.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _workerServiceaccountYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x3c\xca\x31\x8a\xc3\x40\x0c\x05\xd0\x7e\x4e\xa1\x0b\x4c\xb1\xad\xba\x3d\x43\x20\xfd\x67\xfc\x43\x84\xb1\x34\x68\x64\x07\x72\xfa\x34\x21\xf5\x7b\x98\x76\x67\x2e\x0b\x57\xb9\xfe\xda\x6e\xbe\xa9\xdc\x98\x97\x0d\xfe\x8f\x11\xa7\x57\x3b\x58\xd8\x50\xd0\x26\xe2\x38\xa8\x82\x8c\x1e\x93\x89\x8a\xec\xaf\xc8\x9d\xf9\xb5\x35\x31\xa8\x12\x93\xbe\x9e\xf6\xa8\x8e\xf7\x99\xfc\xe5\xf6\x09\x00\x00\xff\xff\xe3\x3c\x43\x66\x70\x00\x00\x00")

func workerServiceaccountYamlBytes() ([]byte, error) {
	return bindataRead(
		_workerServiceaccountYaml,
		"worker/serviceaccount.yaml",
	)
}

func workerServiceaccountYaml() (*asset, error) {
	bytes, err := workerServiceaccountYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "worker/serviceaccount.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"aro.openshift.io_clusters.yaml": aroOpenshiftIo_clustersYaml,
	"master/deployment.yaml":         masterDeploymentYaml,
	"master/rolebinding.yaml":        masterRolebindingYaml,
	"master/service.yaml":            masterServiceYaml,
	"master/serviceaccount.yaml":     masterServiceaccountYaml,
	"namespace.yaml":                 namespaceYaml,
	"worker/deployment.yaml":         workerDeploymentYaml,
	"worker/role.yaml":               workerRoleYaml,
	"worker/rolebinding.yaml":        workerRolebindingYaml,
	"worker/serviceaccount.yaml":     workerServiceaccountYaml,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"aro.openshift.io_clusters.yaml": {aroOpenshiftIo_clustersYaml, map[string]*bintree{}},
	"master": {nil, map[string]*bintree{
		"deployment.yaml":     {masterDeploymentYaml, map[string]*bintree{}},
		"rolebinding.yaml":    {masterRolebindingYaml, map[string]*bintree{}},
		"service.yaml":        {masterServiceYaml, map[string]*bintree{}},
		"serviceaccount.yaml": {masterServiceaccountYaml, map[string]*bintree{}},
	}},
	"namespace.yaml": {namespaceYaml, map[string]*bintree{}},
	"worker": {nil, map[string]*bintree{
		"deployment.yaml":     {workerDeploymentYaml, map[string]*bintree{}},
		"role.yaml":           {workerRoleYaml, map[string]*bintree{}},
		"rolebinding.yaml":    {workerRolebindingYaml, map[string]*bintree{}},
		"serviceaccount.yaml": {workerServiceaccountYaml, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
