package conf

// Config 执行节点配置（business 段，热更新）
type Config struct {
	Shell string `yaml:"shell"` // 执行脚本的 shell 路径，默认 /bin/sh
}

var instance = &Config{}

// Instance 返回当前配置（有状态单例，随 Watch 热更新）
func Instance() *Config { return instance }

// Set 更新配置（main 的 Watch 回调调用）
func Set(c *Config) { instance = c }