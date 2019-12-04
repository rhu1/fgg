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

		case fg.ITypeLit:
			iTypeLit, err := c.convertITypeLit(decl)
			if err != nil {
				return err
			}
			c.fggProg.ds = append(c.fggProg.ds, iTypeLit)

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

	expr, err := c.convertExpr(c.fgProg.GetMain())
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
	for _, f := range s.GetFieldDecls() {
		fd, err := c.convertFieldDecl(f)
		if err != nil {
			return STypeLit{}, err
		}
		fieldDecls = append(fieldDecls, fd)
	}

	return STypeLit{t: typeName, psi: typeFormals, fds: fieldDecls}, nil
}

func (c *fg2fgg) convertITypeLit(i fg.ITypeLit) (ITypeLit, error) {
	typeName, typeFormals := c.convertType(i.GetType())

	var specs []Spec
	for _, s := range i.Specs() {
		sig := s.(fg.Sig)
		var paramDecls []ParamDecl
		for _, p := range sig.MethodParams() {
			pd, err := c.convertParamDecl(p)
			if err != nil {
				return ITypeLit{}, nil
			}
			paramDecls = append(paramDecls, pd)
		}
		retTypeName, _ := c.convertType(sig.ReturnType())

		specs = append(specs, Sig{
			m:   Name(sig.MethodName()),
			psi: TFormals{tfs: nil},
			pds: paramDecls,
			u:   TName{t: retTypeName},
		})
	}

	return ITypeLit{t: typeName, psi: typeFormals, ss: specs}, nil
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
	recvTypeName, recvTypeFormals := c.convertType(md.GetReceiver().GetType())

	var paramDecls []ParamDecl
	for _, p := range md.GetParamDecls() {
		pd, err := c.convertParamDecl(p)
		if err != nil {
			return MDecl{}, err
		}
		paramDecls = append(paramDecls, pd)
	}

	retTypeName, _ := c.convertType(md.GetReturn())
	methImpl, err := c.convertExpr(md.GetBody())
	if err != nil {
		return MDecl{}, err
	}

	return MDecl{
		x_recv:   md.GetReceiver().GetName(),
		t_recv:   recvTypeName,
		psi_recv: recvTypeFormals,
		m:        Name(md.GetName()),
		psi:      TFormals{}, // empty parameter
		pds:      paramDecls,
		u:        TName{t: retTypeName},
		e:        methImpl,
	}, nil
}

func (c *fg2fgg) convertExpr(expr base.Expr) (Expr, error) {
	switch expr := expr.(type) {
	case fg.Variable:
		return Variable{id: expr.String()}, nil

	case fg.StructLit:
		sLitExpr, err := c.convertStructLit(expr)
		if err != nil {
			return nil, err
		}
		return sLitExpr, nil

	case fg.Call:
		callExpr, err := c.convertCall(expr)
		if err != nil {
			return nil, err
		}
		return callExpr, nil

	case fg.Select:
		selExpr, err := c.convertExpr(expr.Expr())
		if err != nil {
			return nil, err
		}
		return Select{e: selExpr, f: Name(expr.FieldName())}, nil

	case fg.Assert:
		assertExpr, err := c.convertExpr(expr.Expr())
		if err != nil {
			return nil, err
		}
		assType, _ := c.convertType(expr.AssertType())
		return Assert{e: assertExpr, u: TName{t: assType}}, nil
	}

	return nil, fmt.Errorf("unknown expression type: %T", expr)
}

func (c *fg2fgg) convertStructLit(sLit fg.StructLit) (StructLit, error) {
	structType, _ := c.convertType(sLit.Type())

	var es []Expr
	for _, expr := range sLit.FieldExprs() {
		fieldExpr, err := c.convertExpr(expr)
		if err != nil {
			return StructLit{}, err
		}
		es = append(es, fieldExpr)
	}

	return StructLit{u: TName{t: structType}, es: es}, nil
}

func (c *fg2fgg) convertCall(call fg.Call) (Call, error) {
	e, err := c.convertExpr(call.Expr())
	if err != nil {
		return Call{}, err
	}

	var args []Expr
	for _, arg := range call.Args() {
		argExpr, err := c.convertExpr(arg)
		if err != nil {
			return Call{}, err
		}
		args = append(args, argExpr)
	}

	return Call{e: e, m: Name(call.MethodName()), args: args}, nil
}
