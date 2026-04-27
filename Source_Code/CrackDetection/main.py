#!/usr/bin/env python3

from __future__ import annotations

import argparse
import json
import math
import time
from functools import partial
from pathlib import Path
from typing import Any

import cv2
import numpy as np
import torch
import torch.nn as nn
import torch.nn.functional as torch_functional
from PIL import Image


APP_ROOT = Path(__file__).resolve().parent
MODEL_PATH = APP_ROOT / 'model' / 'best_epoch_weights.pth'
OUTPUT_ROOT = APP_ROOT / 'outputs'
INPUT_SIZE = 1024
MIN_REGION_AREA = 20
NUM_CLASSES = 2


class Gelu(nn.Module):
    def forward(self, tensor: torch.Tensor) -> torch.Tensor:
        return 0.5 * tensor * (1 + torch.tanh(np.sqrt(2 / np.pi) * (tensor + 0.044715 * torch.pow(tensor, 3))))


class OverlapPatchEmbed(nn.Module):
    def __init__(self, patch_size: int, stride: int, in_channels: int, embed_dim: int) -> None:
        super().__init__()
        self.proj = nn.Conv2d(
            in_channels,
            embed_dim,
            kernel_size=(patch_size, patch_size),
            stride=stride,
            padding=(patch_size // 2, patch_size // 2),
        )
        self.norm = nn.LayerNorm(embed_dim)

    def forward(self, tensor: torch.Tensor) -> tuple[torch.Tensor, int, int]:
        tensor = self.proj(tensor)
        _, _, height, width = tensor.shape
        tensor = tensor.flatten(2).transpose(1, 2)
        tensor = self.norm(tensor)
        return tensor, height, width


class Attention(nn.Module):
    def __init__(
        self,
        dim: int,
        num_heads: int,
        qkv_bias: bool,
        qk_scale: float | None,
        attn_drop: float,
        proj_drop: float,
        sr_ratio: int,
    ) -> None:
        super().__init__()
        head_dim = dim // num_heads
        self.num_heads = num_heads
        self.scale = qk_scale or head_dim ** -0.5
        self.q = nn.Linear(dim, dim, bias=qkv_bias)
        self.sr_ratio = sr_ratio
        if sr_ratio > 1:
            self.sr = nn.Conv2d(dim, dim, kernel_size=sr_ratio, stride=sr_ratio)
            self.norm = nn.LayerNorm(dim)
        self.kv = nn.Linear(dim, dim * 2, bias=qkv_bias)
        self.attn_drop = nn.Dropout(attn_drop)
        self.proj = nn.Linear(dim, dim)
        self.proj_drop = nn.Dropout(proj_drop)

    def forward(self, tensor: torch.Tensor, height: int, width: int) -> torch.Tensor:
        batch_size, point_count, channel_count = tensor.shape
        query = self.q(tensor).reshape(
            batch_size,
            point_count,
            self.num_heads,
            channel_count // self.num_heads,
        ).permute(0, 2, 1, 3)

        if self.sr_ratio > 1:
            source = tensor.permute(0, 2, 1).reshape(batch_size, channel_count, height, width)
            source = self.sr(source).reshape(batch_size, channel_count, -1).permute(0, 2, 1)
            source = self.norm(source)
            key_value = self.kv(source).reshape(
                batch_size,
                -1,
                2,
                self.num_heads,
                channel_count // self.num_heads,
            ).permute(2, 0, 3, 1, 4)
        else:
            key_value = self.kv(tensor).reshape(
                batch_size,
                -1,
                2,
                self.num_heads,
                channel_count // self.num_heads,
            ).permute(2, 0, 3, 1, 4)

        key, value = key_value[0], key_value[1]
        attention = (query @ key.transpose(-2, -1)) * self.scale
        attention = attention.softmax(dim=-1)
        attention = self.attn_drop(attention)
        tensor = (attention @ value).transpose(1, 2).reshape(batch_size, point_count, channel_count)
        tensor = self.proj(tensor)
        return self.proj_drop(tensor)


class DepthwiseConv(nn.Module):
    def __init__(self, dim: int) -> None:
        super().__init__()
        self.dwconv = nn.Conv2d(dim, dim, 3, 1, 1, bias=True, groups=dim)

    def forward(self, tensor: torch.Tensor, height: int, width: int) -> torch.Tensor:
        batch_size, _, channel_count = tensor.shape
        tensor = tensor.transpose(1, 2).view(batch_size, channel_count, height, width)
        tensor = self.dwconv(tensor)
        return tensor.flatten(2).transpose(1, 2)


class FeedForward(nn.Module):
    def __init__(self, in_features: int, hidden_features: int, drop: float) -> None:
        super().__init__()
        self.fc1 = nn.Linear(in_features, hidden_features)
        self.dwconv = DepthwiseConv(hidden_features)
        self.act = Gelu()
        self.fc2 = nn.Linear(hidden_features, in_features)
        self.drop = nn.Dropout(drop)

    def forward(self, tensor: torch.Tensor, height: int, width: int) -> torch.Tensor:
        tensor = self.fc1(tensor)
        tensor = self.dwconv(tensor, height, width)
        tensor = self.act(tensor)
        tensor = self.drop(tensor)
        tensor = self.fc2(tensor)
        return self.drop(tensor)


class TransformerBlock(nn.Module):
    def __init__(
        self,
        dim: int,
        num_heads: int,
        mlp_ratio: int,
        qkv_bias: bool,
        qk_scale: float | None,
        drop: float,
        attn_drop: float,
        norm_layer: type[nn.LayerNorm],
        sr_ratio: int,
    ) -> None:
        super().__init__()
        self.norm1 = norm_layer(dim)
        self.attn = Attention(dim, num_heads, qkv_bias, qk_scale, attn_drop, drop, sr_ratio)
        self.norm2 = norm_layer(dim)
        self.mlp = FeedForward(dim, int(dim * mlp_ratio), drop)
        self.drop_path = nn.Identity()

    def forward(self, tensor: torch.Tensor, height: int, width: int) -> torch.Tensor:
        tensor = tensor + self.drop_path(self.attn(self.norm1(tensor), height, width))
        tensor = tensor + self.drop_path(self.mlp(self.norm2(tensor), height, width))
        return tensor


class MixVisionTransformer(nn.Module):
    def __init__(self) -> None:
        super().__init__()
        embed_dims = [64, 128, 320, 512]
        num_heads = [1, 2, 5, 8]
        mlp_ratios = [4, 4, 4, 4]
        depths = [3, 4, 6, 3]
        sr_ratios = [8, 4, 2, 1]
        norm_layer = partial(nn.LayerNorm, eps=1e-6)

        self.patch_embed1 = OverlapPatchEmbed(7, 4, 3, embed_dims[0])
        self.patch_embed2 = OverlapPatchEmbed(3, 2, embed_dims[0], embed_dims[1])
        self.patch_embed3 = OverlapPatchEmbed(3, 2, embed_dims[1], embed_dims[2])
        self.patch_embed4 = OverlapPatchEmbed(3, 2, embed_dims[2], embed_dims[3])

        self.block1 = self.build_blocks(embed_dims, num_heads, mlp_ratios, depths, sr_ratios, norm_layer, 0)
        self.block2 = self.build_blocks(embed_dims, num_heads, mlp_ratios, depths, sr_ratios, norm_layer, 1)
        self.block3 = self.build_blocks(embed_dims, num_heads, mlp_ratios, depths, sr_ratios, norm_layer, 2)
        self.block4 = self.build_blocks(embed_dims, num_heads, mlp_ratios, depths, sr_ratios, norm_layer, 3)

        self.norm1 = norm_layer(embed_dims[0])
        self.norm2 = norm_layer(embed_dims[1])
        self.norm3 = norm_layer(embed_dims[2])
        self.norm4 = norm_layer(embed_dims[3])

    @staticmethod
    def build_blocks(
        embed_dims: list[int],
        num_heads: list[int],
        mlp_ratios: list[int],
        depths: list[int],
        sr_ratios: list[int],
        norm_layer: type[nn.LayerNorm],
        index: int,
    ) -> nn.ModuleList:
        return nn.ModuleList(
            [
                TransformerBlock(
                    dim=embed_dims[index],
                    num_heads=num_heads[index],
                    mlp_ratio=mlp_ratios[index],
                    qkv_bias=True,
                    qk_scale=None,
                    drop=0.0,
                    attn_drop=0.0,
                    norm_layer=norm_layer,
                    sr_ratio=sr_ratios[index],
                )
                for _ in range(depths[index])
            ]
        )

    def run_stage(
        self,
        tensor: torch.Tensor,
        patch_embed: OverlapPatchEmbed,
        blocks: nn.ModuleList,
        norm: nn.LayerNorm,
    ) -> tuple[torch.Tensor, int, int]:
        tensor, height, width = patch_embed(tensor)
        for block in blocks:
            tensor = block(tensor, height, width)
        tensor = norm(tensor)
        return tensor, height, width

    def forward(self, tensor: torch.Tensor) -> list[torch.Tensor]:
        batch_size = tensor.shape[0]
        outputs = []

        tensor, height, width = self.run_stage(tensor, self.patch_embed1, self.block1, self.norm1)
        tensor = tensor.reshape(batch_size, height, width, -1).permute(0, 3, 1, 2).contiguous()
        outputs.append(tensor)

        tensor, height, width = self.run_stage(tensor, self.patch_embed2, self.block2, self.norm2)
        tensor = tensor.reshape(batch_size, height, width, -1).permute(0, 3, 1, 2).contiguous()
        outputs.append(tensor)

        tensor, height, width = self.run_stage(tensor, self.patch_embed3, self.block3, self.norm3)
        tensor = tensor.reshape(batch_size, height, width, -1).permute(0, 3, 1, 2).contiguous()
        outputs.append(tensor)

        tensor, height, width = self.run_stage(tensor, self.patch_embed4, self.block4, self.norm4)
        tensor = tensor.reshape(batch_size, height, width, -1).permute(0, 3, 1, 2).contiguous()
        outputs.append(tensor)
        return outputs


class LinearEmbedding(nn.Module):
    def __init__(self, input_dim: int, embed_dim: int) -> None:
        super().__init__()
        self.proj = nn.Linear(input_dim, embed_dim)

    def forward(self, tensor: torch.Tensor) -> torch.Tensor:
        tensor = tensor.flatten(2).transpose(1, 2)
        return self.proj(tensor)


class ConvModule(nn.Module):
    def __init__(self, in_channels: int, out_channels: int) -> None:
        super().__init__()
        self.conv = nn.Conv2d(in_channels, out_channels, 1, 1, 0, groups=1, bias=False)
        self.bn = nn.BatchNorm2d(out_channels, eps=0.001, momentum=0.03)
        self.act = nn.ReLU()

    def forward(self, tensor: torch.Tensor) -> torch.Tensor:
        return self.act(self.bn(self.conv(tensor)))


class SegFormerHead(nn.Module):
    def __init__(self) -> None:
        super().__init__()
        in_channels = [64, 128, 320, 512]
        embedding_dim = 768
        self.linear_c4 = LinearEmbedding(in_channels[3], embedding_dim)
        self.linear_c3 = LinearEmbedding(in_channels[2], embedding_dim)
        self.linear_c2 = LinearEmbedding(in_channels[1], embedding_dim)
        self.linear_c1 = LinearEmbedding(in_channels[0], embedding_dim)
        self.linear_fuse = ConvModule(embedding_dim * 4, embedding_dim)
        self.linear_pred = nn.Conv2d(embedding_dim, NUM_CLASSES, kernel_size=1)
        self.dropout = nn.Dropout2d(0.1)

    def forward(self, inputs: list[torch.Tensor]) -> torch.Tensor:
        c1, c2, c3, c4 = inputs
        batch_size = c4.shape[0]

        c4 = self.linear_c4(c4).permute(0, 2, 1).reshape(batch_size, -1, c4.shape[2], c4.shape[3])
        c4 = torch_functional.interpolate(c4, size=c1.size()[2:], mode='bilinear', align_corners=False)

        c3 = self.linear_c3(c3).permute(0, 2, 1).reshape(batch_size, -1, c3.shape[2], c3.shape[3])
        c3 = torch_functional.interpolate(c3, size=c1.size()[2:], mode='bilinear', align_corners=False)

        c2 = self.linear_c2(c2).permute(0, 2, 1).reshape(batch_size, -1, c2.shape[2], c2.shape[3])
        c2 = torch_functional.interpolate(c2, size=c1.size()[2:], mode='bilinear', align_corners=False)

        c1 = self.linear_c1(c1).permute(0, 2, 1).reshape(batch_size, -1, c1.shape[2], c1.shape[3])
        tensor = self.linear_fuse(torch.cat([c4, c3, c2, c1], dim=1))
        tensor = self.dropout(tensor)
        return self.linear_pred(tensor)


class SegFormer(nn.Module):
    def __init__(self) -> None:
        super().__init__()
        self.backbone = MixVisionTransformer()
        self.decode_head = SegFormerHead()

    def forward(self, tensor: torch.Tensor) -> torch.Tensor:
        height, width = tensor.size(2), tensor.size(3)
        tensor = self.backbone(tensor)
        tensor = self.decode_head(tensor)
        return torch_functional.interpolate(tensor, size=(height, width), mode='bilinear', align_corners=True)


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description='本地幕墙裂缝检测')
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


def load_model(model_path: Path, device: torch.device) -> SegFormer:
    validate_model_file(model_path)
    model = SegFormer()
    state_dict = torch.load(model_path, map_location=device)
    if any(key.startswith('module.') for key in state_dict):
        state_dict = {key.replace('module.', '', 1): value for key, value in state_dict.items()}
    model.load_state_dict(state_dict)
    model.to(device)
    model.eval()
    return model


def convert_to_rgb(image: Image.Image) -> Image.Image:
    if len(np.shape(image)) == 3 and np.shape(image)[2] == 3:
        return image
    return image.convert('RGB')


def resize_image(image: Image.Image, size: tuple[int, int]) -> tuple[Image.Image, int, int]:
    image_width, image_height = image.size
    target_width, target_height = size
    scale = min(target_width / image_width, target_height / image_height)
    new_width = int(image_width * scale)
    new_height = int(image_height * scale)
    image = image.resize((new_width, new_height), Image.BICUBIC)
    resized_image = Image.new('RGB', size, (128, 128, 128))
    resized_image.paste(image, ((target_width - new_width) // 2, (target_height - new_height) // 2))
    return resized_image, new_width, new_height


def preprocess_image(image: np.ndarray) -> np.ndarray:
    image -= np.array([123.675, 116.28, 103.53], np.float32)
    image /= np.array([58.395, 57.12, 57.375], np.float32)
    return image


def read_image(input_path: Path) -> Image.Image:
    if not input_path.exists():
        raise FileNotFoundError(f'未找到输入图片: {input_path}')
    return convert_to_rgb(Image.open(input_path))


def predict_mask(model: SegFormer, image: Image.Image, device: torch.device) -> np.ndarray:
    original_height, original_width = np.array(image).shape[:2]
    resized_image, new_width, new_height = resize_image(image, (INPUT_SIZE, INPUT_SIZE))
    image_data = preprocess_image(np.array(resized_image, np.float32))
    image_data = np.expand_dims(np.transpose(image_data, (2, 0, 1)), 0)

    with torch.no_grad():
        tensor = torch.from_numpy(image_data).to(device)
        prediction = model(tensor)[0]
        prediction = torch_functional.softmax(prediction.permute(1, 2, 0), dim=-1).cpu().numpy()

    top = int((INPUT_SIZE - new_height) // 2)
    left = int((INPUT_SIZE - new_width) // 2)
    prediction = prediction[top:top + new_height, left:left + new_width]
    prediction = cv2.resize(prediction, (original_width, original_height), interpolation=cv2.INTER_LINEAR)
    return prediction.argmax(axis=-1).astype(np.uint8)


def extract_regions(mask: np.ndarray) -> tuple[np.ndarray, list[dict[str, Any]]]:
    binary_mask = (mask == 1).astype(np.uint8)
    region_count, labels, stats, centroids = cv2.connectedComponentsWithStats(binary_mask, connectivity=8)
    cleaned_mask = np.zeros_like(binary_mask)
    regions = []

    for label in range(1, region_count):
        x, y, width, height, area = stats[label]
        if int(area) < MIN_REGION_AREA:
            continue
        cleaned_mask[labels == label] = 1
        regions.append(
            {
                'id': len(regions) + 1,
                'area_px': int(area),
                'bbox': {
                    'x': int(x),
                    'y': int(y),
                    'width': int(width),
                    'height': int(height),
                },
                'bbox_diagonal_px': round(math.sqrt(float(width * width + height * height)), 2),
                'centroid': {
                    'x': round(float(centroids[label][0]), 2),
                    'y': round(float(centroids[label][1]), 2),
                },
            }
        )

    regions.sort(key=lambda item: item['area_px'], reverse=True)
    for rank, region in enumerate(regions, start=1):
        region['rank'] = rank
    return cleaned_mask, regions


def build_report(input_path: Path, cleaned_mask: np.ndarray, regions: list[dict[str, Any]]) -> dict[str, Any]:
    image_height, image_width = cleaned_mask.shape[:2]
    total_pixels = int(image_height * image_width)
    crack_pixels = int(cleaned_mask.sum())
    crack_ratio = crack_pixels / total_pixels if total_pixels else 0.0

    if crack_pixels == 0:
        severity = 'none'
        suggestion = '未检出裂缝像素。建议结合原图质量进行人工复核。'
    elif crack_ratio < 0.0005:
        severity = 'low'
        suggestion = '检出少量裂缝像素。建议人工复核并纳入日常巡检。'
    elif crack_ratio < 0.005:
        severity = 'medium'
        suggestion = '检出较明显裂缝区域。建议结合现场尺度标定，安排近距离复查。'
    else:
        severity = 'high'
        suggestion = '裂缝像素占比较高。建议优先现场复核，并评估维修或加固方案。'

    return {
        'input': str(input_path),
        'image_size': {'width': image_width, 'height': image_height},
        'has_crack': bool(crack_pixels > 0),
        'crack_pixels': crack_pixels,
        'crack_ratio': round(crack_ratio, 6),
        'crack_ratio_percent': round(crack_ratio * 100, 4),
        'region_count': len(regions),
        'largest_region': regions[0] if regions else None,
        'severity': severity,
        'suggestion': suggestion,
        'regions': regions,
    }


def save_outputs(image: Image.Image, cleaned_mask: np.ndarray, output_dir: Path, input_path: Path) -> dict[str, str]:
    output_dir.mkdir(parents=True, exist_ok=True)
    file_stem = input_path.stem
    rgb_image = np.array(image).astype(np.uint8)
    mask_image = (cleaned_mask * 255).astype(np.uint8)
    overlay_image = rgb_image.copy()
    red_layer = np.zeros_like(rgb_image)
    red_layer[:, :, 0] = 255
    overlay_image = np.where(
        cleaned_mask[:, :, None] > 0,
        (0.55 * rgb_image + 0.45 * red_layer).astype(np.uint8),
        overlay_image,
    )
    summary_image = np.concatenate([rgb_image, cv2.cvtColor(mask_image, cv2.COLOR_GRAY2RGB), overlay_image], axis=1)

    mask_path = output_dir / f'{file_stem}-mask.png'
    overlay_path = output_dir / f'{file_stem}-overlay.png'
    summary_path = output_dir / f'{file_stem}-summary.png'
    report_path = output_dir / f'{file_stem}-report.json'

    cv2.imwrite(str(mask_path), mask_image)
    cv2.imwrite(str(overlay_path), cv2.cvtColor(overlay_image, cv2.COLOR_RGB2BGR))
    cv2.imwrite(str(summary_path), cv2.cvtColor(summary_image, cv2.COLOR_RGB2BGR))
    return {
        'mask': str(mask_path),
        'overlay': str(overlay_path),
        'summary': str(summary_path),
        'report_json': str(report_path),
    }


def print_report(report: dict[str, Any], outputs: dict[str, str]) -> None:
    image_size = report['image_size']
    has_crack = '是' if report['has_crack'] else '否'
    input_path = report['input']
    image_width = image_size['width']
    image_height = image_size['height']
    crack_pixels = report['crack_pixels']
    crack_ratio_percent = report['crack_ratio_percent']
    region_count = report['region_count']
    severity = report['severity']
    suggestion = report['suggestion']

    print('\n==== 本地裂缝检测报告 ====')
    print(f'输入图片: {input_path}')
    print(f'图片尺寸: {image_width} x {image_height} px')
    print(f'是否检出裂缝: {has_crack}')
    print(f'裂缝像素数: {crack_pixels}')
    print(f'裂缝面积占比: {crack_ratio_percent}%')
    print(f'裂缝连通区域数: {region_count}')
    print(f'风险等级: {severity}')
    print(f'建议: {suggestion}')

    if report['largest_region']:
        largest_region = report['largest_region']
        bbox = largest_region['bbox']
        centroid = largest_region['centroid']
        area_px = largest_region['area_px']
        bbox_x = bbox['x']
        bbox_y = bbox['y']
        bbox_width = bbox['width']
        bbox_height = bbox['height']
        bbox_diagonal_px = largest_region['bbox_diagonal_px']
        centroid_x = centroid['x']
        centroid_y = centroid['y']
        print('\n最大裂缝区域:')
        print(f'  面积: {area_px} px')
        print(f'  边界框: x={bbox_x}, y={bbox_y}, w={bbox_width}, h={bbox_height}')
        print(f'  边界框对角线: {bbox_diagonal_px} px')
        print(f'  中心点: ({centroid_x}, {centroid_y})')

    print('\n输出文件:')
    for name, path in outputs.items():
        print(f'  {name}: {path}')

def main() -> int:
    args = parse_args()
    start_time = time.time()
    input_path = Path(args.input).expanduser().resolve()
    device = get_device()
    image = read_image(input_path)
    model = load_model(MODEL_PATH, device)
    raw_mask = predict_mask(model, image, device)
    cleaned_mask, regions = extract_regions(raw_mask)
    report = build_report(input_path, cleaned_mask, regions)
    report['model'] = str(MODEL_PATH)
    report['device'] = str(device)
    report['input_size'] = INPUT_SIZE
    report['runtime_seconds'] = round(time.time() - start_time, 3)

    timestamp = time.strftime('%Y%m%d-%H%M%S')
    run_dir = OUTPUT_ROOT / f'{input_path.stem}-{timestamp}'
    outputs = save_outputs(image, cleaned_mask, run_dir, input_path)
    report['outputs'] = outputs
    Path(outputs['report_json']).write_text(json.dumps(report, ensure_ascii=False, indent=2), encoding='utf-8')
    print_report(report, outputs)
    return 0


if __name__ == '__main__':
    raise SystemExit(main())
