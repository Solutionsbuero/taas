# MQTT API

This is the MQTT API for the ttrn project. All changes have to be coordinated with all team members. The messages are specified as [JSON Schema](http://json-schema.org/understanding-json-schema/).

## Train

There are _three_ trains, identified with the numbers 0, 1 and 2.

### Train Speed

**Description:** Sets the speed of the train. Speed is a integer ranging between -5 (reverse in full speed) 0 (stop) and 5 (full speed ahead).

**Topic:** `/train/{:id}/speed` 

**Message:**

```json
{
	"type": "object",
	"properties" : {
		"speed": {
			"type": "integer",
			"minimum": -5,
			"maximum": 5
		}
	}
}
```


### Train Position

**Topic:** `/train/{:id}/position/{:position}` 



