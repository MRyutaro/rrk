package tree

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/MRyutaro/rrk/internal/history"
)

// DirectoryNode ディレクトリツリーのノードを表現
type DirectoryNode struct {
	Path     string
	Commands []string
	Children map[string]*DirectoryNode
}

// NewDirectoryNode 新しいディレクトリノードを作成
func NewDirectoryNode(path string) *DirectoryNode {
	return &DirectoryNode{
		Path:     path,
		Commands: make([]string, 0),
		Children: make(map[string]*DirectoryNode),
	}
}

// AddCommand ディレクトリにコマンドを追加
func (node *DirectoryNode) AddCommand(command string) {
	node.Commands = append(node.Commands, command)
}

// TreeBuilder ディレクトリツリー構築器
type TreeBuilder struct {
	root *DirectoryNode
}

// NewTreeBuilder 新しいツリー構築器を作成
func NewTreeBuilder() *TreeBuilder {
	return &TreeBuilder{
		root: NewDirectoryNode(""),
	}
}

// BuildTree 履歴エントリからディレクトリツリーを構築
func (tb *TreeBuilder) BuildTree(entries []history.Entry, limit int) *DirectoryNode {
	// ディレクトリごとにコマンドをグループ化
	dirCommands := make(map[string][]string)
	
	for _, entry := range entries {
		if entry.CWD != "" && entry.Command != "" {
			dirCommands[entry.CWD] = append(dirCommands[entry.CWD], entry.Command)
		}
	}
	
	// 各ディレクトリで重複を除去し、制限を適用
	for dir, commands := range dirCommands {
		uniqueCommands := removeDuplicates(commands)
		if limit > 0 && len(uniqueCommands) > limit {
			uniqueCommands = uniqueCommands[len(uniqueCommands)-limit:]
		}
		dirCommands[dir] = uniqueCommands
	}
	
	// ツリー構造を構築
	return tb.buildDirectoryTree(dirCommands)
}

// buildDirectoryTree ディレクトリマップからツリー構造を構築
func (tb *TreeBuilder) buildDirectoryTree(dirCommands map[string][]string) *DirectoryNode {
	root := NewDirectoryNode("")
	
	// すべてのディレクトリパスを処理
	for dirPath, commands := range dirCommands {
		tb.addDirectoryToTree(root, dirPath, commands)
	}
	
	return root
}

// addDirectoryToTree ツリーにディレクトリとコマンドを追加
func (tb *TreeBuilder) addDirectoryToTree(root *DirectoryNode, dirPath string, commands []string) {
	if dirPath == "" {
		return
	}
	
	// パスを正規化
	cleanPath := filepath.Clean(dirPath)
	if cleanPath == "." {
		cleanPath = ""
	}
	
	// ルートディレクトリの場合
	if cleanPath == "" || cleanPath == "/" {
		for _, cmd := range commands {
			root.AddCommand(cmd)
		}
		return
	}
	
	// パスコンポーネントに分割
	parts := strings.Split(cleanPath, string(filepath.Separator))
	if parts[0] == "" {
		parts = parts[1:] // 絶対パスの先頭の空文字列を除去
	}
	
	current := root
	currentPath := ""
	
	// パスの各部分を辿ってノードを作成
	for _, part := range parts {
		if part == "" {
			continue
		}
		
		if currentPath == "" {
			currentPath = "/" + part
		} else {
			currentPath = filepath.Join(currentPath, part)
		}
		
		if current.Children[part] == nil {
			current.Children[part] = NewDirectoryNode(currentPath)
		}
		current = current.Children[part]
	}
	
	// 最終ノードにコマンドを追加
	for _, cmd := range commands {
		current.AddCommand(cmd)
	}
}

// PrintTree ツリーを表示
func PrintTree(root *DirectoryNode, rootPath string, maxCommands int) {
	if root == nil {
		return
	}
	
	// ルートパスが指定されている場合、そのパス以下のみを表示
	if rootPath != "" {
		targetNode := findNodeByPath(root, rootPath)
		if targetNode != nil {
			printNode(targetNode, "", true, maxCommands)
		} else {
			fmt.Printf("No history found for path: %s\n", rootPath)
		}
		return
	}
	
	// 全体を表示
	printTreeRecursive(root, "", maxCommands)
}

// findNodeByPath パスで指定されたノードを検索
func findNodeByPath(root *DirectoryNode, targetPath string) *DirectoryNode {
	if root == nil {
		return nil
	}
	
	targetPath = filepath.Clean(targetPath)
	if targetPath == "." || targetPath == "" {
		return root
	}
	
	// ルートノードの子ノードから検索
	parts := strings.Split(targetPath, string(filepath.Separator))
	if parts[0] == "" {
		parts = parts[1:]
	}
	
	current := root
	for _, part := range parts {
		if part == "" {
			continue
		}
		if current.Children[part] != nil {
			current = current.Children[part]
		} else {
			return nil
		}
	}
	
	return current
}

// printTreeRecursive ツリー全体を再帰的に表示
func printTreeRecursive(node *DirectoryNode, prefix string, maxCommands int) {
	if node == nil {
		return
	}
	
	// ルートレベルのディレクトリを特定（/で始まるパス）
	rootDirs := make(map[string]*DirectoryNode)
	
	// ルートディレクトリを特定
	for name, child := range node.Children {
		if strings.HasPrefix(child.Path, "/") {
			rootPath := "/" + name
			rootDirs[rootPath] = child
		}
	}
	
	// ルートディレクトリをソート
	var sortedRoots []string
	for root := range rootDirs {
		sortedRoots = append(sortedRoots, root)
	}
	sort.Strings(sortedRoots)
	
	// 各ルートディレクトリを表示
	for i, rootPath := range sortedRoots {
		rootNode := rootDirs[rootPath]
		
		// ルートディレクトリ名を表示
		fmt.Printf("%s\n", rootPath)
		
		// ルートディレクトリのコマンドを表示
		if len(rootNode.Commands) > 0 {
			printCommandsWithTree(rootNode.Commands, "", true, maxCommands)
		}
		
		// 子ディレクトリを表示
		printDirectoryChildren(rootNode, "", maxCommands)
		
		// 最後のルートディレクトリでなければ空行を追加
		if i < len(sortedRoots)-1 {
			fmt.Println()
		}
	}
}

// printNode 特定のノードを表示
func printNode(node *DirectoryNode, prefix string, isRoot bool, maxCommands int) {
	if node == nil {
		return
	}
	
	// パスを表示（ルートでない場合）
	if !isRoot && node.Path != "" {
		fmt.Printf("%s/\n", node.Path)
	}
	
	// コマンドを表示
	if len(node.Commands) > 0 {
		printCommands(node.Commands, "├── ", maxCommands)
	}
	
	// 子ディレクトリをソート
	var childNames []string
	for name := range node.Children {
		childNames = append(childNames, name)
	}
	sort.Strings(childNames)
	
	// 子ディレクトリを表示
	for i, name := range childNames {
		child := node.Children[name]
		isLast := i == len(childNames)-1
		
		var childPrefix string
		if isLast {
			fmt.Printf("└── %s/\n", name)
			childPrefix = "    "
		} else {
			fmt.Printf("├── %s/\n", name)
			childPrefix = "│   "
		}
		
		// 子のコマンドを表示
		if len(child.Commands) > 0 {
			printCommands(child.Commands, childPrefix+"├── ", maxCommands)
		}
		
		// 孫ディレクトリがある場合は再帰的に処理
		if len(child.Children) > 0 {
			printChildNodes(child, childPrefix, maxCommands)
		}
	}
}

// printChildNodes 子ノードを再帰的に表示
func printChildNodes(node *DirectoryNode, prefix string, maxCommands int) {
	var childNames []string
	for name := range node.Children {
		childNames = append(childNames, name)
	}
	sort.Strings(childNames)
	
	for i, name := range childNames {
		child := node.Children[name]
		isLast := i == len(childNames)-1
		
		var childPrefix string
		if isLast {
			fmt.Printf("%s└── %s/\n", prefix, name)
			childPrefix = prefix + "    "
		} else {
			fmt.Printf("%s├── %s/\n", prefix, name)
			childPrefix = prefix + "│   "
		}
		
		// コマンドを表示
		if len(child.Commands) > 0 {
			printCommands(child.Commands, childPrefix+"├── ", maxCommands)
		}
		
		// 再帰的に子ノードを処理
		if len(child.Children) > 0 {
			printChildNodes(child, childPrefix, maxCommands)
		}
	}
}

// printDirectoryChildren ディレクトリの子ノードを表示
func printDirectoryChildren(node *DirectoryNode, prefix string, maxCommands int) {
	if node == nil || len(node.Children) == 0 {
		return
	}
	
	// 子ディレクトリをソート
	var childNames []string
	for name := range node.Children {
		childNames = append(childNames, name)
	}
	sort.Strings(childNames)
	
	// 各子ディレクトリを表示
	for i, name := range childNames {
		child := node.Children[name]
		isLast := i == len(childNames)-1
		
		// ディレクトリ名を表示
		var treeChar string
		if isLast {
			treeChar = "└── "
		} else {
			treeChar = "├── "
		}
		fmt.Printf("%s%s%s/\n", prefix, treeChar, name)
		
		// 子ディレクトリのコマンドを表示
		var childPrefix string
		if isLast {
			childPrefix = prefix + "    "
		} else {
			childPrefix = prefix + "│   "
		}
		
		if len(child.Commands) > 0 {
			printCommandsWithTree(child.Commands, childPrefix, false, maxCommands)
		}
		
		// 孫ディレクトリを再帰的に表示
		if len(child.Children) > 0 {
			printDirectoryChildren(child, childPrefix, maxCommands)
		}
	}
}

// printCommandsWithTree ツリー形式でコマンドを表示
func printCommandsWithTree(commands []string, prefix string, isRoot bool, maxCommands int) {
	if len(commands) == 0 {
		return
	}
	
	displayCommands := commands
	if maxCommands > 0 && len(commands) > maxCommands {
		displayCommands = commands[len(commands)-maxCommands:]
	}
	
	for i, cmd := range displayCommands {
		isLast := i == len(displayCommands)-1
		if isLast {
			fmt.Printf("%s└── %s\n", prefix, cmd)
		} else {
			fmt.Printf("%s├── %s\n", prefix, cmd)
		}
	}
}

// printCommands コマンドリストを表示
func printCommands(commands []string, prefix string, maxCommands int) {
	if len(commands) == 0 {
		return
	}
	
	displayCommands := commands
	if maxCommands > 0 && len(commands) > maxCommands {
		displayCommands = commands[len(commands)-maxCommands:]
	}
	
	for i, cmd := range displayCommands {
		isLast := i == len(displayCommands)-1
		if isLast && strings.HasSuffix(prefix, "├── ") {
			// 最後のコマンドの場合は└──を使用
			modifiedPrefix := strings.Replace(prefix, "├── ", "└── ", 1)
			fmt.Printf("%s%s\n", modifiedPrefix, cmd)
		} else {
			fmt.Printf("%s%s\n", prefix, cmd)
		}
	}
}

// removeDuplicates 重複を除去（順序を保持）
func removeDuplicates(slice []string) []string {
	seen := make(map[string]bool)
	var result []string
	
	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	
	return result
}