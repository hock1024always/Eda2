package models

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

type JsonErrStruct struct {
	Code int         `json:"code"`
	Msg  interface{} `json:"msg"`
}

func ReturnError(c *gin.Context, code int, msg string) {
	//code：响应码，msg：错误信息
	json := &JsonErrStruct{Code: code, Msg: msg}
	c.JSON(http.StatusOK, json)
}

// verilogToDot 将Verilog代码转换为DOT语言的图形表示
func VerilogToDot(verilogCode string) string {
	// 初始化DOT代码
	dotCode := "digraph VerilogDiagram {\n"
	dotCode += "  rankdir=LR;\n"         // 从左到右布局
	dotCode += "  node [shape=box];\n"   // 模块用方框表示
	dotCode += "  edge [fontsize=10];\n" // 边上的字体大小

	// 解析Verilog代码
	lines := strings.Split(verilogCode, "\n")
	modules := make(map[string]string) // 模块名到节点ID的映射
	wires := make(map[string]string)   // 线网名到节点ID的映射
	assigns := make([]string, 0)       // 存储所有assign语句
	nextNodeID := 1

	// 第一次遍历：提取模块、端口和assign语句
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module") {
			// 提取模块名称
			parts := strings.Fields(line)
			if len(parts) < 2 {
				continue
			}
			moduleName := strings.TrimSuffix(parts[1], ";")
			modules[moduleName] = strconv.Itoa(nextNodeID)
			nextNodeID++
		} else if strings.HasPrefix(line, "assign") {
			// 存储assign语句供后续处理
			assigns = append(assigns, line)
		}
	}

	// 第二次遍历：处理assign语句，提取线网连接
	for _, assign := range assigns {
		parts := strings.SplitN(assign, "=", 2)
		if len(parts) != 2 {
			continue
		}
		lhs := strings.TrimSpace(strings.TrimPrefix(parts[0], "assign"))
		rhs := strings.TrimSpace(strings.TrimSuffix(parts[1], ";"))

		// 为每个线网创建节点（如果还不存在）
		if _, exists := wires[lhs]; !exists {
			wires[lhs] = strconv.Itoa(nextNodeID)
			nextNodeID++
		}
		if _, exists := wires[rhs]; !exists {
			wires[rhs] = strconv.Itoa(nextNodeID)
			nextNodeID++
		}
	}

	// 添加模块节点
	for moduleName, nodeID := range modules {
		dotCode += fmt.Sprintf("  %s [label=\"%s\"];\n", nodeID, moduleName)
	}

	// 添加线网节点（用椭圆形表示）
	for wireName, nodeID := range wires {
		dotCode += fmt.Sprintf("  %s [label=\"%s\", shape=ellipse];\n", nodeID, wireName)
	}

	// 添加连接关系
	for _, assign := range assigns {
		parts := strings.SplitN(assign, "=", 2)
		if len(parts) != 2 {
			continue
		}
		lhs := strings.TrimSpace(strings.TrimPrefix(parts[0], "assign"))
		rhs := strings.TrimSpace(strings.TrimSuffix(parts[1], ";"))

		// 查找源和目标节点
		var srcNode, dstNode string
		if nodeID, ok := modules[lhs]; ok {
			srcNode = nodeID
		} else if nodeID, ok := wires[lhs]; ok {
			srcNode = nodeID
		}

		if nodeID, ok := modules[rhs]; ok {
			dstNode = nodeID
		} else if nodeID, ok := wires[rhs]; ok {
			dstNode = nodeID
		}

		if srcNode != "" && dstNode != "" {
			dotCode += fmt.Sprintf("  %s -> %s;\n", srcNode, dstNode)
		}
	}

	// 结束DOT代码
	dotCode += "}"

	return dotCode
}

// parseDotCode 解析 DOT 代码，提取节点和边信息
func ParseDotCode(dotCode string) (string, error) {
	// 使用正则表达式提取节点和边信息
	nodeRegex := regexp.MustCompile(`\b(\w+)\b`)
	edgeRegex := regexp.MustCompile(`\b(\w+)\s*--\s*(\w+)\b`)

	// 提取所有节点
	nodes := nodeRegex.FindAllString(dotCode, -1)
	uniqueNodes := uniqueStrings(nodes)

	// 提取所有边
	edges := edgeRegex.FindAllStringSubmatch(dotCode, -1)

	// 构建电路网表图内容
	var netlist strings.Builder
	netlist.WriteString("电路网表图:\n")
	netlist.WriteString("节点:\n")
	for _, node := range uniqueNodes {
		netlist.WriteString(node + "\n")
	}
	netlist.WriteString("边:\n")
	for _, edge := range edges {
		netlist.WriteString(edge[1] + " -- " + edge[2] + "\n")
	}

	return netlist.String(), nil
}

// uniqueStrings 去除重复的字符串
func uniqueStrings(strs []string) []string {
	unique := make(map[string]bool)
	var result []string
	for _, str := range strs {
		if _, exists := unique[str]; !exists {
			unique[str] = true
			result = append(result, str)
		}
	}
	return result
}

//// constantPropagation 实现常数传播优化
//func ConstantPropagation(dotCode string) string {
//	lines := strings.Split(dotCode, ";")
//	constants := make(map[string]string)
//	var result []string
//
//	for _, line := range lines {
//		line = strings.TrimSpace(line)
//		if line == "" {
//			continue
//		}
//
//		if strings.Contains(line, "=") {
//			parts := strings.SplitN(line, "=", 2)
//			varName := strings.TrimSpace(parts[0])
//			value := strings.TrimSpace(parts[1])
//
//			// 替换已知常数
//			for k, v := range constants {
//				value = strings.ReplaceAll(value, k, v)
//			}
//
//			// 检测是否为常数赋值
//			if isConstant(value) {
//				constants[varName] = value
//			}
//
//			result = append(result, varName+" = "+value)
//		} else {
//			result = append(result, line)
//		}
//	}
//
//	return strings.Join(result, "; ")
//}
//
//// commonSubexpressionElimination 实现共享子表达式消除优化
//func CommonSubexpressionElimination(dotCode string) string {
//	lines := strings.Split(dotCode, ";")
//	subexpressions := make(map[string]string)
//	var newLines []string
//	nextTempVar := 1
//
//	// 第一次遍历：收集所有表达式
//	for _, line := range lines {
//		line = strings.TrimSpace(line)
//		if line == "" {
//			continue
//		}
//
//		if strings.Contains(line, "=") {
//			parts := strings.SplitN(line, "=", 2)
//			varName := strings.TrimSpace(parts[0])
//			expression := strings.TrimSpace(parts[1])
//
//			// 跳过已经生成的临时变量
//			if strings.HasPrefix(varName, "TEMP") {
//				newLines = append(newLines, line)
//				continue
//			}
//
//			// 检查是否已经处理过这个表达式
//			if tempVar, exists := subexpressions[expression]; exists {
//				newLines = append(newLines, varName+" = "+tempVar)
//			} else {
//				// 创建新的临时变量
//				tempVar := "TEMP" + strconv.Itoa(nextTempVar)
//				nextTempVar++
//				subexpressions[expression] = tempVar
//				newLines = append(newLines, tempVar+" = "+expression)
//				newLines = append(newLines, varName+" = "+tempVar)
//			}
//		} else {
//			newLines = append(newLines, line)
//		}
//	}
//
//	return strings.Join(newLines, "; ")
//}
//
//// isConstant 检测字符串是否为常数
//func isConstant(s string) bool {
//	for _, r := range s {
//		if !unicode.IsDigit(r) && r != '.' && r != '-' {
//			return false
//		}
//	}
//	return true
//}

// constantPropagation 实现常数传播优化
func ConstantPropagation(dotCode string) string {
	lines := strings.Split(dotCode, ";")
	constants := make(map[string]string)
	var result []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			varName := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			// 替换已知常数
			for k, v := range constants {
				value = strings.ReplaceAll(value, k, v)
			}

			// 检测是否为常数赋值
			if isConstant(value) {
				constants[varName] = value
			}

			// Escape quotes in the output
			varName = escapeQuotes(varName)
			value = escapeQuotes(value)
			result = append(result, varName+" = "+value)
		} else {
			// Escape quotes in non-assignment lines too
			line = escapeQuotes(line)
			result = append(result, line)
		}
	}

	return strings.Join(result, "; ")
}

// commonSubexpressionElimination 实现共享子表达式消除优化
func CommonSubexpressionElimination(dotCode string) string {
	lines := strings.Split(dotCode, ";")
	subexpressions := make(map[string]string)
	var newLines []string
	nextTempVar := 1

	// 第一次遍历：收集所有表达式
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			varName := strings.TrimSpace(parts[0])
			expression := strings.TrimSpace(parts[1])

			// 跳过已经生成的临时变量
			if strings.HasPrefix(varName, "TEMP") {
				newLines = append(newLines, escapeQuotes(line))
				continue
			}

			// 检查是否已经处理过这个表达式
			if tempVar, exists := subexpressions[expression]; exists {
				newLines = append(newLines, escapeQuotes(varName)+" = "+escapeQuotes(tempVar))
			} else {
				// 创建新的临时变量
				tempVar := "TEMP" + strconv.Itoa(nextTempVar)
				nextTempVar++
				subexpressions[expression] = tempVar
				newLines = append(newLines, escapeQuotes(tempVar)+" = "+escapeQuotes(expression))
				newLines = append(newLines, escapeQuotes(varName)+" = "+escapeQuotes(tempVar))
			}
		} else {
			newLines = append(newLines, escapeQuotes(line))
		}
	}

	return strings.Join(newLines, "; ")
}

// isConstant 检测字符串是否为常数
func isConstant(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) && r != '.' && r != '-' {
			return false
		}
	}
	return true
}

// escapeQuotes 在引号前添加转义字符
func escapeQuotes(s string) string {
	return strings.ReplaceAll(s, `"`, `\"`)
}
