package config

import (
	"github.com/Si40Code/kit/config/provider"
)

type Option func(*options)

type options struct {
	filePaths      []string
	useEnv         bool
	envPrefix      string
	watchFile      bool
	remoteProvider provider.RemoteProvider
	defaults       map[string]interface{}
}

func newOptions(opts ...Option) *options {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// WithFile 加载单个配置文件
func WithFile(path string) Option {
	return func(o *options) {
		o.filePaths = append(o.filePaths, path)
	}
}

// WithFiles 加载多个配置文件（按顺序加载，后面的覆盖前面的）
func WithFiles(paths ...string) Option {
	return func(o *options) {
		o.filePaths = append(o.filePaths, paths...)
	}
}

func WithEnv(prefix string) Option {
	return func(o *options) {
		o.useEnv = true
		o.envPrefix = prefix
	}
}

func WithFileWatcher() Option {
	return func(o *options) {
		o.watchFile = true
	}
}

func WithRemote(p provider.RemoteProvider) Option {
	return func(o *options) {
		o.remoteProvider = p
	}
}

// WithDefaults 设置默认配置值
func WithDefaults(defaults map[string]interface{}) Option {
	return func(o *options) {
		o.defaults = defaults
	}
}

// WithDefaultStruct 从结构体设置默认配置值
func WithDefaultStruct(defaultStruct interface{}) Option {
	return func(o *options) {
		o.defaults = map[string]interface{}{"_struct": defaultStruct}
	}
}
