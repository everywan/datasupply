# datasupply

datasupply 是反作弊的补数 SDK, 由反作弊补数模块重构而来, 目标是提供一个高性能, 节省资源, 通用且易懂的补数 SDK. 所谓节省资源, 是指尽可能减少无效的函数调用.

## Install
```Bash
go get git.in.zhihu.com/antispam/datasupply
```

## 实现
datasupply 的核心思想是基于 DAG 构建字段依赖顺序, 并按依赖顺序进行节点调度.

datasupply 主要分为 **构建和执行** 两部分, 具体的流程和功能如下
````
-> build
    -> build field node: 每个字段构建为一个节点.
    -> merge field node: 合并相同函数调用的字段为一个节点.
    -> analysis dep: 分析依赖关系, 构建依赖图
    -> optomize:
        -> 移除孤儿节点.
        -> 节点优先级调整.
    -> return dag.root

-> run
    -> dag 中间件执行
    -> 构建 root.runtime
    -> 执行 runtime
        -> 节点中间件执行 -> 节点参数校验 -> 节点执行 -> 返回值校验 -> 节点中间件执行
    -> 下游节点状态检测
    -> dag 阶段检测
    -> 结果保存
    -> dag 中间件执行
````

其他模块
1. supplier. 补数只是根据配置生成 dag, 需要有真正的函数去执行才能获取字段. supplier 就负责提供真正的函数执行.

## Example
See [example](./datasupply_test.go)
