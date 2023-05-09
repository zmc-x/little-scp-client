package main

import (
	"archive/zip"
	"fmt"
	"little-scp-client/util"
	"little-scp-client/view"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pkg/sftp"
)

// 控制台输出信息
type info struct {
	filename, filesize, status string
}

var (
	wg        sync.WaitGroup
	zipwriter *zip.Writer
	lock      sync.Mutex
	msg		  chan info
)

func main() {
	msg = make(chan info)
	// 初始化输出界面
	m := view.InitialModel()
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
	// 初始化&建立连接
	sshuserinfo := m.Value()
	localaddr, remoteaddr := m.Getaddr()
	sshclient := sshuserinfo.SshConnect()
	defer sshclient.Close()
	// 建立sftp客户端
	sftpClient := sshuserinfo.SftpConnect(sshclient)
	defer sftpClient.Close()
	var file *os.File
	var err error
	flag := true
	// 判断是否为目录
	if filepath.Ext(remoteaddr) == "" {
		// 本地创建压缩包
		file, err = os.Create(filepath.Join(localaddr, filepath.Base(remoteaddr)+".zip"))
		fmt.Println(filepath.Join(localaddr, filepath.Base(remoteaddr)+".zip"))
		util.CheckErr(err, "failed to create compressed file")
	} else {
		flag = false
	}
	defer file.Close()
	// 初始化zip writer
	zipwriter = zip.NewWriter(file)
	defer zipwriter.Close()
	t := time.Now()
	if flag {
		visitfiles(remoteaddr, sftpClient)
	} else {
		size := util.DownloadFile(filepath.Join(localaddr, filepath.Base(remoteaddr)), remoteaddr, sftpClient)
		fmt.Println(remoteaddr, util.Changefilesize(int(size)), "🚀")
	}
	// closer
	go func() {
		wg.Wait()
		close(msg)
	}()
	for v := range msg {
		fmt.Println(v.filename, v.filesize, v.status)
	}
	log.Println(time.Since(t))
}

// 远程主机的目录遍历 & 文件的并发写入
func visitfiles(str string, client *sftp.Client) {
	files, _ := client.ReadDir(str)
	// dfs
	for _, file := range files {
		temp := str + "/" + file.Name()
		if file.IsDir() {
			visitfiles(temp, client)
		} else {
			// 并发
			wg.Add(1)
			go func() {
				defer wg.Done()
				filename := filepath.Join(strings.Split(temp, "/")[3:]...)
				// 开始写入文件
				lock.Lock()
				defer lock.Unlock()
				size := util.DownloadFoler(filename, temp, zipwriter, client)
				msg <- info{filename: temp, filesize: util.Changefilesize(int(size)), status: "✔"}
			}()
		}
	}
}
