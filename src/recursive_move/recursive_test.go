package recursive_move

import "testing"

func TestRecursive(t *testing.T) {
	targetDir := "E:/예능/골목식당"

	recursive(targetDir, targetDir, true)
}
