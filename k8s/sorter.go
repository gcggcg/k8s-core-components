package k8s

/**
 *    Description: 自定义排序器
 *    Date: 2024/12/23
 */

type MemAscSorter []*StatInfo

func (p MemAscSorter) Len() int           { return len(p) }
func (p MemAscSorter) Less(i, j int) bool { return p[i].MemLoad.Used < p[j].MemLoad.Used }
func (p MemAscSorter) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type MemDescSorter []*StatInfo

func (p MemDescSorter) Len() int           { return len(p) }
func (p MemDescSorter) Less(i, j int) bool { return p[i].MemLoad.Used > p[j].MemLoad.Used }
func (p MemDescSorter) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type CpuAscSorter []*StatInfo

func (p CpuAscSorter) Len() int           { return len(p) }
func (p CpuAscSorter) Less(i, j int) bool { return p[i].CpuLoad.Ratio < p[j].CpuLoad.Ratio }
func (p CpuAscSorter) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type CpuDescSorter []*StatInfo

func (p CpuDescSorter) Len() int           { return len(p) }
func (p CpuDescSorter) Less(i, j int) bool { return p[i].CpuLoad.Ratio > p[j].CpuLoad.Ratio }
func (p CpuDescSorter) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
