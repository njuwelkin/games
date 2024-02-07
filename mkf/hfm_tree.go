package main

type YJ1TreeNode struct {
	value byte
	leaf  bool
	//level  uint16
	//weight uint32
	//parent *YJ1TreeNode
	left  *YJ1TreeNode
	right *YJ1TreeNode
}

func makeHFMTree(src []byte, flags []byte, tree_len int) *YJ1TreeNode {
	br := NewBitReader(flags)
	tree := make([]YJ1TreeNode, tree_len)
	root := &tree[0]
	root.leaf = false
	root.value = 0
	root.left = &tree[1]
	root.right = &tree[2]
	for i := 1; i < tree_len; i++ {
		tree[i].leaf = br.Read(1) == 0
		tree[i].value = src[i]
		if root.leaf {
			tree[i].left = nil
			tree[i].right = nil
		} else {
			tree[i].left = &tree[int(root.value)*2+1]
			tree[i].right = &tree[int(root.value)*2+2]
		}
	}
	return root
}
