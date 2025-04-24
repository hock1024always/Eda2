package controllers

//
//import (
//	"github.com/gin-gonic/gin"
//	"net/http"
//	"your-app/models"
//)
//
//type SchedulingController struct{}
//
//// MLRCS 实现最大延迟-资源约束调度
//func (sc *SchedulingController) MLRCS(c *gin.Context) {
//	var req models.SchedulingRequest
//	if err := c.ShouldBindJSON(&req); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
//		return
//	}
//
//	result := models.MaxLatencyResourceConstrainedScheduling(req.Graph, req.Resources, req.MaxResource)
//	c.JSON(http.StatusOK, result)
//}
//
//// MRLCS 实现最小资源-延迟约束调度
//func (sc *SchedulingController) MRLCS(c *gin.Context) {
//	var req models.SchedulingRequest
//	if err := c.ShouldBindJSON(&req); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
//		return
//	}
//
//	result := models.MinResourceLatencyConstrainedScheduling(req.Graph, req.Resources, req.MaxLatency)
//	c.JSON(http.StatusOK, result)
//}
//
//// MLRCSWithILP 使用整数线性规划求解ML-RCS问题
//func (sc *SchedulingController) MLRCSWithILP(c *gin.Context) {
//	var req models.SchedulingRequest
//	if err := c.ShouldBindJSON(&req); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
//		return
//	}
//
//	result := models.MLRCSWithIntegerLinearProgramming(req.Graph, req.Resources, req.MaxResource)
//	c.JSON(http.StatusOK, result)
//}
