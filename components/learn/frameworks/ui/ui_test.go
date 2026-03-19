//nolint:cyclop // Test file with multiple test cases for learning purposes
package ui

import (
	"bytes"
	"context"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Templ UI Components", func() {

	Context("Basic", func() {
		It("1.1 should render a simple greeting component", func() {
			By("Creating a greeting component")
			component := Greeting("Alice")

			By("Rendering the component to HTML")
			var buf bytes.Buffer
			err := component.Render(context.Background(), &buf)
			Expect(err).ToNot(HaveOccurred())

			By("Verifying the rendered HTML contains expected content")
			html := buf.String()
			Expect(html).To(ContainSubstring("Hello, Alice!"))
			Expect(html).To(ContainSubstring("Welcome to Templ learning"))
			Expect(html).To(ContainSubstring(`class="greeting"`))
		})

		It("1.2 should render component with conditional logic", func() {
			By("Rendering active user card")
			activeComponent := UserCard("John", true)
			var activeBuf bytes.Buffer
			err := activeComponent.Render(context.Background(), &activeBuf)
			Expect(err).ToNot(HaveOccurred())

			activeHTML := activeBuf.String()
			Expect(activeHTML).To(ContainSubstring("John"))
			Expect(activeHTML).To(ContainSubstring("badge active"))
			Expect(activeHTML).To(ContainSubstring("Active"))

			By("Rendering inactive user card")
			inactiveComponent := UserCard("Jane", false)
			var inactiveBuf bytes.Buffer
			err = inactiveComponent.Render(context.Background(), &inactiveBuf)
			Expect(err).ToNot(HaveOccurred())

			inactiveHTML := inactiveBuf.String()
			Expect(inactiveHTML).To(ContainSubstring("Jane"))
			Expect(inactiveHTML).To(ContainSubstring("badge inactive"))
			Expect(inactiveHTML).To(ContainSubstring("Inactive"))
		})

		It("1.3 should render component with loops", func() {
			todos := []string{"Learn Templ", "Build UI", "Test Components"}

			By("Creating todo list component")
			component := TodoList(todos)

			By("Rendering the component")
			var buf bytes.Buffer
			err := component.Render(context.Background(), &buf)
			Expect(err).ToNot(HaveOccurred())

			By("Verifying all todo items are rendered")
			html := buf.String()
			Expect(html).To(ContainSubstring("Todo Items"))
			for _, todo := range todos {
				Expect(html).To(ContainSubstring(todo))
			}
			Expect(html).To(ContainSubstring("<li>"))
		})

		It("1.4 should render button with attributes", func() {
			By("Rendering enabled button")
			enabledBtn := Button("Click Me", false)
			var enabledBuf bytes.Buffer
			err := enabledBtn.Render(context.Background(), &enabledBuf)
			Expect(err).ToNot(HaveOccurred())

			enabledHTML := enabledBuf.String()
			Expect(enabledHTML).To(ContainSubstring("Click Me"))
			Expect(enabledHTML).To(ContainSubstring(`type="button"`))
			Expect(enabledHTML).ToNot(ContainSubstring("disabled"))

			By("Rendering disabled button")
			disabledBtn := Button("Disabled", true)
			var disabledBuf bytes.Buffer
			err = disabledBtn.Render(context.Background(), &disabledBuf)
			Expect(err).ToNot(HaveOccurred())

			disabledHTML := disabledBuf.String()
			Expect(disabledHTML).To(ContainSubstring("Disabled"))
			Expect(disabledHTML).To(ContainSubstring("disabled"))
		})

		It("1.5 should render empty list gracefully", func() {
			By("Creating empty todo list")
			component := TodoList([]string{})

			By("Rendering the component")
			var buf bytes.Buffer
			err := component.Render(context.Background(), &buf)
			Expect(err).ToNot(HaveOccurred())

			By("Verifying structure exists but no items")
			html := buf.String()
			Expect(html).To(ContainSubstring("Todo Items"))
			Expect(html).To(ContainSubstring("<ul>"))
			itemCount := strings.Count(html, "<li>")
			Expect(itemCount).To(Equal(0))
		})
	})

	Context("Medium", func() {
		It("2.1 should render nested components", func() {
			By("Creating nested page with greeting")
			greeting := Greeting("Bob")
			page := Page("Welcome Page", greeting)

			By("Rendering the full page")
			var buf bytes.Buffer
			err := page.Render(context.Background(), &buf)
			Expect(err).ToNot(HaveOccurred())

			By("Verifying page structure and nested content")
			html := buf.String()
			Expect(html).To(ContainSubstring("<!doctype html>"))
			Expect(html).To(ContainSubstring("<title>Welcome Page</title>"))
			Expect(html).To(ContainSubstring("Hello, Bob!"))
			Expect(html).To(ContainSubstring("font-family: Arial"))
		})

		It("2.2 should handle counter with different states", func() {
			testCases := []struct {
				count           int
				expectedMessage string
			}{
				{0, "Counter is at zero"},
				{5, "Counter is positive: 5"},
				{-3, "Counter is negative: -3"},
				{100, "Counter is positive: 100"},
			}

			for _, tc := range testCases {
				By("Testing counter with value: " + string(rune(tc.count)))
				component := Counter(tc.count)

				var buf bytes.Buffer
				err := component.Render(context.Background(), &buf)
				Expect(err).ToNot(HaveOccurred())

				html := buf.String()
				Expect(html).To(ContainSubstring("Counter Value"))
				Expect(html).To(ContainSubstring(tc.expectedMessage))
			}
		})

		It("2.3 should render data table with multiple rows", func() {
			rows := []TableRow{
				{ID: 1, Name: "Alice", Age: 25},
				{ID: 2, Name: "Bob", Age: 30},
				{ID: 3, Name: "Charlie", Age: 35},
			}

			By("Creating data table component")
			component := DataTable(rows)

			By("Rendering the table")
			var buf bytes.Buffer
			err := component.Render(context.Background(), &buf)
			Expect(err).ToNot(HaveOccurred())

			By("Verifying table structure")
			html := buf.String()
			Expect(html).To(ContainSubstring("<table"))
			Expect(html).To(ContainSubstring("<thead>"))
			Expect(html).To(ContainSubstring("<tbody>"))
			Expect(html).To(ContainSubstring("<th>ID</th>"))
			Expect(html).To(ContainSubstring("<th>Name</th>"))
			Expect(html).To(ContainSubstring("<th>Age</th>"))

			By("Verifying all row data is present")
			for _, row := range rows {
				Expect(html).To(ContainSubstring(row.Name))
			}
		})

		It("2.4 should compose multiple components together", func() {
			By("Creating multiple components")
			greeting := Greeting("Team")
			todos := []string{"Review code", "Deploy app", "Write docs"}
			todoList := TodoList(todos)

			By("Rendering greeting")
			var greetingBuf bytes.Buffer
			err := greeting.Render(context.Background(), &greetingBuf)
			Expect(err).ToNot(HaveOccurred())

			By("Rendering todo list")
			var todoBuf bytes.Buffer
			err = todoList.Render(context.Background(), &todoBuf)
			Expect(err).ToNot(HaveOccurred())

			By("Verifying both components rendered independently")
			greetingHTML := greetingBuf.String()
			todoHTML := todoBuf.String()

			Expect(greetingHTML).To(ContainSubstring("Hello, Team!"))
			Expect(todoHTML).To(ContainSubstring("Review code"))
			Expect(todoHTML).To(ContainSubstring("Deploy app"))
		})

		It("2.5 should handle special characters in content", func() {
			By("Creating component with special characters")
			specialName := "<script>alert('xss')</script>"
			component := Greeting(specialName)

			By("Rendering the component")
			var buf bytes.Buffer
			err := component.Render(context.Background(), &buf)
			Expect(err).ToNot(HaveOccurred())

			By("Verifying HTML escaping is applied")
			html := buf.String()
			Expect(html).ToNot(ContainSubstring("<script>"))
			Expect(html).To(ContainSubstring("&lt;script&gt;"))
		})

		It("2.6 should render empty table gracefully", func() {
			By("Creating empty data table")
			component := DataTable([]TableRow{})

			By("Rendering the table")
			var buf bytes.Buffer
			err := component.Render(context.Background(), &buf)
			Expect(err).ToNot(HaveOccurred())

			By("Verifying table structure exists")
			html := buf.String()
			Expect(html).To(ContainSubstring("<table"))
			Expect(html).To(ContainSubstring("<thead>"))
			Expect(html).To(ContainSubstring("<tbody>"))

			By("Verifying no data rows")
			bodyRows := strings.Count(html, "<tbody>")
			Expect(bodyRows).To(Equal(1))
		})
	})
})
