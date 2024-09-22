package factory

import (
	"bufio"
	"fmt"
	"os"

	json "github.com/json-iterator/go"

	"github.com/limes-cloud/cron/internal/client/service"
)

type store struct {
}

const (
	filePrefix = ".db/%s"
)

type scanner struct {
	scanner *bufio.Scanner
}

func (s *scanner) Scan() bool {
	return s.scanner.Scan()
}

func (s *scanner) Data() (*service.ExecTaskReply, error) {
	res := service.ExecTaskReply{}
	if err := json.Unmarshal(s.scanner.Bytes(), &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *store) filename(uuid string) string {
	return fmt.Sprintf(filePrefix, uuid)
}

func (s *store) add(uuid string, data *service.ExecTaskReply) error {
	file, err := os.OpenFile(s.filename(uuid), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	res, _ := json.MarshalToString(data)
	if _, err = file.WriteString(res); err != nil {
		return err
	}
	return nil
}

func (s *store) exist(uuid string) bool {
	_, err := os.Stat(s.filename(uuid))
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

func (s *store) scanner(uuid string) (*scanner, error) {
	srcFile, err := os.Open(s.filename(uuid))
	if err != nil {
		return nil, err
	}
	defer srcFile.Close()
	return &scanner{scanner: bufio.NewScanner(srcFile)}, nil
}

func (s *store) remove(uuid string) error {
	return os.Remove(uuid)
}

// 删除文件的前n行，然后将其余行写入同一文件
func (s *store) removeN(uuid string, n int) error {
	// 打开源文件
	srcFile, err := os.Open(s.filename(uuid))
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "tmp")
	if err != nil {
		return err
	}
	defer tmpFile.Close()

	scanner := bufio.NewScanner(srcFile)
	lineCount := 0

	// 逐行读取，跳过前n行
	for scanner.Scan() {
		if lineCount >= n {
			_, err = tmpFile.WriteString(scanner.Text() + "\n")
			if err != nil {
				return err
			}
		}
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// 关闭源文件和临时文件
	_ = srcFile.Close()
	_ = tmpFile.Close()

	// 替换原始文件
	err = os.Rename(tmpFile.Name(), s.filename(uuid))
	if err != nil {
		return err
	}

	return nil
}
