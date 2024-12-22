package main

import (
	"moefile/internal/viteutil"
)

func generatePlayerProd() {
	println("Generating BUILD_ITEM=player MODE=prod")
	viteutil.WriteAppSrc("player", false, "")
}

func generatePlayerDev() {
	println("Generating BUILD_ITEM=player MODE=dev")
	viteutil.WriteAppSrc("player", true, "")
}
