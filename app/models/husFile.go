package models

import (
	"os"
	"io"
	"mime/multipart"
	"github.com/revel/revel"
	"github.com/disintegration/imaging"
	"path"
)

type HusFile struct {
	UploadDir string
	Files []*multipart.FileHeader
}

func (husFile HusFile) UploadFile(filename string) bool {
	for i, _ := range husFile.Files {
		//For each file header, get a handle to the actual file
		file, err := husFile.Files[i].Open()
		defer file.Close() //close the source file handle on function return
		if err != nil {
			return false
		}
		//Create destination file
		destination := revel.BasePath + husFile.UploadDir + "/" + filename
		dst, err := os.Create(destination)
		defer dst.Close() //Close the destination file handle on function return
		defer os.Chmod(destination, (os.FileMode)(0644)) //Limit access restrictions
		if err != nil {
			return false
		}

		//Copy the uploaded file to the destination file
		if _, err := io.Copy(dst, file); err != nil {
			return false
		}
	}

	return true
}

func (husFile HusFile) ThumbnailImage(filename string, w, h int) bool {
	//Load and resize file
	file := path.Join(revel.BasePath + husFile.UploadDir, filename)
	img, err := imaging.Open(file)
	if err != nil {
		return false
	}
	destImg := imaging.Thumbnail(img, w, h, imaging.Lanczos)

	//Save the combined image to file
	imaging.Save(destImg, file)
	return true
}

func (husFile HusFile) ResizeImage(filename string, w, h int) bool {
	//Load and resize file
	file := path.Join(revel.BasePath + husFile.UploadDir, filename)
	img, err := imaging.Open(file)
	if err != nil {
		return false
	}
	destImg := imaging.Resize(img, w, h, imaging.Lanczos)

	//Save the combined image to file
	imaging.Save(destImg, file)
	return true
}

func (husFile HusFile) ScaleImage(filename string, w, h int) bool {
	//Load and resize file
	file := path.Join(revel.BasePath + husFile.UploadDir, filename)
	img, err := imaging.Open(file)
	if err != nil {
		return false
	}
	destImg := imaging.Fit(img, w, h, imaging.Lanczos)

	//Save the combined image to file
	imaging.Save(destImg, file)
	return true
}