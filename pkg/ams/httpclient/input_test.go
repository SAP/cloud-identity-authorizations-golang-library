package httpclient

import (
	"reflect"
	"testing"
)

func TestInsertCustomInput(t *testing.T) {
	t.Run("map[string]string", func(t *testing.T) {
		input := map[string]string{
			"key1": "value1",
			"key2": "value2",
		}
		result := make(reqInput)
		insertCustomInput(result, reflect.ValueOf(input), []string{"$app"})
		expected := reqInput{
			"$app.key1": "value1",
			"$app.key2": "value2",
		}
		if !reflect.DeepEqual(result, expected) {
			t.Fatalf("expected %v, got %v", expected, result)
		}
	})
}
