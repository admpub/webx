package agent

import (
	xMW "github.com/admpub/webx/application/middleware"
	modelAgent "github.com/admpub/webx/application/model/official/agent"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

// AgentApply 申请成为代理
func AgentApply(c echo.Context) error {
	customer := xMW.Customer(c)
	m := modelAgent.NewAgentProfile(c)
	cond := db.Cond{`customer_id`: customer.Id}
	err := m.Get(nil, cond)
	if err != nil {
		if err != db.ErrNoMoreRows {
			return err
		}
		err = nil
	}
	rvM := modelAgent.NewAgentRecv(c)
	rvM.Get(nil, cond)
	isNew := m.CustomerId <= 0
	if c.IsPost() {
		if !isNew && m.Status != modelAgent.ProfileStatusReject {
			return c.E(`您已经申请过了，当前状态: %s`, modelAgent.ProfileStatus.Get(m.Status, c.T(`未知`)))
		}
		m.CustomerId = customer.Id
		c.Begin()
		rvM.RecvMoneyMethod = c.Formx(`recvMoneyMethod`).String()
		rvM.RecvMoneyBranch = c.Formx(`recvMoneyBranch`).String()
		rvM.RecvMoneyOwner = c.Formx(`recvMoneyOwner`).String()
		rvM.RecvMoneyAccount = c.Formx(`recvMoneyAccount`).String()
		rvM.CustomerId = m.CustomerId
		if isNew {
			_, err = rvM.Add()
		} else {
			err = rvM.Edit(nil, cond)
		}
		if err != nil {
			c.Rollback()
			return err
		}

		m.ApplyLevel = c.Formx(`applyLevel`).Uint()
		m.Remark = c.Formx(`remark`).String()
		m.Use(rvM.Trans())
		if isNew {
			_, err = m.Add()
		} else {
			err = m.Edit(nil, cond)
		}
		if err != nil {
			c.Rollback()
			return err
		}
		c.Commit()
		return c.Redirect(xMW.URLFor(`/user/agent/apply`))
	}
	editable := true
	var statusName, statusColor, statusDescription string
	if !isNew {
		if m.Status != modelAgent.ProfileStatusReject {
			editable = false
		}
		item := modelAgent.ProfileStatus.GetItem(m.Status)
		if item != nil {
			statusName = item.V
			statusColor = item.H.String(`color`)
			statusDescription = item.H.String(`applyDescrition`)
		}
		echo.StructToForm(c, m.OfficialCustomerAgentProfile, ``, echo.LowerCaseFirstLetter)
		if rvM.CustomerId > 0 {
			echo.StructToForm(c, rvM.OfficialCustomerAgentRecv, ``, echo.LowerCaseFirstLetter)
		}
	}
	c.Set(`activeURL`, `/user/agent`)
	c.Set(`recvMoneyMethods`, modelAgent.RecvMoneyMethod.Slice())
	levelM := modelAgent.NewAgentLevel(c)
	levelM.ListByOffset(nil, func(r db.Result) db.Result {
		return r.OrderBy(`id`)
	}, 0, -1)
	agentLevelList := levelM.Objects()
	c.Set(`agentLevelList`, agentLevelList)
	c.Set(`editable`, editable)
	c.Set(`status`, m.Status)
	c.Set(`statusName`, statusName)
	c.Set(`statusColor`, statusColor)
	c.Set(`statusDescription`, statusDescription)
	return c.Render(`/user/agent/apply`, err)
}
