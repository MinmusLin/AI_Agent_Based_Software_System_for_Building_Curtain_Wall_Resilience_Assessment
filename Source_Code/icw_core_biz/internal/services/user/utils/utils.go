package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"strconv"
	"strings"

	"icw_core_biz/internal/services/auth/utils"
)

// NormalizeEmailHash 对标准化邮箱地址做 SHA-256 哈希
func NormalizeEmailHash(email string) (string, error) {
	normalizedEmail, err := utils.NormalizeEmailAddress(email)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256([]byte(normalizedEmail))
	return hex.EncodeToString(sum[:]), nil
}

// BuildDefaultAvatarSVG 生成 GitHub 样式 5×5 镜像默认 SVG 头像
func BuildDefaultAvatarSVG(emailHash string) []byte {
	const (
		avatarSize      = 64
		avatarGridSize  = 5
		avatarMarginPct = 0.08
	)

	foreground := avatarColor(emailHash)
	baseMargin := int(math.Floor(float64(avatarSize) * avatarMarginPct))
	cellSize := int(math.Floor(float64(avatarSize-baseMargin*2) / avatarGridSize))
	margin := int(math.Floor(float64(avatarSize-cellSize*avatarGridSize) / 2))
	strokeWidth := float64(avatarSize) * 0.005

	var blocks strings.Builder

	for i := 0; i < 15 && i < len(emailHash); i++ {
		if hexDigitValue(emailHash[i])%2 != 0 {
			continue
		}
		if i < 5 {
			writeCell(&blocks, 2, i, cellSize, margin)
		} else if i < 10 {
			row := i - 5
			writeCell(&blocks, 1, row, cellSize, margin)
			writeCell(&blocks, 3, row, cellSize, margin)
		} else {
			row := i - 10
			writeCell(&blocks, 0, row, cellSize, margin)
			writeCell(&blocks, 4, row, cellSize, margin)
		}
	}

	return []byte(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d"><rect width="%d" height="%d" fill="#F0F0F0"/><g fill="%s" stroke="%s" stroke-width="%.2f">%s</g></svg>`,
		avatarSize,
		avatarSize,
		avatarSize,
		avatarSize,
		avatarSize,
		avatarSize,
		foreground,
		foreground,
		strokeWidth,
		blocks.String(),
	))
}

// avatarColor 使用标准化邮箱地址哈希末 7 位生成 GitHub 样式前景色
func avatarColor(emailHash string) string {
	if len(emailHash) < 7 {
		return "#26734D"
	}
	hueRaw, err := strconv.ParseUint(emailHash[len(emailHash)-7:], 16, 64)
	if err != nil {
		return "#26734D"
	}
	hue := float64(hueRaw) / 0xFFFFFFF
	red, green, blue := hslToRGB(hue, 0.7, 0.5)
	return fmt.Sprintf("#%02X%02X%02X", red, green, blue)
}

// hslToRGB 将 HSL 颜色格式转换为 RGB 颜色格式
func hslToRGB(hue, saturation, lightness float64) (int, int, int) {
	if saturation == 0 {
		value := colorComponent(lightness)
		return value, value, value
	}

	var q float64
	if lightness < 0.5 {
		q = lightness * (1 + saturation)
	} else {
		q = lightness + saturation - lightness*saturation
	}
	p := 2*lightness - q

	red := huwToRGB(p, q, hue+1.0/3.0)
	green := huwToRGB(p, q, hue)
	blue := huwToRGB(p, q, hue-1.0/3.0)

	return colorComponent(red), colorComponent(green), colorComponent(blue)
}

// huwToRGB 计算单 RGB 通道值
func huwToRGB(p, q, hue float64) float64 {
	if hue < 0 {
		hue++
	}
	if hue > 1 {
		hue--
	}
	switch {
	case hue < 1.0/6.0:
		return p + (q-p)*6*hue
	case hue < 1.0/2.0:
		return q
	case hue < 2.0/3.0:
		return p + (q-p)*(2.0/3.0-hue)*6
	default:
		return p
	}
}

// colorComponent 将 0-1 浮点颜色值转为 0-255 整数颜色值
func colorComponent(value float64) int {
	return int(math.Round(value * 255))
}

// hexDigitValue 返回单个十六进制字符的数值
func hexDigitValue(char byte) int {
	switch {
	case char >= '0' && char <= '9':
		return int(char - '0')
	case char >= 'a' && char <= 'f':
		return int(char-'a') + 10
	case char >= 'A' && char <= 'F':
		return int(char-'A') + 10
	default:
		return 1
	}
}

// writeCell 写入 SVG 色块
func writeCell(blocks *strings.Builder, col, row, cellSize, margin int) {
	x := col*cellSize + margin
	y := row*cellSize + margin
	blocks.WriteString(fmt.Sprintf(`<rect x="%d" y="%d" width="%d" height="%d"/>`, x, y, cellSize, cellSize))
}
