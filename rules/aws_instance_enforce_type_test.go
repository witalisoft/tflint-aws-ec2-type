package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsInstanceEnforceType(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Config   string
		Expected helper.Issues
	}{
		{
			Name: "wrong instance",
			Content: `
resource "aws_instance" "web" {
    instance_type = "t2.nano"
}`,
			Config: `
rule "aws_instance_enforce_type" {
	enabled = true
	types   = ["t2.micro", "t3.micro"]
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsInstanceEnforceTypeRule(),
					Message: "wrong instance type t2.nano, should be on of: [t2.micro t3.micro]",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 21},
						End:      hcl.Pos{Line: 3, Column: 30},
					},
				},
			},
		},
		{
			Name: "no issue",
			Content: `
resource "aws_instance" "web" {
    instance_type = "t2.micro"
}`,
			Config: `
rule "aws_instance_enforce_type" {
	enabled = true
	types   = ["t2.micro", "t3.micro"]
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsInstanceEnforceTypeRule()

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content, ".tflint.hcl": tc.Config})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssues(t, tc.Expected, runner.Issues)
		})
	}
}
