package main

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed web/*
var uiAssets embed.FS

func AssetFile() http.FileSystem {
	subfs, err := fs.Sub(uiAssets, "web")
	if err != nil {
		panic(err)
	}
	return http.FS(subfs)
}
