package router

import (
	"circuit_optimizer/controllers"
	"circuit_optimizer/pkg/logger"
	"github.com/gin-gonic/gin"
)

// 路由 函数的名字要大写，这样才可以被其他包访问！
func Router() *gin.Engine {
	//创建一个路由的实例
	r := gin.Default()

	//日志中间件
	r.Use(gin.LoggerWithConfig(logger.LoggerToFile()))
	r.Use(logger.Recover)
	// 设置最大上传文件大小为10MB
	r.MaxMultipartMemory = 10 << 20

	dot := r.Group("/dot")
	{
		// 输入verilog代码，返回dot代码
		dot.POST("/verilog", controllers.DotController{}.DotVerilog)
		// 输入dot代码，返回产生的
		dot.POST("/netlist", controllers.DotController{}.DotNetlist)
		// 输入dot代码，返回优化后的dot代码
		dot.POST("/netlist/optimize", controllers.DotController{}.OptimizeDotNetlist)
		//设计和实现ML – RCS和MR – LCS算法，完成对逻辑门的调度。

		//使用ILP(Integer Linear Programming)求解ML – RCS调度问题

	}

	return r
}
