package customer

import "github.com/admpub/webx/application/dbschema"

type ComplaintExt struct {
	*dbschema.OfficialCommonComplaint
	Customer       *dbschema.OfficialCustomer `db:"-,relation=id:customer_id"`
	TypeName       string                     `db:"-"`
	TargetTypeName string                     `db:"-"`
	urlFormat      func(*dbschema.OfficialCommonComplaint) string
}

func (c *ComplaintExt) SetURLFormat(fn func(*dbschema.OfficialCommonComplaint) string) *ComplaintExt {
	c.urlFormat = fn
	return c
}

func (c *ComplaintExt) URLFormat() string {
	if c.urlFormat == nil {
		return ``
	}
	return c.urlFormat(c.OfficialCommonComplaint)
}
