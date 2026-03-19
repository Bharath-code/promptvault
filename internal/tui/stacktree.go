package tui

import (
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type TreeNode struct {
	Name     string
	Path     string
	Parent   *TreeNode
	Children map[string]*TreeNode
	Expanded bool
	Count    int
}

var (
	treeStyle = lipgloss.NewStyle().
			Foreground(color("#64748B"))

	treeSelectedStyle = lipgloss.NewStyle().
				Background(color("#334155")).
				Foreground(color("#E2E8F0")).
				Bold(true)

	stackIconStyle  = lipgloss.NewStyle().Foreground(color("#7C3AED"))
	countBadgeStyle = lipgloss.NewStyle().Foreground(color("#64748B"))
	expandIconStyle = lipgloss.NewStyle().Foreground(color("#64748B"))
)

type StackTree struct {
	root      *TreeNode
	nodes     []*TreeNode
	cursor    int
	available []string
	width     int
}

func NewStackTree(availableStacks []string, width int) *StackTree {
	st := &StackTree{
		width:     width,
		available: availableStacks,
	}
	st.buildTree()
	return st
}

func (st *StackTree) buildTree() {
	st.root = &TreeNode{
		Name:     "All Stacks",
		Path:     "",
		Children: make(map[string]*TreeNode),
		Expanded: true,
	}

	for _, stack := range st.available {
		parts := strings.Split(stack, "/")
		current := st.root

		for i, part := range parts {
			if current.Children == nil {
				current.Children = make(map[string]*TreeNode)
			}

			if child, exists := current.Children[part]; exists {
				current = child
			} else {
				path := strings.Join(parts[:i+1], "/")
				newNode := &TreeNode{
					Name:     part,
					Path:     path,
					Parent:   current,
					Children: make(map[string]*TreeNode),
					Expanded: i == 0,
				}
				current.Children[part] = newNode
				current = newNode
			}
		}
	}

	st.rebuildNodeList()
}

func (st *StackTree) rebuildNodeList() {
	st.nodes = nil
	st.flattenTree(st.root, 0)
}

func (st *StackTree) flattenTree(node *TreeNode, depth int) {
	st.nodes = append(st.nodes, node)
	for _, child := range st.sortedChildren(node) {
		if node.Expanded {
			st.flattenTree(child, depth+1)
		}
	}
}

func (st *StackTree) sortedChildren(node *TreeNode) []*TreeNode {
	var children []*TreeNode
	for _, child := range node.Children {
		children = append(children, child)
	}
	sort.Slice(children, func(i, j int) bool {
		return children[i].Name < children[j].Name
	})
	return children
}

func (st *StackTree) Current() *TreeNode {
	if st.cursor >= 0 && st.cursor < len(st.nodes) {
		return st.nodes[st.cursor]
	}
	return nil
}

func (st *StackTree) MoveUp() bool {
	if st.cursor > 0 {
		st.cursor--
		return true
	}
	return false
}

func (st *StackTree) MoveDown() bool {
	if st.cursor < len(st.nodes)-1 {
		st.cursor++
		return true
	}
	return false
}

func (st *StackTree) ToggleExpand() {
	node := st.Current()
	if node != nil && len(node.Children) > 0 {
		node.Expanded = !node.Expanded
		st.rebuildNodeList()
	}
}

func (st *StackTree) Expand() {
	node := st.Current()
	if node != nil && !node.Expanded {
		node.Expanded = true
		st.rebuildNodeList()
	}
}

func (st *StackTree) Collapse() {
	node := st.Current()
	if node != nil && node.Expanded {
		node.Expanded = false
		st.rebuildNodeList()
	}
}

func (st *StackTree) Render() string {
	var lines []string

	// Header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(color("#7C3AED")).
		PaddingBottom(1)

	lines = append(lines, headerStyle.Render("📁 Stack Browser"))

	// Build tree lines
	for i, node := range st.nodes {
		selected := i == st.cursor
		line := st.renderNode(node, selected)
		lines = append(lines, line)
	}

	// Footer
	footerStyle := lipgloss.NewStyle().
		Foreground(color("#64748B")).
		Italic(true)

	lines = append(lines, "")
	lines = append(lines, footerStyle.Render("← → expand/collapse  •  Enter select  •  Esc close"))

	return lipgloss.NewStyle().
		Width(st.width).
		Height(st.getHeight()).
		Render(strings.Join(lines, "\n"))
}

func (st *StackTree) renderNode(node *TreeNode, selected bool) string {
	depth := st.getDepth(node)
	indent := strings.Repeat("  ", depth)

	// Expand/collapse icon
	var expandIcon string
	if len(node.Children) > 0 {
		if node.Expanded {
			expandIcon = expandIconStyle.Render("▼")
		} else {
			expandIcon = expandIconStyle.Render("▶")
		}
	} else {
		expandIcon = "  "
	}

	// Node icon
	var nodeIcon string
	if len(node.Children) > 0 {
		nodeIcon = stackIconStyle.Render("📂")
	} else {
		nodeIcon = stackIconStyle.Render("📄")
	}

	// Name
	name := node.Name
	if node.Path == "" {
		name = "All Stacks"
	}

	// Count badge
	countStr := ""
	if node.Count > 0 {
		countStr = countBadgeStyle.Render(" (" + itoa(node.Count) + ")")
	}

	// Selected marker
	marker := "  "
	if selected {
		marker = "►"
	}

	fullLine := indent + expandIcon + " " + marker + " " + nodeIcon + " " + name + countStr

	if selected {
		return treeSelectedStyle.Render(fullLine)
	}
	return treeStyle.Render(fullLine)
}

func (st *StackTree) getDepth(node *TreeNode) int {
	if node.Path == "" {
		return 0
	}
	return strings.Count(node.Path, "/") + 1
}

func (st *StackTree) getHeight() int {
	h := len(st.nodes) + 4 // header + footer
	if h < 10 {
		h = 10
	}
	if h > 30 {
		h = 30
	}
	return h
}

func (st *StackTree) IsSelectable() bool {
	node := st.Current()
	return node != nil && node.Path != ""
}

// UpdateCounts updates the count for each node based on stack usage
func (st *StackTree) UpdateCounts(stackCounts map[string]int) {
	// Reset counts
	st.resetCounts(st.root)

	// Update counts
	for stack, count := range stackCounts {
		node := st.findNode(st.root, stack)
		if node != nil {
			node.Count = count
			st.incrementParentCounts(node)
		}
	}
}

func (st *StackTree) findNode(node *TreeNode, path string) *TreeNode {
	if node.Path == path {
		return node
	}
	for _, child := range node.Children {
		if found := st.findNode(child, path); found != nil {
			return found
		}
	}
	return nil
}

func (st *StackTree) resetCounts(node *TreeNode) {
	node.Count = 0
	for _, child := range node.Children {
		st.resetCounts(child)
	}
}

func (st *StackTree) incrementParentCounts(node *TreeNode) {
	node.Count++
	if node.Parent != nil {
		st.incrementParentCounts(node.Parent)
	}
}

// Parent is a helper to track parent nodes
type ParentNode struct {
	*TreeNode
	Parent *TreeNode
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	result := ""
	for n > 0 {
		result = string(rune('0'+n%10)) + result
		n /= 10
	}
	return result
}
