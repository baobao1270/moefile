package main

import (
	"flag"
)

func main() {
	dev := flag.Bool("dev", false, "generate test data for development")
	prod := flag.Bool("prod", false, "generate test data for production")
	flag.Parse()
	println("Generating test data...")

	if *dev && *prod {
		panic("Cannot generate both DEV and PROD data")
	}

	if *dev {
		println("Generating MODE=dev")
		generateIndexDev()
		generatePlayerDev()
	}

	if *prod {
		println("Generating MODE=prod")
		generateIndexProd()
		generatePlayerProd()
	}
}
