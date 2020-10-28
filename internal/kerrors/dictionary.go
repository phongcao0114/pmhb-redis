package kerrors

// ErrorMessage is the content of the error dictionary
type ErrorMessage string

// ErrorCode is the key of the error dictionary
type ErrorCode string

func (m ErrorMessage) String() string {
	return string(m)
}
func (c ErrorCode) String() string {
	return string(c)
}

// Dictionary the set of items
type Dictionary struct {
	Items map[ErrorCode]ErrorMessage
}

// default error message
const (
	NotFoundErrorMessageInDictionary = ErrorMessage("the internal code did not match to error message in dictionary")
	CannotGetFunctionNumberErrorText = "Can't get function_number"
)

// setup error code
const (
	NoneError    ErrorCode = "0000"
	UnknownError ErrorCode = "0001"

	/*-------- BadRequest start with = 1 ----------*/
	ValidateFailed           ErrorCode = "1000"
	CannotDecodeInputRequest ErrorCode = "1001"
	MarshalFail              ErrorCode = "1002"
	UnMarshalFail            ErrorCode = "1003"
	DuplicateRequestInfo     ErrorCode = "1004"

	/*-------- InternalServerError start with = 2 ----------*/
	DatabaseServerError        ErrorCode = "2000"
	CannotRemoveRedis          ErrorCode = "2001"
	CannotUpdateDataInDB       ErrorCode = "2002"
	CannotGetDataFromDB        ErrorCode = "2003"
	CardlinkResponseCodeFailed ErrorCode = "2004"
	CannotGetRedis             ErrorCode = "2005"
	CannotSetRedis             ErrorCode = "2006"
	DatabaseScanErr            ErrorCode = "2007"

	/*-------- TimeOutServer start with = 3 ----------*/

	/*-------- BusinessValidationError start with = 4 ----------*/
	NotFoundItemInQuery         ErrorCode = "4000"
	CannotDecodeDatabaseRequest ErrorCode = "4001"
	CursorDatabaseErr           ErrorCode = "4002"
	StoreTransactionErr         ErrorCode = "4003"
	TransactionExisting         ErrorCode = "4004"

	/*-------- SecurityAuthorizeOnlyBackOffice start with = 5 ----------*/
	AuthenticationRedisFailed ErrorCode = "5000"

	/*-------- StatusCode for header only ----------*/
	StatusNoneError   ErrorCode = "00"
	StatusErrorFailed ErrorCode = "01"
)

// ErrDictionary contains all error description
var ErrDictionary = Dictionary{Items: map[ErrorCode]ErrorMessage{
	NoneError:    "Success",
	UnknownError: "Cannot Recognize Error (out update error dictionary)",

	/*-------- BadRequest start with = 1 ----------*/
	ValidateFailed:           "Validate failed",
	CannotDecodeInputRequest: "Can't decode input request",
	MarshalFail:              "Marshal failed",
	UnMarshalFail:            "Unmarshal failed",
	DuplicateRequestInfo:     "Duplicate request information",

	/*-------- InternalServerError start with = 2 ----------*/
	DatabaseServerError:        "Database server error",
	CannotRemoveRedis:          "Can't remove data from redis",
	CannotUpdateDataInDB:       "Can't update data in database",
	CannotGetDataFromDB:        "Can't get data from database",
	CardlinkResponseCodeFailed: "Cardlink response code failed",
	CannotGetRedis:             "Can't get data redis",
	CannotSetRedis:             "Can't set data to redis",
	DatabaseScanErr:            "Database scanning error",

	/*-------- TimeOutServer start with = 3 ----------*/

	/*-------- BusinessValidationError start with = 4 ----------*/
	NotFoundItemInQuery:         "Not found item in query",
	CannotDecodeDatabaseRequest: "Can't decode database request",
	CursorDatabaseErr:           "Cursor database error",
	StoreTransactionErr:         "Store transaction error",
	TransactionExisting:         "The information has existed",

	/*-------- SecurityAuthorizeOnlyBackOffice start with = 5 ----------*/
	AuthenticationRedisFailed: "Redis authentication failed",

	/*-------- StatusCode for header only ----------*/
	StatusNoneError:   "nil",
	StatusErrorFailed: "status failed",
}}

// Get return item in dictionary, but if not found return default not found
func (d *Dictionary) Get(code ErrorCode) ErrorMessage {
	msg, ok := d.Items[code]
	if !ok {
		return NotFoundErrorMessageInDictionary
	}
	return msg
}

// Set adds a new item to the dictionary
func (d *Dictionary) Set(key ErrorCode, value ErrorMessage) {
	d.Items[key] = value
}

// Has return true if the key exist in dictionary
func (d *Dictionary) Has(key ErrorCode) bool {
	_, ok := d.Items[key]
	return ok
}

// Clear removes all items in dictionary
func (d *Dictionary) Clear() {
	d.Items = make(map[ErrorCode]ErrorMessage)
}

// Size returns length of dictionary
func (d *Dictionary) Size() int {
	return len(d.Items)
}
