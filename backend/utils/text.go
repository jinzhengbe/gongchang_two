package utils

import (
	"regexp"
	"strings"
)

// HighlightText 高亮文本中的关键词
func HighlightText(text, keyword string) string {
	if keyword == "" || text == "" {
		return text
	}

	// 转义正则表达式特殊字符
	escapedKeyword := regexp.QuoteMeta(keyword)
	
	// 创建不区分大小写的正则表达式
	re, err := regexp.Compile("(?i)" + escapedKeyword)
	if err != nil {
		return text
	}

	// 替换匹配的文本为高亮格式
	highlighted := re.ReplaceAllString(text, "<em>$0</em>")
	
	return highlighted
}

// ExtractKeywords 从文本中提取关键词
func ExtractKeywords(text string, maxKeywords int) []string {
	if text == "" {
		return []string{}
	}

	// 简单的关键词提取：按空格分割，过滤短词
	words := strings.Fields(text)
	keywords := make([]string, 0)
	seen := make(map[string]bool)

	for _, word := range words {
		// 清理单词（去除标点符号）
		cleanWord := cleanWord(word)
		
		// 过滤条件：长度大于1，不是常见停用词
		if len(cleanWord) > 1 && !isStopWord(cleanWord) && !seen[cleanWord] {
			keywords = append(keywords, cleanWord)
			seen[cleanWord] = true
			
			if len(keywords) >= maxKeywords {
				break
			}
		}
	}

	return keywords
}

// cleanWord 清理单词
func cleanWord(word string) string {
	// 去除标点符号
	re := regexp.MustCompile(`[^\p{Han}\w]`)
	clean := re.ReplaceAllString(word, "")
	
	return strings.ToLower(clean)
}

// isStopWord 判断是否为停用词
func isStopWord(word string) bool {
	stopWords := map[string]bool{
		"的": true, "了": true, "在": true, "是": true, "我": true, "有": true, "和": true, "就": true,
		"不": true, "人": true, "都": true, "一": true, "一个": true, "上": true, "也": true, "很": true,
		"到": true, "说": true, "要": true, "去": true, "你": true, "会": true, "着": true, "没有": true,
		"看": true, "好": true, "自己": true, "这": true, "那": true, "什么": true, "怎么": true,
		"the": true, "a": true, "an": true, "and": true, "or": true, "but": true, "in": true, "on": true,
		"at": true, "to": true, "for": true, "of": true, "with": true, "by": true, "is": true, "are": true,
		"was": true, "were": true, "be": true, "been": true, "being": true, "have": true, "has": true,
		"had": true, "do": true, "does": true, "did": true, "will": true, "would": true, "could": true,
		"should": true, "may": true, "might": true, "can": true, "this": true, "that": true, "these": true,
		"those": true, "i": true, "you": true, "he": true, "she": true, "it": true, "we": true, "they": true,
	}

	return stopWords[word]
}

// CalculateRelevance 计算相关性分数
func CalculateRelevance(text, keyword string) float64 {
	if keyword == "" || text == "" {
		return 0.0
	}

	text = strings.ToLower(text)
	keyword = strings.ToLower(keyword)

	// 精确匹配权重最高
	if strings.Contains(text, keyword) {
		return 1.0
	}

	// 部分匹配
	keywords := strings.Fields(keyword)
	score := 0.0
	totalKeywords := float64(len(keywords))

	for _, kw := range keywords {
		if strings.Contains(text, kw) {
			score += 0.5
		}
	}

	return score / totalKeywords
} 