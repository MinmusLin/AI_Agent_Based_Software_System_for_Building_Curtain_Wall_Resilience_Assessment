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
  "corrosion_pixels": 51465,
  "corrosion_ratio": 0.125647,
  "max_confidence": 0.961351,
  "average_confidence": 0.640132,
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
| `crack_pixels` | 裂缝像素数 | 裂缝掩码覆盖的像素数量 | 推荐落库；核心量化指标 |
| `crack_ratio` | 裂缝面积比例 | 裂缝像素数占整图像素数比例，取值范围 0 到 1 | 推荐落库；核心量化指标 |
| `region_count` | 裂缝区域数量 | 连通裂缝区域数量 | 推荐落库；核心统计字段 |
| `largest_region.id` | 最大裂缝区域序号 | 面积最大的裂缝区域序号 | 可选落库；可直接落最大区域摘要字段 |
| `largest_region.area_px` | 最大裂缝区域面积 | 最大裂缝区域面积，单位像素 | 推荐落库；严重程度辅助指标 |
| `largest_region.bbox.x` | 最大区域边界框 X 坐标 | 最大裂缝区域边界框左上角 X 坐标 | 可选落库；定位字段 |
| `largest_region.bbox.y` | 最大区域边界框 Y 坐标 | 最大裂缝区域边界框左上角 Y 坐标 | 可选落库；定位字段 |
| `largest_region.bbox.width` | 最大区域边界框宽度 | 最大裂缝区域边界框宽度 | 可选落库；定位字段 |
| `largest_region.bbox.height` | 最大区域边界框高度 | 最大裂缝区域边界框高度 | 可选落库；定位字段 |
| `largest_region.bbox_diagonal_px` | 最大区域边界框对角线 | 最大裂缝区域边界框对角线长度 | 可选落库；尺度字段 |
| `largest_region.centroid.x` | 最大区域中心 X 坐标 | 最大裂缝区域质心 X 坐标 | 可选落库；定位字段 |
| `largest_region.centroid.y` | 最大区域中心 Y 坐标 | 最大裂缝区域质心 Y 坐标 | 可选落库；定位字段 |
| `largest_region.rank` | 最大区域排序 | 按面积排序后的名次 | 不建议单独落库；最大区域固定为第 1 名 |
| `severity` | 裂缝严重程度 | 根据裂缝像素比例得到的 `none`、`low`、`medium`、`high` | 推荐落库；核心等级字段 |
| `suggestion` | 处理建议 | 根据严重程度生成的检查建议 | 推荐落库；可直接用于报告展示 |
| `regions[].id` | 裂缝区域序号 | 单个裂缝连通区域序号 | 可选落库；如需目标级明细表则落库 |
| `regions[].area_px` | 裂缝区域面积 | 单个裂缝区域面积，单位像素 | 可选落库；目标级明细字段 |
| `regions[].bbox.x` | 区域边界框 X 坐标 | 单个裂缝区域边界框左上角 X 坐标 | 可选落库；定位字段 |
| `regions[].bbox.y` | 区域边界框 Y 坐标 | 单个裂缝区域边界框左上角 Y 坐标 | 可选落库；定位字段 |
| `regions[].bbox.width` | 区域边界框宽度 | 单个裂缝区域边界框宽度 | 可选落库；定位字段 |
| `regions[].bbox.height` | 区域边界框高度 | 单个裂缝区域边界框高度 | 可选落库；定位字段 |
| `regions[].bbox_diagonal_px` | 区域边界框对角线 | 单个裂缝区域边界框对角线长度 | 可选落库；尺度字段 |
| `regions[].centroid.x` | 区域中心 X 坐标 | 单个裂缝区域质心 X 坐标 | 可选落库；定位字段 |
| `regions[].centroid.y` | 区域中心 Y 坐标 | 单个裂缝区域质心 Y 坐标 | 可选落库；定位字段 |
| `regions[].rank` | 区域面积排序 | 单个裂缝区域按面积从大到小的排序 | 可选落库；目标级明细字段 |
| `runtime_seconds` | 执行耗时 | 原子检测运行耗时，单位秒，数值格式保存 | 推荐落库；用于性能统计 |

## 三、石材污渍检测 StainDetection

| JSON 英文名称 | 中文名称 | 意义 | 备注 |
| --- | --- | --- | --- |
| `status` | 检测状态 | `success` 表示成功，`failed` 表示未形成有效污渍检测结果 | 推荐落库；核心状态字段 |
| `total_blocks` | 有污渍块数量 | 当前报告中保留下来的污渍块数量 | 推荐落库；核心统计字段 |
| `average_stain_percentage` | 平均污渍百分比 | 所有污渍块的平均污渍占比 | 推荐落库；核心量化指标 |
| `max_stain_percentage` | 最大污渍百分比 | 所有污渍块中的最大污渍占比 | 推荐落库；核心量化指标 |
| `detections[].block_number` | 块编号 | 污渍块在图像中的编号，和输出文件 `block_{no}.png` 对应 | 可选落库；目标级明细字段 |
| `detections[].confidence` | 检测置信度 | 图像块检测模型对该块的置信度 | 可选落库；目标级明细字段 |
| `detections[].class_id` | 类别 ID | 图像块检测模型类别 ID | 可选落库；目标级明细字段 |
| `detections[].class_name` | 类别名称 | 图像块检测模型类别名称 | 可选落库；目标级明细字段 |
| `detections[].bbox_xyxy` | 边界框坐标 | 图像块检测框坐标，格式为 `[x1,y1,x2,y2]` | 可选落库；定位字段 |
| `detections[].corners[].x` | 透视角点 X 坐标 | 图像块四角点的 X 坐标 | 不建议单独落库；适合原始 JSON |
| `detections[].corners[].y` | 透视角点 Y 坐标 | 图像块四角点的 Y 坐标 | 不建议单独落库；适合原始 JSON |
| `detections[].block_size.width` | 块宽度 | 透视矫正后图像块宽度 | 可选落库；目标级明细字段 |
| `detections[].block_size.height` | 块高度 | 透视矫正后图像块高度 | 可选落库；目标级明细字段 |
| `detections[].block_area_px` | 块面积 | 透视矫正后图像块面积，单位像素 | 可选落库；目标级明细字段 |
| `detections[].has_stain` | 是否有污渍 | 当前块是否判定存在污渍 | 推荐落库；如建明细表则为核心字段 |
| `detections[].stain_ratio` | 污渍比例 | 当前块污渍面积占比，取值范围 0 到 1 | 推荐落库；如建明细表则为核心字段 |
| `detections[].stain_percentage` | 污渍百分比 | 当前块污渍比例的百分比表达 | 可选落库；如只保留比例可不单独落库 |
| `detections[].variance` | 块亮度方差 | 分块统计得到的方差，用于跳过低变化区域 | 可选落库；算法解释字段 |
| `detections[].threshold` | 二值化阈值 | 当前块污渍分割使用的阈值 | 可选落库；算法复现字段 |
| `detections[].status` | 块处理状态 | 当前块处理状态，如 `processed` | 可选落库；明细表状态字段 |
| `runtime_seconds` | 执行耗时 | 原子检测运行耗时，单位秒，数值格式保存 | 推荐落库；用于性能统计 |

## 四、玻璃平整度检测 FlatnessDetection

| JSON 英文名称 | 中文名称 | 意义 | 备注 |
| --- | --- | --- | --- |
| `flatness_result` | 平整度结论 | 字符串枚举，仅包含 `平整`、`不平整`、`非玻璃` | 推荐落库；核心筛选字段 |
| `glass_count` | 玻璃区域数量 | 检测出的玻璃区域总数 | 推荐落库；核心统计字段 |
| `flat_glass_count` | 平整玻璃区域数量 | 被判定为平整的玻璃区域数量 | 推荐落库；核心统计字段 |
| `uneven_glass_count` | 不平整玻璃区域数量 | 被判定为不平整的玻璃区域数量 | 推荐落库；核心统计字段 |
| `glass_regions[].id` | 玻璃区域序号 | 单个玻璃区域编号 | 可选落库；如建区域明细表则落库 |
| `glass_regions[].bbox.x` | 区域边界框 X 坐标 | 单个玻璃区域边界框左上角 X 坐标 | 可选落库；区域定位字段 |
| `glass_regions[].bbox.y` | 区域边界框 Y 坐标 | 单个玻璃区域边界框左上角 Y 坐标 | 可选落库；区域定位字段 |
| `glass_regions[].bbox.width` | 区域边界框宽度 | 单个玻璃区域边界框宽度 | 可选落库；区域定位字段 |
| `glass_regions[].bbox.height` | 区域边界框高度 | 单个玻璃区域边界框高度 | 可选落库；区域定位字段 |
| `glass_regions[].area_px` | 区域面积 | 单个玻璃区域面积，单位像素 | 推荐落库；如建区域明细表则为核心字段 |
| `glass_regions[].flatness_result` | 区域平整度结果 | `1` 表示平整，`0` 表示不平整 | 推荐落库；如建区域明细表则为核心字段 |
| `glass_regions[].is_flat` | 区域是否平整 | 单个玻璃区域是否平整 | 推荐落库；如建区域明细表则为核心字段 |
| `glass_regions[].edge_result` | 边缘分析结果 | 边缘清晰度分析得到的平整性判断 | 可选落库；算法解释字段 |
| `glass_regions[].line_result` | 直线分析结果 | Hough 直线角度分布得到的平整性判断 | 可选落库；算法解释字段 |
| `glass_regions[].gradient_result` | 梯度分析结果 | 梯度变化得到的平整性判断 | 可选落库；算法解释字段 |
| `glass_regions[].frequency_result` | 频谱分析结果 | 频域变化得到的平整性判断 | 可选落库；算法解释字段 |
| `glass_regions[].edge_count` | 边缘数量 | Canny 边缘点数量 | 可选落库；算法解释字段 |
| `glass_regions[].laplacian_variance` | 拉普拉斯方差 | 衡量边缘清晰度和模糊程度 | 可选落库；算法解释字段 |
| `glass_regions[].line_count` | 直线数量 | Hough 检测出的直线数量 | 可选落库；算法解释字段 |
| `glass_regions[].angle_std` | 直线角度标准差 | 直线角度离散程度 | 可选落库；算法解释字段 |
| `glass_regions[].gradient_mean` | 梯度均值 | 玻璃区域梯度强度平均值 | 可选落库；算法解释字段 |
| `glass_regions[].gradient_std` | 梯度标准差 | 玻璃区域梯度强度离散程度 | 可选落库；算法解释字段 |
| `glass_regions[].frequency_min` | 频谱最小值 | 频谱幅值最小值 | 可选落库；算法解释字段 |
| `glass_regions[].frequency_max` | 频谱最大值 | 频谱幅值最大值 | 可选落库；算法解释字段 |
| `glass_regions[].frequency_diff` | 频谱范围差 | 频谱最大值与最小值差值 | 可选落库；算法解释字段 |
| `runtime_seconds` | 执行耗时 | 原子检测运行耗时，单位秒，数值格式保存 | 推荐落库；用于性能统计 |

## 五、玻璃爆裂检测 SpallingDetection

| JSON 英文名称 | 中文名称 | 意义 | 备注 |
| --- | --- | --- | --- |
| `has_spalling` | 是否存在爆裂 | 布尔值，表示是否判定为玻璃爆裂 | 推荐落库；核心筛选字段 |
| `confidence` | 检测置信度 | 分类预测置信度 | 推荐落库；核心可信度字段 |
| `predicted_index` | 预测类别索引 | 模型输出的类别索引 | 可选落库；模型明细字段 |
| `logits` | 原始输出值 | 模型分类层原始 logits | 不建议单独落库；算法调试字段 |
| `probabilities` | 类别概率列表 | 各类别 softmax 概率列表 | 可选落库；如果需要完整概率解释可落 JSON 字段 |
| `probabilities_by_class` | 按类别映射的概率 | 类别名称到概率的映射 | 可选落库；比 `probabilities` 更适合展示 |
| `runtime_seconds` | 执行耗时 | 原子检测运行耗时，单位秒，数值格式保存 | 推荐落库；用于性能统计 |
