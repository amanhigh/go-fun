package libraries

import (
	"fmt"

	"github.com/andygrunwald/go-jira"
)

func JiraFun() {
	var err error
	var jiraClient *jira.Client
	var issue *jira.Issue

	//tp := jira.BasicAuthTransport{
	//	Username: "aman",
	//	Password: "token",
	//}

	// Pass tp.Client() for issue creation
	if jiraClient, err = jira.NewClient(nil, "https://issues.apache.org/jira/"); err == nil {
		if issue, _, err = jiraClient.Issue.Get("MESOS-3325", nil); err == nil {
			fmt.Printf("IssueId: %s\n", issue.Key)
			fmt.Printf("%s: %+v\n", issue.Key, issue.Fields.Summary)
			fmt.Printf("Type: %s\n", issue.Fields.Type.Name)
			fmt.Printf("Priority: %s\n", issue.Fields.Priority.Name)
		}
	}

	/* Create */
	//newIssue := jira.Issue{
	//	Fields: &jira.IssueFields{
	//		Description: "Demo",
	//		Type: jira.IssueType{
	//			Name: "Task",
	//		},
	//		Project: jira.Project{
	//			Key: "DEMO_PROJECT",
	//		},
	//		Summary: "Demo Issue",
	//	},
	//}
	//if issue, _, err = jiraClient.Issue.Create(&newIssue); err == nil {
	//	fmt.Println(issue)
	//}
	//fmt.Println(issue.ID)

	fmt.Println("Error: ", err)

}
