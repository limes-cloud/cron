package factory

import (
	"bufio"
	"errors"
	"os"
	"os/exec"

	"github.com/limes-cloud/kratosx"

	"github.com/limes-cloud/cron/internal/client/biz"
)

func (f *Factory) shell(ctx kratosx.Context, task *biz.Task) (int, error) {
	shell := f.conf.Shell
	if shell == "" {
		shell = defaultShell
	}

	code := defaultErrorCode

	tpFile, err := os.CreateTemp("", "*.sh")
	if err != nil {
		return code, err
	}
	if _, err = tpFile.WriteString(task.Value); err != nil {
		return code, err
	}
	if err := tpFile.Sync(); err != nil {
		return code, err
	}

	defer func() {
		_ = tpFile.Close()
		_ = os.Remove(tpFile.Name())
	}()

	cmd := exec.CommandContext(ctx, shell, tpFile.Name())
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return code, err
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return code, err
	}

	go func() {
		stdout := bufio.NewScanner(stdoutPipe)
		for stdout.Scan() {
			f.reply(task.Uuid, logInfo, stdout.Text())
		}
	}()

	go func() {
		stderr := bufio.NewScanner(stderrPipe)
		for stderr.Scan() {
			f.reply(task.Uuid, logError, stderr.Text())
		}
	}()

	if err := cmd.Start(); err != nil {
		return code, err
	}

	if err := cmd.Wait(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode(), nil
		}
		return defaultErrorCode, err
	}
	return 0, nil
}
