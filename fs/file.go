package fs

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Locatable interface {
	Loc() string
}

func fileAttr(loc Locatable, a *fuse.Attr) error {
	info, err := os.Lstat(loc.Loc())
	if err != nil {
		panic(fmt.Errorf("failed to get file info for %s: %w", loc.Loc(), err))
	}
	a.Mode = info.Mode()
	a.Size = uint64(info.Size())
	return nil
}

// RealFile implements both Node and Handle for a file that actually exists in our file tree
type File struct{
	mooFS *FS

	path string
	fromPath string
	writeable bool
	attr *fuse.Attr

	openFile *os.File
}

func (mooFS *FS) newFile(parent *Dir, name string, mode os.FileMode) *File {
	return &File {
		path: filepath.Join(parent.path, name),
		writeable: true,
		attr: &fuse.Attr{
			Mode: mode,
			Inode: fs.GenerateDynamicInode(parent.attr.Inode, name),
		},
		mooFS: mooFS,
	}
}

func (f *File) Loc() string {
	if f.writeable {
		return filepath.Join(tmpDir, f.mooFS.mountPoint, f.path)
	}
	return f.fromPath
}

func (f *File) Attr(ctx context.Context, a *fuse.Attr) error {
	if f.attr != nil {
		a.Mode = f.attr.Mode
		a.Size = f.attr.Size
		return nil
	}
	f.attr = a

	return fileAttr(f, a)
}

func (f *File) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error {
	// TODO this should be on a file handle and the client should've already opened the file by this point
	osFile, err := os.Open(f.Loc())
	if err != nil {
		return err
	}

	resp.Data = make([]byte, req.Size)
	if _, err := osFile.ReadAt(resp.Data, req.Offset); err != io.EOF {
		return err
	}

	return nil
}

func (f *File) openForWrite() (*os.File, error){
	if !f.writeable {
		f.writeable = true
		writeFile, err := os.Create(f.Loc())
		if err != nil {
			return nil, err
		}

		origFile, err := os.Open(f.fromPath)
		if err != nil {
			return nil, err
		}
		defer origFile.Close()

		if _, err := io.Copy(origFile, writeFile); err != nil {
			return nil, err
		}
		f.openFile = writeFile
		return writeFile, nil
	}
	return f.openFile, nil
}

func (f *File) Write(ctx context.Context, req *fuse.WriteRequest, resp *fuse.WriteResponse) error {
	osFile, err := f.openForWrite()
	if err != nil {
		return err
	}

	n, err := osFile.WriteAt(req.Data, req.Offset)
	resp.Size = n
	f.attr.Size += uint64(n)
	return err
}

func (f *File) Fsync(ctx context.Context, req *fuse.FsyncRequest) error {
	return f.openFile.Sync()
}

func (f *File) Flush(ctx context.Context, req *fuse.FlushRequest) error {
	if f.openFile == nil {
		return nil
	}

	err := f.openFile.Close()
	if err != nil {
		return err
	}
	f.openFile = nil
	return nil
}