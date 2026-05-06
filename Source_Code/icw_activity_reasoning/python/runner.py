import argparse
import subprocess
import sys
from pathlib import Path


ERROR_FILE_NAME = 'error.log'


# 解析命令行输入参数
def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser()
    parser.add_argument('--task-code', required=True)
    parser.add_argument('--detector-path', required=True)
    parser.add_argument('--image-uuid', required=True)
    parser.add_argument('--runtime-root', required=True)
    return parser.parse_args()


# 调用指定脚本并校验报告输出
def run_detector(task_code: str, detector_path: Path, image_path: Path, task_runtime_dir: Path) -> Path:
    if not detector_path.exists():
        raise FileNotFoundError(f'detector directory not found: {detector_path}')
    script = detector_path / 'main.py'
    if not script.exists():
        raise FileNotFoundError(f'detector script not found: {script}')

    completed = subprocess.run(
        [sys.executable, str(script), '--input', str(image_path)],
        cwd=str(detector_path),
        check=False,
        capture_output=True,
        text=True,
    )
    combined_output = '\n'.join(item for item in [completed.stdout, completed.stderr] if item)
    if completed.returncode != 0:
        message = combined_output.strip()
        if message:
            raise RuntimeError(f'detector {task_code} failed with exit code {completed.returncode}: {message}')
        raise RuntimeError(f'detector {task_code} failed with exit code {completed.returncode}')

    report_path = task_runtime_dir / 'report.json'
    if not report_path.is_file():
        raise FileNotFoundError(f'report json not found: {report_path}')
    return report_path


# 写入捕获的错误信息
def write_error(task_runtime_dir: Path, err: Exception) -> None:
    task_runtime_dir.mkdir(parents=True, exist_ok=True)
    message = str(err).strip() or err.__class__.__name__
    (task_runtime_dir / ERROR_FILE_NAME).write_text(message, encoding='utf-8')


# 校验路径片段是否安全可用
def safe_path_part(value: str, name: str) -> str:
    value = value.strip()
    if not value or value in {'.', '..'} or '/' in value or '\\' in value:
        raise ValueError(f'{name} must be a single path segment: {value}')
    return value


# 执行检测
def main() -> int:
    args = parse_args()
    task_code = safe_path_part(args.task_code, 'task_code')
    detector_path = Path(args.detector_path).expanduser().resolve()
    image_uuid = safe_path_part(args.image_uuid, 'image_uuid')
    runtime_root = Path(args.runtime_root).expanduser().resolve()
    task_runtime_dir = runtime_root / task_code / image_uuid
    task_runtime_dir.mkdir(parents=True, exist_ok=True)
    error_path = task_runtime_dir / ERROR_FILE_NAME
    if error_path.exists():
        error_path.unlink()
    image_path = task_runtime_dir / 'original.png'
    if not image_path.is_file():
        raise FileNotFoundError(f'original image not found: {image_path}')

    try:
        run_detector(task_code, detector_path, image_path, task_runtime_dir)
    except Exception as err:
        write_error(task_runtime_dir, err)
        return 1
    return 0


# 执行主流程并兜底异常
def run() -> int:
    try:
        return main()
    except Exception:
        return 1


raise SystemExit(run())
