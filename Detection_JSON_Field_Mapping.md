# 原子检测 JSON 报告字段落库建议

## 一、金属锈蚀检测 CorrosionDetection

| JSON 英文名称 | 中文名称 | 意义 | 备注 |
| --- | --- | --- | --- |
| `has_corrosion` | 是否存在锈蚀 | 布尔值，表示图像是否检测到锈蚀目标 | 推荐落库；核心筛选字段 |
| `corrosion_count` | 锈蚀目标数量 | 检测出的锈蚀实例数量 | 推荐落库；核心统计字段 |
| `max_confidence` | 最高置信度 | 所有锈蚀目标中的最高模型置信度 | 推荐落库；用于结果可信度排序 |
| `average_confidence` | 平均置信度 | 所有锈蚀目标的平均模型置信度 | 推荐落库；用于整体可信度评估 |
| `corrosion_pixels` | 锈蚀像素数 | 锈蚀掩码覆盖的总像素数量 | 推荐落库；核心量化指标 |
| `corrosion_ratio` | 锈蚀面积比例 | 锈蚀像素数占整图像素数的比例，取值范围 0 到 1 | 推荐落库；核心量化指标 |
| `regions[].id` | 锈蚀目标序号 | 过滤 `class_id = 0` 且 `class_name = rust` 后重新从 1 顺延编号 | 可选落库；如果需要目标级明细表则落库 |
| `regions[].confidence` | 目标置信度 | 单个锈蚀目标的模型置信度 | 可选落库；目标级明细字段 |
| `regions[].bbox_xyxy` | 边界框坐标 | 左上角和右下角坐标，格式为 `[x1,y1,x2,y2]` | 可选落库；目标级定位字段 |
| `regions[].mask_pixels` | 目标掩码像素数 | 单个目标掩码覆盖像素数量 | 可选落库；目标级明细字段 |
| `regions[].mask_ratio` | 目标掩码比例 | 单个目标掩码面积占整图面积比例，取值范围 0 到 1 | 可选落库；目标级明细字段 |
| `runtime` | 执行时间 | 生成报告时的 Linux Unix 秒级时间戳 | 推荐落库；用于任务时间追踪 |

```
{
  "has_corrosion": true,
  "corrosion_count": 8,
  "max_confidence": 0.961351,
  "average_confidence": 0.640132,
  "corrosion_pixels": 51465,
  "corrosion_ratio": 0.125647,
  "regions": [
    {
      "id": 1,
      "confidence": 0.961351,
      "bbox_xyxy": [179.3,160.69,530.77,291.78],
      "mask_pixels": 23086,
      "mask_ratio": 0.056362
    }
  ],
  "runtime": 1122334
}
```

## 二、石材裂缝检测 CrackDetection

| JSON 英文名称 | 中文名称 | 意义 | 备注 |
| --- | --- | --- | --- |
| `has_crack` | 是否存在裂缝 | 布尔值，表示是否检测到裂缝像素 | 推荐落库；核心筛选字段 |
| `crack_count` | 裂缝区域数量 | 按面积排序后的裂缝区域数量 | 推荐落库；核心统计字段 |
| `crack_pixels` | 裂缝像素数 | 裂缝掩码覆盖的像素数量 | 推荐落库；核心量化指标 |
| `crack_ratio` | 裂缝面积比例 | 裂缝像素数占整图像素数比例，取值范围 0 到 1 | 推荐落库；核心量化指标 |
| `regions[].id` | 裂缝区域序号 | 裂缝区域按 `mask_pixels` 从大到小排序后，从 1 开始重新编号 | 可选落库；如需目标级明细表则落库 |
| `regions[].bbox_xyxy` | 裂缝区域边界框 | 左上角和右下角坐标，格式为 `[x1,y1,x2,y2]` | 可选落库；定位字段 |
| `regions[].mask_pixels` | 裂缝区域像素数 | 单个裂缝区域的掩码像素数量 | 可选落库；目标级明细字段 |
| `regions[].mask_ratio` | 裂缝区域面积比例 | 单个裂缝区域像素数占整图像素数比例，取值范围 0 到 1 | 可选落库；目标级明细字段 |
| `runtime` | 执行时间 | 生成报告时的 Linux Unix 秒级时间戳 | 推荐落库；用于任务时间追踪 |

```
{
  "has_crack": true,
  "crack_count": 1,
  "crack_pixels": 8450,
  "crack_ratio": 0.042102,
  "regions": [
    {
      "id": 1,
      "bbox_xyxy": [179,160,530,291],
      "mask_pixels": 8450,
      "mask_ratio": 0.056362
    }
  ],
  "runtime": 12312312
}
```

## 三、石材污渍检测 StainDetection

| JSON 英文名称 | 中文名称 | 意义 | 备注 |
| --- | --- | --- | --- |
| `has_stain` | 是否存在污渍 | 布尔值，表示是否存在有效污渍区域 | 推荐落库；核心筛选字段 |
| `stain_count` | 污渍区域数量 | 有污渍的区域数量，不统计无污渍块 | 推荐落库；核心统计字段 |
| `average_stain_ratio` | 平均污渍比例 | 所有污渍区域污渍占比的平均值，按百分比数值保存 | 推荐落库；核心量化指标 |
| `max_stain_ratio` | 最大污渍比例 | 所有污渍区域污渍占比的最大值，按百分比数值保存 | 推荐落库；核心量化指标 |
| `regions[].id` | 污渍区域序号 | 只对有污渍区域从 1 开始顺序编号，不跳号 | 可选落库；如需目标级明细表则落库 |
| `regions[].confidence` | 检测置信度 | 图像块检测模型对该区域的置信度 | 可选落库；目标级明细字段 |
| `regions[].bbox_xyxy` | 边界框坐标 | 图像块检测框坐标，格式为 `[x1,y1,x2,y2]` | 可选落库；定位字段 |
| `regions[].region_width` | 区域宽度 | 透视矫正后区域宽度 | 可选落库；目标级明细字段 |
| `regions[].region_height` | 区域高度 | 透视矫正后区域高度 | 可选落库；目标级明细字段 |
| `regions[].stain_pixels` | 污渍像素数 | 污渍在矫正后区域中的像素数量 | 推荐落库；如建明细表则为核心字段 |
| `regions[].stain_ratio` | 污渍面积比例 | 污渍像素数占矫正后区域像素数比例，取值范围 0 到 1 | 推荐落库；如建明细表则为核心字段 |
| `runtime` | 执行时间 | 生成报告时的 Linux Unix 秒级时间戳 | 推荐落库；用于任务时间追踪 |

```
{
  "has_stain": true,
  "stain_count": 5,
  "average_stain_ratio": 12.0048,
  "max_stain_ratio": 38.3384,
  "regions": [
    {
      "id": 1,
      "confidence": 0.953774,
      "bbox_xyxy": [21.41,1000.95,2362.81,1784.65],
      "region_width": 2327,
      "region_height": 783,
      "stain_pixels": 3833,
      "stain_ratio": 0.383384
    }
  ],
  "runtime": 12312332
}
```

## 四、玻璃平整度检测 FlatnessDetection

| JSON 英文名称 | 中文名称 | 意义 | 备注 |
| --- | --- | --- | --- |
| `result` | 平整度结论 | 字符串枚举，仅包含 `平整`、`不平整`、`非玻璃` | 推荐落库；核心筛选字段 |
| `is_flat` | 是否整体平整 | 仅当结论为 `平整` 时为 `true` | 推荐落库；核心筛选字段 |
| `uneven_count` | 不平整玻璃块数量 | 不平整玻璃区域数量 | 推荐落库；核心统计字段 |
| `regions[].id` | 不平整区域序号 | 只对不平整玻璃区域从 1 开始顺序编号，不跳号 | 可选落库；如需目标级明细表则落库 |
| `regions[].bbox_xyxy` | 区域边界框坐标 | 左上角和右下角坐标，格式为 `[x1,y1,x2,y2]` | 可选落库；区域定位字段 |
| `regions[].edge_uneven_detected` | 边缘判据是否发现不平整 | `true` 表示边缘判据认为不平整，`false` 表示平整 | 可选落库；算法解释字段 |
| `regions[].line_uneven_detected` | 直线判据是否发现不平整 | `true` 表示直线判据认为不平整，`false` 表示平整 | 可选落库；算法解释字段 |
| `regions[].gradient_uneven_detected` | 梯度判据是否发现不平整 | `true` 表示梯度判据认为不平整，`false` 表示平整 | 可选落库；算法解释字段 |
| `regions[].frequency_uneven_detected` | 频谱判据是否发现不平整 | `true` 表示频谱判据认为不平整，`false` 表示平整 | 可选落库；算法解释字段 |
| `regions[].edge_count` | 边缘数量 | Canny 边缘点数量 | 可选落库；算法解释字段 |
| `regions[].laplacian_variance` | 拉普拉斯方差 | 衡量边缘清晰度和模糊程度 | 可选落库；算法解释字段 |
| `regions[].line_count` | 直线数量 | Hough 检测出的直线数量 | 可选落库；算法解释字段 |
| `regions[].angle_std` | 直线角度标准差 | 直线角度离散程度 | 可选落库；算法解释字段 |
| `regions[].gradient_mean` | 梯度均值 | 玻璃区域梯度强度平均值 | 可选落库；算法解释字段 |
| `regions[].gradient_std` | 梯度标准差 | 玻璃区域梯度强度离散程度 | 可选落库；算法解释字段 |
| `regions[].frequency_min` | 频谱最小值 | 频谱幅值最小值 | 可选落库；算法解释字段 |
| `regions[].frequency_max` | 频谱最大值 | 频谱幅值最大值 | 可选落库；算法解释字段 |
| `runtime` | 执行时间 | 生成报告时的 Linux Unix 秒级时间戳 | 推荐落库；用于任务时间追踪 |

```
{
  "result": "不平整",
  "is_flat": false,
  "uneven_count": 1,
  "regions": [
    {
      "id": 1,
      "bbox_xyxy": [0,0,0,0],
      "edge_uneven_detected": false,
      "line_uneven_detected": true,
      "gradient_uneven_detected": true,
      "frequency_uneven_detected": false,
      "edge_count": 60438,
      "laplacian_variance": 4254.554801,
      "line_count": 587,
      "angle_std": 73.2174,
      "gradient_mean": 137.417486,
      "gradient_std": 155.640067,
      "frequency_min": 35.819298,
      "frequency_max": 343.158295
    }
  ],
  "runtime": 1141
}
```

## 五、玻璃爆裂检测 SpallingDetection

| JSON 英文名称 | 中文名称 | 意义 | 备注 |
| --- | --- | --- | --- |
| `has_spalling` | 是否存在爆裂 | 布尔值，表示是否判定为玻璃爆裂 | 推荐落库；核心筛选字段 |
| `confidence` | 检测置信度 | 分类预测置信度 | 推荐落库；核心可信度字段 |
| `runtime` | 执行时间 | 生成报告时的 Linux Unix 秒级时间戳 | 推荐落库；用于任务时间追踪 |

```
{
  "has_spalling": true,
  "confidence": 0.989133,
  "runtime": 256
}
```
