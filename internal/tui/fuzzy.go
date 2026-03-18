package tui

import (
	"sort"
	"strings"
	"unicode"

	"github.com/Bharath-code/promptvault/internal/model"
)

// FuzzyMatch performs fuzzy matching and returns match status and score
func FuzzyMatch(query, text string) (bool, int) {
	if query == "" {
		return true, 100
	}

	query = strings.ToLower(query)
	text = strings.ToLower(text)

	// Exact match gets highest score
	if text == query {
		return true, 100
	}

	// Contains match gets high score
	if strings.Contains(text, query) {
		return true, 80
	}

	// Fuzzy match with scoring
	score := 0
	queryIdx := 0
	consecutiveMatches := 0
	maxConsecutive := 0
	lastMatchIdx := -10

	for i, char := range text {
		if queryIdx >= len(query) {
			break
		}

		if char == rune(query[queryIdx]) {
			// Base score for match
			score += 10

			// Bonus for consecutive matches
			consecutiveMatches++
			if consecutiveMatches > 1 {
				score += consecutiveMatches * 2
			}
			if consecutiveMatches > maxConsecutive {
				maxConsecutive = consecutiveMatches
			}

			// Bonus for word boundary matches
			if i == 0 || !unicode.IsLetter(rune(text[i-1])) {
				score += 5
			}

			// Bonus for matches at start
			if i == 0 {
				score += 10
			}

			lastMatchIdx = i
			queryIdx++
		} else {
			consecutiveMatches = 0
		}
	}

	// All query characters matched
	if queryIdx == len(query) {
		// Bonus for longer consecutive matches
		score += maxConsecutive * 3

		// Penalty for spreading matches across text
		if lastMatchIdx > len(text)/2 {
			score -= 5
		}

		// Cap score at 100
		if score > 100 {
			score = 100
		}

		return true, score
	}

	return false, 0
}

// FuzzySearch filters prompts by fuzzy matching on multiple fields
func FuzzySearch(query string, prompts []*model.Prompt) ([]*model.Prompt, []int) {
	if query == "" {
		scores := make([]int, len(prompts))
		for i := range scores {
			scores[i] = 100
		}
		return prompts, scores
	}

	type scoredPrompt struct {
		prompt *model.Prompt
		score  int
		index  int
	}

	var scored []scoredPrompt

	for i, p := range prompts {
		// Match on title (highest weight)
		if matched, score := FuzzyMatch(query, p.Title); matched {
			scored = append(scored, scoredPrompt{p, score, i})
			continue
		}

		// Match on stack
		if matched, score := FuzzyMatch(query, p.Stack); matched {
			scored = append(scored, scoredPrompt{p, score - 10, i})
			continue
		}

		// Match on tags
		for _, tag := range p.Tags {
			if matched, score := FuzzyMatch(query, tag); matched {
				scored = append(scored, scoredPrompt{p, score - 15, i})
				break
			}
		}

		// Match on content (lower weight, expensive)
		if len(query) > 2 {
			if matched, score := FuzzyMatch(query, p.Content); matched {
				scored = append(scored, scoredPrompt{p, score - 30, i})
			}
		}
	}

	// Sort by score (descending)
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	// Extract sorted prompts and scores
	result := make([]*model.Prompt, len(scored))
	scores := make([]int, len(scored))
	for i, sp := range scored {
		result[i] = sp.prompt
		scores[i] = sp.score
	}

	return result, scores
}

// HighlightMatch adds markdown formatting to matched characters
func HighlightMatch(text, query string) string {
	if query == "" || text == "" {
		return text
	}

	query = strings.ToLower(query)
	textLower := strings.ToLower(text)

	var result strings.Builder
	queryIdx := 0
	inMatch := false

	for i, char := range text {
		if queryIdx < len(query) && char == rune(query[queryIdx]) {
			if !inMatch {
				result.WriteString("**")
				inMatch = true
			}
			result.WriteRune(char)
			queryIdx++
		} else {
			if inMatch {
				result.WriteString("**")
				inMatch = false
			}
			result.WriteRune(char)
		}

		// Prevent index out of bounds
		_ = textLower
		_ = i
	}

	if inMatch {
		result.WriteString("**")
	}

	return result.String()
}
