/* package supplier 为 node/datasupply 提供获取一系列外部数据的能力.
将外部能力抽象为一个 Plugin, suppleir 通过保存一系列 Plugin 来实现不同的能力.

Supplier 负责管理和执行 Plugins, 如注册, 获取, 在 Plugin 真正运行前执行自定义处理,
如 singleflight, middlerwares 等. Plugin 负责一次函数调用的处理, 包括 ttl, cache,
CircuitBreaker 等.

supplier 是对 antispam.plugin-framework 的封装, 修改了如下功能
1. 插件注册/管理: 面向函数=>面向对象. 原因如下
	1. 接口易于 mock, 方便测试.
	2. 接口易于替换实现. 项目中有多种实现, 直接使用函数式较难满足.
2. Call 函数签名: Call(ctx, request, response)=>Call(ctx, params), 目的是参数中
	只传递函数必要的信息, 原环境信息如 traceid 放到 ctx 中传递.
	req/resp 加一层封装或许有利于之后的扩展, 但目前没想到有需要扩展的场景.
*/
package supplier

import (
	"context"
	"sync"

	"git.in.zhihu.com/antispam/datasupply/constant"
)

//go:generate mockgen -package mock -destination ./mock/supplier.go -source=supplier.go
type ISupplier interface {
	GetName() string
	GetPlugin(pluginName string) (plugin IPlugin, isExist bool)
	GetAllPlugin() map[string]IPlugin
	RegisterPlugin(IPlugin)

	// Supply 是 Plugin.Func 的封装, 提供快速入口, 返回值封装, 以及配置处理.
	Supply(ctx context.Context, pluginName string, params []interface{}) (map[string]interface{}, error)
}

type DefaultSupplier struct {
	name      string
	pluginMap map[string]IPlugin
	locker    sync.Locker
}

var _ ISupplier = new(DefaultSupplier)

func NewDefaultSupplier(name string, plugins []IPlugin, options ...Option) *DefaultSupplier {
	pluginMap := make(map[string]IPlugin)
	for _, plugin := range plugins {
		pluginMap[plugin.GetName()] = plugin
	}
	return &DefaultSupplier{
		name:      name,
		pluginMap: pluginMap,
		locker:    &sync.Mutex{},
	}
}

func (supplier *DefaultSupplier) GetName() string {
	return supplier.name
}

func (supplier *DefaultSupplier) GetPlugin(pluginName string) (plugin IPlugin, isExist bool) {
	return supplier.getPlugin(pluginName)
}

func (supplier *DefaultSupplier) getPlugin(pluginName string) (plugin IPlugin, isExist bool) {
	plugin, isExist = supplier.pluginMap[pluginName]
	return plugin, isExist
}

func (supplier *DefaultSupplier) GetAllPlugin() map[string]IPlugin {
	return supplier.pluginMap
}

func (supplier *DefaultSupplier) RegisterPlugin(plugin IPlugin) {
	supplier.locker.Lock()
	supplier.pluginMap[plugin.GetName()] = plugin
	supplier.locker.Unlock()
}

func (supplier *DefaultSupplier) Supply(ctx context.Context, pluginName string,
	params []interface{}) (map[string]interface{}, error) {
	plugin, isExist := supplier.getPlugin(pluginName)
	if !isExist {
		return map[string]interface{}{}, constant.NotFoundError
	}
	return plugin.Call(ctx, params...)
}
