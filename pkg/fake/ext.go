package fake

import (
	"moefile/pkg/randcjk"
)

var AvailExt = []string{
	"md", "txt", "rst", "rtf",
	"doc", "docx", "xls", "xlsx", "ppt", "pptx", "pdf",
	"odt", "ods", "odp", "odg", "odf",
	"epub", "mobi", "djvu", "fb2",

	"jpg", "jpeg", "png", "gif", "bmp", "tif", "tiff", "svg", "ico",
	"webp", "avif", "heif", "heic",
	"raw", "dng", "nef", "arw", "cr2", "cr3",

	"wav", "flac", "alac", "dsd", "ape",
	"mp3", "aac", "ogg", "m4a", "opus",
	"wma", "aiff", "amr", "mka", "mks",
	"mid", "midi",

	"mp4", "mkv", "webm", "avi", "mov", "wmv", "flv", "f4v", "f4p", "f4a", "f4b",
	"m4v", "3gp", "3g2", "ogv", "ogg", "rm", "rmvb", "m2v", "m4p", "m4b",
	"mpg", "mpeg", "m2ts", "mts", "vob",

	"c", "cpp", "h", "hpp", "cs", "java", "js", "ts", "jsx", "tsx", "html", "css", "scss", "sass", "less",
	"py", "rb", "php", "go", "rs", "swift", "kt", "sh", "bash", "zsh", "fish", "ps1",
	"pl", "lua", "r", "dart", "scala", "groovy", "tsv", "csv", "yaml", "yml", "toml", "sql",

	"json", "xml", "yaml", "yml", "toml", "ini", "cfg", "conf", "properties", "csv", "sqlite", "db",

	"zip", "tar", "gz", "bz2", "xz", "7z", "rar", "zst",
	"pkg", "deb", "rpm", "apk", "ipa", "vhd", "vmdk", "qcow2",

	"ova", "ovf", "iso", "img",
	"exe", "msi", "appx", "elf",
	"dll", "so", "dylib", "a", "lib", "bin",

	"unknown",
}

func FakeExt() string {
	return AvailExt[randcjk.RRange(0, len(AvailExt))]
}
