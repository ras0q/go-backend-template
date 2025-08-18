package frontend

import (
	"embed"
	"io/fs"
)

// NOTE: go:embed でフロントエンドをGoのバイナリに埋め込んで配信している
// 同様にembed.FSを増やすことで複数のフロントエンドプロジェクトを同時に埋め込むことが可能

//go:embed app-ui/dist
var uiDist embed.FS

var UI, _ = fs.Sub(uiDist, "app-ui/dist")
