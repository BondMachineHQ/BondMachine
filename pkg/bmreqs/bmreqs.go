package bmreqs

type bmReqObj struct {
	bmReqMap
}

type bmReqSet interface {

	// Initialization
	init()

	// Handling of the single requirement node
	getName() string
	getType() uint8
	setName(string) error
	setType(uint8) error

	// Insert and remove requirements from/to the current node
	insertReq(string) error
	removeReq(string) error

	// Check if the current node contains the given requirement
	checkReq(string) (string, error)

	// Exporting requirements for the current node
	getReqs() string
	importReqs(*ReqRoot, string, string, string) error

	// SubRequirements
	supportSub() bool
	listSub() []string
	getSub(string) (*bmReqObj, error)
}

type bmReqMap map[string]bmReqSet

func (o *bmReqObj) init() {
	if o.bmReqMap == nil {
		o.bmReqMap = make(map[string]bmReqSet)
	}
}

func (o *bmReqObj) getMap() bmReqMap {
	return o.bmReqMap
}

type ExportedReq struct {
	Node string
	Type uint8
	Name string
	Req  string
}

type ExportedReqs []ExportedReq
