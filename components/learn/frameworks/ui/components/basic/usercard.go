package basic

import (
	"context"
	"io"

	"github.com/a-h/templ"
	"github.com/amanhigh/go-fun/components/learn/frameworks/ui/components"
)

// UserCardComponent demonstrates conditional rendering with active/inactive badge
type UserCardComponent struct {
	*components.BaseComponent
	username string
	isActive bool
}

var _ components.Component = (*UserCardComponent)(nil)

// NewUserCardComponent creates a new user card component
func NewUserCardComponent(username string, isActive bool) *UserCardComponent {
	activeStr := "inactive"
	if isActive {
		activeStr = "active"
	}
	c := &UserCardComponent{username: username, isActive: isActive}
	c.BaseComponent = components.NewBaseComponent(
		"usercard",
		"User card with conditional active/inactive badge rendering",
		"/basic/usercard/"+username+"/"+activeStr,
		components.LevelBasic,
		2,
		c.render,
	)
	return c
}

func (c *UserCardComponent) render() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, _ = io.WriteString(w, `<div class="user-card"><h2>`+c.username+`</h2>`)
		if c.isActive {
			_, _ = io.WriteString(w, `<span class="badge active">Active</span>`)
		} else {
			_, _ = io.WriteString(w, `<span class="badge inactive">Inactive</span>`)
		}
		_, _ = io.WriteString(w, `</div>`)
		return nil
	})
}

// DefaultUserCardComponent returns the default user card component for demo
func DefaultUserCardComponent() *UserCardComponent {
	return NewUserCardComponent("John", true)
}
