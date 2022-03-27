package rules

import (
	"fmt"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsInstanceEnforceTypeRule checks whether ...
type AwsInstanceEnforceTypeRule struct{}

// AwsInstanceEnforceTypeRuleConfig is the config structure for the AwsInstanceEnforceTypeRule rule
type AwsInstanceEnforceTypeRuleConfig struct {
	Types  []string `hcl:"types"`
	Remain hcl.Body `hcl:",remain"`
}

// NewAwsInstanceEnforceTypeRule returns a new rule
func NewAwsInstanceEnforceTypeRule() *AwsInstanceEnforceTypeRule {
	return &AwsInstanceEnforceTypeRule{}
}

// Name returns the rule name
func (r *AwsInstanceEnforceTypeRule) Name() string {
	return "aws_instance_enforce_type"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsInstanceEnforceTypeRule) Enabled() bool {
	return false
}

// Severity returns the rule severity
func (r *AwsInstanceEnforceTypeRule) Severity() string {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsInstanceEnforceTypeRule) Link() string {
	return ""
}

// Check checks whether ...
func (r *AwsInstanceEnforceTypeRule) Check(runner tflint.Runner) error {
	return runner.WalkResourceAttributes("aws_instance", "instance_type", func(attribute *hcl.Attribute) error {
		var instanceType string
		err := runner.EvaluateExpr(attribute.Expr, &instanceType, nil)

		config := AwsInstanceEnforceTypeRuleConfig{}
		if err := runner.DecodeRuleConfig(r.Name(), &config); err != nil {
			return err
		}
		return runner.EnsureNoError(err, func() error {
			if !checkContains(config.Types, instanceType) {
				return runner.EmitIssueOnExpr(
					r,
					fmt.Sprintf("wrong instance type %s, should be one of: %v", instanceType, config.Types),
					attribute.Expr,
				)
			} else {
				return nil
			}
		})
	})
}

func checkContains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
