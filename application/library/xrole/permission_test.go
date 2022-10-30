package xrole

import (
	"testing"

	"github.com/admpub/nging/v5/application/registry/navigate"
	"github.com/admpub/webx/application/dbschema"
	"github.com/stretchr/testify/assert"
	"github.com/webx-top/com"
	"github.com/webx-top/echo/defaults"
)

func TestRolePermission(t *testing.T) {
	p := NewRolePermission()
	roleList := []*CustomerRoleWithPermissions{
		{
			OfficialCustomerRole: &dbschema.OfficialCustomerRole{
				Id:   1,
				Name: `test role`,
			},
			Permissions: []*dbschema.OfficialCustomerRolePermission{
				{
					RoleId:     1,
					Type:       `page`,
					Permission: `test/*`,
				},
			},
		},
	}
	p.Init(roleList)
	navList := &navigate.List{
		{
			Display: true,
			Name:    `test group`,
			Action:  `test`,
			Children: &navigate.List{
				{
					Display: true,
					Name:    `test item 1`,
					Action:  `1`,
				},
				{
					Display: true,
					Name:    `test subgroup 1`,
					Action:  `subtest`,
					Children: &navigate.List{
						{
							Display: true,
							Name:    `test subgroup item 1`,
							Action:  `1`,
						},
					},
				},
			},
		},
		{
			Display: true,
			Name:    `deny group`,
			Action:  `deny`,
			Children: &navigate.List{
				{
					Display: true,
					Name:    `deny item 1`,
					Action:  `1`,
				},
			},
		},
	}
	ctx := defaults.NewMockContext()
	actualNavList := p.FilterNavigate(ctx, navList)
	com.Dump(actualNavList)

	excepted := navigate.List{
		{
			Display: true,
			Name:    `test group`,
			Action:  `test`,
			Children: &navigate.List{
				{
					Display:  true,
					Name:     `test item 1`,
					Action:   `1`,
					Children: &navigate.List{},
				},
				{
					Display: true,
					Name:    `test subgroup 1`,
					Action:  `subtest`,
					Children: &navigate.List{
						{
							Display:  true,
							Name:     `test subgroup item 1`,
							Action:   `1`,
							Children: &navigate.List{},
						},
					},
				},
			},
		},
	}
	assert.Equal(t, excepted, actualNavList)
}
