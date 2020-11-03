package fs

import (
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const mountDir = "moofs"

type mooFSSuite struct {
	suite.Suite
}

func (s *mooFSSuite) SetupSuite(){
	if err := os.Chdir("fs/test_data"); err != nil {
		panic(err)
	}
}

func (s *mooFSSuite) SetupTest(){
	go MountAndServe(mountDir, []string{"foo.go"})
	time.Sleep(100 * time.Millisecond)
}

func (s *mooFSSuite) TearDownTest(){
	if err := Unmount(mountDir); err != nil {
		panic(err)
	}
}

func TestMooFSSuite(t *testing.T) {
	s := new(mooFSSuite)
	suite.Run(t, s)
}

func (s *mooFSSuite) TestSrcFileMounted() {
	mountInfo, err := os.Lstat(filepath.Join(mountDir, "foo.go"))
	s.Require().NoError(err)

	fooInfo, err := os.Lstat("foo.go")
	s.Require().NoError(err)

	s.Equal(fooInfo.Size(), mountInfo.Size())
	s.Equal(fooInfo.Mode(), mountInfo.Mode())
}


func (s *mooFSSuite) TestEditSrcFile() {
	origContents, err := ioutil.ReadFile("foo.go")
	s.Require().NoError(err)

	f, err := os.OpenFile(filepath.Join(mountDir, "foo.go"), os.O_RDWR, 0)
	s.Require().NoError(err)


	comment := []byte("//this is a comment")
	n, err := f.WriteAt(comment, 0)
	s.Require().NoError(err)

	s.Equal(len(comment), n)

	f.Close()

	newContents, err := ioutil.ReadFile("foo.go")
	s.Require().NoError(err)

	s.Equal(origContents, newContents)

	cowFSContent, err := ioutil.ReadFile(filepath.Join(mountDir, "foo.go"))
	s.Require().NoError(err)

	s.Equal(append(comment, origContents...), cowFSContent)

}