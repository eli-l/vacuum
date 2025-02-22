// Copyright 2023 Princess B33f Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package owasp

import (
	"github.com/daveshanley/vacuum/model"
	"github.com/daveshanley/vacuum/utils"
	"github.com/pb33f/doctor/model/high/base"
	"gopkg.in/yaml.v3"
	"slices"
)

type AdditionalPropertiesConstrained struct{}

// GetSchema returns a model.RuleFunctionSchema defining the schema of the DefineError rule.
func (ad AdditionalPropertiesConstrained) GetSchema() model.RuleFunctionSchema {
	return model.RuleFunctionSchema{Name: "additional_properties_constrained"}
}

// RunRule will execute the DefineError rule, based on supplied context and a supplied []*yaml.Node slice.
func (ad AdditionalPropertiesConstrained) RunRule(_ []*yaml.Node, context model.RuleFunctionContext) []model.RuleFunctionResult {

	var results []model.RuleFunctionResult

	if context.DrDocument == nil {
		return results
	}

	for _, schema := range context.DrDocument.Schemas {
		if slices.Contains(schema.Value.Type, "object") {
			if schema.Value.AdditionalProperties != nil {

				node := schema.Value.GoLow().Type.KeyNode
				result := model.RuleFunctionResult{
					Message: utils.SuppliedOrDefault(context.Rule.Message,
						"schema should also define `maxProperties` when `additionalProperties` is an object"),
					StartNode: node,
					EndNode:   node,
					Path:      schema.GenerateJSONPath(),
					Rule:      context.Rule,
				}

				if schema.Value.AdditionalProperties.IsA() {
					if schema.Value.MaxProperties == nil {

						schema.AddRuleFunctionResult(base.ConvertRuleResult(&result))
						results = append(results, result)
						continue
					}
				}
				if schema.Value.AdditionalProperties.IsB() && schema.Value.AdditionalProperties.B {

					if schema.Value.MaxProperties == nil {

						schema.AddRuleFunctionResult(base.ConvertRuleResult(&result))
						results = append(results, result)
						continue
					}
				}

			}
		}
	}
	return results
}
