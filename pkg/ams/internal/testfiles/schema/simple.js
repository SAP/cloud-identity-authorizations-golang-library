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
					"NumberValue": { "attribute": "Number" },
					"BoolArrayValue": { "attribute": "Boolean[]" },
					"NumberArrayValue": { "attribute": "Number[]" }
				}
			},
			"$env": {
				"attribute": "Structure",
				"nested": {
					"EnvN": { "attribute": "Number" }
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