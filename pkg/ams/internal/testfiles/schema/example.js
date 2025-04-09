[{
	"schema": [
		"schema"
	],
	"definition": {
		"attribute": "Structure",
		"nested": {
			"$app": {
				"attribute": "Structure",
				"nested": {
					"string_value": { "attribute": "String" },
					"bool_value": { "attribute": "Boolean" },
					"number_value": { "attribute": "Number" },
					"string_array_value": { "attribute": "String[]" },
					"bool_array_value": { "attribute": "Boolean[]" },
					"number_array_value": { "attribute": "Number[]" },
					"deeper_nested": {
						"attribute": "Structure",
						"nested": {
							"nested_number_array_value": { "attribute": "Number[]" },
							"dot.in.name": {
								"attribute": "Structure",
								"nested": {
									"number_value": { "attribute": "Number" }
								}
							}
						}
					}
				}
			}
		}
	}
},
{

	"schema": [
		"_dcltentant_",
		"random",
		"package",
		"schema"
	],
	"tenant": "tenant_id1"
}]