package fun

type Inner struct {
	Name string
}

type Metadata struct {
	PackageName Inner
	Imports     []string
	Type        string
}

type Entry struct {
	Name string
	Done bool
}

type ToDo struct {
	User string
	List []Entry
}
