package compilerlabs

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type RandomProgramGenerator struct {
	r *rand.Rand

	varNames []string

	declaredVars []string

	mathOps    []string
	compareOps []string
	logicOps   []string
}

// NewRandomProgramGenerator создает генератор.
func NewRandomProgramGenerator(seed int64) *RandomProgramGenerator {
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	return &RandomProgramGenerator{
		r: rand.New(rand.NewSource(seed)),
		varNames: []string{
			"x", "y", "z", "alpha", "beta", "count", "total", "index", "sum",
		},
		declaredVars: nil,
		mathOps:      []string{"+", "-", "*", "/"},
		compareOps:   []string{"==", "!=", "<", ">", "<=", ">="},
		logicOps:     []string{"&&", "||"},
	}
}

// Generate генерирует случайную программу.
// statementCount — количество инструкций на верхнем уровне (по умолчанию 10, если <= 0).
func (g *RandomProgramGenerator) Generate(statementCount int) string {
	if statementCount <= 0 {
		statementCount = 10
	}

	g.declaredVars = g.declaredVars[:0]
	var b strings.Builder

	// Обязательно объявляем пару переменных в начале.
	for i := 0; i < 3; i++ {
		b.WriteString(g.generateVarDeclaration(0))
		b.WriteByte('\n')
	}

	g.generateBlock(&b, statementCount, 0)
	return b.String()
}

func (g *RandomProgramGenerator) generateBlock(b *strings.Builder, count int, indentLevel int) {
	indent := strings.Repeat(" ", indentLevel*4)

	for i := 0; i < count; i++ {
		statementType := g.r.Intn(5)

		if indentLevel > 2 && statementType > 2 {
			statementType = g.r.Intn(3)
		}

		switch statementType {
		case 0:
			b.WriteString(g.generateVarDeclaration(indentLevel))
			b.WriteByte('\n')

		case 1:
			if len(g.declaredVars) > 0 {
				fmt.Fprintf(b, "%s%s = %s;\n", indent, g.getRandomVar(), g.generateExpression())
			} else {
				b.WriteString(g.generateVarDeclaration(indentLevel))
				b.WriteByte('\n')
			}

		case 2:
			fmt.Fprintf(b, "%sprint %s;\n", indent, g.generateExpression())

		case 3:
			fmt.Fprintf(b, "%sif (%s) {\n", indent, g.generateCondition())
			g.generateBlock(b, g.r.Intn(3)+1, indentLevel+1) // 1..3

			if g.r.Float64() > 0.5 { // 50% else
				fmt.Fprintf(b, "%s} else {\n", indent)
				g.generateBlock(b, g.r.Intn(2)+1, indentLevel+1) // 1..2
			}
			fmt.Fprintf(b, "%s}\n", indent)

		case 4:
			fmt.Fprintf(b, "%swhile (%s) {\n", indent, g.generateCondition())
			g.generateBlock(b, g.r.Intn(3)+1, indentLevel+1) // 1..3
			fmt.Fprintf(b, "%s}\n", indent)
		}
	}
}

func (g *RandomProgramGenerator) generateVarDeclaration(indentLevel int) string {
	indent := strings.Repeat(" ", indentLevel*4)

	// Берем случайное имя; если уже есть, просто "переобъявим" (как в C# версии)
	varName := g.varNames[g.r.Intn(len(g.varNames))]
	if !contains(g.declaredVars, varName) {
		g.declaredVars = append(g.declaredVars, varName)
	}

	return fmt.Sprintf("%svar %s = %s;", indent, varName, g.generateExpression())
}

func (g *RandomProgramGenerator) generateExpression() string {
	// Простые числа или переменные
	if g.r.Float64() > 0.6 || len(g.declaredVars) == 0 {
		return fmt.Sprintf("%d", g.r.Intn(99)+1) // 1..99
	}

	if g.r.Float64() > 0.5 {
		return g.getRandomVar()
	}

	// Составное мат. выражение: x + 42
	var left string
	if g.r.Float64() > 0.5 && len(g.declaredVars) > 0 {
		left = g.getRandomVar()
	} else {
		left = fmt.Sprintf("%d", g.r.Intn(99)+1)
	}

	var right string
	if g.r.Float64() > 0.5 && len(g.declaredVars) > 0 {
		right = g.getRandomVar()
	} else {
		right = fmt.Sprintf("%d", g.r.Intn(99)+1)
	}

	op := g.mathOps[g.r.Intn(len(g.mathOps))]
	return fmt.Sprintf("%s %s %s", left, op, right)
}

func (g *RandomProgramGenerator) generateCondition() string {
	// Например: x <= 10
	left := g.getRandomVarOrNumber()
	right := g.getRandomVarOrNumber()
	comp := g.compareOps[g.r.Intn(len(g.compareOps))]

	condition := fmt.Sprintf("%s %s %s", left, comp, right)

	// 30% шанс усложнить условие логическим оператором
	if g.r.Float64() > 0.7 {
		logic := g.logicOps[g.r.Intn(len(g.logicOps))]
		extraLeft := g.getRandomVarOrNumber()
		extraRight := g.getRandomVarOrNumber()
		extraComp := g.compareOps[g.r.Intn(len(g.compareOps))]

		condition = fmt.Sprintf("(%s) %s (%s %s %s)", condition, logic, extraLeft, extraComp, extraRight)
	}

	return condition
}

func (g *RandomProgramGenerator) getRandomVar() string {
	if len(g.declaredVars) == 0 {
		return "1" // Fallback
	}
	return g.declaredVars[g.r.Intn(len(g.declaredVars))]
}

func (g *RandomProgramGenerator) getRandomVarOrNumber() string {
	if len(g.declaredVars) > 0 && g.r.Float64() > 0.5 {
		return g.getRandomVar()
	}
	return fmt.Sprintf("%d", g.r.Intn(99)+1)
}

func contains(ss []string, s string) bool {
	for _, x := range ss {
		if x == s {
			return true
		}
	}
	return false
}
