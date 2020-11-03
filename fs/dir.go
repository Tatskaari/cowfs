package fs

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

//TODO this should probably use constructors and such. Raw struct init here is risky.
func (mooFS *FS) fromPaths(paths []string) *Dir {
	d := new(Dir)
	d.path = ""
	d.entries = make(map[string]fs.Node, len(paths))
	d.mooFS = mooFS
	d.attr = &fuse.Attr{
		Inode: fs.GenerateDynamicInode(0, ""),
	}

	for _, p := range paths {
		f := &File{
			path:      filepath.Base(p),
			fromPath:  p,
			writeable: false,
			mooFS: mooFS,
		}
		f.attr = new(fuse.Attr)
		err := fileAttr(f, f.attr)
		f.attr.Inode = fs.GenerateDynamicInode(d.attr.Inode, f.path)
		if err != nil {
			panic(err)
		}
		d.entries[filepath.Base(p)] = f
	}

	return d
}

// Dir implements both Node and Handle for the root directory.
type Dir struct{
	mooFS *FS
	attr *fuse.Attr
	path string
	entries map[string]fs.Node
}

func (mooFS *FS) newDir(parent *Dir, name string, mode os.FileMode) *Dir {
	return &Dir{
		path: filepath.Join(parent.path, name),
		attr: &fuse.Attr{
			Inode: fs.GenerateDynamicInode(parent.attr.Inode, name),
			Mode: mode,
		},
		entries: map[string]fs.Node{},
		mooFS: mooFS,
	}
}

func (d *Dir) Loc() string {
	return filepath.Join(tmpDir, d.path)
}

func (d *Dir) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Mode = os.ModeDir | 0775
	d.attr = a
	return nil
}

func (d *Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	if f, ok := d.entries[name]; ok {
		return f, nil
	}
	return nil, syscall.ENOENT
}

func (d *Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	// TODO cache this?
	entries := make([]fuse.Dirent, 0, len(d.entries))

	for k, v := range d.entries {
		switch f := v.(type) {
		case *File:
			entries = append(entries, fuse.Dirent{
				Type:  fuse.DT_File,
				Inode: f.attr.Inode,
				Name:  k,
			})
		case *Dir:
			entries = append(entries, fuse.Dirent{
				Type:  fuse.DT_Dir,
				Inode: f.attr.Inode,
				Name:  k,
			})
		default:
			panic(fmt.Errorf("unhandled file type %v", f))
		}
	}
	return entries, nil
}

func (d *Dir) Remove(ctx context.Context, req *fuse.RemoveRequest) error {
	if e, ok := d.entries[req.Name]; !ok {
		return syscall.ENOENT
	} else {
		if _, ok := e.(*Dir); ok && !req.Dir {
			return syscall.EISDIR
		}
		delete(d.entries, req.Name)
		return nil
	}
}

func (d *Dir) Create(ctx context.Context, req *fuse.CreateRequest, resp *fuse.CreateResponse) (fs.Node, fs.Handle, error) {
	child := d.mooFS.newFile(d, req.Name, req.Mode)
	path := child.Loc()

	f, err := os.OpenFile(path, int(req.Flags), req.Mode.Perm())
	if err != nil {
		return nil, nil, err
	}
	child.openFile = f

	d.entries[req.Name] = child
	return child, child, nil
}

func (d *Dir) Mkdir(ctx context.Context, req *fuse.MkdirRequest) (fs.Node, error) {
	path := filepath.Join(tmpDir, d.mooFS.mountPoint, d.path, req.Name)
	if err := os.Mkdir(path, req.Mode.Perm()); err != nil {
		return nil, err
	}
	child := d.mooFS.newDir(d, req.Name, req.Mode)
	d.entries[req.Name] = child
	return child, nil
}
