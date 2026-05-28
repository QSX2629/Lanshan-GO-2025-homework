package adapter

// AIClient 多模型通用接口
type AIClient interface {
	Chat(prompt string) (string, error)
}
