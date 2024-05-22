package server

import (
	"bytes"
	"httphere/asset"
	"httphere/utils"
	"io"
	"net/http"
)

func ServeQr(w http.ResponseWriter, r *http.Request) {
	fi, _ := asset.Dist.Open("dist/qrcode.html")
	content, _ := io.ReadAll(fi)

	viewUrl, uploadUrl, _ := utils.GetUrls()

	cnt1 := bytes.Replace(content, []byte("---viewUrl---"), []byte(viewUrl), -1)
	cnt2 := bytes.Replace(cnt1, []byte("---uploadUrl---"), []byte(uploadUrl), -1)
	w.Write(cnt2)
	return
}
