from __future__ import annotations

import argparse
import json
import math
import os
import time
from pathlib import Path
from typing import Any

CACHE_ROOT = Path('/tmp/stain-detection-cache')
os.environ.setdefault('MPLCONFIGDIR', str(CACHE_ROOT / 'matplotlib'))
os.environ.setdefault('YOLO_CONFIG_DIR', str(CACHE_ROOT / 'ultralytics'))
for cache_path in (os.environ['MPLCONFIGDIR'], os.environ['YOLO_CONFIG_DIR']):
    Path(cache_path).mkdir(parents=True, exist_ok=True)

import cv2
import numpy as np
from ultralytics import YOLO

APP_ROOT = Path(__file__).resolve().parent
MODEL_PATH = APP_ROOT / 'model' / 'best_weights_model.pt'
IMAGE_SIZE = 640
CONFIDENCE_THRESHOLD = 0.3
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


def sort_corners(points: np.ndarray) -> np.ndarray:
    points = points.reshape(-1, 2).astype(np.float32)
    ordered = np.zeros((4, 2), dtype=np.float32)
    point_sum = points.sum(axis=1)
    point_diff = np.diff(points, axis=1).reshape(-1)
    ordered[0] = points[np.argmin(point_sum)]
    ordered[2] = points[np.argmax(point_sum)]
    ordered[1] = points[np.argmin(point_diff)]
    ordered[3] = points[np.argmax(point_diff)]
    return ordered


def polygon_to_corners(polygon: np.ndarray) -> np.ndarray | None:
    if polygon is None or len(polygon) < 4:
        return None
    contour = polygon.reshape(-1, 1, 2).astype(np.float32)
    hull = cv2.convexHull(contour)
    perimeter = cv2.arcLength(hull, True)
    approx = cv2.approxPolyDP(hull, 0.02 * perimeter, True)
    if len(approx) == 4:
        return sort_corners(approx.reshape(4, 2))
    rect = cv2.minAreaRect(hull)
    return sort_corners(cv2.boxPoints(rect))


def warp_block(image: np.ndarray, corners: np.ndarray) -> tuple[np.ndarray, int, int]:
    width_top = cv2.norm(corners[0], corners[1])
    width_bottom = cv2.norm(corners[3], corners[2])
    height_left = cv2.norm(corners[0], corners[3])
    height_right = cv2.norm(corners[1], corners[2])
    output_width = int(max(width_top, width_bottom))
    output_height = int(max(height_left, height_right))
    if output_width <= 0 or output_height <= 0:
        raise ValueError('curtain wall block size is invalid')

    target = np.float32(
        [[0, 0], [output_width - 1, 0], [output_width - 1, output_height - 1], [0, output_height - 1]]
    )
    matrix = cv2.getPerspectiveTransform(corners, target)
    warped_image = cv2.warpPerspective(image, matrix, (output_width, output_height))
    return warped_image, output_width, output_height


def partition_image(image: np.ndarray, num_rows: int, num_cols: int) -> list[np.ndarray]:
    height, width = image.shape[:2]
    block_height = height // num_rows
    block_width = width // num_cols
    return [
        image[row * block_height:(row + 1) * block_height, col * block_width:(col + 1) * block_width]
        for row in range(num_rows)
        for col in range(num_cols)
    ]


def calculate_stain(block_image: np.ndarray, output_path: Path) -> dict[str, Any]:
    blocks = partition_image(block_image, 4, 4)
    mean_values = [float(np.mean(block[:, :, 0])) for block in blocks if block.size > 0]
    variance = float(np.var(mean_values)) if mean_values else 0.0

    if variance <= 40:
        return {
            'has_stain': False,
            'stain_ratio': 0.0,
            'stain_percentage': 0.0,
            'variance': round(variance, 6),
            'result_image': None,
            'status': 'low_variance_skipped',
        }

    gray_image = cv2.cvtColor(block_image, cv2.COLOR_BGR2GRAY)
    denoised_image = cv2.fastNlMeansDenoising(gray_image, None, h=10, templateWindowSize=7, searchWindowSize=21)
    threshold_value = float(cv2.mean(denoised_image)[0])
    _, binary_image = cv2.threshold(denoised_image, threshold_value, 255, cv2.THRESH_BINARY)

    binary_bgr = cv2.cvtColor(binary_image, cv2.COLOR_GRAY2BGR)
    result_image = cv2.addWeighted(binary_bgr, 0.5, block_image, 0.5, 0)
    white_pixels = int(cv2.countNonZero(binary_image))
    total_pixels = int(binary_image.size)
    stain_ratio = 1 - white_pixels / total_pixels if total_pixels else 0.0

    output_path.parent.mkdir(parents=True, exist_ok=True)
    cv2.imwrite(str(output_path), result_image)
    return {
        'has_stain': bool(stain_ratio > 0),
        'stain_ratio': round(float(stain_ratio), 6),
        'stain_percentage': round(float(stain_ratio * 100), 4),
        'variance': round(variance, 6),
        'threshold': round(threshold_value, 4),
        'result_image': str(output_path),
        'status': 'processed',
    }


def predict(model: YOLO, image_path: Path) -> Any:
    return model.predict(
        source=str(image_path),
        imgsz=IMAGE_SIZE,
        conf=CONFIDENCE_THRESHOLD,
        iou=IOU_THRESHOLD,
        show=False,
        save=False,
        save_txt=False,
        save_conf=False,
        save_crop=False,
        show_labels=False,
        show_conf=True,
        retina_masks=True,
        show_boxes=False,
        verbose=False,
    )


def build_outputs(input_path: Path, image: np.ndarray, results: Any, output_dir: Path) -> dict[str, Any]:
    if not results or len(results) == 0:
        return {
            'status': 'failed',
            'input': str(input_path),
            'error_code': 'NO_RESULTS',
            'message': 'model returned no result.',
        }

    result = results[0]
    masks = result.masks
    if masks is None or masks.xy is None or len(masks.xy) == 0:
        return {
            'status': 'failed',
            'input': str(input_path),
            'error_code': 'NO_MASKS_DETECTED',
            'message': 'no curtain wall block detected in image.',
        }

    block_dir = output_dir / 'blocks'
    result_dir = output_dir / 'results'
    block_dir.mkdir(parents=True, exist_ok=True)
    result_dir.mkdir(parents=True, exist_ok=True)

    annotated_image = image.copy()
    detections = []
    boxes = result.boxes
    names = result.names if isinstance(result.names, dict) else {}

    for index, polygon in enumerate(masks.xy, start=1):
        corners = polygon_to_corners(np.asarray(polygon, dtype=np.float32))
        if corners is None:
            continue

        try:
            warped_image, output_width, output_height = warp_block(image, corners)
        except ValueError:
            continue

        block_path = block_dir / f'block-{index}.jpg'
        stain_result_path = result_dir / f'stain-{index}.jpg'
        cv2.imwrite(str(block_path), warped_image)
        stain_metrics = calculate_stain(warped_image, stain_result_path)

        cv2.polylines(annotated_image, [corners.astype(np.int32)], isClosed=True, color=(0, 255, 0), thickness=2)
        cv2.putText(
            annotated_image,
            str(index),
            tuple(corners[0].astype(int)),
            cv2.FONT_HERSHEY_SIMPLEX,
            0.6,
            (255, 255, 255),
            2,
        )

        box_data: dict[str, Any] = {}
        if boxes is not None and len(boxes) >= index:
            box = boxes[index - 1]
            cls_id = int(box.cls.item()) if box.cls is not None else None
            box_data = {
                'confidence': round(float(box.conf.item()), 6) if box.conf is not None else None,
                'class_id': cls_id,
                'class_name': names.get(cls_id, str(cls_id)) if cls_id is not None else None,
                'bbox_xyxy': [round(float(value), 2) for value in box.xyxy[0].tolist()],
            }

        detections.append(
            {
                'block_number': index,
                **box_data,
                'corners': [{'x': round(float(x), 2), 'y': round(float(y), 2)} for x, y in corners],
                'block_size': {'width': output_width, 'height': output_height},
                'block_area_px': int(output_width * output_height),
                'block_image': str(block_path),
                **stain_metrics,
            }
        )

    if not detections:
        return {
            'status': 'failed',
            'input': str(input_path),
            'error_code': 'PROCESSING_FAILED',
            'message': 'model detected masks but no block was processed.',
        }

    annotated_path = output_dir / 'annotated.png'
    cv2.imwrite(str(annotated_path), annotated_image)
    stain_percentages = [item['stain_percentage'] for item in detections]

    return {
        'status': 'success',
        'input': str(input_path),
        'image_size': {'width': int(image.shape[1]), 'height': int(image.shape[0])},
        'total_blocks': len(detections),
        'average_stain_percentage': round(float(np.mean(stain_percentages)), 4),
        'max_stain_percentage': round(float(np.max(stain_percentages)), 4),
        'annotated_image': str(annotated_path),
        'detections': detections,
    }


def main() -> int:
    args = parse_args()
    start_time = time.time()
    input_path = Path(args.input).expanduser().resolve()
    validate_model_file(MODEL_PATH)

    image = read_image(input_path)
    model = YOLO(str(MODEL_PATH))
    raw_results = predict(model, input_path)

    output_dir = input_path.parent
    report = build_outputs(input_path, image, raw_results, output_dir)
    output_dir.mkdir(parents=True, exist_ok=True)
    report_path = output_dir / 'report.json'
    report['model'] = str(MODEL_PATH)
    report['image_size_for_model'] = IMAGE_SIZE
    report['confidence_threshold'] = CONFIDENCE_THRESHOLD
    report['iou_threshold'] = IOU_THRESHOLD
    report['runtime_seconds'] = round(time.time() - start_time, 3)
    report['outputs'] = {'report_json': str(report_path)}
    report_path.write_text(json.dumps(report, ensure_ascii=False, indent=2), encoding='utf-8')
    return 0


if __name__ == '__main__':
    raise SystemExit(main())
