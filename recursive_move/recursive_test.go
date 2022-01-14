package recursive_move

import "testing"

func TestRecursive(t *testing.T) {
	targetDir := "D:\\download\\_new\\[Moozzi2] Jujutsu Kaisen [ x265-10Bit Ver. ] - TV + SP"

	MoveRecursive(targetDir, targetDir, true)
}

func TestMoveSubtitle(t *testing.T) {
	MoveSubtitle("D:\\download\\_new\\[Moozzi2] Jujutsu Kaisen [ x265-10Bit Ver. ] - TV + SP")
}
