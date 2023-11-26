# param
param 的修改和困惑源自于 `node 要支持 middlewares 判断, 且 middlewares 可以使用 dag 内产生的 field, 并构建依赖顺序. 但 middleware 失败不一定需要剪枝`.

有两种方案
1. 方案一: param{id,kv,field_name,onerror}, 缺点在于会有重复的 param.kv
2. 方案二: param{id,kv,field_name,callers[]}, 优点是没有重复了, 但是 callers 会有不同的值. 后续如果要实现 Default, 会很难.

可以确定的是, node 对外展示和接收的 params 格式是固定的. 对外展示 []Params, 表示需要外界传入的参数合集, 然后接受这一系列参数的 `[]interface{}`.

目前的处理是方案一, 
1. Param 添加 OnError 处理, 分别支持 prune(supplier 默认)/skip(middleware 默认).
    1. node.runtime 参数判断失败时, 自行判断处理.
2. Param 修改 ID(添加 caller_id), 添加字段 field_name, 作为变量时的 field_name(拆分出来了 id 的功能)

1. Param.OnError == prune, 则对下游进行全部剪枝, 默认处理.
2. Param.OnError == default, 则当作正常参数处理. (目前不加这个, 预留功能)
3. Param.OnError == skip, 如果是 middleware, 则跳过该 middleware, 如果是函数, 则执行 field.OnError
