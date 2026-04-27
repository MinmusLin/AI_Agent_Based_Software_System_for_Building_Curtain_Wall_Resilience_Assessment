from __future__ import annotations

import argparse
import json
import time
from pathlib import Path
from typing import Any

import torch
import torch.nn as nn
import torchvision.models as models
from PIL import Image
from torchvision import transforms


APP_ROOT = Path(__file__).resolve().parent
MODEL_PATH = APP_ROOT / 'model' / 'best_weights_model.pt'
OUTPUT_ROOT = APP_ROOT / 'outputs'
INPUT_SIZE = 224
CLASS_NAMES = ['defect', 'undefect']


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description='本地幕墙爆裂检测')
    parser.add_argument('--input', required=True, help='输入图片路径')
    return parser.parse_args()


def validate_model_file(model_path: Path) -> None:
    if not model_path.exists():
        raise FileNotFoundError(f'未找到模型文件: {model_path}')
    if model_path.stat().st_size < 1024 * 1024:
        head = model_path.read_bytes()[:128]
        if b'git-lfs.github.com/spec' in head:
            raise RuntimeError(f'模型文件是 Git LFS 指针，不是真实权重: {model_path}')
        raise RuntimeError(f'模型文件大小异常: {model_path}')


def get_device() -> torch.device:
    return torch.device('cuda' if torch.cuda.is_available() else 'cpu')


def build_model() -> nn.Module:
    model = models.resnet34(weights=None)
    num_features = model.fc.in_features
    model.fc = nn.Linear(num_features, len(CLASS_NAMES))
    return model


def load_model(model_path: Path, device: torch.device) -> nn.Module:
    validate_model_file(model_path)
    model = build_model()
    state_dict = torch.load(model_path, map_location=device)
    if any(key.startswith('module.') for key in state_dict):
        state_dict = {key.replace('module.', '', 1): value for key, value in state_dict.items()}
    model.load_state_dict(state_dict)
    model.to(device)
    model.eval()
    return model


def read_image(input_path: Path) -> Image.Image:
    if not input_path.exists():
        raise FileNotFoundError(f'未找到输入图片: {input_path}')
    return Image.open(input_path).convert('RGB')


def preprocess_image(image: Image.Image) -> torch.Tensor:
    transform = transforms.Compose(
        [
            transforms.Resize((INPUT_SIZE, INPUT_SIZE)),
            transforms.ToTensor(),
            transforms.Normalize(mean=[0.485, 0.456, 0.406], std=[0.229, 0.224, 0.225]),
        ]
    )
    return transform(image).unsqueeze(0)


def predict(model: nn.Module, image: Image.Image, device: torch.device) -> dict[str, Any]:
    tensor = preprocess_image(image).to(device)
    with torch.no_grad():
        logits_tensor = model(tensor)[0]
        probabilities_tensor = torch.softmax(logits_tensor, dim=0)
        confidence_tensor, predicted_tensor = torch.max(probabilities_tensor, dim=0)

    logits = [round(float(value), 6) for value in logits_tensor.cpu().tolist()]
    probabilities = [round(float(value), 6) for value in probabilities_tensor.cpu().tolist()]
    predicted_index = int(predicted_tensor.cpu().item())
    predicted_class = CLASS_NAMES[predicted_index]

    return {
        'predicted_index': predicted_index,
        'predicted_class': predicted_class,
        'is_spalling': bool(predicted_class == 'defect'),
        'confidence': round(float(confidence_tensor.cpu().item()), 6),
        'class_names': CLASS_NAMES,
        'logits': logits,
        'probabilities': probabilities,
        'probabilities_by_class': dict(zip(CLASS_NAMES, probabilities)),
    }


def build_report(input_path: Path, image: Image.Image, prediction: dict[str, Any]) -> dict[str, Any]:
    image_width, image_height = image.size
    return {
        'input': str(input_path),
        'image_size': {'width': image_width, 'height': image_height},
        'model_output': prediction,
        'result': prediction['predicted_class'],
        'has_spalling': prediction['is_spalling'],
        'confidence': prediction['confidence'],
    }


def build_output_paths(output_dir: Path, input_path: Path) -> dict[str, str]:
    output_dir.mkdir(parents=True, exist_ok=True)
    report_path = output_dir / f'{input_path.stem}-report.json'
    return {'report_json': str(report_path)}


def print_report(report: dict[str, Any]) -> None:
    model_output = report['model_output']
    image_size = report['image_size']
    has_spalling = '是' if report['has_spalling'] else '否'

    print('\n==== 本地爆裂检测报告 ====')
    print(f"输入图片: {report['input']}")
    print(f"图片尺寸: {image_size['width']} x {image_size['height']} px")
    print(f'是否检出爆裂: {has_spalling}')
    print(f"预测类别: {model_output['predicted_class']}")
    print(f"置信度: {model_output['confidence']}")
    print(f"JSON: {report['outputs']['report_json']}")


def main() -> int:
    args = parse_args()
    start_time = time.time()
    input_path = Path(args.input).expanduser().resolve()
    device = get_device()
    image = read_image(input_path)
    model = load_model(MODEL_PATH, device)
    prediction = predict(model, image, device)
    report = build_report(input_path, image, prediction)
    report['model'] = str(MODEL_PATH)
    report['device'] = str(device)
    report['input_size'] = INPUT_SIZE
    report['runtime_seconds'] = round(time.time() - start_time, 3)

    timestamp = time.strftime('%Y%m%d-%H%M%S')
    run_dir = OUTPUT_ROOT / f'{input_path.stem}-{timestamp}'
    report['outputs'] = build_output_paths(run_dir, input_path)
    Path(report['outputs']['report_json']).write_text(json.dumps(report, ensure_ascii=False, indent=2), encoding='utf-8')
    print_report(report)
    return 0


if __name__ == '__main__':
    raise SystemExit(main())
