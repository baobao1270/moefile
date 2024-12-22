package fake

import (
	"fmt"
	"io/fs"
	"time"

	"moefile/pkg/randcjk"
)

const (
	FakeIsDir     = 0b001
	FakeIsFile    = 0b010
	FakeRandomDir = 0b100
)

type FakeFile struct {
	fs.FileInfo
	NameLength int
	NameFlag   int
	MinSize    int
	MaxSize    int
	Dir        int
	BackYears  int
}

func (f *FakeFile) Name() string {
	name := randcjk.RString(f.NameLength, f.NameFlag)
	if !f.IsDir() {
		return fmt.Sprintf("%s.%s", name, FakeExt())
	}
	return name
}

func (f *FakeFile) Size() int64 {
	return int64(randcjk.RRange(f.MinSize, f.MaxSize))
}

func (f *FakeFile) IsDir() bool {
	if f.Dir&FakeIsDir != 0 {
		return true
	}
	if f.Dir&FakeIsFile != 0 {
		return false
	}
	return randcjk.RRange(0, 2) == 0
}

func (f *FakeFile) ModTime() time.Time {
	return FakeYearsBefore(f.BackYears)
}

func (f *FakeFile) Mode() fs.FileMode {
	return 0
}

func (f *FakeFile) Sys() interface{} {
	return nil
}

func (f *FakeFile) As() fs.FileInfo {
	return f.FileInfo
}
