// Package simple provides a simple example.
package main

import (
	"github.com/iZIVer/imagemaker/imagemaker-service"
	imageserver_http "github.com/iZIVer/imagemaker/imagemaker-service/http"
	imageserver_http_crop "github.com/iZIVer/imagemaker/imagemaker-service/http/crop"
	imageserver_http_gamma "github.com/iZIVer/imagemaker/imagemaker-service/http/gamma"
	imageserver_http_gift "github.com/iZIVer/imagemaker/imagemaker-service/http/gift"
	imageserver_http_image "github.com/iZIVer/imagemaker/imagemaker-service/http/image"
	imageserver_image "github.com/iZIVer/imagemaker/imagemaker-service/image"
	imageserver_image_crop "github.com/iZIVer/imagemaker/imagemaker-service/image/crop"
	imageserver_image_gamma "github.com/iZIVer/imagemaker/imagemaker-service/image/gamma"
	_ "github.com/iZIVer/imagemaker/imagemaker-service/image/gif"
	imageserver_image_gift "github.com/iZIVer/imagemaker/imagemaker-service/image/gift"
	_ "github.com/iZIVer/imagemaker/imagemaker-service/image/jpeg"
	_ "github.com/iZIVer/imagemaker/imagemaker-service/image/png"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	errors "errors"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	pathfile "path/filepath"
	"time"

	"github.com/disintegration/gift"
)

func saveFile(w http.ResponseWriter, file multipart.File, handle *multipart.FileHeader) {
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	var fileName = handle.Filename
	var ext = pathfile.Ext(fileName)
	algorithm := sha1.New()
	algorithm.Write([]byte(fileName + time.Now().Format("2006-01-02 15:04:05")))
	var hashString = hex.EncodeToString(algorithm.Sum(nil))
	err = ioutil.WriteFile("./src/"+hashString+ext, data, 0666)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
}

func jsonResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprint(w, message)
}
func getFile(params imageserver.Params) (image.Image, error) {
	files, err := ioutil.ReadDir("./src")
	if err != nil {
		return nil, errors.New("file not found")
	}
	source, err := params.GetString("source")
	if err != nil {
		return nil, errors.New("file not found")
	}
	for _, f := range files {
		var fileName = f.Name()

		fName := pathfile.Base(fileName)
		extName := pathfile.Ext(fileName)
		bName := fName[:len(fName)-len(extName)]

		if source == bName {
			extWithoutDot := extName[1:len(extName)]
			var image = loadImage(fName, extWithoutDot)
			if image != nil {
				return image, nil
			}
		}
	}
	return nil, errors.New("file not found")
}
func loadImage(filename string, format string) image.Image {
	filePath := pathfile.Join("./src", filename)
	reader, _ := os.Open(filePath)
	defer reader.Close()
	im, _, err := image.Decode(reader)
	if err != nil {
		return nil
	}
	return im

}
func newImageHTTPHandler() http.Handler {
	var handler http.Handler = &imageserver_http.Handler{
		Parser: imageserver_http.ListParser([]imageserver_http.Parser{
			//	&imageserver_http.SourcePathParser{},
			&imageserver_http.SourceParser{},
			&imageserver_http_crop.Parser{},
			&imageserver_http_gift.RotateParser{},
			&imageserver_http_gift.ResizeParser{},
			&imageserver_http_image.FormatParser{},
			&imageserver_http_image.QualityParser{},
			&imageserver_http_gamma.CorrectionParser{},
		}),
		Server:   newServerImage(),
		ETagFunc: imageserver_http.NewParamsHashETagFunc(sha256.New),
	}
	return handler
}
func newServerImage() imageserver.Server {
	basicHdr := &imageserver_image.Handler{
		Processor: imageserver_image_gamma.NewCorrectionProcessor(
			imageserver_image.ListProcessor([]imageserver_image.Processor{
				&imageserver_image_crop.Processor{},
				&imageserver_image_gift.RotateProcessor{
					DefaultInterpolation: gift.CubicInterpolation,
				},
				&imageserver_image_gift.ResizeProcessor{
					DefaultResampling: gift.LanczosResampling,
					MaxWidth:          2048,
					MaxHeight:         2048,
				},
			}),
			true,
		),
	}
	/*
		gifHdr := &imageserver_image_gif.FallbackHandler{
			Handler: &imageserver_image_gif.Handler{
				Processor: &imageserver_image_gif.SimpleProcessor{
					Processor: imageserver_image.ListProcessor([]imageserver_image.Processor{
						&imageserver_image_crop.Processor{},
						&imageserver_image_gift.RotateProcessor{
							DefaultInterpolation: gift.NearestNeighborInterpolation,
						},
						&imageserver_image_gift.ResizeProcessor{
							DefaultResampling: gift.NearestNeighborResampling,
							MaxWidth:          1024,
							MaxHeight:         1024,
						},
					}),
				},
			},
			Fallback: basicHdr,
		}*/
	return &imageserver.HandlerServer{
		Server: &imageserver_image.Server{

			Provider: imageserver_image.ProviderFunc(func(params imageserver.Params) (image.Image, error) {
				return getFile(params)
			}),
			DefaultFormat: "png",
		},
		Handler: basicHdr,
	}
}
func main() {
	path, err1 := os.Getwd()
	if err1 != nil {
		if _, err1 := os.Stat(path); os.IsNotExist(err1) {
			os.Mkdir(path, os.ModePerm)
		}
	}
	http.Handle("/src", newImageHTTPHandler())
	//http.Handle("/src", http.StripPrefix("/", newImageHTTPHandler()))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var path = r.URL.Path[1:]
		if path == "" {
			path = "index.html"
		}
		dir, err := os.Getwd()
		if err == nil {
			fmt.Println(dir)
		}
		fmt.Println(path)
		http.ServeFile(w, r, dir+"/ui/"+path)
		//		http.ServeFile(w, r, "ui/")
	})
	http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		files, err := ioutil.ReadDir("./src")

		if err != nil {
			log.Fatal(err)
		}
		var filesCount = len(files)
		var filesList []string
		filesList = make([]string, filesCount)
		i := 0
		for _, f := range files {
			filesList[i] = f.Name()
			i++
		}
		fileData, err := json.Marshal(filesList)
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		w.Write(fileData)
	})
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {

		fileData, err := json.Marshal("hello")
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		w.Write(fileData)
	})
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		file, handle, err := r.FormFile("file")
		if err != nil {
			fmt.Fprintf(w, "%v", err)
			return
		}
		defer file.Close()

		mimeType := handle.Header.Get("Content-Type")
		switch mimeType {
		case "image/jpeg":
			saveFile(w, file, handle)
			http.Redirect(w, r, "/", 301)
		case "image/png":
			saveFile(w, file, handle)
			http.Redirect(w, r, "/", 301)
		default:
			jsonResponse(w, http.StatusBadRequest, "The format file is not valid.")
		}
	})

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		panic(err)
	}
}
