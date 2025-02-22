package openapi

import (
	"github.com/daveshanley/vacuum/model"
	"github.com/pb33f/libopenapi/datamodel"
	"github.com/pb33f/libopenapi/index"
	"github.com/pb33f/libopenapi/utils"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestUnusedComponent_GetSchema(t *testing.T) {
	def := UnusedComponent{}
	assert.Equal(t, "unused_component", def.GetSchema().Name)
}

func TestUnusedComponent_RunRule(t *testing.T) {
	def := UnusedComponent{}
	res := def.RunRule(nil, model.RuleFunctionContext{})
	assert.Len(t, res, 0)
}

func TestUnusedComponent_RunRule_Success(t *testing.T) {

	yml := `paths:
  /naughty/{puppy}:
    parameters:
      - $ref: '#/components/parameters/Chewy'
    get:
      responses:
      "200":
        description: The naughty pup
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/Puppy'
components:
  schemas:
    Puppy:
      type: string
      description: pup
  parameters:
    Chewy:
      description: chewy
      in: query
      name: chewy`

	path := "$"

	var rootNode yaml.Node
	mErr := yaml.Unmarshal([]byte(yml), &rootNode)
	assert.NoError(t, mErr)

	nodes, _ := utils.FindNodes([]byte(yml), path)

	rule := buildOpenApiTestRuleAction(path, "unused_component", "", nil)
	ctx := buildOpenApiTestContext(model.CastToRuleAction(rule.Then), nil)
	config := index.CreateOpenAPIIndexConfig()
	ctx.Index = index.NewSpecIndexWithConfig(&rootNode, config)

	def := UnusedComponent{}
	res := def.RunRule(nodes, ctx)

	assert.Len(t, res, 0)
}

func TestUnusedComponent_RunRule_SuccessSwaggerSecurity(t *testing.T) {

	yml := `swagger: 2.0
securityDefinitions:
  basicAuth:
    type: basic
  sessionAuth:
    type: apiKey
    in: header
    name: X-API-Key
paths:
  "/store/inventory":
    get:
      security:
        - basicAuth: []
  "/store/inventory/doSomething":
    get:
      security:
        - sessionAuth: []`

	path := "$"

	var rootNode yaml.Node
	mErr := yaml.Unmarshal([]byte(yml), &rootNode)
	assert.NoError(t, mErr)

	nodes, _ := utils.FindNodes([]byte(yml), path)

	rule := buildOpenApiTestRuleAction(path, "unused_component", "", nil)
	ctx := buildOpenApiTestContext(model.CastToRuleAction(rule.Then), nil)
	config := index.CreateOpenAPIIndexConfig()
	ctx.Index = index.NewSpecIndexWithConfig(&rootNode, config)
	info, _ := datamodel.ExtractSpecInfo([]byte(yml))
	ctx.SpecInfo = info
	ctx.Rule = &rule

	def := UnusedComponent{}
	res := def.RunRule(nodes, ctx)

	assert.Len(t, res, 0)
}

func TestUnusedComponent_RunRule_SuccessOpenAPISecurity(t *testing.T) {

	yml := `openapi: 3.0.1
info:
  description: A test spec with a security def that is not a ref!
security:
  - SomeSecurity: []
components:
  securitySchemes:
    SomeSecurity:
      description: A secure way to do things and stuff.`
	path := "$"

	var rootNode yaml.Node
	mErr := yaml.Unmarshal([]byte(yml), &rootNode)
	assert.NoError(t, mErr)

	nodes, _ := utils.FindNodes([]byte(yml), path)

	rule := buildOpenApiTestRuleAction(path, "unused_component", "", nil)
	ctx := buildOpenApiTestContext(model.CastToRuleAction(rule.Then), nil)
	config := index.CreateOpenAPIIndexConfig()
	ctx.Index = index.NewSpecIndexWithConfig(&rootNode, config)
	info, _ := datamodel.ExtractSpecInfo([]byte(yml))
	ctx.SpecInfo = info

	def := UnusedComponent{}
	res := def.RunRule(nodes, ctx)

	assert.Len(t, res, 0)
}

func TestUnusedComponent_RunRule_Success_Fail_TwoMissing_Two_Undefined(t *testing.T) {

	yml := `parameters:
  Chewy:
    description: chewy
    in: query
    name: chewy
paths:
  /naughty/{puppy}:
    parameters:
      - $ref: '#/parameters/Nothing'
    get:
      responses:
      "200":
        description: The naughty pup
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/Cupcakes_And_Sugar'
components:
  schemas:
    Puppy:
      type: string
      description: pup
    Kitty:
      $ref: '#/components/schemas/Puppy' `

	path := "$"

	var rootNode yaml.Node
	mErr := yaml.Unmarshal([]byte(yml), &rootNode)
	assert.NoError(t, mErr)

	nodes, _ := utils.FindNodes([]byte(yml), path)

	rule := buildOpenApiTestRuleAction(path, "unused_component", "", nil)
	ctx := buildOpenApiTestContext(model.CastToRuleAction(rule.Then), nil)
	config := index.CreateOpenAPIIndexConfig()
	ctx.Index = index.NewSpecIndexWithConfig(&rootNode, config)

	def := UnusedComponent{}
	res := def.RunRule(nodes, ctx)

	assert.Len(t, res, 4)
}

func TestUnusedComponent_RunRule_Success_Fail_Four_Undefined(t *testing.T) {

	yml := `paths:
  /naughty/{puppy}:
    parameters:
      - $ref: '#/components/parameters/Chewy'
    get:
      responses:
      "200":
        description: The naughty pup
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/Puppy'
components:
  responses:
    Chappy:
      type: string
      description: Chappy
  schemas:  
    Chippy:
      type: string
      description: chippy
    Puppy:
      type: string
      description: pup
    Kitty:
      $ref: '#/components/schemas/Puppy'
  parameters:
    Minty:
      description: minty
      in: header
      name: minty
    Chewy:
      description: chewy
      in: query
      name: chewy`

	path := "$"

	var rootNode yaml.Node
	mErr := yaml.Unmarshal([]byte(yml), &rootNode)
	assert.NoError(t, mErr)

	nodes, _ := utils.FindNodes([]byte(yml), path)

	rule := buildOpenApiTestRuleAction(path, "unused_component", "", nil)
	ctx := buildOpenApiTestContext(model.CastToRuleAction(rule.Then), nil)
	config := index.CreateOpenAPIIndexConfig()
	ctx.Index = index.NewSpecIndexWithConfig(&rootNode, config)

	def := UnusedComponent{}
	res := def.RunRule(nodes, ctx)

	assert.Len(t, res, 4)
}

func TestUnusedComponent_RunRule_Success_PolymorphicCheck(t *testing.T) {

	yml := `paths:
  /naughty/{puppy}:
    get:
      responses:
      "200":
        description: The naughty pup
        content:
          application/json:
            schema:
              oneOf:
                - $ref: '#/components/schemas/Puppy'
      "404":
        description: The naughty kitty
        content:
          application/json:
            schema:
              anyOf:
                - $ref: '#/components/schemas/Kitty'
      "500":
        description: The naughty bunny
        content:
          application/json:
            schema:
              allOf:
                - $ref: '#/components/schemas/Bunny'
components:
  schemas:
    Puppy:
      type: string
      description: pup
    Kitty:
      type: string
      description: kitty
    Bunny:
      type: string
      description: bunny`

	path := "$"

	var rootNode yaml.Node
	mErr := yaml.Unmarshal([]byte(yml), &rootNode)
	assert.NoError(t, mErr)

	nodes, _ := utils.FindNodes([]byte(yml), path)

	rule := buildOpenApiTestRuleAction(path, "unused_component", "", nil)
	ctx := buildOpenApiTestContext(model.CastToRuleAction(rule.Then), nil)
	config := index.CreateOpenAPIIndexConfig()
	ctx.Index = index.NewSpecIndexWithConfig(&rootNode, config)

	def := UnusedComponent{}
	res := def.RunRule(nodes, ctx)

	assert.Len(t, res, 0)
}
