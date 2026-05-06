import argparse
import json
import time
from pathlib import Path
from typing import Any

import cv2
import numpy as np
from ultralytics import YOLO


APP_ROOT = Path(__file__).resolve().parent
MODEL_PATH = APP_ROOT / 'model' / 'best_weights_model.pt'
IMAGE_SIZE = 640
CONFIDENCE_THRESHOLD = 0.25
IOU_THRESHOLD = 0.45


# 解析命令行输入参数
def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser()
    parser.add_argument('--input', required=True)
    return parser.parse_args()


# 读取并校验输入图像
def read_image(input_path: Path) -> np.ndarray:
    if not input_path.exists():
        raise FileNotFoundError(f'input image not found: {input_path}')
    image = cv2.imread(str(input_path), cv2.IMREAD_COLOR)
    if image is None:
        raise RuntimeError(f'input image cannot be read: {input_path}')
    return image


# 执行模型推理并返回预测结果
def predict(model: YOLO, image: np.ndarray) -> Any:
    return model.predict(
        source=image,
        imgsz=IMAGE_SIZE,
        conf=CONFIDENCE_THRESHOLD,
        iou=IOU_THRESHOLD,
        show=False,
        save=False,
        save_txt=False,
        save_conf=False,
        save_crop=False,
        show_labels=True,
        show_conf=True,
        retina_masks=True,
        verbose=False,
    )


# 将掩码转换为像素数量和占比指标
def mask_to_metrics(mask: np.ndarray, image_area: int) -> dict[str, Any]:
    binary_mask = (mask > 0.5).astype(np.uint8)
    mask_pixels = int(binary_mask.sum())
    return {
        'mask_pixels': mask_pixels,
        'mask_ratio': round(mask_pixels / image_area, 6) if image_area else 0.0,
    }


# 提取锈蚀检测区域结果
def extract_detections(result: Any, image_shape: tuple[int, int, int]) -> list[dict[str, Any]]:
    image_height, image_width = image_shape[:2]
    image_area = int(image_height * image_width)
    names = result.names if isinstance(result.names, dict) else {}
    boxes = result.boxes
    masks = result.masks
    detections = []

    if boxes is None or len(boxes) == 0:
        return detections

    mask_data = masks.data.cpu().numpy() if masks is not None and masks.data is not None else None

    for box_index, box in enumerate(boxes):
        cls_id = int(box.cls.item()) if box.cls is not None else None
        class_name = names.get(cls_id, str(cls_id)) if cls_id is not None else None
        if cls_id != 0 or str(class_name or '').strip().lower() != 'rust':
            continue

        xyxy = [round(float(value), 2) for value in box.xyxy[0].tolist()]
        detection: dict[str, Any] = {
            'id': len(detections) + 1,
            'confidence': round(float(box.conf.item()), 6) if box.conf is not None else None,
            'bbox_xyxy': xyxy,
            'mask_pixels': 0,
            'mask_ratio': 0.0,
        }

        if mask_data is not None and box_index < len(mask_data):
            detection.update(mask_to_metrics(mask_data[box_index], image_area))

        detections.append(detection)

    return detections


# 保存带检测标注的图像
def save_annotated_image(image: np.ndarray, result: Any | None, output_dir: Path) -> None:
    output_dir.mkdir(parents=True, exist_ok=True)
    annotated_image = result.plot() if result is not None else image
    output_path = output_dir / 'annotated.png'
    cv2.imwrite(str(output_path), annotated_image)


# 构建检测结果报告数据
def build_report(image: np.ndarray, detections: list[dict[str, Any]]) -> dict[str, Any]:
    image_height, image_width = image.shape[:2]
    rust_pixels = int(sum(item.get('mask_pixels', 0) for item in detections))
    image_area = int(image_height * image_width)
    confidences = [item['confidence'] for item in detections if item.get('confidence') is not None]
    return {
        'has_corrosion': bool(detections),
        'corrosion_count': len(detections),
        'max_confidence': round(float(max(confidences)), 6) if confidences else 0.0,
        'average_confidence': round(float(np.mean(confidences)), 6) if confidences else 0.0,
        'corrosion_pixels': rust_pixels,
        'corrosion_ratio': round(rust_pixels / image_area, 6) if image_area else 0.0,
        'regions': detections,
    }


# 执行金属锈蚀检测
def main() -> int:
    args = parse_args()
    start_time = time.time()
    input_path = Path(args.input).expanduser().resolve()

    image = read_image(input_path)
    model = YOLO(str(MODEL_PATH))
    raw_results = predict(model, image)
    result = raw_results[0] if raw_results and len(raw_results) > 0 else None
    detections = extract_detections(result, image.shape) if result is not None else []
    output_dir = input_path.parent
    save_annotated_image(image, result, output_dir)
    report = build_report(image, detections)
    report_path = output_dir / 'report.json'
    report['runtime_seconds'] = round(time.time() - start_time, 3)
    report_path.write_text(json.dumps(report, ensure_ascii=False, indent=2), encoding='utf-8')
    return 0


raise SystemExit(main())
