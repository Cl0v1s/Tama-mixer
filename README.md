
## Bodyparts

A bodypart is labelled as follows: 

`type@name-subname-frame`

**type**: 
* eye
* mouth
* arm1
* arm2
* leg1
* leg2
* leg3

**name**:
The identifier for this particular bodyshape appearance.

**subname**:
An optional identifier for this particular bodyshape state. For eyes it actually represents the current expression.

**frame**:
Allows to create animations. Frame 0 is always Idle.

## Generated 

Generated frames are labelled as follows:

`bodyshape_bodyframe@bodypart_type=bodypartname-bodypartsubname-bodypartframe`