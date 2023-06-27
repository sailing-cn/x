package tree

type TreeResult struct {
	Nodes []*TreeNode `json:"nodes"`
}

type TreeNode struct {
	Key       string      `json:"key"`
	Title     string      `json:"title"`
	Level     int32       `json:"level"`
	Checked   bool        `json:"checked"`
	IsLeaf    bool        `json:"is_leaf"`
	Origin    interface{} `json:"origin"`
	Children  []*TreeNode `json:"children"`
	ParentKey string      `json:"parent_key"`
}

type ByLevel []*TreeNode

func (l ByLevel) Len() int {
	return len(l)
}
func (l ByLevel) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l ByLevel) Less(i, j int) bool { return l[i].Level < l[j].Level }

func (t *TreeResult) ToTree() *TreeResult {
	result := &TreeResult{}
	temp := make(map[string]*TreeNode)
	topLevel := getTopLevel(t.Nodes)
	for len(temp) < len(t.Nodes) {
		for _, node := range t.Nodes {
			if _, ok := temp[node.Key]; ok {
				continue
			}
			node.HandleLeaf(t.Nodes)
			if node.Level == topLevel {
				result.Nodes = append(result.Nodes, node)
				temp[node.Key] = node
			} else if v, ok := temp[node.ParentKey]; ok {
				if v.Children == nil {
					v.Children = make([]*TreeNode, 0)
				}
				v.Children = append(v.Children, node)
				temp[node.Key] = node
			}
		}
	}
	return result
}

func getTopLevel(list []*TreeNode) int32 {
	var level = list[0].Level
	for _, region := range list {
		_level := region.Level
		if level > _level {
			level = _level
		}
	}
	return level
}

func (node *TreeNode) HandleLeaf(list []*TreeNode) {
	if len(node.ParentKey) == 0 {
		node.IsLeaf = false
		return
	}
	for _, item := range list {
		if item.ParentKey != node.Key && len(item.Key) > 0 {
			node.IsLeaf = true
		} else {
			node.IsLeaf = false
		}
	}
}
