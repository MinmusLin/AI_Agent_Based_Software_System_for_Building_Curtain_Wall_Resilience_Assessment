from __future__ import annotations

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


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser()
    parser.add_argument('--input', required=True)
    return parser.parse_args()


def validate_model_file(model_path: Path) -> None:
    if not model_path.exists():
        raise FileNotFoundError(f'model file not found: {model_path}')
    if model_path.stat().st_size < 1024 * 1024:
        head = model_path.read_bytes()[:128]
        if b'git-lfs.github.com/spec' in head:
            raise RuntimeError(f'model file is a Git LFS pointer, not real weights: {model_path}')
        raise RuntimeError(f'model file size is invalid: {model_path}')


def read_image(input_path: Path) -> np.ndarray:
    if not input_path.exists():
        raise FileNotFoundError(f'input image not found: {input_path}')
    image = cv2.imread(str(input_path), cv2.IMREAD_COLOR)
    if image is None:
        raise RuntimeError(f'input image cannot be read: {input_path}')
    return image


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


def mask_to_metrics(mask: np.ndarray, image_area: int) -> dict[str, Any]:
    binary_mask = (mask > 0.5).astype(np.uint8)
    mask_pixels = int(binary_mask.sum())
    contours, _ = cv2.findContours(binary_mask, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
    contour_area = float(sum(cv2.contourArea(contour) for contour in contours))
    return {
        'mask_pixels': mask_pixels,
        'mask_ratio': round(mask_pixels / image_area, 6) if image_area else 0.0,
        'mask_ratio_percent': round(mask_pixels / image_area * 100, 4) if image_area else 0.0,
        'contour_area_px': round(contour_area, 2),
        'contour_count': len(contours),
    }


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
    mask_polygons = masks.xy if masks is not None and masks.xy is not None else []

    for index, box in enumerate(boxes, start=1):
        cls_id = int(box.cls.item()) if box.cls is not None else None
        xyxy = [round(float(value), 2) for value in box.xyxy[0].tolist()]
        xywh = [round(float(value), 2) for value in box.xywh[0].tolist()]
        x1, y1, x2, y2 = xyxy
        bbox_width = max(0.0, x2 - x1)
        bbox_height = max(0.0, y2 - y1)
        detection: dict[str, Any] = {
            'id': index,
            'confidence': round(float(box.conf.item()), 6) if box.conf is not None else None,
            'class_id': cls_id,
            'class_name': names.get(cls_id, str(cls_id)) if cls_id is not None else None,
            'bbox_xyxy': xyxy,
            'bbox_xywh': xywh,
            'bbox_area_px': round(float(bbox_width * bbox_height), 2),
            'bbox_ratio': round(float(bbox_width * bbox_height) / image_area, 6) if image_area else 0.0,
        }

        if mask_data is not None and index - 1 < len(mask_data):
            detection.update(mask_to_metrics(mask_data[index - 1], image_area))
        if index - 1 < len(mask_polygons):
            polygon = np.asarray(mask_polygons[index - 1], dtype=np.float32)
            detection['polygon_point_count'] = int(len(polygon))
            detection['polygon'] = [
                {'x': round(float(x), 2), 'y': round(float(y), 2)}
                for x, y in polygon
            ]

        detections.append(detection)

    return detections


def save_annotated_image(result: Any, output_dir: Path, input_path: Path) -> str:
    output_dir.mkdir(parents=True, exist_ok=True)
    annotated_image = result.plot()
    output_path = output_dir / 'annotated.png'
    cv2.imwrite(str(output_path), annotated_image)
    return str(output_path)


def build_report(input_path: Path, image: np.ndarray, result: Any, detections: list[dict[str, Any]]) -> dict[str, Any]:
    image_height, image_width = image.shape[:2]
    rust_pixels = int(sum(item.get('mask_pixels', 0) for item in detections))
    image_area = int(image_height * image_width)
    confidences = [item['confidence'] for item in detections if item.get('confidence') is not None]
    return {
        'input': str(input_path),
        'image_size': {'width': image_width, 'height': image_height},
        'result': 'corrosion_detected' if detections else 'no_corrosion_detected',
        'has_corrosion': bool(detections),
        'detection_count': len(detections),
        'max_confidence': round(float(max(confidences)), 6) if confidences else None,
        'average_confidence': round(float(np.mean(confidences)), 6) if confidences else None,
        'corrosion_pixels': rust_pixels,
        'corrosion_ratio': round(rust_pixels / image_area, 6) if image_area else 0.0,
        'corrosion_ratio_percent': round(rust_pixels / image_area * 100, 4) if image_area else 0.0,
        'class_names': result.names if isinstance(result.names, dict) else {},
        'detections': detections,
    }


def main() -> int:
    args = parse_args()
    start_time = time.time()
    input_path = Path(args.input).expanduser().resolve()
    validate_model_file(MODEL_PATH)

    image = read_image(input_path)
    model = YOLO(str(MODEL_PATH))
    raw_results = predict(model, image)
    if not raw_results or len(raw_results) == 0:
        raise RuntimeError('model returned no result')

    result = raw_results[0]
    detections = extract_detections(result, image.shape)
    output_dir = input_path.parent
    annotated_path = save_annotated_image(result, output_dir, input_path)
    report = build_report(input_path, image, result, detections)
    report_path = output_dir / 'report.json'
    report['model'] = str(MODEL_PATH)
    report['image_size_for_model'] = IMAGE_SIZE
    report['confidence_threshold'] = CONFIDENCE_THRESHOLD
    report['iou_threshold'] = IOU_THRESHOLD
    report['runtime_seconds'] = round(time.time() - start_time, 3)
    report['outputs'] = {
        'annotated_image': annotated_path,
        'report_json': str(report_path),
    }
    report_path.write_text(json.dumps(report, ensure_ascii=False, indent=2), encoding='utf-8')
    return 0


if __name__ == '__main__':
    raise SystemExit(main())
