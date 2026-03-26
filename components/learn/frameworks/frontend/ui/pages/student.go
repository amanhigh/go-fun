package pages

import (
	"github.com/a-h/templ"
	"github.com/amanhigh/go-fun/components/learn/frameworks/frontend/ui/components"
)

// StudentListComponent renders the student management CRUD showcase page.
type StudentListComponent struct {
	components.BaseComponent
}

// NewStudentListComponent creates the student management showcase component.
func NewStudentListComponent() *StudentListComponent {
	return &StudentListComponent{
		BaseComponent: components.NewBaseComponent(
			"student-management",
			"👥 Student Management - Search, filter, and complete CRUD workflows with modals",
			"/students",
			components.LevelMedium,
			2,
		),
	}
}

var _ components.Component = (*StudentListComponent)(nil)

func (c *StudentListComponent) Render() templ.Component {
	return StudentPage()
}
