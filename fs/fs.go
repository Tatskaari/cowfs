package fs

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"gopkg.in/op/go-logging.v1"
	"os"
	"path/filepath"
)

var log = logging.MustGetLogger("moofs/fs")

const tmpDir = "/tmp/moofs"

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
	_ = os.MkdirAll(filepath.Join("/tmp/cowfs", mountPoint), 0755)

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
		log.Infof("%v", msg)
	}

	err = fs.Serve(c, mooFS)
	if err != nil {
		log.Fatal(err)
	}
}

func Unmount(mountPoint string) error {
	err := fuse.Unmount(mountPoint)
	_ = os.RemoveAll(filepath.Join("/tmp/cowfs", mountPoint))
	return err
}