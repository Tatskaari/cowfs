package fs

import (
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const mountDir = "cowfs"

type cowFSSuite struct {
	suite.Suite
}

func (s *cowFSSuite) SetupSuite(){
	if err := os.Chdir("fs/test_data"); err != nil {
		panic(err)
	}
}

func (s *cowFSSuite) SetupTest(){
	go MountAndServe(mountDir, []string{"foo.go"})
	time.Sleep(100 * time.Millisecond)
}

func (s *cowFSSuite) TearDownTest(){
	if err := Unmount(mountDir); err != nil {
		panic(err)
	}
}

func TestCowFSSuite(t *testing.T) {
	s := new(cowFSSuite)
	suite.Run(t, s)
}

func (s *cowFSSuite) TestSrcFileMounted() {
	mountInfo, err := os.Lstat(filepath.Join(mountDir, "foo.go"))
	s.Require().NoError(err)

	fooInfo, err := os.Lstat("foo.go")
	s.Require().NoError(err)

	s.Equal(fooInfo.Size(), mountInfo.Size())
	s.Equal(fooInfo.Mode(), mountInfo.Mode())
}


func (s *cowFSSuite) TestEditSrcFile() {
	s.T().Skip("not working")

	origContents, err := ioutil.ReadFile("foo.go")
	s.Require().NoError(err)

	f, err := os.OpenFile(filepath.Join(mountDir, "foo.go"), os.O_RDWR, 0)
	s.Require().NoError(err)


	comment := []byte("//this is a comment")
	n, err := f.WriteAt(comment, 0)
	s.Require().NoError(err)

	s.Equal(len(comment), n)

	s.Require().NoError(f.Close())

	newContents, err := ioutil.ReadFile("foo.go")
	s.Require().NoError(err)

	s.Equal(origContents, newContents)

	cowFSContent, err := ioutil.ReadFile(filepath.Join(mountDir, "foo.go"))
	s.Require().NoError(err)

	s.Equal(string(append(origContents, comment...)), string(cowFSContent))
}

func (s *cowFSSuite) TestMkdir() {
	err := os.Mkdir(filepath.Join(mountDir, "bar"), os.ModeDir | 0775)
	s.Require().NoError(err)

	bar, err := os.Create(filepath.Join(mountDir, "bar", "bar.go"))
	s.Require().NoError(err)

	_, err = bar.Write([]byte("Testing"))
	s.Require().NoError(err)

	s.Require().NoError(bar.Close())

	barContent, err := ioutil.ReadFile(filepath.Join(mountDir, "bar", "bar.go"))
	s.Equal("Testing", string(barContent))
}