package common

import (
	"crypto/tls"
	"math"
	"strconv"

	"gopkg.in/gomail.v2"

	"net"
	"os"
	"strings"
)

const MaxUnixTimeInt64 int64 = math.MaxInt64

// image formats and magic numbers
var magicTable = map[string]string{
	"\xff\xd8\xff":      "image/jpeg",
	"\x89PNG\r\n\x1a\n": "image/png",
	"GIF87a":            "image/gif",
	"GIF89a":            "image/gif",
}

//檢查 PATH 狀態 0: 不存在, 1: 是檔案, 2: 是目錄, -1: 其他錯誤
func CheckPath(path string) (int, error) {
	if fInfo, err := os.Stat(path); err == nil {
		if fInfo.IsDir() {
			return 2, nil
		} else {
			return 1, nil
		}
	} else if os.IsNotExist(err) {
		return 0, err
	} else {
		return -1, err
	}

}

// mimeFromIncipit returns the mime type of an image file from its first few
// bytes or the empty string if the file does not look like a known file type
func ImageFileType(incipit []byte) string {
	incipitStr := string(incipit)
	for magic, mime := range magicTable {
		if strings.HasPrefix(incipitStr, magic) {
			return mime
		}
	}

	return ""
}

func CheckEmailLogin(smtpServerAddress string, smtpUserName string, smtpPassword string) error {
	host, portStr, _ := net.SplitHostPort(smtpServerAddress)
	port, _ := strconv.Atoi(portStr)
	d := gomail.NewDialer(host, port, smtpUserName, smtpPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	var err error
	var s gomail.SendCloser
	s, err = d.Dial()
	if err == nil {
		s.Close()
	}

	return err
}

func GomailSender(from string, to string, subject string, body string, smtpServerAddress string, smtpUserName string, smtpPassword string) error {
	host, portStr, _ := net.SplitHostPort(smtpServerAddress)
	port, _ := strconv.Atoi(portStr)
	d := gomail.NewDialer(host, port, smtpUserName, smtpPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	m := gomail.NewMessage(gomail.SetEncoding(gomail.Base64))
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	//m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	//m.Attach("/home/Alex/lolcat.jpg")

	err := d.DialAndSend(m)
	return err

}
