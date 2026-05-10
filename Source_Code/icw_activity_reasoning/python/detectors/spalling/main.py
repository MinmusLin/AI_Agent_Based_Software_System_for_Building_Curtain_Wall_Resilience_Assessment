import json
import os
import time
from pathlib import Path

import torch
import torch.nn as nn
import torchvision.models as models
from PIL import Image
from torchvision import transforms


MODEL_PATH = Path(os.environ['REASONING_MODEL_PATH']).expanduser().resolve()
INPUT_SIZE = 224
CLASS_NAMES = ['defect', 'undefect']
DETECTOR: 'SpallingDetector | None' = None


# 获取当前可用的推理设备
def get_device() -> torch.device:
    return torch.device('cuda' if torch.cuda.is_available() else 'cpu')


# 构建模型网络结构
def build_model() -> nn.Module:
    model = models.resnet34(weights=None)
    num_features = model.fc.in_features
    model.fc = nn.Linear(num_features, len(CLASS_NAMES))
    return model


# 加载模型权重并切换为推理模式
def load_model(model_path: Path, device: torch.device) -> nn.Module:
    model = build_model()
    state_dict = torch.load(model_path, map_location=device)
    if any(key.startswith('module.') for key in state_dict):
        state_dict = {key.replace('module.', '', 1): value for key, value in state_dict.items()}
    model.load_state_dict(state_dict)
    model.to(device)
    model.eval()
    return model


# 读取并校验输入图像
def read_image(input_path: Path) -> Image.Image:
    if not input_path.exists():
        raise FileNotFoundError(f'input image not found: {input_path}')
    return Image.open(input_path).convert('RGB')


# 对输入图像执行模型预处理
def preprocess_image(image: Image.Image) -> torch.Tensor:
    transform = transforms.Compose(
        [
            transforms.Resize((INPUT_SIZE, INPUT_SIZE)),
            transforms.ToTensor(),
            transforms.Normalize(mean=[0.485, 0.456, 0.406], std=[0.229, 0.224, 0.225]),
        ]
    )
    return transform(image).unsqueeze(0)


# 执行模型推理并返回预测结果
def predict(model: nn.Module, image: Image.Image, device: torch.device) -> dict[str, float | bool]:
    tensor = preprocess_image(image).to(device)
    with torch.no_grad():
        logits_tensor = model(tensor)[0]
        probabilities_tensor = torch.softmax(logits_tensor, dim=0)
        confidence_tensor, predicted_tensor = torch.max(probabilities_tensor, dim=0)

    predicted_index = int(predicted_tensor.cpu().item())
    predicted_class = CLASS_NAMES[predicted_index]
    return {
        'has_spalling': bool(predicted_class == 'defect'),
        'confidence': round(float(confidence_tensor.cpu().item()), 6),
    }


# 构建检测结果报告数据
def build_report(prediction: dict[str, float | bool]) -> dict[str, float | bool]:
    return {
        'has_spalling': prediction['has_spalling'],
        'confidence': prediction['confidence'],
    }


class SpallingDetector:
    def __init__(self) -> None:
        self.device = get_device()
        self.model = load_model(MODEL_PATH, self.device)

    # 执行玻璃爆裂检测
    def detect(self, input_path: Path) -> None:
        start_time = time.time()
        input_path = input_path.expanduser().resolve()
        image = read_image(input_path)
        prediction = predict(self.model, image, self.device)
        report = build_report(prediction)
        report['runtime_seconds'] = round(time.time() - start_time, 3)
        run_dir = input_path.parent
        run_dir.mkdir(parents=True, exist_ok=True)
        report_path = run_dir / 'report.json'
        report_path.write_text(json.dumps(report, ensure_ascii=False, indent=2), encoding='utf-8')


# 获取可复用的玻璃爆裂检测器
def get_detector() -> SpallingDetector:
    global DETECTOR
    if DETECTOR is None:
        DETECTOR = SpallingDetector()
    return DETECTOR


# 执行玻璃爆裂检测
def detect(input_path: Path) -> None:
    get_detector().detect(input_path)
