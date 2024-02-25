package fun

type Inner struct {
	Name string
}

type Metadata struct {
	PackageName Inner
	Imports     []string
	Type        string
}
