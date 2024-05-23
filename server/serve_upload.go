package server

import (
	"fmt"
	"github.com/youwen21/httphere/asset"
	"github.com/youwen21/httphere/conf"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
)

func ServeUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fi, _ := asset.Dist.Open("dist/upload.html")
		content, _ := io.ReadAll(fi)
		w.Write(content)
		return
	} else {
		// Multipart form
		err := r.ParseMultipartForm(32 << 20) // maxMemory 32MB
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		multipartForm := r.MultipartForm
		if multipartForm == nil {
			w.Write([]byte("未找到上传的文件"))
			return
		}

		files := multipartForm.File["file"]
		fmt.Println(files)
		for _, file := range files {
			dst := path.Join(conf.GetRoot(), file.Filename)

			// Upload the file to specific dst.
			err := SaveUploadedFile(file, dst)
			if err != nil {
				fmt.Println("add chapter SaveUploadedFile error: ", err)
				continue
			}
		}

		msg := fmt.Sprintf("上传完成, 共上传 %v个文件", len(files))
		w.Write([]byte(msg))
		return
	}
}

func SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
