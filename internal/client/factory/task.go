package factory

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/limes-cloud/kratosx"

	"github.com/limes-cloud/cron/internal/client/biz"
	"github.com/limes-cloud/cron/internal/client/conf"
)

type watcher struct {
	cancel context.CancelFunc
	mutex  sync.RWMutex
	buf    []*biz.ExecTaskReply
	reply  biz.ExecTaskReplyFunc
	closer chan error
	close  atomic.Bool
	time   time.Time
	err    error
}

type Factory struct {
	conf  *conf.Config
	ws    map[string]*watcher
	mutex sync.RWMutex
}

var (
	_f   *Factory
	once sync.Once
)

const (
	defaultErrorCode = 400
	defaultShell     = "/bin/sh"

	logInfo  = "info"
	logError = "error"

	startInfo = "开始进行第%d次执行任务，任务索引:%s"
	errorInfo = "执行任务失败:"
)

func New(conf *conf.Config) *Factory {
	once.Do(func() {
		_f = &Factory{
			conf:  conf,
			ws:    make(map[string]*watcher),
			mutex: sync.RWMutex{},
		}
	})
	return _f
}

func (f *Factory) CancelExecTask(uuid string) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	wtr := f.ws[uuid]
	if wtr != nil && !wtr.close.Load() {
		wtr.cancel()
	}
}

func (f *Factory) ExecTask(ctx kratosx.Context, task *biz.Task, fn biz.ExecTaskReplyFunc) error {
	defer func() {
		f.mutex.Lock()
		if f.ws[task.Uuid] != nil && f.ws[task.Uuid].close.Load() && len(f.ws[task.Uuid].buf) == 0 {
			close(f.ws[task.Uuid].closer)
			delete(f.ws, task.Uuid)
		}
		f.mutex.Unlock()
	}()

	childCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	f.mutex.Lock()
	if wtr, ok := f.ws[task.Uuid]; ok {
		f.mutex.Unlock()

		closer := make(chan error, 1)
		if !wtr.close.Load() {
			wtr.mutex.Lock()
			close(wtr.closer)
			wtr.reply = fn
			wtr.closer = closer
			wtr.mutex.Unlock()
		}

		var i = 0
		defer func(index *int) {
			f.mutex.RLock()
			itr, ok := f.ws[task.Uuid]
			f.mutex.RUnlock()

			itr.mutex.Lock()
			if *index == len(itr.buf) {
				wtr.close.Store(true)
			}

			if *index != 0 && ok {
				itr.buf = itr.buf[*index:]
			}
			itr.mutex.Unlock()
		}(&i)

		var tmp = make([]*biz.ExecTaskReply, len(wtr.buf))
		copy(tmp, wtr.buf)

		for ; i < len(tmp); i++ {
			if err := fn(tmp[i]); err != nil {
				return err
			}
		}

		if !wtr.close.Load() {
			return <-wtr.closer
		}
		return wtr.err
	}

	closer := make(chan error, 1)
	wtr := &watcher{
		buf:    make([]*biz.ExecTaskReply, 0),
		reply:  fn,
		closer: closer,
		close:  atomic.Bool{},
		cancel: cancel,
		time:   time.Now(),
	}
	f.ws[task.Uuid] = wtr
	f.mutex.Unlock()

	f.exec(kratosx.MustContext(childCtx), task, wtr)
	wtr.err = <-closer
	return wtr.err
}

func (f *Factory) reply(uuid string, tp string, text string) {
	f.mutex.RLock()
	wtr := f.ws[uuid]
	f.mutex.RUnlock()

	res := &biz.ExecTaskReply{
		Type:    tp,
		Content: text,
		Time:    uint32(time.Now().Unix()),
	}
	wtr.mutex.Lock()
	defer wtr.mutex.Unlock()
	if err := wtr.reply(res); err != nil {
		wtr.buf = append(wtr.buf, res)
	}
}

func (f *Factory) exec(ctx kratosx.Context, task *biz.Task, wtr *watcher) {
	var (
		err   error
		code  = defaultErrorCode
		count = int(task.RetryCount) + 1
	)

	for i := 0; i < count; i++ {
		f.reply(task.Uuid, logInfo, fmt.Sprintf(startInfo, i+1, task.Uuid))
		childCtx := ctx
		if task.MaxExecTime != 0 {
			c, _ := context.WithTimeout(ctx, time.Duration(task.MaxExecTime)*time.Second)
			childCtx = kratosx.MustContext(c)
		}

		var ()
		switch task.Type {
		case "shell":
			code, err = f.shell(childCtx, task)
		case "http":
			code, err = f.http(childCtx, task)
		default:
			err = errors.New("不支持的任务类型")
		}

		if err == nil && int(task.ExpectCode) == code {
			break
		}

		if err == nil && code != int(task.ExpectCode) {
			err = fmt.Errorf("code %d not eq expect_code %d", code, task.ExpectCode)
		}
		f.reply(task.Uuid, logError, errorInfo+err.Error())

		if task.RetryWaitTime != 0 {
			time.Sleep(time.Duration(task.RetryWaitTime) * time.Second)
		}
	}

	wtr.mutex.Lock()
	defer func() {
		wtr.close.Store(true)
		wtr.mutex.Unlock()
	}()
	wtr.closer <- err
}
