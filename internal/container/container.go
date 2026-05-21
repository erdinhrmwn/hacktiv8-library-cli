package container

import (
	"context"
	"fmt"
)

type Container struct {
	AuthMenu    *AuthMenu
	VisitorMenu *VisitorMenu
	StaffMenu   *StaffMenu
}

func New() *Container {
	return &Container{
		AuthMenu:    &AuthMenu{},
		VisitorMenu: &VisitorMenu{},
		StaffMenu:   &StaffMenu{},
	}
}

func (c *Container) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("\n👋 Selamat tinggal!")
			return
		default:
		}

		user := c.AuthMenu.Main(ctx)
		if user == nil {
			fmt.Println("\n👋 Selamat tinggal!")
			return
		}

		switch user.Role {
		case "staff":
			c.StaffMenu.Dashboard(ctx)
		case "visitor":
			c.VisitorMenu.Dashboard(ctx)
		}
	}
}
