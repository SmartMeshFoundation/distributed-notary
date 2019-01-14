package userapi

import (
	"fmt"

	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/nkbai/log"
)

// APIName :
type APIName string

/* #nosec */
const (
	APIAdminNamePrefix             = "Admin-"
	APIAdminNameCreatePrivateKey   = APIAdminNamePrefix + "CreatePrivateKey"   // 发起一次私钥协商
	APIAdminNameGetPrivateKeyList  = APIAdminNamePrefix + "GetPrivateKeyList"  // 私钥片列表查询
	APIAdminNameRegisterNewSCToken = APIAdminNamePrefix + "RegisterNewSCToken" // 注册一个新的侧链token

	APIUserNamePrefix           = "User-"
	APIUserNameGetNotaryList    = APIUserNamePrefix + "GetNotaryList"  // 公证人列表查询
	APIUserNameGetSCTokenList   = APIUserNamePrefix + "GetSCTokenList" // 当前支持的SCToken列表查询
	APIUserNameGetLockinStatus  = APIUserNamePrefix + "GetLockinStatus"
	APIUserNameSCPrepareLockin  = APIUserNamePrefix + "SCPrepareLockin"
	APIUserNameGetLockoutStatus = APIUserNamePrefix + "GetLockoutStatus"
	APIUserNameMCPrepareLockout = APIUserNamePrefix + "MCPrepareLockout"

	APIDebugNamePrefix            = "Debug-"
	APIDebugNameTransferToAccount = APIDebugNamePrefix + "TransferToAccount" // 给某个账户在所有链上转10eth,为了测试
	APIDebugNameGetAllLockinInfo  = APIDebugNamePrefix + "GetAllLockinInfo"
	APIDebugNameClearSCTokenList  = APIDebugNamePrefix + "ClearSCTokenList" // 清空所有的SCToken
)

// APIName2URLMap :
var APIName2URLMap map[string]string

func init() {
	APIName2URLMap = make(map[string]string)
	/*
		admin
	*/
	APIName2URLMap[APIAdminNameCreatePrivateKey] = "/api/1/admin/private-key"
	APIName2URLMap[APIAdminNameGetPrivateKeyList] = "/api/1/admin/private-keys"
	APIName2URLMap[APIAdminNameRegisterNewSCToken] = "/api/1/admin/sctoken"
	/*
		user
	*/
	APIName2URLMap[APIUserNameGetNotaryList] = "/api/1/user/notaries"
	APIName2URLMap[APIUserNameGetSCTokenList] = "/api/1/user/sctokens"
	APIName2URLMap[APIUserNameGetLockinStatus] = "/api/1/user/lockin/:sctoken/:secrethash"
	/*
		debug
	*/
	APIName2URLMap[APIDebugNameGetAllLockinInfo] = "/api/1/debug/lockin"
	APIName2URLMap[APIDebugNameTransferToAccount] = "/api/1/debug/transfer-to-account/:account"
}

/*
UserAPI :
提供给用户的API
暂时把SystemRequest和NotaryRequest都放在UserAPI
*/
type UserAPI struct {
	api.BaseAPI
}

// NewUserAPI :
func NewUserAPI(host string) *UserAPI {
	var userAPI UserAPI
	router, err := rest.MakeRouter(
		/*
			api about private key
		*/
		rest.Put(APIName2URLMap[APIAdminNameCreatePrivateKey], userAPI.createPrivateKey),
		/*
			user api
		*/
		rest.Get(APIName2URLMap[APIUserNameGetNotaryList], userAPI.getNotaryList),
		rest.Get(APIName2URLMap[APIUserNameGetSCTokenList], userAPI.getSCTokenList),
		rest.Get(APIName2URLMap[APIUserNameGetLockinStatus], userAPI.getLockinStatus),
		/*
			admin api
		*/
		rest.Get(APIName2URLMap[APIAdminNameGetPrivateKeyList], userAPI.getPrivateKeyList),
		rest.Put(APIName2URLMap[APIAdminNameRegisterNewSCToken], userAPI.registerNewSCToken),
		/*
			debug api
		*/
		rest.Get(APIName2URLMap[APIDebugNameGetAllLockinInfo], userAPI.getAllLockinInfo),
		rest.Get(APIName2URLMap[APIDebugNameTransferToAccount], userAPI.transferToAccount),
	)
	if err != nil {
		log.Crit(fmt.Sprintf("maker router :%s", err))
	}
	userAPI.BaseAPI = api.NewBaseAPI("UserAPI-Server", host, router)
	return &userAPI
}
