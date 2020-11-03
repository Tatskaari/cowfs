package fs

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)


const tmpDir = "/tmp/cowfs"

// FS implements the hello world file system.
type FS struct{
	mountPoint string
	root *Dir
}

func (mooFS *FS) Root() (fs.Node, error) {
	return mooFS.root, nil
}

func MountAndServe(mountPoint string, paths []string) {
	_ = os.MkdirAll(mountPoint, 0755)
	_ = os.MkdirAll(filepath.Join(tmpDir, mountPoint), 0755)

	c, err := fuse.Mount(
		mountPoint,
		fuse.FSName("cowfs"),
		fuse.Subtype("cowfs"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	mooFS := &FS{mountPoint: mountPoint}
	mooFS.root = mooFS.fromPaths(paths)

	fuse.Debug = func(msg interface{}) {
		logLine := fmt.Sprintf("%v", msg)
		if !strings.Contains(logLine, "attr") {
			fmt.Println(logLine)
		}
	}

	err = fs.Serve(c, mooFS)
	if err != nil {
		log.Fatal(err)
	}
}

func Unmount(mountPoint string) error {
	err := fuse.Unmount(mountPoint)
	_ = os.RemoveAll(filepath.Join(tmpDir, mountPoint))
	return err
}