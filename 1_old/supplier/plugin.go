package supplier

import (
	"context"
)

// Plugin 的概念来自于 bazooka, 表示一种可以动态扩展的能力.
type (
	IPlugin interface {
		GetName() string
		// TODO [optimize] 如果我直接使用 reflect.Method 呢?
		Call(ctx context.Context, args ...interface{}) (map[string]interface{}, error)
	}
	PluginFunc func(ctx context.Context, args ...interface{}) (map[string]interface{}, error)
)

// type Plugin = plugin.Plugin

type DefaultPlugin struct {
	name string
	fn   PluginFunc
}

var _ IPlugin = new(DefaultPlugin)

func NewDefaultPlugin(name string, fn PluginFunc) *DefaultPlugin {
	return &DefaultPlugin{
		name: name,
		fn:   fn,
	}
}

func (p *DefaultPlugin) GetName() string {
	return p.name
}

func (p *DefaultPlugin) Call(ctx context.Context, args ...interface{}) (map[string]interface{}, error) {
	return p.fn(ctx, args...)
}
