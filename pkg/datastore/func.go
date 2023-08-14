package datastore

const (
	KFuncTableName  = "FUNCTION"
	KFuncKey        = "PRIMARY_KEY"
	KFuncSdModel    = "SD_MODEL"
	KFuncSdVae      = "SD_VAE"
	KFuncEndPoint   = "END_POINT"
	KCreateTime     = "FUNC_CREATE_TIME"
	kLastModifyTime = "FUNC_LAST_MODIFY_TIME"
)

type FuncInterface interface {
	Put(key string, data map[string]interface{}) error
	Get(key string, fields []string) (map[string]interface{}, error)
	ListAll(fields []string) ([]map[string]interface{}, error)
}

func NewFuncDataStore(dbType DatastoreType) FuncInterface {
	switch dbType {
	case SQLite:

	}
	return nil
}