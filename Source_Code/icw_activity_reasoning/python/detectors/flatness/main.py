import json
import time
from functools import reduce
from pathlib import Path
from typing import Any

import cv2
import numpy as np
import torch
import torch.nn as nn
import torch.nn.functional as F
from PIL import Image, ImageDraw
from torchvision import transforms


APP_ROOT = Path(__file__).resolve().parent
MODEL_PATH = APP_ROOT / 'model' / 'best_weights_model.pt'
INPUT_SIZE = 416
MASK_THRESHOLD = 200
MIN_SEGMENT_AREA = 200
DETECTOR: 'FlatnessDetector | None' = None


class LambdaBase(nn.Sequential):

    # 初始化当前模型模块的参数
    def __init__(self, lambda_func, *args):
        super(LambdaBase, self).__init__(*args)
        self.lambda_func = lambda_func

    # 准备 Lambda 模块的前向输入
    def forward_prepare(self, input_tensor):
        output = []
        for module in self._modules.values():
            output.append(module(input_tensor))
        return output if output else input_tensor


class Lambda(LambdaBase):

    # 执行当前模型模块的前向传播
    def forward(self, input_tensor):
        return self.lambda_func(self.forward_prepare(input_tensor))


class LambdaMap(LambdaBase):

    # 执行当前模型模块的前向传播
    def forward(self, input_tensor):
        return list(map(self.lambda_func, self.forward_prepare(input_tensor)))


class LambdaReduce(LambdaBase):

    # 执行当前模型模块的前向传播
    def forward(self, input_tensor):
        return reduce(self.lambda_func, self.forward_prepare(input_tensor))


RESNEXT_101_32X4D = nn.Sequential(
    nn.Conv2d(3, 64, (7, 7), (2, 2), (3, 3), 1, 1, bias=False),
    nn.BatchNorm2d(64),
    nn.ReLU(),
    nn.MaxPool2d((3, 3), (2, 2), (1, 1)),
    nn.Sequential(
        nn.Sequential(
            LambdaMap(lambda x: x,
                      nn.Sequential(
                          nn.Sequential(
                              nn.Conv2d(64, 128, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                              nn.BatchNorm2d(128),
                              nn.ReLU(),
                              nn.Conv2d(128, 128, (3, 3), (1, 1), (1, 1), 1, 32, bias=False),
                              nn.BatchNorm2d(128),
                              nn.ReLU(),
                          ),
                          nn.Conv2d(128, 256, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                          nn.BatchNorm2d(256),
                      ),
                      nn.Sequential(
                          nn.Conv2d(64, 256, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                          nn.BatchNorm2d(256),
                      ),
                      ),
            LambdaReduce(lambda x, y: x + y),
            nn.ReLU(),
        ),
        nn.Sequential(
            LambdaMap(lambda x: x,
                      nn.Sequential(
                          nn.Sequential(
                              nn.Conv2d(256, 128, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                              nn.BatchNorm2d(128),
                              nn.ReLU(),
                              nn.Conv2d(128, 128, (3, 3), (1, 1), (1, 1), 1, 32, bias=False),
                              nn.BatchNorm2d(128),
                              nn.ReLU(),
                          ),
                          nn.Conv2d(128, 256, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                          nn.BatchNorm2d(256),
                      ),
                      Lambda(lambda x: x),
                      ),
            LambdaReduce(lambda x, y: x + y),
            nn.ReLU(),
        ),
        nn.Sequential(
            LambdaMap(lambda x: x,
                      nn.Sequential(
                          nn.Sequential(
                              nn.Conv2d(256, 128, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                              nn.BatchNorm2d(128),
                              nn.ReLU(),
                              nn.Conv2d(128, 128, (3, 3), (1, 1), (1, 1), 1, 32, bias=False),
                              nn.BatchNorm2d(128),
                              nn.ReLU(),
                          ),
                          nn.Conv2d(128, 256, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                          nn.BatchNorm2d(256),
                      ),
                      Lambda(lambda x: x),
                      ),
            LambdaReduce(lambda x, y: x + y),
            nn.ReLU(),
        ),
    ),
    nn.Sequential(
        nn.Sequential(
            LambdaMap(lambda x: x,
                      nn.Sequential(
                          nn.Sequential(
                              nn.Conv2d(256, 256, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                              nn.BatchNorm2d(256),
                              nn.ReLU(),
                              nn.Conv2d(256, 256, (3, 3), (2, 2), (1, 1), 1, 32, bias=False),
                              nn.BatchNorm2d(256),
                              nn.ReLU(),
                          ),
                          nn.Conv2d(256, 512, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                          nn.BatchNorm2d(512),
                      ),
                      nn.Sequential(
                          nn.Conv2d(256, 512, (1, 1), (2, 2), (0, 0), 1, 1, bias=False),
                          nn.BatchNorm2d(512),
                      ),
                      ),
            LambdaReduce(lambda x, y: x + y),
            nn.ReLU(),
        ),
        nn.Sequential(
            LambdaMap(lambda x: x,
                      nn.Sequential(
                          nn.Sequential(
                              nn.Conv2d(512, 256, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                              nn.BatchNorm2d(256),
                              nn.ReLU(),
                              nn.Conv2d(256, 256, (3, 3), (1, 1), (1, 1), 1, 32, bias=False),
                              nn.BatchNorm2d(256),
                              nn.ReLU(),
                          ),
                          nn.Conv2d(256, 512, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                          nn.BatchNorm2d(512),
                      ),
                      Lambda(lambda x: x),
                      ),
            LambdaReduce(lambda x, y: x + y),
            nn.ReLU(),
        ),
        nn.Sequential(
            LambdaMap(lambda x: x,
                      nn.Sequential(
                          nn.Sequential(
                              nn.Conv2d(512, 256, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                              nn.BatchNorm2d(256),
                              nn.ReLU(),
                              nn.Conv2d(256, 256, (3, 3), (1, 1), (1, 1), 1, 32, bias=False),
                              nn.BatchNorm2d(256),
                              nn.ReLU(),
                          ),
                          nn.Conv2d(256, 512, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                          nn.BatchNorm2d(512),
                      ),
                      Lambda(lambda x: x),
                      ),
            LambdaReduce(lambda x, y: x + y),
            nn.ReLU(),
        ),
        nn.Sequential(
            LambdaMap(lambda x: x,
                      nn.Sequential(
                          nn.Sequential(
                              nn.Conv2d(512, 256, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                              nn.BatchNorm2d(256),
                              nn.ReLU(),
                              nn.Conv2d(256, 256, (3, 3), (1, 1), (1, 1), 1, 32, bias=False),
                              nn.BatchNorm2d(256),
                              nn.ReLU(),
                          ),
                          nn.Conv2d(256, 512, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                          nn.BatchNorm2d(512),
                      ),
                      Lambda(lambda x: x),
                      ),
            LambdaReduce(lambda x, y: x + y),
            nn.ReLU(),
        ),
    ),
    nn.Sequential(
        nn.Sequential(
            LambdaMap(lambda x: x,
                      nn.Sequential(
                          nn.Sequential(
                              nn.Conv2d(512, 512, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                              nn.BatchNorm2d(512),
                              nn.ReLU(),
                              nn.Conv2d(512, 512, (3, 3), (2, 2), (1, 1), 1, 32, bias=False),
                              nn.BatchNorm2d(512),
                              nn.ReLU(),
                          ),
                          nn.Conv2d(512, 1024, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                          nn.BatchNorm2d(1024),
                      ),
                      nn.Sequential(
                          nn.Conv2d(512, 1024, (1, 1), (2, 2), (0, 0), 1, 1, bias=False),
                          nn.BatchNorm2d(1024),
                      ),
                      ),
            LambdaReduce(lambda x, y: x + y),
            nn.ReLU(),
        ),
        nn.Sequential(
            LambdaMap(lambda x: x,
                      nn.Sequential(
                          nn.Sequential(
                              nn.Conv2d(1024, 512, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                              nn.BatchNorm2d(512),
                              nn.ReLU(),
                              nn.Conv2d(512, 512, (3, 3), (1, 1), (1, 1), 1, 32, bias=False),
                              nn.BatchNorm2d(512),
                              nn.ReLU(),
                          ),
                          nn.Conv2d(512, 1024, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                          nn.BatchNorm2d(1024),
                      ),
                      Lambda(lambda x: x),
                      ),
            LambdaReduce(lambda x, y: x + y),
            nn.ReLU(),
        ),
        nn.Sequential(
            LambdaMap(lambda x: x,
                      nn.Sequential(
                          nn.Sequential(
                              nn.Conv2d(1024, 512, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                              nn.BatchNorm2d(512),
                              nn.ReLU(),
                              nn.Conv2d(512, 512, (3, 3), (1, 1), (1, 1), 1, 32, bias=False),
                              nn.BatchNorm2d(512),
                              nn.ReLU(),
                          ),
                          nn.Conv2d(512, 1024, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                          nn.BatchNorm2d(1024),
                      ),
                      Lambda(lambda x: x),
                      ),
            LambdaReduce(lambda x, y: x + y),
            nn.ReLU(),
        ),
        nn.Sequential(
            LambdaMap(lambda x: x,
                      nn.Sequential(
                          nn.Sequential(
                              nn.Conv2d(1024, 512, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                              nn.BatchNorm2d(512),
                              nn.ReLU(),
                              nn.Conv2d(512, 512, (3, 3), (1, 1), (1, 1), 1, 32, bias=False),
                              nn.BatchNorm2d(512),
                              nn.ReLU(),
                          ),
                          nn.Conv2d(512, 1024, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                          nn.BatchNorm2d(1024),
                      ),
                      Lambda(lambda x: x),
                      ),
            LambdaReduce(lambda x, y: x + y),
            nn.ReLU(),
        ),
        nn.Sequential(
            LambdaMap(lambda x: x,
                      nn.Sequential(
                          nn.Sequential(
                              nn.Conv2d(1024, 512, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                              nn.BatchNorm2d(512),
                              nn.ReLU(),
                              nn.Conv2d(512, 512, (3, 3), (1, 1), (1, 1), 1, 32, bias=False),
                              nn.BatchNorm2d(512),
                              nn.ReLU(),
                          ),
                          nn.Conv2d(512, 1024, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                          nn.BatchNorm2d(1024),
                      ),
                      Lambda(lambda x: x),
                      ),
            LambdaReduce(lambda x, y: x + y),
            nn.ReLU(),
        ),
        nn.Sequential(
            LambdaMap(lambda x: x,
                      nn.Sequential(
                          nn.Sequential(
                              nn.Conv2d(1024, 512, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                              nn.BatchNorm2d(512),
                              nn.ReLU(),
                              nn.Conv2d(512, 512, (3, 3), (1, 1), (1, 1), 1, 32, bias=False),
                              nn.BatchNorm2d(512),
                              nn.ReLU(),
                          ),
                          nn.Conv2d(512, 1024, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                          nn.BatchNorm2d(1024),
                      ),
                      Lambda(lambda x: x),
                      ),
            LambdaReduce(lambda x, y: x + y),
            nn.ReLU(),
        ),
        nn.Sequential(
            LambdaMap(lambda x: x,
                      nn.Sequential(
                          nn.Sequential(
                              nn.Conv2d(1024, 512, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                              nn.BatchNorm2d(512),
                              nn.ReLU(),
                              nn.Conv2d(512, 512, (3, 3), (1, 1), (1, 1), 1, 32, bias=False),
                              nn.BatchNorm2d(512),
                              nn.ReLU(),
                          ),
                          nn.Conv2d(512, 1024, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                          nn.BatchNorm2d(1024),
                      ),
                      Lambda(lambda x: x),
                      ),
            LambdaReduce(lambda x, y: x + y),
            nn.ReLU(),
        ),
        nn.Sequential(
            LambdaMap(lambda x: x,
                      nn.Sequential(
                          nn.Sequential(
                              nn.Conv2d(1024, 512, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                              nn.BatchNorm2d(512),
                              nn.ReLU(),
                              nn.Conv2d(512, 512, (3, 3), (1, 1), (1, 1), 1, 32, bias=False),
                              nn.BatchNorm2d(512),
                              nn.ReLU(),
                          ),
                          nn.Conv2d(512, 1024, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                          nn.BatchNorm2d(1024),
                      ),
                      Lambda(lambda x: x),
                      ),
            LambdaReduce(lambda x, y: x + y),
            nn.ReLU(),
        ),
        nn.Sequential(
            LambdaMap(lambda x: x,
                      nn.Sequential(
                          nn.Sequential(
                              nn.Conv2d(1024, 512, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                              nn.BatchNorm2d(512),
                              nn.ReLU(),
                              nn.Conv2d(512, 512, (3, 3), (1, 1), (1, 1), 1, 32, bias=False),
                              nn.BatchNorm2d(512),
                              nn.ReLU(),
                          ),
                          nn.Conv2d(512, 1024, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                          nn.BatchNorm2d(1024),
                      ),
                      Lambda(lambda x: x),
                      ),
            LambdaReduce(lambda x, y: x + y),
            nn.ReLU(),
        ),
        nn.Sequential(
            LambdaMap(lambda x: x,
                      nn.Sequential(
                          nn.Sequential(
                              nn.Conv2d(1024, 512, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                              nn.BatchNorm2d(512),
                              nn.ReLU(),
                              nn.Conv2d(512, 512, (3, 3), (1, 1), (1, 1), 1, 32, bias=False),
                              nn.BatchNorm2d(512),
                              nn.ReLU(),
                          ),
                          nn.Conv2d(512, 1024, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                          nn.BatchNorm2d(1024),
                      ),
                      Lambda(lambda x: x),
                      ),
            LambdaReduce(lambda x, y: x + y),
            nn.ReLU(),
        ),
        nn.Sequential(
            LambdaMap(lambda x: x,
                      nn.Sequential(
                          nn.Sequential(
                              nn.Conv2d(1024, 512, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                              nn.BatchNorm2d(512),
                              nn.ReLU(),
                              nn.Conv2d(512, 512, (3, 3), (1, 1), (1, 1), 1, 32, bias=False),
                              nn.BatchNorm2d(512),
                              nn.ReLU(),
                          ),
                          nn.Conv2d(512, 1024, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                          nn.BatchNorm2d(1024),
                      ),
                      Lambda(lambda x: x),
                      ),
            LambdaReduce(lambda x, y: x + y),
            nn.ReLU(),
        ),
        nn.Sequential(
            LambdaMap(lambda x: x,
                      nn.Sequential(
                          nn.Sequential(
                              nn.Conv2d(1024, 512, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                              nn.BatchNorm2d(512),
                              nn.ReLU(),
                              nn.Conv2d(512, 512, (3, 3), (1, 1), (1, 1), 1, 32, bias=False),
                              nn.BatchNorm2d(512),
                              nn.ReLU(),
                          ),
                          nn.Conv2d(512, 1024, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                          nn.BatchNorm2d(1024),
                      ),
                      Lambda(lambda x: x),
                      ),
            LambdaReduce(lambda x, y: x + y),
            nn.ReLU(),
        ),
        nn.Sequential(
            LambdaMap(lambda x: x,
                      nn.Sequential(
                          nn.Sequential(
                              nn.Conv2d(1024, 512, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                              nn.BatchNorm2d(512),
                              nn.ReLU(),
                              nn.Conv2d(512, 512, (3, 3), (1, 1), (1, 1), 1, 32, bias=False),
                              nn.BatchNorm2d(512),
                              nn.ReLU(),
                          ),
                          nn.Conv2d(512, 1024, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                          nn.BatchNorm2d(1024),
                      ),
                      Lambda(lambda x: x),
                      ),
            LambdaReduce(lambda x, y: x + y),
            nn.ReLU(),
        ),
        nn.Sequential(
            LambdaMap(lambda x: x,
                      nn.Sequential(
                          nn.Sequential(
                              nn.Conv2d(1024, 512, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                              nn.BatchNorm2d(512),
                              nn.ReLU(),
                              nn.Conv2d(512, 512, (3, 3), (1, 1), (1, 1), 1, 32, bias=False),
                              nn.BatchNorm2d(512),
                              nn.ReLU(),
                          ),
                          nn.Conv2d(512, 1024, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                          nn.BatchNorm2d(1024),
                      ),
                      Lambda(lambda x: x),
                      ),
            LambdaReduce(lambda x, y: x + y),
            nn.ReLU(),
        ),
        nn.Sequential(
            LambdaMap(lambda x: x,
                      nn.Sequential(
                          nn.Sequential(
                              nn.Conv2d(1024, 512, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                              nn.BatchNorm2d(512),
                              nn.ReLU(),
                              nn.Conv2d(512, 512, (3, 3), (1, 1), (1, 1), 1, 32, bias=False),
                              nn.BatchNorm2d(512),
                              nn.ReLU(),
                          ),
                          nn.Conv2d(512, 1024, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                          nn.BatchNorm2d(1024),
                      ),
                      Lambda(lambda x: x),
                      ),
            LambdaReduce(lambda x, y: x + y),
            nn.ReLU(),
        ),
        nn.Sequential(
            LambdaMap(lambda x: x,
                      nn.Sequential(
                          nn.Sequential(
                              nn.Conv2d(1024, 512, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                              nn.BatchNorm2d(512),
                              nn.ReLU(),
                              nn.Conv2d(512, 512, (3, 3), (1, 1), (1, 1), 1, 32, bias=False),
                              nn.BatchNorm2d(512),
                              nn.ReLU(),
                          ),
                          nn.Conv2d(512, 1024, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                          nn.BatchNorm2d(1024),
                      ),
                      Lambda(lambda x: x),
                      ),
            LambdaReduce(lambda x, y: x + y),
            nn.ReLU(),
        ),
        nn.Sequential(
            LambdaMap(lambda x: x,
                      nn.Sequential(
                          nn.Sequential(
                              nn.Conv2d(1024, 512, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                              nn.BatchNorm2d(512),
                              nn.ReLU(),
                              nn.Conv2d(512, 512, (3, 3), (1, 1), (1, 1), 1, 32, bias=False),
                              nn.BatchNorm2d(512),
                              nn.ReLU(),
                          ),
                          nn.Conv2d(512, 1024, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                          nn.BatchNorm2d(1024),
                      ),
                      Lambda(lambda x: x),
                      ),
            LambdaReduce(lambda x, y: x + y),
            nn.ReLU(),
        ),
        nn.Sequential(
            LambdaMap(lambda x: x,
                      nn.Sequential(
                          nn.Sequential(
                              nn.Conv2d(1024, 512, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                              nn.BatchNorm2d(512),
                              nn.ReLU(),
                              nn.Conv2d(512, 512, (3, 3), (1, 1), (1, 1), 1, 32, bias=False),
                              nn.BatchNorm2d(512),
                              nn.ReLU(),
                          ),
                          nn.Conv2d(512, 1024, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                          nn.BatchNorm2d(1024),
                      ),
                      Lambda(lambda x: x),
                      ),
            LambdaReduce(lambda x, y: x + y),
            nn.ReLU(),
        ),
        nn.Sequential(
            LambdaMap(lambda x: x,
                      nn.Sequential(
                          nn.Sequential(
                              nn.Conv2d(1024, 512, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                              nn.BatchNorm2d(512),
                              nn.ReLU(),
                              nn.Conv2d(512, 512, (3, 3), (1, 1), (1, 1), 1, 32, bias=False),
                              nn.BatchNorm2d(512),
                              nn.ReLU(),
                          ),
                          nn.Conv2d(512, 1024, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                          nn.BatchNorm2d(1024),
                      ),
                      Lambda(lambda x: x),
                      ),
            LambdaReduce(lambda x, y: x + y),
            nn.ReLU(),
        ),
        nn.Sequential(
            LambdaMap(lambda x: x,
                      nn.Sequential(
                          nn.Sequential(
                              nn.Conv2d(1024, 512, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                              nn.BatchNorm2d(512),
                              nn.ReLU(),
                              nn.Conv2d(512, 512, (3, 3), (1, 1), (1, 1), 1, 32, bias=False),
                              nn.BatchNorm2d(512),
                              nn.ReLU(),
                          ),
                          nn.Conv2d(512, 1024, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                          nn.BatchNorm2d(1024),
                      ),
                      Lambda(lambda x: x),
                      ),
            LambdaReduce(lambda x, y: x + y),
            nn.ReLU(),
        ),
        nn.Sequential(
            LambdaMap(lambda x: x,
                      nn.Sequential(
                          nn.Sequential(
                              nn.Conv2d(1024, 512, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                              nn.BatchNorm2d(512),
                              nn.ReLU(),
                              nn.Conv2d(512, 512, (3, 3), (1, 1), (1, 1), 1, 32, bias=False),
                              nn.BatchNorm2d(512),
                              nn.ReLU(),
                          ),
                          nn.Conv2d(512, 1024, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                          nn.BatchNorm2d(1024),
                      ),
                      Lambda(lambda x: x),
                      ),
            LambdaReduce(lambda x, y: x + y),
            nn.ReLU(),
        ),
        nn.Sequential(
            LambdaMap(lambda x: x,
                      nn.Sequential(
                          nn.Sequential(
                              nn.Conv2d(1024, 512, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                              nn.BatchNorm2d(512),
                              nn.ReLU(),
                              nn.Conv2d(512, 512, (3, 3), (1, 1), (1, 1), 1, 32, bias=False),
                              nn.BatchNorm2d(512),
                              nn.ReLU(),
                          ),
                          nn.Conv2d(512, 1024, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                          nn.BatchNorm2d(1024),
                      ),
                      Lambda(lambda x: x),
                      ),
            LambdaReduce(lambda x, y: x + y),
            nn.ReLU(),
        ),
        nn.Sequential(
            LambdaMap(lambda x: x,
                      nn.Sequential(
                          nn.Sequential(
                              nn.Conv2d(1024, 512, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                              nn.BatchNorm2d(512),
                              nn.ReLU(),
                              nn.Conv2d(512, 512, (3, 3), (1, 1), (1, 1), 1, 32, bias=False),
                              nn.BatchNorm2d(512),
                              nn.ReLU(),
                          ),
                          nn.Conv2d(512, 1024, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                          nn.BatchNorm2d(1024),
                      ),
                      Lambda(lambda x: x),
                      ),
            LambdaReduce(lambda x, y: x + y),
            nn.ReLU(),
        ),
    ),
    nn.Sequential(
        nn.Sequential(
            LambdaMap(lambda x: x,
                      nn.Sequential(
                          nn.Sequential(
                              nn.Conv2d(1024, 1024, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                              nn.BatchNorm2d(1024),
                              nn.ReLU(),
                              nn.Conv2d(1024, 1024, (3, 3), (2, 2), (1, 1), 1, 32, bias=False),
                              nn.BatchNorm2d(1024),
                              nn.ReLU(),
                          ),
                          nn.Conv2d(1024, 2048, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                          nn.BatchNorm2d(2048),
                      ),
                      nn.Sequential(
                          nn.Conv2d(1024, 2048, (1, 1), (2, 2), (0, 0), 1, 1, bias=False),
                          nn.BatchNorm2d(2048),
                      ),
                      ),
            LambdaReduce(lambda x, y: x + y),
            nn.ReLU(),
        ),
        nn.Sequential(
            LambdaMap(lambda x: x,
                      nn.Sequential(
                          nn.Sequential(
                              nn.Conv2d(2048, 1024, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                              nn.BatchNorm2d(1024),
                              nn.ReLU(),
                              nn.Conv2d(1024, 1024, (3, 3), (1, 1), (1, 1), 1, 32, bias=False),
                              nn.BatchNorm2d(1024),
                              nn.ReLU(),
                          ),
                          nn.Conv2d(1024, 2048, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                          nn.BatchNorm2d(2048),
                      ),
                      Lambda(lambda x: x),
                      ),
            LambdaReduce(lambda x, y: x + y),
            nn.ReLU(),
        ),
        nn.Sequential(
            LambdaMap(lambda x: x,
                      nn.Sequential(
                          nn.Sequential(
                              nn.Conv2d(2048, 1024, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                              nn.BatchNorm2d(1024),
                              nn.ReLU(),
                              nn.Conv2d(1024, 1024, (3, 3), (1, 1), (1, 1), 1, 32, bias=False),
                              nn.BatchNorm2d(1024),
                              nn.ReLU(),
                          ),
                          nn.Conv2d(1024, 2048, (1, 1), (1, 1), (0, 0), 1, 1, bias=False),
                          nn.BatchNorm2d(2048),
                      ),
                      Lambda(lambda x: x),
                      ),
            LambdaReduce(lambda x, y: x + y),
            nn.ReLU(),
        ),
    ),
    nn.AvgPool2d((7, 7), (1, 1)),
    Lambda(lambda x: x.view(x.size(0), -1)),
    nn.Sequential(Lambda(lambda x: x.view(1, -1) if 1 == len(x.size()) else x), nn.Linear(2048, 1000)),
)


class ResNeXt101(nn.Module):

    # 初始化当前模型模块的参数
    def __init__(self, backbone_path):
        super(ResNeXt101, self).__init__()
        net = RESNEXT_101_32X4D
        if backbone_path is not None:
            weights = torch.load(backbone_path)
            net.load_state_dict(weights, strict=True)

        net = list(net.children())
        self.layer0 = nn.Sequential(*net[:3])
        self.layer1 = nn.Sequential(*net[3:5])
        self.layer2 = net[5]
        self.layer3 = net[6]
        self.layer4 = net[7]

    # 执行当前模型模块的前向传播
    def forward(self, x):
        layer0 = self.layer0(x)
        layer1 = self.layer1(layer0)
        layer2 = self.layer2(layer1)
        layer3 = self.layer3(layer2)
        layer4 = self.layer4(layer3)
        return layer4

class BasicConv(nn.Module):

    # 初始化当前模型模块的参数
    def __init__(self, in_planes, out_planes, kernel_size, stride=1, padding=0, dilation=1, groups=1, relu=True,
                 bn=True, bias=False):
        super(BasicConv, self).__init__()
        self.out_channels = out_planes
        self.conv = nn.Conv2d(in_planes, out_planes, kernel_size=kernel_size, stride=stride, padding=padding,
                              dilation=dilation, groups=groups, bias=bias)
        self.bn = nn.BatchNorm2d(out_planes, eps=1e-5, momentum=0.01, affine=True) if bn else None
        self.relu = nn.ReLU() if relu else None

    # 执行当前模型模块的前向传播
    def forward(self, x):
        x = self.conv(x)
        if self.bn is not None:
            x = self.bn(x)
        if self.relu is not None:
            x = self.relu(x)
        return x


class Flatten(nn.Module):

    # 执行当前模型模块的前向传播
    def forward(self, x):
        return x.view(x.size(0), -1)


class ChannelGate(nn.Module):

    # 初始化当前模型模块的参数
    def __init__(self, gate_channels, reduction_ratio=16, pool_types=['avg']):
        super(ChannelGate, self).__init__()
        self.gate_channels = gate_channels
        self.mlp = nn.Sequential(
            Flatten(),
            nn.Linear(gate_channels, gate_channels // reduction_ratio),
            nn.ReLU(),
            nn.Linear(gate_channels // reduction_ratio, gate_channels)
        )
        self.pool_types = pool_types

    # 执行当前模型模块的前向传播
    def forward(self, x):
        channel_att_sum = None
        for pool_type in self.pool_types:
            if pool_type == 'avg':
                avg_pool = F.avg_pool2d(x, (x.size(2), x.size(3)), stride=(x.size(2), x.size(3)))
                channel_att_raw = self.mlp(avg_pool)
            elif pool_type == 'max':
                max_pool = F.max_pool2d(x, (x.size(2), x.size(3)), stride=(x.size(2), x.size(3)))
                channel_att_raw = self.mlp(max_pool)
            elif pool_type == 'lp':
                lp_pool = F.lp_pool2d(x, 2, (x.size(2), x.size(3)), stride=(x.size(2), x.size(3)))
                channel_att_raw = self.mlp(lp_pool)
            elif pool_type == 'lse':
                lse_pool = logsumexp_2d(x)
                channel_att_raw = self.mlp(lse_pool)

            if channel_att_sum is None:
                channel_att_sum = channel_att_raw
            else:
                channel_att_sum = channel_att_sum + channel_att_raw

        scale = F.sigmoid(channel_att_sum).unsqueeze(2).unsqueeze(3).expand_as(x)
        return x * scale


# 计算二维特征图的 log-sum-exp 池化结果
def logsumexp_2d(tensor):
    tensor_flatten = tensor.view(tensor.size(0), tensor.size(1), -1)
    s, _ = torch.max(tensor_flatten, dim=2, keepdim=True)
    outputs = s + (tensor_flatten - s).exp().sum(dim=2, keepdim=True).log()
    return outputs


class ChannelPool(nn.Module):

    # 执行当前模型模块的前向传播
    def forward(self, x):
        return torch.mean(x, 1).unsqueeze(1)


class SpatialGate(nn.Module):

    # 初始化当前模型模块的参数
    def __init__(self):
        super(SpatialGate, self).__init__()
        kernel_size = 7
        self.compress = ChannelPool()
        self.spatial = BasicConv(1, 1, kernel_size, stride=1, padding=(kernel_size - 1) // 2, relu=False)

    # 执行当前模型模块的前向传播
    def forward(self, x):
        x_compress = self.compress(x)
        x_out = self.spatial(x_compress)
        scale = F.sigmoid(x_out)
        return x * scale


class CBAM(nn.Module):

    # 初始化当前模型模块的参数
    def __init__(self, gate_channels=128, reduction_ratio=16, pool_types=['avg'], no_spatial=False):
        super(CBAM, self).__init__()
        self.ChannelGate = ChannelGate(gate_channels, reduction_ratio, pool_types)
        self.no_spatial = no_spatial
        if not no_spatial:
            self.SpatialGate = SpatialGate()

    # 执行当前模型模块的前向传播
    def forward(self, x):
        x_out = self.ChannelGate(x)
        if not self.no_spatial:
            x_out = self.SpatialGate(x_out)
        return x_out


class LCFI(nn.Module):

    # 初始化当前模型模块的参数
    def __init__(self, input_channels, dr1=1, dr2=2, dr3=3, dr4=4):
        super(LCFI, self).__init__()
        self.input_channels = input_channels
        self.channels_single = int(input_channels / 4)
        self.channels_double = int(input_channels / 2)
        self.dr1 = dr1
        self.dr2 = dr2
        self.dr3 = dr3
        self.dr4 = dr4
        self.padding1 = 1 * dr1
        self.padding2 = 2 * dr2
        self.padding3 = 3 * dr3
        self.padding4 = 4 * dr4

        self.p1_channel_reduction = nn.Sequential(
            nn.Conv2d(self.input_channels, self.channels_single, 3, 1, 1, dilation=1),
            nn.BatchNorm2d(self.channels_single), nn.ReLU())
        self.p2_channel_reduction = nn.Sequential(
            nn.Conv2d(self.input_channels, self.channels_single, 3, 1, 1, dilation=1),
            nn.BatchNorm2d(self.channels_single), nn.ReLU())
        self.p3_channel_reduction = nn.Sequential(
            nn.Conv2d(self.input_channels, self.channels_single, 3, 1, 1, dilation=1),
            nn.BatchNorm2d(self.channels_single), nn.ReLU())
        self.p4_channel_reduction = nn.Sequential(
            nn.Conv2d(self.input_channels, self.channels_single, 3, 1, 1, dilation=1),
            nn.BatchNorm2d(self.channels_single), nn.ReLU())

        self.p1_d1 = nn.Sequential(
            nn.Conv2d(self.channels_single, self.channels_single, (3, 1), 1, padding=(self.padding1, 0),
                      dilation=(self.dr1, 1)),
            nn.Conv2d(self.channels_single, self.channels_single, (1, 3), 1, padding=(0, self.padding1),
                      dilation=(1, self.dr1)),
            nn.BatchNorm2d(self.channels_single), nn.ReLU())
        self.p1_d2 = nn.Sequential(
            nn.Conv2d(self.channels_single, self.channels_single, (1, 3), 1, padding=(0, self.padding1),
                      dilation=(1, self.dr1)),
            nn.Conv2d(self.channels_single, self.channels_single, (3, 1), 1, padding=(self.padding1, 0),
                      dilation=(self.dr1, 1)),
            nn.BatchNorm2d(self.channels_single), nn.ReLU())
        self.p1_fusion = nn.Sequential(nn.Conv2d(self.channels_double, self.channels_single, 3, 1, 1, dilation=1),
                                       nn.BatchNorm2d(self.channels_single), nn.ReLU())

        self.p2_d1 = nn.Sequential(
            nn.Conv2d(self.channels_double, self.channels_single, (5, 1), 1, padding=(self.padding2, 0),
                      dilation=(self.dr2, 1)),
            nn.Conv2d(self.channels_single, self.channels_single, (1, 5), 1, padding=(0, self.padding2),
                      dilation=(1, self.dr2)),
            nn.BatchNorm2d(self.channels_single), nn.ReLU())
        self.p2_d2 = nn.Sequential(
            nn.Conv2d(self.channels_double, self.channels_single, (1, 5), 1, padding=(0, self.padding2),
                      dilation=(1, self.dr2)),
            nn.Conv2d(self.channels_single, self.channels_single, (5, 1), 1, padding=(self.padding2, 0),
                      dilation=(self.dr2, 1)),
            nn.BatchNorm2d(self.channels_single), nn.ReLU())
        self.p2_fusion = nn.Sequential(nn.Conv2d(self.channels_double, self.channels_single, 3, 1, 1, dilation=1),
                                       nn.BatchNorm2d(self.channels_single), nn.ReLU())

        self.p3_d1 = nn.Sequential(
            nn.Conv2d(self.channels_double, self.channels_single, (7, 1), 1, padding=(self.padding3, 0),
                      dilation=(self.dr3, 1)),
            nn.Conv2d(self.channels_single, self.channels_single, (1, 7), 1, padding=(0, self.padding3),
                      dilation=(1, self.dr3)),
            nn.BatchNorm2d(self.channels_single), nn.ReLU())
        self.p3_d2 = nn.Sequential(
            nn.Conv2d(self.channels_double, self.channels_single, (1, 7), 1, padding=(0, self.padding3),
                      dilation=(1, self.dr3)),
            nn.Conv2d(self.channels_single, self.channels_single, (7, 1), 1, padding=(self.padding3, 0),
                      dilation=(self.dr3, 1)),
            nn.BatchNorm2d(self.channels_single), nn.ReLU())
        self.p3_fusion = nn.Sequential(nn.Conv2d(self.channels_double, self.channels_single, 3, 1, 1, dilation=1),
                                       nn.BatchNorm2d(self.channels_single), nn.ReLU())

        self.p4_d1 = nn.Sequential(
            nn.Conv2d(self.channels_double, self.channels_single, (9, 1), 1, padding=(self.padding4, 0),
                      dilation=(self.dr4, 1)),
            nn.Conv2d(self.channels_single, self.channels_single, (1, 9), 1, padding=(0, self.padding4),
                      dilation=(1, self.dr4)),
            nn.BatchNorm2d(self.channels_single), nn.ReLU())
        self.p4_d2 = nn.Sequential(
            nn.Conv2d(self.channels_double, self.channels_single, (1, 9), 1, padding=(0, self.padding4),
                      dilation=(1, self.dr4)),
            nn.Conv2d(self.channels_single, self.channels_single, (9, 1), 1, padding=(self.padding4, 0),
                      dilation=(self.dr4, 1)),
            nn.BatchNorm2d(self.channels_single), nn.ReLU())
        self.p4_fusion = nn.Sequential(nn.Conv2d(self.channels_double, self.channels_single, 3, 1, 1, dilation=1),
                                       nn.BatchNorm2d(self.channels_single), nn.ReLU())

        self.cbam = CBAM(self.input_channels)

        self.channel_reduction = nn.Sequential(
            nn.Conv2d(self.input_channels, self.channels_single, 3, 1, 1, dilation=1),
            nn.BatchNorm2d(self.channels_single),
            nn.ReLU())

    # 执行当前模型模块的前向传播
    def forward(self, x):
        p1_input = self.p1_channel_reduction(x)
        p1 = self.p1_fusion(torch.cat((self.p1_d1(p1_input), self.p1_d2(p1_input)), 1))

        p2_input = torch.cat((self.p2_channel_reduction(x), p1), 1)
        p2 = self.p2_fusion(torch.cat((self.p2_d1(p2_input), self.p2_d2(p2_input)), 1))

        p3_input = torch.cat((self.p3_channel_reduction(x), p2), 1)
        p3 = self.p3_fusion(torch.cat((self.p3_d1(p3_input), self.p3_d2(p3_input)), 1))

        p4_input = torch.cat((self.p4_channel_reduction(x), p3), 1)
        p4 = self.p4_fusion(torch.cat((self.p4_d1(p4_input), self.p4_d2(p4_input)), 1))

        channel_reduction = self.channel_reduction(self.cbam(torch.cat((p1, p2, p3, p4), 1)))

        return channel_reduction


class GDNet(nn.Module):

    # 初始化当前模型模块的参数
    def __init__(self, backbone_path=None):
        super(GDNet, self).__init__()

        resnext = ResNeXt101(backbone_path)
        self.layer0 = resnext.layer0
        self.layer1 = resnext.layer1
        self.layer2 = resnext.layer2
        self.layer3 = resnext.layer3
        self.layer4 = resnext.layer4

        self.h5_conv = LCFI(2048, 1, 2, 3, 4)
        self.h4_conv = LCFI(1024, 1, 2, 3, 4)
        self.h3_conv = LCFI(512, 1, 2, 3, 4)
        self.l2_conv = LCFI(256, 1, 2, 3, 4)

        self.h5_up = nn.UpsamplingBilinear2d(scale_factor=2)
        self.h3_down = nn.AvgPool2d((2, 2), stride=2)
        self.h_fusion = CBAM(896)
        self.h_fusion_conv = nn.Sequential(nn.Conv2d(896, 896, 3, 1, 1), nn.BatchNorm2d(896), nn.ReLU())

        self.l_fusion_conv = nn.Sequential(nn.Conv2d(64, 64, 3, 1, 1), nn.BatchNorm2d(64), nn.ReLU())
        self.h2l = nn.ConvTranspose2d(896, 1, 8, 4, 2)

        self.h_up_for_final_fusion = nn.ConvTranspose2d(896, 256, 8, 4, 2)
        self.final_fusion = CBAM(320)
        self.final_fusion_conv = nn.Sequential(nn.Conv2d(320, 320, 3, 1, 1), nn.BatchNorm2d(320), nn.ReLU())

        self.h_predict = nn.Conv2d(896, 1, 3, 1, 1)
        self.l_predict = nn.Conv2d(64, 1, 3, 1, 1)
        self.final_predict = nn.Conv2d(320, 1, 3, 1, 1)

        for m in self.modules():
            if isinstance(m, nn.ReLU):
                m.inplace = True

    # 执行当前模型模块的前向传播
    def forward(self, x):
        layer0 = self.layer0(x)
        layer1 = self.layer1(layer0)
        layer2 = self.layer2(layer1)
        layer3 = self.layer3(layer2)
        layer4 = self.layer4(layer3)

        h5_conv = self.h5_conv(layer4)
        h4_conv = self.h4_conv(layer3)
        h3_conv = self.h3_conv(layer2)
        l2_conv = self.l2_conv(layer1)

        h5_up = self.h5_up(h5_conv)
        h3_down = self.h3_down(h3_conv)
        h_fusion = self.h_fusion(torch.cat((h5_up, h4_conv, h3_down), 1))
        h_fusion = self.h_fusion_conv(h_fusion)

        l_fusion = self.l_fusion_conv(l2_conv)
        h2l = self.h2l(h_fusion)
        l_fusion = F.sigmoid(h2l) * l_fusion

        h_up_for_final_fusion = self.h_up_for_final_fusion(h_fusion)
        final_fusion = self.final_fusion(torch.cat((h_up_for_final_fusion, l_fusion), 1))
        final_fusion = self.final_fusion_conv(final_fusion)

        h_predict = self.h_predict(h_fusion)

        l_predict = self.l_predict(l_fusion)

        final_predict = self.final_predict(final_fusion)

        h_predict = F.upsample(h_predict, size=x.size()[2:], mode='bilinear', align_corners=True)
        l_predict = F.upsample(l_predict, size=x.size()[2:], mode='bilinear', align_corners=True)
        final_predict = F.upsample(final_predict, size=x.size()[2:], mode='bilinear', align_corners=True)

        return torch.sigmoid(h_predict), torch.sigmoid(l_predict), torch.sigmoid(final_predict)


# 获取当前可用的推理设备
def get_device() -> torch.device:
    return torch.device('cuda' if torch.cuda.is_available() else 'cpu')


# 加载模型权重并切换为推理模式
def load_model(model_path: Path, device: torch.device) -> nn.Module:
    model = GDNet()
    checkpoint = torch.load(model_path, map_location=device)
    if isinstance(checkpoint, dict) and 'state_dict' in checkpoint:
        checkpoint = checkpoint['state_dict']
    if isinstance(checkpoint, dict) and any(key.startswith('module.') for key in checkpoint):
        checkpoint = {key.replace('module.', '', 1): value for key, value in checkpoint.items()}
    model.load_state_dict(checkpoint)
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
            transforms.Normalize([0.485, 0.456, 0.406], [0.229, 0.224, 0.225]),
        ]
    )
    return transform(image).unsqueeze(0)


# 执行分割模型推理并生成掩码
def predict_mask(model: nn.Module, image: Image.Image, device: torch.device) -> np.ndarray:
    width, height = image.size
    tensor = preprocess_image(image).to(device)
    with torch.no_grad():
        _, _, final_prediction = model(tensor)

    final_tensor = final_prediction.data.squeeze(0).cpu()
    mask = np.array(transforms.Resize((height, width))(transforms.ToPILImage()(final_tensor)))
    return np.where(mask > MASK_THRESHOLD, 255, 0).astype(np.uint8)


# 从玻璃掩码中提取白色候选区域
def extract_white_regions(binary_mask: np.ndarray) -> list[np.ndarray]:
    contours, _ = cv2.findContours(binary_mask, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
    return list(contours)


# 按比例扩展候选区域边界框
def expand_bounding_box(x: int, y: int, width: int, height: int, image_shape: tuple[int, int], padding_ratio: float = 0.1) -> tuple[int, int, int, int]:
    image_height, image_width = image_shape
    padding_width = int(width * padding_ratio)
    padding_height = int(height * padding_ratio)
    new_x = max(x - padding_width, 0)
    new_y = max(y - padding_height, 0)
    new_width = min(width + 2 * padding_width, image_width - new_x)
    new_height = min(height + 2 * padding_height, image_height - new_y)
    return new_x, new_y, new_width, new_height


# 从原图中裁剪候选玻璃区域
def segment_from_original_image(original_image: np.ndarray, contours: list[np.ndarray]) -> list[dict[str, Any]]:
    segments = []
    image_height, image_width = original_image.shape[:2]
    for contour in contours:
        x, y, width, height = cv2.boundingRect(contour)
        x, y, width, height = expand_bounding_box(x, y, width, height, (image_height, image_width))
        area = int(width * height)
        if area < MIN_SEGMENT_AREA:
            continue
        crop = original_image[y:y + height, x:x + width]
        segments.append({'image': crop, 'bbox': {'x': x, 'y': y, 'width': width, 'height': height}, 'area_px': area})
    segments.sort(key=lambda item: item['area_px'], reverse=True)
    return segments


# 裁剪玻璃区域边缘以保留主体区域
def crop_glass_region(image: np.ndarray, border_ratio: float = 0.1) -> np.ndarray:
    height, width = image.shape[:2]
    top = int(height * border_ratio)
    bottom = int(height * (1 - border_ratio))
    left = int(width * border_ratio)
    right = int(width * (1 - border_ratio))
    return image[top:bottom, left:right]


# 保存频谱分析可视化图像
def save_frequency_image(path: Path, magnitude_spectrum: np.ndarray) -> None:
    normalized = cv2.normalize(magnitude_spectrum, None, 0, 255, cv2.NORM_MINMAX)
    cv2.imwrite(str(path), normalized.astype(np.uint8))


# 检测单个玻璃区域的平整度指标
def detect_glass_flatness(segment_image: np.ndarray, region_id: int, output_dir: Path) -> dict[str, Any]:
    output_dir.mkdir(parents=True, exist_ok=True)
    cropped_image = crop_glass_region(segment_image, border_ratio=0.1)
    if cropped_image.size == 0:
        cropped_image = segment_image

    gray = cv2.cvtColor(cropped_image, cv2.COLOR_RGB2GRAY)
    edges = cv2.Canny(gray, 50, 150)
    edge_count = int(np.sum(edges > 0))
    laplacian_variance = float(cv2.Laplacian(gray, cv2.CV_64F).var())
    edge_is_flat = laplacian_variance >= 10

    lines = cv2.HoughLinesP(edges, 1, np.pi / 180, 100, minLineLength=50, maxLineGap=10)
    line_image = cropped_image.copy()
    angles = []
    if lines is not None:
        for line in lines:
            x1, y1, x2, y2 = line[0]
            cv2.line(line_image, (x1, y1), (x2, y2), (0, 255, 0), 2)
            angles.append(float(np.arctan2(y2 - y1, x2 - x1) * 180 / np.pi))
    angle_std = float(np.std(angles)) if angles else 0.0
    line_is_flat = angle_std < 50

    grad_x = cv2.Sobel(gray, cv2.CV_64F, 1, 0, ksize=3)
    grad_y = cv2.Sobel(gray, cv2.CV_64F, 0, 1, ksize=3)
    grad_magnitude = cv2.magnitude(grad_x, grad_y)
    gradient_mean = float(np.mean(grad_magnitude))
    gradient_std = float(np.std(grad_magnitude))
    gradient_is_flat = gradient_std < 100

    dft = cv2.dft(np.float32(gray), flags=cv2.DFT_COMPLEX_OUTPUT)
    dft_shift = np.fft.fftshift(dft)
    magnitude = cv2.magnitude(dft_shift[:, :, 0], dft_shift[:, :, 1])
    magnitude_spectrum = 20 * np.log(magnitude + 1)
    freq_max = float(np.max(magnitude_spectrum))
    freq_min = float(np.min(magnitude_spectrum))
    freq_diff = freq_max - freq_min
    frequency_is_flat = freq_diff < 400

    is_flat = sum([line_is_flat, gradient_is_flat, frequency_is_flat]) >= 2
    if not is_flat:
        region_path = output_dir / f'region_{region_id}.png'
        lines_path = output_dir / f'lines_{region_id}.png'
        gradient_path = output_dir / f'gradient_{region_id}.png'
        frequency_path = output_dir / f'frequency_{region_id}.png'
        cv2.imwrite(str(region_path), cv2.cvtColor(segment_image, cv2.COLOR_RGB2BGR))
        cv2.imwrite(str(lines_path), cv2.cvtColor(line_image, cv2.COLOR_RGB2BGR))
        cv2.imwrite(str(gradient_path), np.uint8(np.clip(grad_magnitude, 0, 255)))
        save_frequency_image(frequency_path, magnitude_spectrum)
    return {
        'is_flat': bool(is_flat),
        'edge_uneven_detected': not edge_is_flat,
        'line_uneven_detected': not line_is_flat,
        'gradient_uneven_detected': not gradient_is_flat,
        'frequency_uneven_detected': not frequency_is_flat,
        'edge_count': edge_count,
        'laplacian_variance': round(laplacian_variance, 6),
        'line_count': int(len(lines)) if lines is not None else 0,
        'angle_std': round(angle_std, 6),
        'gradient_mean': round(gradient_mean, 6),
        'gradient_std': round(gradient_std, 6),
        'frequency_min': round(freq_min, 6),
        'frequency_max': round(freq_max, 6),
    }


# 保存玻璃分割掩码图像
def save_mask(binary_mask: np.ndarray, output_dir: Path) -> None:
    mask_path = output_dir / 'mask.png'
    cv2.imwrite(str(mask_path), binary_mask)


# 保存玻璃区域叠加可视化图像
def save_overlay(image: Image.Image, glass_reports: list[dict[str, Any]], output_dir: Path) -> None:
    pil_image = image.convert('RGBA')
    overlay = Image.new('RGBA', pil_image.size, (0, 0, 0, 0))
    draw = ImageDraw.Draw(overlay)
    for glass in glass_reports:
        bbox = glass['bbox']
        color = (0, 255, 0, 110) if glass['is_flat'] else (255, 0, 0, 130)
        x, y, width, height = bbox['x'], bbox['y'], bbox['width'], bbox['height']
        draw.rectangle([x, y, x + width, y + height], fill=color, outline=(0, 0, 0, 180))
    result = Image.alpha_composite(pil_image, overlay).convert('RGB')
    overlay_path = output_dir / 'overlay.png'
    result.save(overlay_path)


# 构建检测结果报告数据
def build_report(glass_reports: list[dict[str, Any]]) -> dict[str, Any]:
    uneven_reports = [item for item in glass_reports if not item['is_flat']]
    if not glass_reports:
        result = 'notglass'
    elif not uneven_reports:
        result = 'flat'
    else:
        result = 'uneven'
    return {
        'result': result,
        'uneven_count': len(uneven_reports),
        'regions': [
            {
                'id': index,
                'bbox_xyxy': item['bbox_xyxy'],
                'edge_uneven_detected': item['edge_uneven_detected'],
                'line_uneven_detected': item['line_uneven_detected'],
                'gradient_uneven_detected': item['gradient_uneven_detected'],
                'frequency_uneven_detected': item['frequency_uneven_detected'],
                'edge_count': item['edge_count'],
                'laplacian_variance': item['laplacian_variance'],
                'line_count': item['line_count'],
                'angle_std': item['angle_std'],
                'gradient_mean': item['gradient_mean'],
                'gradient_std': item['gradient_std'],
                'frequency_min': item['frequency_min'],
                'frequency_max': item['frequency_max'],
            }
            for index, item in enumerate(uneven_reports, start=1)
        ],
    }


class FlatnessDetector:
    def __init__(self) -> None:
        self.device = get_device()
        self.model = load_model(MODEL_PATH, self.device)

    # 执行玻璃平整度检测
    def detect(self, input_path: Path) -> None:
        start_time = time.time()
        input_path = input_path.expanduser().resolve()
        output_dir = input_path.parent
        output_dir.mkdir(parents=True, exist_ok=True)
        image = read_image(input_path)
        binary_mask = predict_mask(self.model, image, self.device)
        original_image = np.array(image)
        contours = extract_white_regions(binary_mask)
        segments = segment_from_original_image(original_image, contours)
        glass_reports = []
        for segment in segments:
            region_id = sum(1 for item in glass_reports if not item['is_flat']) + 1
            analysis = detect_glass_flatness(segment['image'], region_id, output_dir)
            bbox = segment['bbox']
            glass_reports.append(
                {
                    'bbox': bbox,
                    'bbox_xyxy': [bbox['x'], bbox['y'], bbox['x'] + bbox['width'], bbox['y'] + bbox['height']],
                    **analysis,
                }
            )
        report = build_report(glass_reports)
        save_mask(binary_mask, output_dir)
        save_overlay(image, glass_reports, output_dir)
        report_path = output_dir / 'report.json'
        report['runtime_seconds'] = round(time.time() - start_time, 3)
        report_path.write_text(json.dumps(report, ensure_ascii=False, indent=2), encoding='utf-8')


# 获取可复用的玻璃平整度检测器
def get_detector() -> FlatnessDetector:
    global DETECTOR
    if DETECTOR is None:
        DETECTOR = FlatnessDetector()
    return DETECTOR


# 执行玻璃平整度检测
def detect(input_path: Path) -> None:
    get_detector().detect(input_path)
