package util

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/pkg/sftp"
)

// 计算文件的大小
func Changefilesize(filesize int) string {
	if filesize < 1024 {
		return fmt.Sprintf("%.3f", float64(filesize)) + "kB"
	} else if filesize >= 1024 && filesize < 1024 * 1024 {
		return fmt.Sprintf("%.3f", float64(filesize) / 1024.0) + "MB"
	} else {
		return fmt.Sprintf("%.3f", float64(filesize) / 1048576.0) + "GB"
	}
}

// 下载文件夹
func DownloadFoler(localname, remotename string, zipwriter *zip.Writer, client *sftp.Client) int64 {
	// 创建文件
	dst, err := zipwriter.Create(localname)
	CheckErr(err, "failed to create the compressed file")
	// 打开文件
	src, err := client.Open(remotename)
	CheckErr(err, "failed to open remote file")
	defer src.Close()
	// 开始写入
	size, err := io.Copy(dst, src)
	CheckErr(err, "copy faile")
	return size
}

// 下载单个文件
func DownloadFile(localname, remotename string, client *sftp.Client) int64 {
	// 创建文件
	dst, err := os.Create(localname)
	CheckErr(err, "failed to create the file")
	defer dst.Close()
	// 打开文件
	src, err := client.Open(remotename)
	CheckErr(err, "failed to open remote file")
	defer src.Close()
	// 开始写入
	size, err := io.Copy(dst, src)
	CheckErr(err, "copy faile")
	return size
}


// 错误检查
func CheckErr(err error, str string) {
	if err != nil {
		log.Fatal(str + " " + err.Error())
	}
}