# MQTT API

This is the MQTT API for the ttrn project. All changes have to be coordinated with all team members. The messages are specified as [JSON Schema](http://json-schema.org/understanding-json-schema/).

## Train

**Description:** There are _three_ trains, identified with the numbers 0, 1 and 2.

**Topic:** `/train/{:id}/speed` 


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

**Description:** There are three (maybe four) checkpoints which update the position of the trains. The checkpoints are numbered starting form 0. If a train passes a checkpoint the topic will be updated with the number of the checkpoint.

**Topic:** `/train/{:id}/position/` 

```json
{
	"type": "object",
	"properties" : {
		"position": {
			"type": "integer",
			"enum": [0, 1, 2],
		}
	}
}
```


## Turnout

There are at max 10 turnouts. Identified by numbering them starting at 0. Each switch can have two different positions: _straight_ (the train will continue it's current direction) and _diverging_ (the train will branch of the straight). Straight is _0_, diverging is _1_.


