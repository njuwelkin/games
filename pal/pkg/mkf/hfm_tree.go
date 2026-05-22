package mkf

import "fmt"

type YJ1TreeNode struct {
	value byte
	leaf  bool
	//level  uint16
	//weight uint32
	parent *YJ1TreeNode
	left   *YJ1TreeNode
	right  *YJ1TreeNode
}

func makeHFMTree(src []byte, flags []byte, tree_len int) *YJ1TreeNode {
	br := NewBitReader(flags)
	tree := make([]YJ1TreeNode, tree_len+1)
	root := &tree[0]
	root.leaf = false
	root.value = 0
	root.left = &tree[1]
	tree[1].parent = root
	root.right = &tree[2]
	tree[2].parent = root
	for i := 1; i <= tree_len; i++ {
		if tree[i].parent == nil {
			panic("")
		}
		tree[i].leaf = br.Read(1) == 0
		tree[i].value = src[i-1]
		if tree[i].leaf {
			tree[i].left = nil
			tree[i].right = nil
		} else {
			fmt.Println(tree[i].value, int(tree[i].value), int(tree[i].value)*2+1)

			tree[i].left = &(tree[int(tree[i].value)*2+1])
			tree[i].left.parent = &tree[i]
			tree[i].right = &tree[int(tree[i].value)*2+2]
			tree[i].right.parent = &tree[i]
		}
	}
	return root
}

func (node *YJ1TreeNode) Print() {
	buf := []byte{}
	var dfs func(*YJ1TreeNode)
	dfs = func(node *YJ1TreeNode) {
		if node.leaf {
			fmt.Printf("%d: %s\n", node.value, string(buf))
			return
		}
		if node.left != nil {
			buf = append(buf, '0')
			dfs(node.left)
			buf = buf[:len(buf)-1]
		}
		if node.right != nil {
			buf = append(buf, '1')
			dfs(node.right)
			buf = buf[:len(buf)-1]
		}
	}
	dfs(node)
}
