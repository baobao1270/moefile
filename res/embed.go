package res

import (
	_ "embed"
)

//go:embed xslt.xml
var XSLT string

//go:embed xml_header.xml
var XMLHeader []byte
