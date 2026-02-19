package galao

type ViewNode struct {
	Kind     string     `json:"kind"`
	ID       string     `json:"id,omitempty"`
	Value    string     `json:"value,omitempty"`
	Label    string     `json:"label,omitempty"`
	Children []ViewNode `json:"children,omitempty"`
}

func VStack(children ...ViewNode) ViewNode {
	return ViewNode{Kind: "vstack", Children: children}
}

func HStack(children ...ViewNode) ViewNode {
	return ViewNode{Kind: "hstack", Children: children}
}

func Text(value string) ViewNode {
	return ViewNode{Kind: "text", Value: value}
}

func Button(id, label string) ViewNode {
	return ViewNode{Kind: "button", ID: id, Label: label}
}
