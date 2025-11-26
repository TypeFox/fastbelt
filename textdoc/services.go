// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package textdoc

// TextdocSrvCont is an interface for service containers which include the textdoc services.
type TextdocSrvCont interface {
	Textdoc() *TextdocSrv
}

// TextdocSrvContBlock is used to define a service container satisfying TextdocSrvCont.
type TextdocSrvContBlock struct {
	textdoc TextdocSrv
}

func (b *TextdocSrvContBlock) Textdoc() *TextdocSrv {
	return &b.textdoc
}

// TextdocSrv contains the services for the textdoc package.
type TextdocSrv struct {
	Store Store
}

// CreateDefaultServices creates the default services for the textdoc package.
// If the services are already set, they are not overwritten.
func CreateDefaultServices(c TextdocSrvCont) {
	s := c.Textdoc()
	if s.Store == nil {
		s.Store = NewDefaultStore()
	}
}
