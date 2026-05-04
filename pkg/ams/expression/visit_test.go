package expression

import (
	"fmt"
	"strings"
	"testing"
)

func TestVisit(t *testing.T) {
	t.Run("Example for creating SQL", func(t *testing.T) {
		exp := And(
			Or(
				Eq(Ref("x"), String("a")),
				Lt(Ref("y"), Number(10)),
			),
			IsNull(Ref("z")),
			Function("customFunc", nil, []Expression{
				Bool(true),
			}),
		)
		sql := Visit(exp,
			func(t string, args []string) string {
				switch t {
				case AND:
					return "(" + strings.Join(args, " AND ") + ")"
				case OR:
					return "(" + strings.Join(args, " OR ") + ")"
				case EQ:
					return args[0] + " = " + args[1]
				case NE:
					return args[0] + " != " + args[1]
				case GT:
					return args[0] + " > " + args[1]
				case LT:
					return args[0] + " < " + args[1]
				case GE:
					return args[0] + " >= " + args[1]
				case LE:
					return args[0] + " <= " + args[1]
				case BETWEEN:
					return args[0] + " BETWEEN " + args[1] + " AND " + args[2]
				case NOT_BETWEEN:
					return args[0] + " NOT BETWEEN " + args[1] + " AND " + args[2]
				case LIKE:
					return args[0] + " LIKE " + args[1]
				case NOT_LIKE:
					return args[0] + " NOT LIKE " + args[1]
				case IN:
					return args[0] + " IN " + args[1]
				case NOT_IN:
					return args[0] + " NOT IN " + args[1]
				case IS_NULL:
					return args[0] + " IS NULL"
				case IS_NOT_NULL:
					return args[0] + " IS NOT NULL"
				case "customFunc":
					return fmt.Sprintf("CUSTOM_FUNC(%s)", strings.Join(args, ", "))
				default:
					panic("unsupported operator: " + t)
				}
			},
			func(r Reference) string {
				return r.GetName()
			},
			func(c Constant) string {
				switch v := c.(type) {
				case String:
					return "'" + string(v) + "'"
				case Number:
					return fmt.Sprintf("%v", v)
				case Bool:
					return fmt.Sprintf("%v", v)
				case StringArray:
					builder := strings.Builder{}
					builder.WriteString("(")
					for i, s := range v {
						if i > 0 {
							builder.WriteString(", ")
						}
						builder.WriteString("'" + string(s) + "'")
					}
					builder.WriteString(")")
					return builder.String()
				case NumberArray:
					builder := strings.Builder{}
					builder.WriteString("(")
					for i, n := range v {
						if i > 0 {
							builder.WriteString(", ")
						}
						fmt.Fprintf(&builder, "%v", n)
					}
					builder.WriteString(")")
					return builder.String()
				case BoolArray:
					builder := strings.Builder{}
					builder.WriteString("(")
					for i, b := range v {
						if i > 0 {
							builder.WriteString(", ")
						}
						fmt.Fprintf(&builder, "%v", b)
					}
					builder.WriteString(")")
					return builder.String()
				default:
					panic("unsupported constant type")
				}
			},
		)

		expected := "((x = 'a' OR y < 10) AND z IS NULL AND CUSTOM_FUNC(true))"
		if sql != expected {
			t.Fatalf("expected SQL to be %s, but was %s", expected, sql)
		}
	})
}
