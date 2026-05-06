import argparse
import json
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


# 按透视变换需要排序四个角点
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


# 将检测多边形转换为四角点
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


# 按四角点透视矫正幕墙块
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


# 将图像切分为固定网格块
def partition_image(image: np.ndarray, num_rows: int, num_cols: int) -> list[np.ndarray]:
    height, width = image.shape[:2]
    block_height = height // num_rows
    block_width = width // num_cols
    return [
        image[row * block_height:(row + 1) * block_height, col * block_width:(col + 1) * block_width]
        for row in range(num_rows)
        for col in range(num_cols)
    ]


# 计算矫正区域中的污渍指标
def calculate_stain(block_image: np.ndarray, output_path: Path) -> dict[str, Any]:
    blocks = partition_image(block_image, 4, 4)
    mean_values = [float(np.mean(block[:, :, 0])) for block in blocks if block.size > 0]
    variance = float(np.var(mean_values)) if mean_values else 0.0

    if variance <= 40:
        return {
            'has_stain': False,
            'stain_pixels': 0,
            'stain_ratio': 0.0,
        }

    gray_image = cv2.cvtColor(block_image, cv2.COLOR_BGR2GRAY)
    denoised_image = cv2.fastNlMeansDenoising(gray_image, None, h=10, templateWindowSize=7, searchWindowSize=21)
    threshold_value = float(cv2.mean(denoised_image)[0])
    _, binary_image = cv2.threshold(denoised_image, threshold_value, 255, cv2.THRESH_BINARY)

    binary_bgr = cv2.cvtColor(binary_image, cv2.COLOR_GRAY2BGR)
    result_image = cv2.addWeighted(binary_bgr, 0.5, block_image, 0.5, 0)
    white_pixels = int(cv2.countNonZero(binary_image))
    total_pixels = int(binary_image.size)
    stain_pixels = max(total_pixels - white_pixels, 0)
    stain_ratio = stain_pixels / total_pixels if total_pixels else 0.0
    if stain_ratio <= 0:
        return {
            'has_stain': False,
            'stain_pixels': 0,
            'stain_ratio': 0.0,
        }

    output_path.parent.mkdir(parents=True, exist_ok=True)
    cv2.imwrite(str(output_path), result_image)
    return {
        'has_stain': True,
        'stain_pixels': stain_pixels,
        'stain_ratio': round(float(stain_ratio), 6),
    }


# 执行模型推理并返回预测结果
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


# 构建无污渍结果的默认报告
def empty_report() -> dict[str, Any]:
    return {
        'has_stain': False,
        'stain_count': 0,
        'average_stain_ratio': 0.0,
        'max_stain_ratio': 0.0,
        'regions': [],
    }


# 生成污渍检测产物和报告数据
def build_outputs(image: np.ndarray, results: Any, output_dir: Path) -> dict[str, Any]:
    annotated_path = output_dir / 'annotated.png'
    if not results or len(results) == 0:
        cv2.imwrite(str(annotated_path), image)
        return empty_report()

    result = results[0]
    masks = result.masks
    if masks is None or masks.xy is None or len(masks.xy) == 0:
        cv2.imwrite(str(annotated_path), image)
        return empty_report()

    annotated_image = image.copy()
    regions = []
    boxes = result.boxes

    for index, polygon in enumerate(masks.xy, start=1):
        corners = polygon_to_corners(np.asarray(polygon, dtype=np.float32))
        if corners is None:
            continue

        try:
            warped_image, output_width, output_height = warp_block(image, corners)
        except ValueError:
            continue

        region_id = len(regions) + 1
        block_path = output_dir / f'block_{region_id}.png'
        stain_result_path = output_dir / f'overlay_{region_id}.png'
        cv2.imwrite(str(block_path), warped_image)
        stain_metrics = calculate_stain(warped_image, stain_result_path)
        if not stain_metrics.get('has_stain') or not stain_result_path.exists():
            if block_path.exists():
                block_path.unlink()
            if stain_result_path.exists():
                stain_result_path.unlink()
            continue

        cv2.polylines(annotated_image, [corners.astype(np.int32)], isClosed=True, color=(0, 255, 0), thickness=2)
        cv2.putText(
            annotated_image,
            str(region_id),
            tuple(corners[0].astype(int)),
            cv2.FONT_HERSHEY_SIMPLEX,
            0.6,
            (255, 255, 255),
            2,
        )

        confidence = 0.0
        bbox_xyxy = [0.0, 0.0, 0.0, 0.0]
        if boxes is not None and len(boxes) >= index:
            box = boxes[index - 1]
            confidence = round(float(box.conf.item()), 6) if box.conf is not None else None
            bbox_xyxy = [round(float(value), 2) for value in box.xyxy[0].tolist()]

        regions.append(
            {
                'id': region_id,
                'confidence': confidence,
                'bbox_xyxy': bbox_xyxy,
                'region_width': output_width,
                'region_height': output_height,
                'stain_pixels': stain_metrics['stain_pixels'],
                'stain_ratio': stain_metrics['stain_ratio'],
            }
        )

    if not regions:
        cv2.imwrite(str(annotated_path), annotated_image)
        return empty_report()

    cv2.imwrite(str(annotated_path), annotated_image)
    stain_ratios = [item['stain_ratio'] for item in regions]

    return {
        'has_stain': True,
        'stain_count': len(regions),
        'average_stain_ratio': round(float(np.mean(stain_ratios) * 100), 4),
        'max_stain_ratio': round(float(np.max(stain_ratios) * 100), 4),
        'regions': regions,
    }


# 执行石材污渍检测
def main() -> int:
    args = parse_args()
    start_time = time.time()
    input_path = Path(args.input).expanduser().resolve()

    image = read_image(input_path)
    model = YOLO(str(MODEL_PATH))
    raw_results = predict(model, input_path)

    output_dir = input_path.parent
    output_dir.mkdir(parents=True, exist_ok=True)
    report = build_outputs(image, raw_results, output_dir)
    report_path = output_dir / 'report.json'
    report['runtime_seconds'] = round(time.time() - start_time, 3)
    report_path.write_text(json.dumps(report, ensure_ascii=False, indent=2), encoding='utf-8')
    return 0


raise SystemExit(main())
