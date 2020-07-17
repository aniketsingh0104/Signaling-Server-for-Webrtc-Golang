package signaling

// TYPES

// offer type
const OFFER string = "OFFER"

// answer type
const ANSWER string = "ANSWER"

// candidate type
const CANDIDATE string = "CANDIDATE"

/*----------------------------------------------------------*/

// REQUEST ACTIONS

// to start a new meeting room
const START string = "START"

// to join room
const JOIN string = "JOIN"

// to end room / meeting - only by "owner"
const END string = "END" // *also a reply message

// to leave the room / meeting
const LEAVE string = "LEAVE"

/*----------------------------------------------------------*/

// REPLY ACTIONS

// action to indicate send "offer"
const READY string = "READY"

// action to indicate wait for pair to join
const WAIT_PAIR string = "WAIT_PAIR"

// action to indicate wait for pair to make offer
const WAIT_OFFER string = "WAIT_OFFER"
