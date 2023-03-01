package trees

type TreeResult struct {
	Nodes []*TreeNode `json:"nodes" mapstructure:"Nodes"`
}
type TreeNode struct {
	Key       string      `json:"key" mapstructure:"Key"`
	Title     string      `json:"title" mapstructure:"Title"`
	Level     int         `json:"level" mapstructure:"Level"`
	Children  []*TreeNode `json:"children" mapstructure:"Children"`
	Checked   bool        `json:"checked" mapstructure:"Checked"`
	IsLeaf    bool        `json:"isLeaf" mapstructure:"IsLeaf"`
	ParentKey string      `json:"-"`
	Origin    interface{} `json:"origin" mapstructure:"Origin"`
}

func BuildTree(nodes []*TreeNode) (*TreeResult, error) {
	result := &TreeResult{}
	for _, node := range nodes {
		if len(node.ParentKey) == 0 {
			node.getChildren(nodes)
			result.Nodes = append(result.Nodes, node)
		}
	}
	return result, nil
}

func (node *TreeNode) getChildren(source []*TreeNode) {
	for _, item := range source {
		if item.ParentKey == node.Key {
			item.getChildren(source)
			node.Children = append(node.Children, item)
		}
	}
	if node.Children == nil || len(node.Children) == 0 {
		node.IsLeaf = true
	} else {
		node.IsLeaf = false
	}
}
