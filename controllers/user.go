package controllers

import (
	"circuit_optimizer/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type DotController struct{}

type VerilogApi struct {
	VerilogCode string `json:"verilog_code"`
}

type DotApi struct {
	DotCode string `json:"dot_code"`
}

// DotOptimizationResult 定义优化结果的结构体
type DotOptimizationResult struct {
	OptimizedDotCode string `json:"optimized_dot_code"`
}

type SchedulingRequest struct {
	Graph       []Gate   `json:"graph"`       // 逻辑门依赖图
	Resources   []string `json:"resources"`   // 可用资源类型
	MaxLatency  int      `json:"maxLatency"`  // 最大延迟约束（用于MR-LCS）
	MaxResource int      `json:"maxResource"` // 最大资源约束（用于ML-RCS）
}

type Gate struct {
	ID           string   `json:"id"`           // 逻辑门ID
	Type         string   `json:"type"`         // 逻辑门类型
	Duration     int      `json:"duration"`     // 执行时间
	Dependencies []string `json:"dependencies"` // 依赖的门
}

type SchedulingResult struct {
	Schedule      map[string]int         `json:"schedule"`      // 门ID到开始时间的映射
	TotalLatency  int                    `json:"totalLatency"`  // 总延迟
	ResourceUsage map[int]map[string]int `json:"resourceUsage"` // 每个时间步的资源使用情况
}

// DotVerilog 实现接口，接收前端传入的Verilog代码并转换为DOT代码
func (d DotController) DotVerilog(c *gin.Context) {
	var verilog VerilogApi
	if err := c.ShouldBindJSON(&verilog); err != nil {
		models.ReturnError(c, 4001, "无效的请求数据")
		return
	}

	// 获取前端传入的Verilog代码
	verilogCode := verilog.VerilogCode

	// 将Verilog代码转换为DOT代码
	dotCode := models.VerilogToDot(verilogCode)

	// 返回DOT代码
	c.JSON(http.StatusOK, DotApi{DotCode: dotCode})
}

func (d DotController) DotNetlist(c *gin.Context) {
	var dot DotApi
	if err := c.ShouldBindJSON(&dot); err != nil {
		models.ReturnError(c, 4001, "无效的请求数据")
		return
	}

	dot_code := dot.DotCode

	// 解析 DOT 代码
	parsedResult, err := models.ParseDotCode(dot_code)
	if err != nil {
		models.ReturnError(c, 4002, err.Error())
		return
	}

	// 返回解析结果
	c.JSON(200, gin.H{
		"message": "成功",
		"data":    parsedResult,
	})
}

// DotNetlist 实现接口，接收前端传入的dot代码并进行优化
func (d DotController) OptimizeDotNetlist(c *gin.Context) {
	var dot DotApi
	if err := c.ShouldBindJSON(&dot); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": 4001, "message": "无效的请求数据"})
		return
	}

	// 获取前端传入的dot代码
	dotCode := dot.DotCode

	// 执行优化操作
	optimizedDotCode := optimizeDotCode(dotCode)

	// 返回优化后的dot代码
	c.JSON(http.StatusOK, DotOptimizationResult{OptimizedDotCode: optimizedDotCode})
}

// optimizeDotCode 对dot代码进行优化
func optimizeDotCode(dotCode string) string {
	// 执行常数传播优化
	dotCode = models.ConstantPropagation(dotCode)

	// 执行共享子表达式消除优化
	dotCode = models.CommonSubexpressionElimination(dotCode)

	return dotCode
}
