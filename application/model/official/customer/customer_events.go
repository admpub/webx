package customer

import "github.com/admpub/webx/application/dbschema"

type Event func(*dbschema.OfficialCustomer) error

type Events []Event

func (e Events) Execute(customer *dbschema.OfficialCustomer) error {
	for _, fn := range e {
		if err := fn(customer); err != nil {
			return err
		}
	}
	return nil
}

var (
	onSignIn  Events
	onSignOut Events
)

func OnSignIn(f Event) {
	onSignIn = append(onSignIn, f)
}

func FireSignIn(customer *dbschema.OfficialCustomer) error {
	return onSignIn.Execute(customer)
}

func OnSignOut(f Event) {
	onSignOut = append(onSignOut, f)
}

func FireSignOut(customer *dbschema.OfficialCustomer) error {
	return onSignOut.Execute(customer)
}
