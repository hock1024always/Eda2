// models/ilp_solver.go
package models

//
//import (
//	"fmt"
//	"math"
//	"time"
//
//	"github.com/crillab/gophersat/solver"
//)
//
//// MLRCSWithIntegerLinearProgramming 使用ILP求解ML-RCS问题
//func MLRCSWithIntegerLinearProgramming(graph []Gate, resources []string, maxResource int) SchedulingResult {
//	startTime := time.Now()
//
//	// 1. 预处理和初始化
//	maxTime := estimateMaxLatency(graph)
//	if maxTime <= 0 {
//		maxTime = 100 // 默认最大值
//	}
//
//	// 2. 创建ILP问题
//	problem := solver.NewProblem()
//
//	// 3. 创建决策变量: x_g,t 表示门g在时间t开始
//	vars := make(map[string]map[int]*solver.Var)
//	for _, gate := range graph {
//		vars[gate.ID] = make(map[int]*solver.Var)
//		// 每个门只能在它的时间帧内调度
//		for t := 0; t <= maxTime; t++ {
//			vars[gate.ID][t] = problem.NewVar()
//		}
//	}
//
//	// 4. 添加约束条件
//
//	// 约束1: 每个门必须被调度到且仅被调度到一个时间步
//	for _, gate := range graph {
//		var sum []*solver.Var
//		for t := 0; t <= maxTime; t++ {
//			sum = append(sum, vars[gate.ID][t])
//		}
//		problem.AddConstraint(problem.Sum(sum...).Eq(1))
//	}
//
//	// 约束2: 依赖关系约束
//	for _, gate := range graph {
//		for _, depID := range gate.Dependencies {
//			dep := findGateByID(graph, depID)
//			if dep == nil {
//				continue
//			}
//
//			// 确保依赖门在依赖它的门之前完成
//			for t1 := 0; t1 <= maxTime; t1++ {
//				for t2 := t1; t2 <= maxTime; t2++ {
//					if t2 >= t1+dep.Duration {
//						continue // 满足依赖关系
//					}
//					// 如果不满足依赖关系，添加约束
//					problem.AddConstraint(
//						problem.Sum(vars[gate.ID][t1], vars[dep.ID][t2]).Le(1))
//				}
//			}
//		}
//	}
//
//	// 约束3: 资源约束
//	for _, res := range resources {
//		for t := 0; t <= maxTime; t++ {
//			var sum []*solver.Var
//			for _, gate := range graph {
//				if gate.Type != res {
//					continue
//				}
//				// 门g在时间t使用时，必须是在[t-g.Duration+1, t]区间内开始的
//				for tStart := max(0, t-gate.Duration+1); tStart <= t; tStart++ {
//					if tStart > maxTime {
//						continue
//					}
//					sum = append(sum, vars[gate.ID][tStart])
//				}
//			}
//			if len(sum) > 0 {
//				problem.AddConstraint(problem.Sum(sum...).Le(maxResource))
//			}
//		}
//	}
//
//	// 5. 设置目标函数: 最小化总延迟
//	// 我们需要找到最后一个完成的门的时间
//	// 为此，我们引入一个辅助变量表示总延迟
//	totalLatencyVar := problem.NewVar()
//	problem.SetObjective(totalLatencyVar, true) // 最小化
//
//	// 添加约束确保totalLatencyVar >= 所有门的完成时间
//	for _, gate := range graph {
//		for t := 0; t <= maxTime; t++ {
//			// totalLatencyVar >= t + gate.Duration 当 x_g,t = 1时
//			// 用大M法表示这个约束
//			M := maxTime + gate.Duration + 1
//			problem.AddConstraint(
//				problem.Sum(totalLatencyVar, vars[gate.ID][t].Mul(-M)).Ge(t + gate.Duration - M))
//		}
//	}
//
//	// 6. 求解ILP问题
//	solution := problem.Solve()
//	if solution == nil {
//		return SchedulingResult{Error: "无可行解"}
//	}
//
//	// 7. 解析解决方案
//	result := SchedulingResult{
//		Schedule:      make(map[string]int),
//		ResourceUsage: make(map[int]map[string]int),
//	}
//
//	// 收集每个门的开始时间
//	for _, gate := range graph {
//		for t := 0; t <= maxTime; t++ {
//			if solution.Value(vars[gate.ID][t]) > 0.5 { // 大于0.5视为1
//				result.Schedule[gate.ID] = t
//				// 更新资源使用情况
//				for tt := t; tt < t+gate.Duration; tt++ {
//					if _, exists := result.ResourceUsage[tt]; !exists {
//						result.ResourceUsage[tt] = make(map[string]int)
//					}
//					result.ResourceUsage[tt][gate.Type]++
//				}
//				break
//			}
//		}
//	}
//
//	// 计算总延迟
//	result.TotalLatency = computeTotalLatency(result.Schedule, graph)
//
//	// 记录求解时间
//	result.SolveTime = time.Since(startTime).Seconds()
//
//	return result
//}
//
//// 辅助函数
//
//func findGateByID(graph []Gate, id string) *Gate {
//	for _, gate := range graph {
//		if gate.ID == id {
//			return &gate
//		}
//	}
//	return nil
//}
//
//func estimateMaxLatency(graph []Gate) int {
//	// 简单估计最大延迟为所有门持续时间的总和
//	maxLatency := 0
//	for _, gate := range graph {
//		maxLatency += gate.Duration
//	}
//	return maxLatency
//}
//
//func computeTotalLatency(schedule map[string]int, graph []Gate) int {
//	maxLatency := 0
//	for _, gate := range graph {
//		if start, ok := schedule[gate.ID]; ok {
//			if start+gate.Duration > maxLatency {
//				maxLatency = start + gate.Duration
//			}
//		}
//	}
//	return maxLatency
//}
//
//func max(a, b int) int {
//	if a > b {
//		return a
//	}
//	return b
//}
