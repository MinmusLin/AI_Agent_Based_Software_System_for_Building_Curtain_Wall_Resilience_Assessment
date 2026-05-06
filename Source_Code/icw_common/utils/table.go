package utils

import (
	"strings"
	"unicode"
)

// TableColumn 终端表格列数据
type TableColumn struct {
	Header string
	Values []string
}

// FormatTable 将表格列数据格式化为终端表格
func FormatTable(columns []*TableColumn) string {
	if len(columns) == 0 {
		return ""
	}
	widths := make([]int, len(columns))
	rowCount := 0
	for columnIndex, column := range columns {
		widths[columnIndex] = displayWidth(column.Header)
		rowCount = max(rowCount, len(column.Values))
		for _, value := range column.Values {
			widths[columnIndex] = max(widths[columnIndex], displayWidth(value))
		}
	}
	border := formatTableBorder(widths)
	var builder strings.Builder
	builder.WriteString(border)
	builder.WriteString("\n")
	builder.WriteString(formatTableRow(formatTableHeaderRow(columns), widths))
	builder.WriteString("\n")
	builder.WriteString(border)
	builder.WriteString("\n")
	for rowIndex := 0; rowIndex < rowCount; rowIndex++ {
		builder.WriteString(formatTableRow(formatTableDataRow(columns, rowIndex), widths))
		builder.WriteString("\n")
	}
	builder.WriteString(border)
	return builder.String()
}

// formatTableRow 格式化表格行数据
func formatTableRow(values []string, widths []int) string {
	var builder strings.Builder
	builder.WriteString("|")
	for index, width := range widths {
		value := ""
		if index < len(values) {
			value = values[index]
		}
		builder.WriteString(" ")
		builder.WriteString(padRight(value, width))
		builder.WriteString(" |")
	}
	return builder.String()
}

// formatTableHeaderRow 格式化标题行数据
func formatTableHeaderRow(columns []*TableColumn) []string {
	headers := make([]string, 0, len(columns))
	for _, column := range columns {
		headers = append(headers, column.Header)
	}
	return headers
}

// formatTableBorder 格式化表格边框线
func formatTableBorder(widths []int) string {
	var builder strings.Builder
	builder.WriteString("+")
	for _, width := range widths {
		builder.WriteString("-")
		builder.WriteString(strings.Repeat("-", width))
		builder.WriteString("-+")
	}
	return builder.String()
}

// formatTableDataRow 格式化数据行数据
func formatTableDataRow(columns []*TableColumn, rowIndex int) []string {
	values := make([]string, 0, len(columns))
	for _, column := range columns {
		if rowIndex >= len(column.Values) {
			values = append(values, "")
			continue
		}
		values = append(values, column.Values[rowIndex])
	}
	return values
}

// padRight 在字符串右侧填充空格至指定宽度
func padRight(value string, width int) string {
	padding := width - displayWidth(value)
	if padding <= 0 {
		return value
	}
	return value + strings.Repeat(" ", padding)
}

// displayWidth 计算字符串终端显示宽度
func displayWidth(value string) int {
	width := 0
	for _, r := range value {
		if unicode.Is(unicode.Han, r) || unicode.Is(unicode.Hiragana, r) ||
			unicode.Is(unicode.Katakana, r) || unicode.Is(unicode.Hangul, r) ||
			(r >= 0xFF01 && r <= 0xFF60) || (r >= 0xFFE0 && r <= 0xFFE6) {
			width += 2
			continue
		}
		width++
	}
	return width
}
