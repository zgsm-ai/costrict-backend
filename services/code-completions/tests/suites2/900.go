package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type DataAnalyzer struct {
	data        [][]string
	headers     []string
	numRows     int
	numCols     int
	numericCols map[int]bool
}

func NewDataAnalyzer() *DataAnalyzer {
	return &DataAnalyzer{
		data:        make([][]string, 0),
		headers:     make([]string, 0),
		numRows:     0,
		numCols:     0,
		numericCols: make(map[int]bool),
	}
}

func (da *DataAnalyzer) LoadFromCSV(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("无法打开文件: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("读取CSV文件失败: %v", err)
	}

	if len(records) == 0 {
		return fmt.Errorf("CSV文件为空")
	}

	da.headers = records[0]
	da.numCols = len(da.headers)
	da.data = records[1:]
	da.numRows = len(da.data)
	da.detectNumericColumns()

	return nil
}

func (da *DataAnalyzer) LoadFromMemory(headers []string, data [][]string) {
	da.headers = headers
	da.data = data
	da.numRows = len(data)
	da.numCols = len(headers)
	da.detectNumericColumns()
}

func (da *DataAnalyzer) detectNumericColumns() {
	for col := 0; col < da.numCols; col++ {
		isNumeric := true
		for row := 0; row < da.numRows; row++ {
			if col >= len(da.data[row]) {
				isNumeric = false
				break
			}
			if _, err := strconv.ParseFloat(da.data[row][col], 64); err != nil {
				isNumeric = false
				break
			}
		}
		da.numericCols[col] = isNumeric
	}
}

func (da *DataAnalyzer) GetColumnIndex(name string) (int, error) {
	for i, header := range da.headers {
		if header == name {
			return i, nil
		}
	}
	return -1, fmt.Errorf("列名 '%s' 不存在", name)
}

func (da *DataAnalyzer) GetColumn(name string) ([]string, error) {
	idx, err := da.GetColumnIndex(name)
	if err != nil {
		return nil, err
	}

	column := make([]string, da.numRows)
	for i := 0; i < da.numRows; i++ {
		if idx < len(da.data[i]) {
			column[i] = da.data[i][idx]
		} else {
			column[i] = ""
		}
	}

	return column, nil
}

func (da *DataAnalyzer) GetNumericColumn(name string) ([]float64, error) {
	column, err := da.GetColumn(name)
	if err != nil {
		return nil, err
	}

	idx, _ := da.GetColumnIndex(name)
	if !da.numericCols[idx] {
		return nil, fmt.Errorf("列 '%s' 不是数值列", name)
	}

	numericColumn := make([]float64, da.numRows)
	for i, val := range column {
		num, err := strconv.ParseFloat(val, 64)
		if err != nil {
			numericColumn[i] = 0
		} else {
			numericColumn[i] = num
		}
	}

	return numericColumn, nil
}

func (da *DataAnalyzer) Mean(name string) (float64, error) {
	column, err := da.GetNumericColumn(name)
	if err != nil {
		return 0, err
	}

	sum := 0.0
	for _, val := range column {
		sum += val
	}

	return sum / float64(len(column)), nil
}

func (da *DataAnalyzer) Median(name string) (float64, error) {
	column, err := da.GetNumericColumn(name)
	if err != nil {
		return 0, err
	}

	sorted := make([]float64, len(column))
	copy(sorted, column)
	
	// 简单排序
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i] > sorted[j] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	n := len(sorted)
	if n%2 == 1 {
		return sorted[n/2], nil
	}
	return (sorted[n/2-1] + sorted[n/2]) / 2.0, nil
}

func (da *DataAnalyzer) StdDev(name string) (float64, error) {
	mean, err := da.Mean(name)
	if err != nil {
		return 0, err
	}

	column, err := da.GetNumericColumn(name)
	if err != nil {
		return 0, err
	}

	variance := 0.0
	for _, val := range column {
		variance += (val - mean) * (val - mean)
	}
	variance /= float64(len(column))

	return variance, nil
}

func (da *DataAnalyzer) Min(name string) (float64, error) {
	column, err := da.GetNumericColumn(name)
	if err != nil {
		return 0, err
	}

	min := column[0]
	for _, val := range column {
		if val < min {
			min = val
		}
	}

	return min, nil
}

func (da *DataAnalyzer) Max(name string) (float64, error) {
	column, err := da.GetNumericColumn(name)
	if err != nil {
		return 0, err
	}

	max := column[0]
	for _, val := range column {
		if val > max {
			max = val
		}
	}

	return max, nil
}

func (da *DataAnalyzer) Filter(columnName string, value string, operator string) (*DataAnalyzer, error) {
	idx, err := da.GetColumnIndex(columnName)
	if err != nil {
		return nil, err
	}

	filteredData := make([][]string, 0)
	for _, row := range da.data {
		if idx < len(row) {
			match := false
			switch operator {
			case "==":
				match = row[idx] == value
			case "!=":
				match = row[idx] != value
			case "contains":
				match = strings.Contains(row[idx], value)
			}

			if match {
				filteredData = append(filteredData, row)
			}
		}
	}

	filtered := NewDataAnalyzer()
	filtered.LoadFromMemory(da.headers, filteredData)
	return filtered, nil
}

func (da *DataAnalyzer) GroupBy(groupBy string, aggCols map[string]string) (map[string]map[string]float64, error) {
	groupIdx, err := da.GetColumnIndex(groupBy)
	if err != nil {
		return nil, err
	}

	result := make(map[string]map[string]float64)
	groups := make(map[string][][]float64)

	// 按组收集数据
	for _, row := range da.data {
		if groupIdx < len(row) {
			group := row[groupIdx]
			if groups[group] == nil {
				groups[group] = make([][]float64, 0)
			}

			groupRow := make([]float64, 0)
			for col := range aggCols {
				idx, _ := da.GetColumnIndex(col)
				if idx < len(row) {
					val, _ := strconv.ParseFloat(row[idx], 64)
					groupRow = append(groupRow, val)
				}
			}

			groups[group] = append(groups[group], groupRow)
		}
	}

	// 计算每个组的聚合值
	for group, data := range groups {
		result[group] = make(map[string]float64)

		i := 0
		for col, aggFunc := range aggCols {
			values := make([]float64, 0)
			for _, row := range data {
				if i < len(row) {
					values = append(values, row[i])
				}
			}

			var aggValue float64
			switch aggFunc {
			case "sum":
				sum := 0.0
				for _, val := range values {
					sum += val
				}
				aggValue = sum
			case "mean":
				sum := 0.0
				for _, val := range values {
					sum += val
				}
				aggValue = sum / float64(len(values))
			case "min":
				if len(values) > 0 {
					min := values[0]
					for _, val := range values {
						if val < min {
							min = val
						}
					}
					aggValue = min
				}
			case "max":
				if len(values) > 0 {
					max := values[0]
					for _, val := range values {
						if val > max {
							max = val
						}
					}
					aggValue = max
				}
			case "count":
				aggValue = float64(len(values))
			}

			result[group][col] = aggValue
			i++
		}
	}

	return result, nil
}

func (da *DataAnalyzer) NumRows() int {
	return da.numRows
}

func (da *DataAnalyzer) NumCols() int {
	return da.numCols
}

func (da *DataAnalyzer) Headers() []string {
	return da.headers
}

func (da *DataAnalyzer) Describe() map[string]map[string]float64 {
	result := make(map[string]map[string]float64)

	for i, header := range da.headers {
		if da.numericCols[i] {
			mean, _ := da.Mean(header)
			median, _ := da.Median(header)
			stdDev, _ := da.StdDev(header)
			min, _ := da.Min(header)
			max, _ := da.Max(header)

			result[header] = map[string]float64{
				"count":   float64(da.numRows),
				"mean":    mean,
				"median":  median,
				"std_dev": stdDev,
				"min":     min,
				"max":     max,
			}
		}
	}

	return result
}

func main() {
	analyzer := NewDataAnalyzer()
	
	// 创建示例数据
	headers := []string{"Name", "Age", "Salary", "Department"}
	data := [][]string{
		{"Alice", "28", "50000", "Engineering"},
		{"Bob", "32", "60000", "Marketing"},
		{"Charlie", "25", "45000", "Engineering"},
		{"David", "35", "75000", "Management"},
		{"Eve", "30", "55000", "Marketing"},
	}
	
	analyzer.LoadFromMemory(headers, data)
	
	fmt.Printf("数据集包含 %d 行，%d 列\n", analyzer.NumRows(), analyzer.NumCols())
	fmt.Printf("列名: %v\n", analyzer.Headers())
	
	// 计算描述性统计
	description := analyzer.Describe()
	fmt.Println("\n描述性统计:")
	for col, stats := range description {
		fmt.Printf("\n%s:\n", col)
		for stat, value := range stats {
			fmt.Printf("  %s: %.2f\n", stat, value)
		}
	}
	
	// 按部门分组并计算平均薪资
	aggCols := map[string]string{"Age": "mean", "Salary": "mean"}
	grouped, _ := analyzer.GroupBy("Department", aggCols)
	
	fmt.Println("\n按部门分组统计:")
	for dept, stats := range grouped {
		fmt.Printf("\n%s:\n", dept)
		for stat, value := range stats {
			fmt.Printf("  %s: %.2f\n", stat, value)
		}
	}
	
	// 过滤数据
	<｜fim▁hole｜>
	
	fmt.Println("\n数据分析演示完成")
}
