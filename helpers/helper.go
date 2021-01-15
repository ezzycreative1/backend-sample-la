package helpers

import (
	"crypto/aes"
	"crypto/cipher"
	cRand "crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"

	"github.com/gin-gonic/gin"
)

func ValidateImage(f multipart.File) bool {
	// maximize CPU for better performance
	runtime.GOMAXPROCS(runtime.NumCPU())

	buff := make([]byte, 512) // why 512 byte ? please read http://golang.org/pkg/net/http/#DetectContentType
	if _, err := f.Read(buff); err != nil {
		fmt.Println("========= level debug ==============")
		fmt.Println("user updaload image")
		fmt.Printf("message : %s\n", err.Error())
		return false
	}

	filetype := http.DetectContentType(buff)
	switch filetype {
	case "image/jpeg", "image/jpg":
		return true
	case "image/png":
		return true
	default:
		return false
	}
}

// FileUpload ..
func FileUpload(c *gin.Context) (error, string) {

	file, handler, err := c.Request.FormFile("filePhoto")
	if err != nil {
		return fmt.Errorf("cannot read request body"), ""
	}
	valid := ValidateImage(file)
	if !valid {
		return fmt.Errorf("format image not valid only jpeg, jpg, and png"), ""
	}
	defer file.Close()
	dir := fmt.Sprintf("%s/src/backend-sample-la/storage/photo", os.Getenv("GOPATH"))
	f, err := os.OpenFile(dir+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return err, ""
	}
	defer f.Close()
	io.Copy(f, file)
	newName := fmt.Sprintf("%s%s", GenerateRandString(40), filepath.Ext(handler.Filename))
	os.Rename(dir+handler.Filename, dir+newName)

	return nil, newName
}

// GenerateRandString ..
func GenerateRandString(lg int) string {

	var letterRunes = []rune("123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, lg)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)

}

// ValidEmail ..
func ValidEmail(email string) bool {

	reg := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	regMatch := reg.MatchString(email)

	return regMatch
}

func GenetateUrlImage(image string) string {

	payload := []byte(image)

	block, err := aes.NewCipher([]byte("inikey"))
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	ciphertext := make([]byte, aes.BlockSize+len(payload))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(cRand.Reader, iv); err != nil {
		fmt.Println(err.Error())
		return ""
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], payload)

	content := base64.URLEncoding.EncodeToString(ciphertext)
	url := os.Getenv("BASE_URL")

	return fmt.Sprintf("%s/user/image/%s", url, content)
}

func DecodeImageContent(cryptoText string) string {

	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher([]byte(""))
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	if len(ciphertext) < aes.BlockSize {
		return ""
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext)
}
