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
				case "and":
					return "(" + strings.Join(args, " AND ") + ")"
				case "or":
					return "(" + strings.Join(args, " OR ") + ")"
				case "eq":
					return args[0] + " = " + args[1]
				case "ne":
					return args[0] + " != " + args[1]
				case "gt":
					return args[0] + " > " + args[1]
				case "lt":
					return args[0] + " < " + args[1]
				case "ge":
					return args[0] + " >= " + args[1]
				case "le":
					return args[0] + " <= " + args[1]
				case "between":
					return args[0] + " BETWEEN " + args[1] + " AND " + args[2]
				case "not_between":
					return args[0] + " NOT BETWEEN " + args[1] + " AND " + args[2]
				case "like":
					return args[0] + " LIKE " + args[1]
				case "not_like":
					return args[0] + " NOT LIKE " + args[1]
				case "in":
					return args[0] + " IN " + args[1]
				case "not_in":
					return args[0] + " NOT IN " + args[1]
				case "is_null":
					return args[0] + " IS NULL"
				case "is_not_null":
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
						builder.WriteString(fmt.Sprintf("%v", n))
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
						builder.WriteString(fmt.Sprintf("%v", b))
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
