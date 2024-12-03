package tree

// Node 파일시스템의 단일 항목(파일/디렉토리)을 표현
type Node struct {
	Path     string
	Name     string
	IsDir    bool
	Selected bool
	Children []*Node
	Parent   *Node
	IsLast   bool
}

// ToggleSelect는 현재 노드와 모든 하위 노드의 선택 상태를 토글하고,
// 상위 노드들의 선택 상태도 자식 노드들의 상태에 따라 업데이트합니다.
func (n *Node) ToggleSelect() {
	newState := !n.Selected
	n.setSelectionState(newState)
	n.updateParentSelection()
}

// setSelectionState는 현재 노드와 모든 하위 노드의 선택 상태를 설정합니다.
func (n *Node) setSelectionState(selected bool) {
	n.Selected = selected
	for _, child := range n.Children {
		child.setSelectionState(selected)
	}
}

// updateParentSelection은 상위 노드들의 선택 상태를
// 자식 노드들의 상태에 따라 업데이트합니다.
func (n *Node) updateParentSelection() {
	if n.Parent == nil {
		return
	}

	parent := n.Parent
	allSelected := true
	noneSelected := true

	for _, child := range parent.Children {
		if child.Selected {
			noneSelected = false
		} else {
			allSelected = false
		}

		// 일부만 선택된 상태를 발견하면 더 이상 검사할 필요 없음
		if !allSelected && !noneSelected {
			break
		}
	}

	// 세 가지 상태 처리:
	// 1. 모든 자식이 선택됨 -> 부모도 선택
	// 2. 모든 자식이 선택 해제됨 -> 부모도 선택 해제
	// 3. 일부 자식만 선택됨 -> 부모는 선택 해제
	if allSelected {
		parent.Selected = true
	} else { // noneSelected이거나 일부만 선택된 경우
		parent.Selected = false
	}

	// 재귀적으로 상위로 올라가며 업데이트
	parent.updateParentSelection()
}

// GetSelected는 선택된 모든 파일 노드를 반환합니다.
func (n *Node) GetSelected() []*Node {
	var selected []*Node
	if n.Selected && !n.IsDir {
		selected = append(selected, n)
	}
	for _, child := range n.Children {
		selected = append(selected, child.GetSelected()...)
	}
	return selected
}
