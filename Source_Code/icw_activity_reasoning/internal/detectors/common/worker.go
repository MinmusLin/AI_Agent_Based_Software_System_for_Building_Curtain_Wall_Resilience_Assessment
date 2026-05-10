package common

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

// pythonWorkerRequest Python Worker 请求结构体
type pythonWorkerRequest struct {
	TaskCode    string `json:"task_code"`
	ImageUuid   string `json:"image_uuid"`
	RuntimeRoot string `json:"runtime_root"`
	ModelPath   string `json:"model_path"`
}

// pythonWorkerResponse Python Worker 响应结构体
type pythonWorkerResponse struct {
	OK    bool   `json:"ok"`
	Error string `json:"error"`
}

// pythonWorkerReadResult Python Worker 标准输出读取结果
type pythonWorkerReadResult struct {
	line string
	err  error
}

// pythonWorker 单个原子检测能力的常驻 Python 进程
type pythonWorker struct {
	code        string
	modelPath   string
	projectRoot string
	workerPath  string
	runtimeRoot string
	mu          sync.Mutex
	cmd         *exec.Cmd
	stdin       io.WriteCloser
	stdout      *bufio.Reader
}

// newPythonWorker 创建单个原子检测能力的常驻 Python 进程
func newPythonWorker(code, runtimeRoot, modelPath string) *pythonWorker {
	return &pythonWorker{
		code:        code,
		modelPath:   modelPath,
		projectRoot: absPath("."),
		workerPath:  absPath("python/worker.py"),
		runtimeRoot: absPath(runtimeRoot),
	}
}

// Run 向 Python Worker 投递检测任务
func (w *pythonWorker) Run(ctx context.Context, imageUuid string) error {
	if w == nil {
		return errors.New("python worker is nil")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	if err := ctx.Err(); err != nil {
		return err
	}
	if err := w.startLocked(); err != nil {
		return err
	}

	req := &pythonWorkerRequest{
		TaskCode:    w.code,
		ImageUuid:   imageUuid,
		RuntimeRoot: w.runtimeRoot,
		ModelPath:   w.modelPath,
	}
	payload, err := json.Marshal(req)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w.stdin, "%s\n", payload); err != nil {
		w.resetLocked()
		return err
	}

	readResult := make(chan pythonWorkerReadResult, 1)
	go func() {
		line, err := w.stdout.ReadString('\n')
		readResult <- pythonWorkerReadResult{
			line: line,
			err:  err,
		}
	}()

	var result pythonWorkerReadResult
	select {
	case <-ctx.Done():
		w.resetLocked()
		return ctx.Err()
	case result = <-readResult:
	}

	if result.err != nil {
		w.resetLocked()
		return result.err
	}

	var resp pythonWorkerResponse
	if err := json.Unmarshal([]byte(result.line), &resp); err != nil {
		w.resetLocked()
		return err
	}

	if !resp.OK {
		if strings.TrimSpace(resp.Error) == "" {
			return fmt.Errorf("run python detector %s failed", w.code)
		}
		return fmt.Errorf("run python detector %s failed: %s", w.code, strings.TrimSpace(resp.Error))
	}

	return nil
}

// startLocked 启动 Python Worker 进程
func (w *pythonWorker) startLocked() error {
	if w.cmd != nil && w.stdin != nil && w.stdout != nil {
		return nil
	}

	cmd := exec.Command(uvBinary(), "run", "--project", w.projectRoot, "python", w.workerPath)
	cmd.Dir = w.projectRoot

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		_ = stdin.Close()
		return err
	}
	if err := cmd.Start(); err != nil {
		_ = stdin.Close()
		return err
	}

	w.cmd = cmd
	w.stdin = stdin
	w.stdout = bufio.NewReader(stdoutPipe)

	go func() {
		_ = cmd.Wait()
	}()

	return nil
}

// uvBinary 获取 uv 可执行文件路径
func uvBinary() string {
	if path, err := exec.LookPath("uv"); err == nil {
		return path
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "uv"
	}
	path := filepath.Join(homeDir, ".local", "bin", "uv")
	if _, err := os.Stat(path); err == nil {
		return path
	}
	return "uv"
}

// resetLocked 关闭并清理 Python Worker 进程
func (w *pythonWorker) resetLocked() {
	if w.stdin != nil {
		_ = w.stdin.Close()
	}
	if w.cmd != nil && w.cmd.Process != nil {
		_ = w.cmd.Process.Kill()
	}
	w.cmd = nil
	w.stdin = nil
	w.stdout = nil
}

// absPath 将相对路径转换为绝对路径
func absPath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return ""
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	return abs
}
