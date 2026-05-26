package prompt

const (
	SummaryPrompt = `请用简洁的语言总结以下聊天记录：
%s
`
	TodoPrompt = `从以下聊天记录中提取待办事项，每条一行：
%s
`
	ReplyPrompt = `根据以下聊天记录，生成3条简短的回复候选，每条一行：
%s
`
)
