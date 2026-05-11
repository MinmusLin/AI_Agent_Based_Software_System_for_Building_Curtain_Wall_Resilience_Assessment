package utils

import (
	"encoding/json"

	"icw_common/utils"
)

// NormalizeSummaryJSON 将智能体输出转换为单行压缩 JSON 字符串
func NormalizeSummaryJSON(output string) (string, error) {
	if compacted, err := utils.CompactJSONObjectString(output); err == nil {
		return compacted, nil
	}
	unquoted := ""
	if err := json.Unmarshal([]byte(output), &unquoted); err != nil {
		return "", err
	}
	return utils.CompactJSONObjectString(unquoted)
}
