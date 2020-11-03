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

// File implements both Node and Handle for a file that actually exists in our file tree
type File struct{
	mooFS *FS

	path string
	fromPath string
	writeable bool
	attr *fuse.Attr

}

type FileHandle struct {
	*os.File

	file *File
	wasWriteable bool

	mode int
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

// prepareForWrite will copy the file to the tmp directory if the file is not currently writeable
func (f *File) prepareForWrite() error {
	// If we're not writable, copy the file to the tmp fs
	if !f.writeable {
		writeFile, err := os.Create(f.Loc())
		if err != nil {
			return err
		}
		defer writeFile.Close()

		origFile, err := os.Open(f.fromPath)
		if err != nil {
			return err
		}
		defer origFile.Close()

		if _, err := io.Copy(origFile, writeFile); err != nil {
			return err
		}
		f.writeable = true
	}
	return nil
}

func (f *File) Fsync(ctx context.Context, req *fuse.FsyncRequest) error {
	// TODO: I think we should be able to ignore this however it might be best to keep track of all file handles
	// and get them to flush here.
	return nil
}

func (f *File) Open(ctx context.Context, req *fuse.OpenRequest, resp *fuse.OpenResponse) (fs.Handle, error) {
	var osFile *os.File
	var err error

	if req.Flags.IsReadWrite() || req.Flags.IsWriteOnly() {
		if err := f.prepareForWrite(); err != nil {
			return nil, err
		}
		osFile, err = os.OpenFile(f.Loc(), os.O_RDWR, 0)
	} else {
		osFile, err = os.Open(f.Loc())
	}

	return &FileHandle{
		File: osFile,
		file: f,
		wasWriteable: f.writeable,
	}, err
}

func (f *FileHandle) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error {
	// If the file has become writable since we opened it, open the new file
	if !f.wasWriteable && f.file.writeable {
		f.Close()
		newFile, err := os.OpenFile(f.file.Loc(), f.mode, 0)
		if err != nil {
			return err
		}
		f.File = newFile
	}
	resp.Data = make([]byte, req.Size)
	if _, err := f.ReadAt(resp.Data, req.Offset); err != io.EOF {
		return err
	}

	return nil
}

func (f *FileHandle) Flush(ctx context.Context, req *fuse.FlushRequest) error {
	// ignore the error as go will complain on multiple closes however FUSE allows this
	f.Close()
	return nil
}

func (f *FileHandle) Write(ctx context.Context, req *fuse.WriteRequest, resp *fuse.WriteResponse) error {
	n, err := f.WriteAt(req.Data, req.Offset)
	if err != nil {
		return err
	}
	resp.Size = n
	f.file.attr.Size += uint64(n)
	return nil
}