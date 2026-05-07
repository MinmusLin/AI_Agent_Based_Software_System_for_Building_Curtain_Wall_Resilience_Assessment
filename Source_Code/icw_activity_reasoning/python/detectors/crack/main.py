import json
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
MODEL_PATH = APP_ROOT / 'model' / 'best_weights_model.pt'
INPUT_SIZE = 1024
MIN_REGION_AREA = 20
NUM_CLASSES = 2
DETECTOR: 'CrackDetector | None' = None


class Gelu(nn.Module):

    # 执行当前模型模块的前向传播
    def forward(self, tensor: torch.Tensor) -> torch.Tensor:
        return 0.5 * tensor * (1 + torch.tanh(np.sqrt(2 / np.pi) * (tensor + 0.044715 * torch.pow(tensor, 3))))


class OverlapPatchEmbed(nn.Module):

    # 初始化当前模型模块的参数
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

    # 执行当前模型模块的前向传播
    def forward(self, tensor: torch.Tensor) -> tuple[torch.Tensor, int, int]:
        tensor = self.proj(tensor)
        _, _, height, width = tensor.shape
        tensor = tensor.flatten(2).transpose(1, 2)
        tensor = self.norm(tensor)
        return tensor, height, width


class Attention(nn.Module):

    # 初始化当前模型模块的参数
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

    # 执行当前模型模块的前向传播
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

    # 初始化当前模型模块的参数
    def __init__(self, dim: int) -> None:
        super().__init__()
        self.dwconv = nn.Conv2d(dim, dim, 3, 1, 1, bias=True, groups=dim)

    # 执行当前模型模块的前向传播
    def forward(self, tensor: torch.Tensor, height: int, width: int) -> torch.Tensor:
        batch_size, _, channel_count = tensor.shape
        tensor = tensor.transpose(1, 2).view(batch_size, channel_count, height, width)
        tensor = self.dwconv(tensor)
        return tensor.flatten(2).transpose(1, 2)


class FeedForward(nn.Module):

    # 初始化当前模型模块的参数
    def __init__(self, in_features: int, hidden_features: int, drop: float) -> None:
        super().__init__()
        self.fc1 = nn.Linear(in_features, hidden_features)
        self.dwconv = DepthwiseConv(hidden_features)
        self.act = Gelu()
        self.fc2 = nn.Linear(hidden_features, in_features)
        self.drop = nn.Dropout(drop)

    # 执行当前模型模块的前向传播
    def forward(self, tensor: torch.Tensor, height: int, width: int) -> torch.Tensor:
        tensor = self.fc1(tensor)
        tensor = self.dwconv(tensor, height, width)
        tensor = self.act(tensor)
        tensor = self.drop(tensor)
        tensor = self.fc2(tensor)
        return self.drop(tensor)


class TransformerBlock(nn.Module):

    # 初始化当前模型模块的参数
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

    # 执行当前模型模块的前向传播
    def forward(self, tensor: torch.Tensor, height: int, width: int) -> torch.Tensor:
        tensor = tensor + self.drop_path(self.attn(self.norm1(tensor), height, width))
        tensor = tensor + self.drop_path(self.mlp(self.norm2(tensor), height, width))
        return tensor


class MixVisionTransformer(nn.Module):

    # 初始化当前模型模块的参数
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

    # 构建 MixVisionTransformer 阶段模块列表
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

    # 执行 MixVisionTransformer 的单个阶段
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

    # 执行当前模型模块的前向传播
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

    # 初始化当前模型模块的参数
    def __init__(self, input_dim: int, embed_dim: int) -> None:
        super().__init__()
        self.proj = nn.Linear(input_dim, embed_dim)

    # 执行当前模型模块的前向传播
    def forward(self, tensor: torch.Tensor) -> torch.Tensor:
        tensor = tensor.flatten(2).transpose(1, 2)
        return self.proj(tensor)


class ConvModule(nn.Module):

    # 初始化当前模型模块的参数
    def __init__(self, in_channels: int, out_channels: int) -> None:
        super().__init__()
        self.conv = nn.Conv2d(in_channels, out_channels, 1, 1, 0, groups=1, bias=False)
        self.bn = nn.BatchNorm2d(out_channels, eps=0.001, momentum=0.03)
        self.act = nn.ReLU()

    # 执行当前模型模块的前向传播
    def forward(self, tensor: torch.Tensor) -> torch.Tensor:
        return self.act(self.bn(self.conv(tensor)))


class SegFormerHead(nn.Module):

    # 初始化当前模型模块的参数
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

    # 执行当前模型模块的前向传播
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

    # 初始化当前模型模块的参数
    def __init__(self) -> None:
        super().__init__()
        self.backbone = MixVisionTransformer()
        self.decode_head = SegFormerHead()

    # 执行当前模型模块的前向传播
    def forward(self, tensor: torch.Tensor) -> torch.Tensor:
        height, width = tensor.size(2), tensor.size(3)
        tensor = self.backbone(tensor)
        tensor = self.decode_head(tensor)
        return torch_functional.interpolate(tensor, size=(height, width), mode='bilinear', align_corners=True)


# 获取当前可用的推理设备
def get_device() -> torch.device:
    return torch.device('cuda' if torch.cuda.is_available() else 'cpu')


# 加载模型权重并切换为推理模式
def load_model(model_path: Path, device: torch.device) -> SegFormer:
    model = SegFormer()
    state_dict = torch.load(model_path, map_location=device)
    if any(key.startswith('module.') for key in state_dict):
        state_dict = {key.replace('module.', '', 1): value for key, value in state_dict.items()}
    model.load_state_dict(state_dict)
    model.to(device)
    model.eval()
    return model


# 将输入图像转换为 RGB 格式
def convert_to_rgb(image: Image.Image) -> Image.Image:
    if len(np.shape(image)) == 3 and np.shape(image)[2] == 3:
        return image
    return image.convert('RGB')


# 按比例缩放图像并填充到目标尺寸
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


# 对输入图像执行模型预处理
def preprocess_image(image: np.ndarray) -> np.ndarray:
    image -= np.array([123.675, 116.28, 103.53], np.float32)
    image /= np.array([58.395, 57.12, 57.375], np.float32)
    return image


# 读取并校验输入图像
def read_image(input_path: Path) -> Image.Image:
    if not input_path.exists():
        raise FileNotFoundError(f'input image not found: {input_path}')
    return convert_to_rgb(Image.open(input_path))


# 执行分割模型推理并生成掩码
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


# 从裂缝掩码中提取连通区域
def extract_regions(mask: np.ndarray) -> tuple[np.ndarray, list[dict[str, Any]]]:
    binary_mask = (mask == 1).astype(np.uint8)
    image_height, image_width = binary_mask.shape[:2]
    image_area = int(image_height * image_width)
    region_count, labels, stats, _ = cv2.connectedComponentsWithStats(binary_mask, connectivity=8)
    cleaned_mask = np.zeros_like(binary_mask)
    regions = []

    for label in range(1, region_count):
        x, y, width, height, area = stats[label]
        mask_pixels = int(area)
        if mask_pixels < MIN_REGION_AREA:
            continue
        cleaned_mask[labels == label] = 1
        regions.append(
            {
                'bbox_xyxy': [int(x), int(y), int(x + width), int(y + height)],
                'mask_pixels': mask_pixels,
                'mask_ratio': round(mask_pixels / image_area, 6) if image_area else 0.0,
            }
        )

    regions.sort(key=lambda item: item['mask_pixels'], reverse=True)
    return cleaned_mask, [
        {
            'id': index,
            'bbox_xyxy': region['bbox_xyxy'],
            'mask_pixels': region['mask_pixels'],
            'mask_ratio': region['mask_ratio'],
        }
        for index, region in enumerate(regions, start=1)
    ]


# 构建检测结果报告数据
def build_report(cleaned_mask: np.ndarray, regions: list[dict[str, Any]]) -> dict[str, Any]:
    image_height, image_width = cleaned_mask.shape[:2]
    total_pixels = int(image_height * image_width)
    crack_pixels = int(cleaned_mask.sum())
    crack_ratio = crack_pixels / total_pixels if total_pixels else 0.0

    return {
        'has_crack': bool(crack_pixels > 0),
        'crack_count': len(regions),
        'crack_pixels': crack_pixels,
        'crack_ratio': round(crack_ratio, 6),
        'regions': regions,
    }


# 保存检测输出图像文件
def save_outputs(image: Image.Image, cleaned_mask: np.ndarray, output_dir: Path) -> None:
    output_dir.mkdir(parents=True, exist_ok=True)
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
    mask_path = output_dir / 'mask.png'
    overlay_path = output_dir / 'overlay.png'

    cv2.imwrite(str(mask_path), mask_image)
    cv2.imwrite(str(overlay_path), cv2.cvtColor(overlay_image, cv2.COLOR_RGB2BGR))


class CrackDetector:
    def __init__(self) -> None:
        self.device = get_device()
        self.model = load_model(MODEL_PATH, self.device)

    # 执行石材裂缝检测
    def detect(self, input_path: Path) -> None:
        start_time = time.time()
        input_path = input_path.expanduser().resolve()
        image = read_image(input_path)
        raw_mask = predict_mask(self.model, image, self.device)
        cleaned_mask, regions = extract_regions(raw_mask)
        report = build_report(cleaned_mask, regions)
        report['runtime_seconds'] = round(time.time() - start_time, 3)
        run_dir = input_path.parent
        save_outputs(image, cleaned_mask, run_dir)
        report_path = run_dir / 'report.json'
        report_path.write_text(json.dumps(report, ensure_ascii=False, indent=2), encoding='utf-8')


# 获取可复用的石材裂缝检测器
def get_detector() -> CrackDetector:
    global DETECTOR
    if DETECTOR is None:
        DETECTOR = CrackDetector()
    return DETECTOR


# 执行石材裂缝检测
def detect(input_path: Path) -> None:
    get_detector().detect(input_path)
