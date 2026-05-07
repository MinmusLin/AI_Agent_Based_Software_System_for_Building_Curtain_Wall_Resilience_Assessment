import contextlib
import importlib.util
import io
import json
import sys
import traceback
from pathlib import Path
from types import ModuleType
from typing import Any


PROJECT_ROOT = Path(__file__).resolve().parents[1]
DETECTOR_CACHE: dict[str, ModuleType] = {}


def is_safe_path_part(value: str) -> bool:
    return bool(value) and value not in {'.', '..'} and '/' not in value and '\\' not in value


def load_detector(task_code: str) -> ModuleType:
    if task_code in DETECTOR_CACHE:
        return DETECTOR_CACHE[task_code]
    detector_path = PROJECT_ROOT / 'python' / 'detectors' / task_code / 'main.py'
    if not detector_path.exists():
        raise FileNotFoundError(f'detector file not found: {detector_path}')
    spec = importlib.util.spec_from_file_location(f'icw_detector_{task_code}', detector_path)
    if spec is None or spec.loader is None:
        raise RuntimeError(f'detector module cannot be loaded: {task_code}')
    module = importlib.util.module_from_spec(spec)
    with contextlib.redirect_stdout(io.StringIO()), contextlib.redirect_stderr(io.StringIO()):
        spec.loader.exec_module(module)
    DETECTOR_CACHE[task_code] = module
    return module


def run_detector(request: dict[str, Any]) -> None:
    task_code = str(request.get('task_code', '')).strip()
    image_uuid = str(request.get('image_uuid', '')).strip()
    runtime_root = Path(str(request.get('runtime_root', '')).strip()).expanduser().resolve()
    if not is_safe_path_part(task_code):
        raise ValueError('task_code is invalid')
    if not is_safe_path_part(image_uuid):
        raise ValueError('image_uuid is invalid')
    image_path = runtime_root / task_code / image_uuid / 'original.png'
    if not image_path.exists():
        raise FileNotFoundError(f'original image not found: {image_path}')
    with contextlib.redirect_stdout(io.StringIO()), contextlib.redirect_stderr(io.StringIO()):
        module = load_detector(task_code)
        detect = getattr(module, 'detect', None)
        if detect is None:
            raise RuntimeError(f'detector detect function is missing: {task_code}')
        detect(image_path)


def write_response(response: dict[str, Any]) -> None:
    sys.stdout.write(json.dumps(response, ensure_ascii=False, separators=(',', ':')) + '\n')
    sys.stdout.flush()


def main() -> int:
    for line in sys.stdin:
        try:
            request = json.loads(line)
            run_detector(request)
            write_response({'ok': True})
        except Exception as exc:
            error_message = ''.join(traceback.format_exception_only(type(exc), exc)).strip()
            write_response({'ok': False, 'error': error_message})
    return 0


if __name__ == '__main__':
    raise SystemExit(main())
