package storage

import (
	"errors"
	"fmt"
	"go-prerender/internal/config"
	"go-prerender/utils"
	"os"
	"time"
)

func GetData(domain string, name string) (string, error) {
	filepath := fmt.Sprintf("data/%s/%s.html", domain, name)
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		return "", err
	}
	cacheTime := config.GetConfig().System.CacheTime
	if cacheTime == 0 {
		return "", errors.New("no cache")
	}

	// 如果文件不是永久不过其，并且已经过期，则删除文件
	if cacheTime != -1 && time.Now().Unix()-fileInfo.ModTime().Unix() >= cacheTime {
		_ = os.Remove(filepath)
		return "", errors.New("the file has expired")
	}

	return utils.FileLib.ReadFileContent(filepath)
}

func SaveData(domain string, name string, content string) error {
	dirPath := fmt.Sprintf("data/%s/", domain)
	_ = utils.FileLib.CreateDir(dirPath)

	filepath := fmt.Sprintf("data/%s/%s.html", domain, name)
	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		return err
	}
	defer f.Close()
	_, _ = f.WriteString(content)
	return nil
}
