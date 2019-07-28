package fgg

import (
	"fmt"

	"github.com/rhu1/fgg/base"
	"github.com/rhu1/fgg/fg"
)

type fg2fgg struct {
	fgProg  fg.FGProgram
	fggProg FGGProgram
}

// FromFG converts a FG program prog into a FGG program
// with empty type parameters
func FromFG(prog fg.FGProgram) (FGGProgram, error) {
	c := &fg2fgg{fgProg: prog}
	if err := c.convert(); err != nil {
		return FGGProgram{}, err
	}
	return c.fggProg, nil
}

func (c *fg2fgg) convert() error {
	// convert toplevel declarations
	for _, decl := range c.fgProg.GetDecls() {
		switch decl := decl.(type) {
		case fg.STypeLit:
			sTypeLit, err := c.convertSTypeLit(decl)
			if err != nil {
				return err
			}
			c.fggProg.ds = append(c.fggProg.ds, sTypeLit)

		case fg.MDecl:
			mDecl, err := c.convertMDecl(decl)
			if err != nil {
				return err
			}
			c.fggProg.ds = append(c.fggProg.ds, mDecl)

		default:
			return fmt.Errorf("unknown declaration type: %T", decl)
		}
	}

	expr, err := c.convertExpr(c.fgProg.GetExpr())
	if err != nil {
		return err
	}
	c.fggProg.e = expr

	return nil
}

// convertType converts a plain type to a parameterised type
func (c *fg2fgg) convertType(t fg.Type) (Name, TFormals) {
	return Name(t.String()), TFormals{tfs: nil} // 0 formal parameters
}

func (c *fg2fgg) convertSTypeLit(s fg.STypeLit) (STypeLit, error) {
	typeName, typeFormals := c.convertType(s.GetType())

	var fieldDecls []FieldDecl
	for _, f := range s.Fields() {
		fd, err := c.convertFieldDecl(f)
		if err != nil {
			return STypeLit{}, err
		}
		fieldDecls = append(fieldDecls, fd)
	}

	return STypeLit{t: typeName, psi: typeFormals, fds: fieldDecls}, nil
}

func (c *fg2fgg) convertFieldDecl(fd fg.FieldDecl) (FieldDecl, error) {
	typeName, _ := c.convertType(fd.GetType())
	return FieldDecl{f: fd.GetName(), u: TName{t: typeName}}, nil
}

func (c *fg2fgg) convertParamDecl(pd fg.ParamDecl) (ParamDecl, error) {
	typeName, _ := c.convertType(pd.GetType())
	return ParamDecl{x: pd.GetName(), u: TName{t: typeName}}, nil
}

func (c *fg2fgg) convertMDecl(md fg.MDecl) (MDecl, error) {
	recvTypeName, recvTypeFormals := c.convertType(md.Receiver().GetType())

	var paramDecls []ParamDecl
	for _, p := range md.MethodParams() {
		pd, err := c.convertParamDecl(p)
		if err != nil {
			return MDecl{}, err
		}
		paramDecls = append(paramDecls, pd)
	}
	retTypeName, _ := c.convertType(md.ReturnType())
	methImpl, err := c.convertExpr(md.Impl())
	if err != nil {
		return MDecl{}, err
	}

	return MDecl{
		x_recv:   md.Receiver().GetName(),
		t_recv:   recvTypeName,
		psi_recv: recvTypeFormals,
		m:        Name(md.MethodName()),
		pds:      paramDecls,
		u:        TName{t: retTypeName},
		e:        methImpl,
	}, nil
}

func (c *fg2fgg) convertExpr(e base.Expr) (Expr, error) {
	return Variable{id: "_TODO_"}, nil
	/*
		switch e := e.(type) {
		case fg.Variable:
			return Variable{id: e.String()}, nil
		case fg.StructLit:
			return StructLit{}, nil
		case fg.Call:
			// TODO: add type param
			return Call{}, nil
		case fg.Select:
			return Select{}, nil
		case fg.Assert:
			return Assert{}, nil
		}
		return nil, fmt.Errorf("unknown expression type: %T", e)
	*/
}
