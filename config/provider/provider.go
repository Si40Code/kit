package provider

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
)

func LoadFile(k *koanf.Koanf, path string) error {
	parser, err := getParser(path)
	if err != nil {
		return err
	}
	return k.Load(file.Provider(path), parser)
}

// getParser 根据文件扩展名返回对应的解析器
func getParser(path string) (koanf.Parser, error) {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".yaml", ".yml":
		return yaml.Parser(), nil
	case ".json":
		return json.Parser(), nil
	case ".toml":
		return toml.Parser(), nil
	default:
		return nil, fmt.Errorf("unsupported config file format: %s (supported: .yaml, .yml, .json, .toml)", ext)
	}
}

func LoadEnv(k *koanf.Koanf, prefix string) error {
	return k.Load(env.Provider(prefix, ".", func(s string) string {
		return strings.ToLower(strings.TrimPrefix(s, prefix))
	}), nil)
}

func LoadDefaults(k *koanf.Koanf, defaults map[string]interface{}) error {
	// 如果 defaults 中包含 _struct 键，说明是从结构体传入的
	if structValue, ok := defaults["_struct"]; ok {
		return k.Load(structs.Provider(structValue, "koanf"), nil)
	}

	// 否则直接加载 map
	for key, value := range defaults {
		k.Set(key, value)
	}
	return nil
}

func MapProvider(data map[string]interface{}) *rawbytes.RawBytes {
	return rawbytes.Provider(nil)
}
