package connect

import (
	"little-scp-client/util"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// 建立ssh连接所需要的信息
type UserInfo struct {
	// 用户名
	UserName string
	// 远程主机密码
	passwd string
	// 远程主机ip地址
	Addr string
	// 下载目的地址
	DestAddr string
	// 原文件地址
	SrcAddr string
}


// 初始化
func Init(userName, passwd, Addr, destAddr, srcaddr string) UserInfo {
	return UserInfo{
		UserName: userName,
		passwd: passwd,
		Addr: Addr + ":22",
		DestAddr: destAddr,
		SrcAddr: srcaddr,
	}
}

// 建立ssh连接
func(u *UserInfo) SshConnect() *ssh.Client {
	config := &ssh.ClientConfig{
		User: u.UserName,
		Auth: []ssh.AuthMethod{
			ssh.Password(u.passwd),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", u.Addr, config)
	util.CheckErr(err, "Failed to create SSH connection")
	return client
}

// 建立sftp连接
func(u *UserInfo) SftpConnect(sshclient *ssh.Client) *sftp.Client {
	sftpclient, err := sftp.NewClient(sshclient)
	util.CheckErr(err, "Failed to create SFTP connection")
	return sftpclient
}
