package main

import (
	"fmt"

	"moefile/internal/viteutil"
	"moefile/pkg/dto"
	"moefile/pkg/fake"
	rcjk "moefile/pkg/randcjk"
)

func generateIndexProd() {
	println("Generating BUILD_ITEM=index MODE=prod")
	viteutil.WriteAppSrc("index", false, "")
}

func generateIndexDev() {
	println("Generating BUILD_ITEM=index MODE=dev")
	fakeName := fmt.Sprintf("MoeFile %s", rcjk.RString(8, rcjk.ASCIIURL))
	fakeFolder := fmt.Sprintf("%s/%s/%s",
		rcjk.RString(8, rcjk.ASCIIURL), rcjk.RString(16, rcjk.ASCII), rcjk.RString(12, rcjk.CJK))
	data := dto.NewFSDirInfo(fakeName, fakeFolder)

	// 10 ~ 20 folders with ASCII URL name
	for i := 0; i < rcjk.RRange(10, 20); i++ {
		data.AddFSFile(&fake.FakeFile{
			NameLength: 8,
			NameFlag:   rcjk.ASCIIURL,
			MaxSize:    4096,
			Dir:        fake.FakeIsDir,
			BackYears:  3,
		})
	}

	// 5 ~ 10 folders with CJK name
	for i := 0; i < rcjk.RRange(5, 10); i++ {
		data.AddFSFile(&fake.FakeFile{
			NameLength: 16,
			NameFlag:   rcjk.CJK,
			MaxSize:    4096,
			Dir:        fake.FakeIsDir,
			BackYears:  5,
		})
	}

	// 1 ~ 5 files with Lang Chinese name
	for i := 0; i < rcjk.RRange(1, 5); i++ {
		data.AddFSFile(&fake.FakeFile{
			NameLength: 255,
			NameFlag:   rcjk.CJKChinese,
			MinSize:    1024,
			MaxSize:    1024 * 1024 * 1024 * 10,
			Dir:        fake.FakeIsFile,
			BackYears:  10,
		})
	}

	// 1 ~ 5 files OR folder with ASCII Symbol name
	for i := 0; i < rcjk.RRange(1, 5); i++ {
		data.AddFSFile(&fake.FakeFile{
			NameLength: 24,
			NameFlag:   rcjk.ASCIISymbol,
			MinSize:    1024,
			MaxSize:    1024 * 1024 * 1024 * 10,
			Dir:        fake.FakeRandomDir,
			BackYears:  10,
		})
	}

	// 10 ~ 20 files with ANY name
	for i := 0; i < rcjk.RRange(10, 20); i++ {
		data.AddFSFile(&fake.FakeFile{
			NameLength: 12,
			NameFlag:   rcjk.CJKASCII,
			MaxSize:    1024 * 1024,
			Dir:        fake.FakeIsFile,
			BackYears:  20,
		})
	}

	dataXML, err := data.ToS3XMLWithoutXSLT(true)
	if err != nil {
		panic(err)
	}

	viteutil.WriteAppSrc("index", true, string(dataXML))
}
