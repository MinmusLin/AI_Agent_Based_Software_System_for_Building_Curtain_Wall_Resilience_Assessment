from __future__ import annotations

import argparse
import subprocess
import sys
from pathlib import Path


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser()
    parser.add_argument('--task-code', required=True)
    parser.add_argument('--detector-path', required=True)
    parser.add_argument('--image-uuid', required=True)
    parser.add_argument('--runtime-root', required=True)
    return parser.parse_args()


def detector_dir(detector_path: str) -> Path:
    path = Path(detector_path).expanduser().resolve()
    if not path.exists():
        raise FileNotFoundError(f'detector directory not found: {path}')
    return path


def run_detector(task_code: str, detector_path: Path, image_path: Path, task_runtime_dir: Path) -> Path:
    target_dir = detector_dir(str(detector_path))
    script = target_dir / 'main.py'
    if not script.exists():
        raise FileNotFoundError(f'detector script not found: {script}')

    completed = subprocess.run(
        [sys.executable, str(script), '--input', str(image_path)],
        cwd=str(target_dir),
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


def safe_path_part(value: str, name: str) -> str:
    value = value.strip()
    if not value or value in {'.', '..'} or '/' in value or '\\' in value:
        raise ValueError(f'{name} must be a single path segment: {value}')
    return value


def main() -> int:
    args = parse_args()
    task_code = safe_path_part(args.task_code, 'task_code')
    detector_path = Path(args.detector_path).expanduser().resolve()
    image_uuid = safe_path_part(args.image_uuid, 'image_uuid')
    runtime_root = Path(args.runtime_root).expanduser().resolve()
    task_runtime_dir = runtime_root / task_code / image_uuid
    task_runtime_dir.mkdir(parents=True, exist_ok=True)
    image_path = task_runtime_dir / 'original.png'
    if not image_path.is_file():
        raise FileNotFoundError(f'original image not found: {image_path}')

    run_detector(task_code, detector_path, image_path, task_runtime_dir)
    return 0


if __name__ == '__main__':
    raise SystemExit(main())
