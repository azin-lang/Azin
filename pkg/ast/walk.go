package ast

// AdjustPositions adds delta to all token positions in the program.
func AdjustPositions(program *Program, delta uint32) {
	for _, stmt := range program.Statements {
		adjustStmt(stmt, delta)
	}
}

func adjustStmt(stmt Stmt, delta uint32) {
	switch s := stmt.(type) {
	case *BadStmt:
		s.Token.Position.Offset += delta
	case *VarStmt:
		s.Token.Position.Offset += delta
		adjustExpr(s.Name, delta)
		if s.SynType != nil {
			adjustExpr(s.SynType, delta)
		}
		adjustExpr(s.Value, delta)
	case *AssignmentStmt:
		s.Token.Position.Offset += delta
		adjustExpr(s.Left, delta)
		adjustExpr(s.Value, delta)
	case *StructStmt:
		s.Token.Position.Offset += delta
		adjustExpr(s.Name, delta)
		for _, f := range s.Fields {
			adjustExpr(f.Name, delta)
			if f.SynType != nil {
				adjustExpr(f.SynType, delta)
			}
		}
	case *EnumStmt:
		s.Token.Position.Offset += delta
		adjustExpr(s.Name, delta)
		for _, v := range s.Variants {
			adjustExpr(v, delta)
		}
	case *FuncStmt:
		s.Token.Position.Offset += delta
		adjustExpr(s.Name, delta)
		for _, p := range s.Params {
			adjustExpr(p.Name, delta)
			if p.SynType != nil {
				adjustExpr(p.SynType, delta)
			}
		}
		if s.SynReturnType != nil {
			adjustExpr(s.SynReturnType, delta)
		}
		for _, bodyStmt := range s.Body {
			adjustStmt(bodyStmt, delta)
		}
	case *ReturnStmt:
		s.Token.Position.Offset += delta
		adjustExpr(s.Value, delta)
	case *IfStmt:
		s.Token.Position.Offset += delta
		adjustExpr(s.Condition, delta)
		for _, bodyStmt := range s.Then {
			adjustStmt(bodyStmt, delta)
		}
		for _, bodyStmt := range s.Else {
			adjustStmt(bodyStmt, delta)
		}
	case *LoopStmt:
		s.Token.Position.Offset += delta
		for _, bodyStmt := range s.Body {
			adjustStmt(bodyStmt, delta)
		}
	case *StopStmt:
		s.Token.Position.Offset += delta
	case *DeferStmt:
		s.Token.Position.Offset += delta
		adjustExpr(s.Call, delta)
	case *ImportCStmt:
		s.Token.Position.Offset += delta
		adjustExpr(s.Path, delta)
	case *ImportStmt:
		s.Token.Position.Offset += delta
		adjustExpr(s.Path, delta)
	case *ExpressionStmt:
		s.Token.Position.Offset += delta
		adjustExpr(s.Expression, delta)
	}
}

func adjustExpr(expr Expr, delta uint32) {
	if expr == nil {
		return
	}
	switch e := expr.(type) {
	case *BadExpr:
		e.Token.Position.Offset += delta
	case *Identifier:
		e.Token.Position.Offset += delta
	case *IntegerLiteral:
		e.Token.Position.Offset += delta
	case *FloatLiteral:
		e.Token.Position.Offset += delta
	case *StringLiteral:
		e.Token.Position.Offset += delta
	case *CharacterLiteral:
		e.Token.Position.Offset += delta
	case *BooleanLiteral:
		e.Token.Position.Offset += delta
	case *CallExpr:
		adjustExpr(e.Callee, delta)
		for _, arg := range e.Args {
			adjustExpr(arg, delta)
		}
	case *BinaryExpr:
		adjustExpr(e.Left, delta)
		e.Operator.Position.Offset += delta
		adjustExpr(e.Right, delta)
	case *MemberExpr:
		adjustExpr(e.Object, delta)
		adjustExpr(e.Property, delta)
	}
}
