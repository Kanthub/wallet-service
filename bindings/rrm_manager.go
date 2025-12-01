// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// ReferralRewardManagerMetaData contains all meta data concerning the ReferralRewardManager contract.
var ReferralRewardManagerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"FUND_MANAGER_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"PAUSER_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"REWARD_MANAGER_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"batchAddRewards\",\"inputs\":[{\"name\":\"users\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"amounts\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"batchClearRewards\",\"inputs\":[{\"name\":\"users\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"batchReduceRewards\",\"inputs\":[{\"name\":\"users\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"amounts\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"batchSetRewards\",\"inputs\":[{\"name\":\"users\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"amounts\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claimReward\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"fundRewardPool\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getConfig\",\"inputs\":[],\"outputs\":[{\"name\":\"_minClaimAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_maxRewardPerUser\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_maxReferralsPerUser\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getReferrals\",\"inputs\":[{\"name\":\"_referrer\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRewardPoolStatus\",\"inputs\":[],\"outputs\":[{\"name\":\"poolAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"totalPending\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"totalClaimed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSystemStats\",\"inputs\":[],\"outputs\":[{\"name\":\"_totalUsers\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_totalReferrers\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_totalPending\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_totalClaimed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getUserReferralInfo\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"_referrer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_referralCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_pendingReward\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_totalClaimed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"admin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_ypusdToken\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"maxReferralsPerUser\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"maxRewardPerUser\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"minClaimAmount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingRewards\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"referralCount\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"referrals\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"referrer\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"callerConfirmation\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setReferrer\",\"inputs\":[{\"name\":\"_referrer\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalClaimedRewards\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalClaimedRewardsGlobal\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalPendingRewards\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalReferrers\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalUsers\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateConfig\",\"inputs\":[{\"name\":\"_minClaimAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_maxRewardPerUser\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_maxReferralsPerUser\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"withdrawFunds\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"ypusdToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIyPUSD\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"ConfigUpdated\",\"inputs\":[{\"name\":\"minClaimAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"maxRewardPerUser\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"maxReferralsPerUser\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ReferrerSet\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"referrer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RewardAdded\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"manager\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RewardClaimed\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RewardCleared\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RewardPoolFunded\",\"inputs\":[{\"name\":\"funder\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RewardReduced\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"manager\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RewardSet\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"oldAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessControlBadConfirmation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"neededRole\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedInnerCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	Bin: "0x60a080604052346100cc57306080527ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a009081549060ff8260401c166100bd57506001600160401b036002600160401b031982821601610078575b60405161262490816100d182396080518181816111f301526112d10152f35b6001600160401b031990911681179091556040519081527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d290602090a15f8080610059565b63f92ee8a960e01b8152600490fd5b5f80fdfe6080604090808252600480361015610015575f80fd5b5f3560e01c91826301ffc9a714611cb557508163155dd5ee14611b755781631d583e0d14611ab7578163248a9ca314611a815781632cf003c214611a485781632f2ff15d14611a205781632f83a6bc146119e657816330eae166146118e157816331d7a262146118aa57816336568abe1461186657816338196bc7146118445781633f4ba83a146117cf57816341a0894d14611725578163485cc955146114ad5781634f1ef2861461125757816352d1902d146111e05781635c975abb146111b157816362d03cb7146111935781636983e74d146110355781636d750d8014610ffb5781636fdebdc014610e8057816373ad494614610d1e5781638456cb5914610cb957816391d1485414610c695781639e6c612414610c4b578163a110b93f14610c14578163a18a7bfc146109d2578163a217fddf146109b8578163ad3cb1cc1461091a578163ad8d6f2c146108fc578163adc25bde146108de578163b88a802f14610625578163bff1f9e114610607578163c3f909d4146105d3578163cb05347a14610526578163d1ce59a7146104ee578163d355271214610471578163d547741f14610429578163db74559b146103f4578163e3e9dfba146102dc57508063e63ab1e9146102a2578063e697b5d81461024c578063ec1371f21461022e5763f339ae7814610204575f80fd5b3461022a575f36600319011261022a575f5490516001600160a01b039091168152602090f35b5f80fd5b503461022a575f36600319011261022a576020906006549051908152f35b503461022a578060031936011261022a57610265611d07565b906024359160018060a01b038091165f526005602052815f20805484101561022a5760209361029391611e01565b92905490519260031b1c168152f35b503461022a575f36600319011261022a57602090517f65d7a28e3265b37a6474929f336521b332c1681b933f6cb9f3376673440d862a8152f35b823461022a576102eb36611db6565b90929391946102f861204c565b610303828714611f4b565b85156103c257505f5b85811061031557005b8061032b6103266001938989611ec3565b611ed3565b610336828588611ec3565b3590838060a01b031661034a811515611ee7565b610355821515611e2a565b805f527f54dba80f86e498df5a2fcdb5c089d051d97ce8ca9ddb708c70da66557a954de78460209080825261038d858a5f2054611fdb565b9061039c600b54831115611f8f565b845f528252885f205560066103b2858254611fdb565b905587519384523393a30161030c565b606490602084519162461bcd60e51b8352820152600c60248201526b456d7074792061727261797360a01b6044820152fd5b823461022a57602036600319011261022a576020916001600160a01b03610419611d07565b165f528252805f20549051908152f35b823461022a578060031936011261022a5761046f913561046a600161044c611d1d565b93835f525f805160206125af8339815191526020525f200154612120565b612441565b005b90503461022a57602036600319011261022a576001600160a01b0380610495611d07565b165f526003602052825f205416906020526104ea825f2054926001602052805f20546002602052815f205491519485948590949392606092608083019660018060a01b03168352602083015260408201520152565b0390f35b823461022a575f36600319011261022a5760809060085490600954906006549060075492815194855260208501528301526060820152f35b90503461022a575f36600319011261022a575f5482516370a0823160e01b81523092810192909252602090829060249082906001600160a01b03165afa9081156105c9575f91610597575b506006546007549251308152602081019290925260408201526060810191909152608090f35b90506020813d6020116105c1575b816105b260209383611d64565b8101031261022a57515f610571565b3d91506105a5565b82513d5f823e3d90fd5b823461022a575f36600319011261022a57600a54600b54600c549251918252602082015261ffff9091166040820152606090f35b823461022a575f36600319011261022a576020906008549051908152f35b823461022a575f36600319011261022a577f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f009060028254146108cf576002825561066d6124c0565b335f5260209060018252805f205493600a54851061088e575f5482516370a0823160e01b815230838201526001600160a01b039185908290602490829086165afa908115610884579087915f91610853575b50106107eb57335f52600184525f838120556106dd86600654611f2a565b600655335f5260028452825f206106f5878254611fdb565b905561070386600754611fdb565b6007555f8054845163a9059cbb60e01b815233858201908152602081018a905290938793859392849003604001928492165af19081156107e1575f916107b4575b501561077a57505192835260019233917f106f923f993c2149d49b4255ff723acafa1f2d94393f561d3eda32ae348f724191a255005b82606492519162461bcd60e51b835282015260166024820152753ca82aa9a2103a3930b739b332b9103330b4b632b21760511b6044820152fd5b6107d49150843d86116107da575b6107cc8183611d64565b810190611e67565b86610744565b503d6107c2565b83513d5f823e3d90fd5b5082608492519162461bcd60e51b8352820152603760248201527f496e73756666696369656e742062616c616e636520696e20636f6e747261637460448201527f2e20506c6561736520636f6e746163742061646d696e2e0000000000000000006064820152fd5b809250868092503d831161087d575b61086c8183611d64565b8101031261022a57869051886106bf565b503d610862565b84513d5f823e3d90fd5b82606492519162461bcd60e51b8352820152601a60248201527f42656c6f77206d696e696d756d20636c61696d20616d6f756e740000000000006044820152fd5b51633ee5aeb560e01b81529050fd5b823461022a575f36600319011261022a57602090600a549051908152f35b823461022a575f36600319011261022a57602090600b549051908152f35b823461022a575f36600319011261022a578051918183019083821067ffffffffffffffff8311176109a557508152600582526020640352e302e360dc1b60208401528151928391602083528151918260208501525f5b83811061098f5750505f83830185015250601f01601f19168101030190f35b8181018301518782018701528694508201610970565b604190634e487b7160e01b5f525260245ffd5b823461022a575f36600319011261022a57602090515f8152f35b823461022a576020908160031936011261022a576109ee611d07565b926109f76124c0565b6001600160a01b03938416928315610bd457338414610b9c57335f526003815284835f205416610b6557835f52818152825f205461ffff600c54161115610b1a57335f5260038152825f20846bffffffffffffffffffffffff60a01b825416179055835f52818152825f20610a6c8154611fe8565b9055835f5260058152825f209485549568010000000000000000871015610b075786610a9f916001809899018155611e01565b819291549060031b9133831b921b1916179055610abd600854611fe8565b600855845f52525f205414610af4575b337f5f7165288eef601591cf549e15ff19ef9060b7f71b9c115be946fa1fe7ebf68a5f80a3005b610aff600954611fe8565b600955610acd565b604184634e487b7160e01b5f525260245ffd5b608492519162461bcd60e51b8352820152602260248201527f5265666572726572206861732072656163686564206d617820726566657272616044820152616c7360f01b6064820152fd5b606492519162461bcd60e51b83528201526014602482015273149959995c9c995c88185b1c9958591e481cd95d60621b6044820152fd5b606492519162461bcd60e51b8352820152601560248201527421b0b73737ba103932b332b9103cb7bab939b2b63360591b6044820152fd5b606492519162461bcd60e51b8352820152601860248201527f496e76616c6964207265666572726572206164647265737300000000000000006044820152fd5b823461022a57602036600319011261022a576020906001600160a01b03610c39611d07565b165f5260028252805f20549051908152f35b823461022a575f36600319011261022a576020906007549051908152f35b823461022a578060031936011261022a57602091610c85611d1d565b90355f525f805160206125af8339815191528352815f209060018060a01b03165f52825260ff815f20541690519015158152f35b823461022a575f36600319011261022a5760207f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a25891610cf66120c5565b610cfe6124c0565b5f805160206125cf833981519152805460ff1916600117905551338152a1005b90503461022a57606036600319011261022a5780356024359160443590610d43611ff6565b8215610e3d57828410610e055781151580610df9575b15610dbe5750610db9907f6d05db271e19f930af71c4765de54ef86294762644c20f4d6fd2609d057d3c7b9483600a5584600b5561ffff821661ffff19600c541617600c5551938493846040919493926060820195825260208201520152565b0390a1005b606490602086519162461bcd60e51b83528201526015602482015274496e76616c6964206d617820726566657272616c7360581b6044820152fd5b5061ffff821115610d59565b606490602086519162461bcd60e51b83528201526012602482015271125b9d985b1a59081b585e081c995dd85c9960721b6044820152fd5b606490602086519162461bcd60e51b8352820152601860248201527f496e76616c6964206d696e20636c61696d20616d6f756e7400000000000000006044820152fd5b90503461022a57610e9036611db6565b610e9e95929491939561204c565b610ea9818614611f4b565b5f5b858110610eb457005b6001600160a01b03610ed681610ece610326858b8d611ec3565b161515611ee7565b610eec610ee4838589611ec3565b351515611e2a565b80610efb610326848a8c611ec3565b165f52600190602091808352865f2054610f1685878b611ec3565b3511610fb857907f4f98a5b5c72400c0bdafc1fe75f6069e95d1355e66133b9e2e517a00fde2b6bc610f9a610326868c8e610f8c838f8f90848f9260019f9e9d90610f85918f8f610f776103268a8e610f70828d8d611ec3565b3599611ec3565b165f52525f20918254611f2a565b9055611ec3565b35610f856006918254611f2a565b92610fa685888c611ec3565b3592895193845233941692a301610eab565b865162461bcd60e51b8152808701849052601c60248201527f496e73756666696369656e742070656e64696e672072657761726473000000006044820152606490fd5b823461022a575f36600319011261022a57602090517f0f51adb3f49e4a9bbb17b3783f025995eaf8c24be2c8eefff214bdfda05ef94d8152f35b823461022a5761104436611db6565b909361105193929361204c565b61105c828514611f4b565b5f5b84811061106757005b80807fbd0682fae90263f394bf341cc4c207fc356246f627963b5fdf3580867cb4a74e84888a89611157868b61115161032660019c8d8060a01b03986110b58a610ece61032685858d611ec3565b6110ce6110c383878c611ec3565b35600b541015611f8f565b896110dd61032684848c611ec3565b165f528e97602098808a528c5f20549e8f888d826110fc858484611ec3565b35111561116a57926111119161111794611ec3565b35611f2a565b6111246006918254611fdb565b90555b61113284888d611ec3565b35908c611143610326878787611ec3565b165f528a528c5f2055611ec3565b94611ec3565b359084519687528601521692a20161105e565b926111789161117f94611ec3565b3590611f2a565b61118c6006918254611f2a565b9055611127565b823461022a575f36600319011261022a576020906009549051908152f35b823461022a575f36600319011261022a5760209060ff5f805160206125cf833981519152541690519015158152f35b823461022a575f36600319011261022a577f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316300361124a57602090517f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc8152f35b5163703e46dd60e11b8152fd5b90508160031936011261022a5761126c611d07565b602492833567ffffffffffffffff811161022a573660238201121561022a578084013561129881611d9a565b936112a584519586611d64565b818552602091828601933689838301011161022a57815f928a8693018737870101526001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000811630811490811561147f575b5061146f5761130a611ff6565b81169484516352d1902d60e01b8152838189818a5afa5f9181611440575b506113435750505050505191634c9c8ce360e01b8352820152fd5b9087878794938b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc9182810361142b5750853b15611417575080546001600160a01b031916821790558451907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b5f80a28251156113fd5750505f61046f9581925190845af4913d156113f3573d6113e56113dc82611d9a565b92519283611d64565b81525f81943d92013e61252b565b506060925061252b565b945094505050503461140b57005b63b398979f60e01b8152fd5b8651634c9c8ce360e01b8152808501849052fd5b8751632a87526960e21b815280860191909152fd5b9091508481813d8311611468575b6114588183611d64565b8101031261022a5751905f611328565b503d61144e565b845163703e46dd60e11b81528790fd5b9050817f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc541614155f6112fd565b90503461022a578160031936011261022a576114c7611d07565b906114d0611d1d565b907ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a009283549260ff84871c16159367ffffffffffffffff81168015908161171d575b6001149081611713575b15908161170a575b506116fa5767ffffffffffffffff1981166001178655846116db575b506001600160a01b0390828216156116a0571691821561166557506115ec906115676124ea565b61156f6124ea565b6115776124ea565b60017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f00556115a36124ea565b6115ab6124ea565b5f805160206125cf833981519152805460ff191690556115c96124ea565b6115d28161214e565b506115dc816121d8565b506115e681612285565b5061232b565b506bffffffffffffffffffffffff60a01b5f5416175f55620f4240600a556402540be400600b556103e861ffff19600c541617600c5561162857005b805468ff00000000000000001916905551600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d290602090a1005b606490602087519162461bcd60e51b83528201526015602482015274496e76616c6964207950555344206164647265737360581b6044820152fd5b865162461bcd60e51b81526020818601526015602482015274496e76616c69642061646d696e206164647265737360581b6044820152606490fd5b68ffffffffffffffffff1916680100000000000000011785555f611540565b865163f92ee8a960e01b81528490fd5b9050155f611524565b303b15915061151c565b869150611512565b823461022a576020908160031936011261022a576001600160a01b03918261174b611d07565b165f5260058152815f209082518082845491828152019081945f52835f20905f5b8181106117b95750505081611782910382611d64565b8351938285019183865251809252840192915f5b8281106117a35785850386f35b8351871685529381019392810192600101611796565b825489168452928501926001928301920161176c565b823461022a575f36600319011261022a576117e86120c5565b5f805160206125cf8339815191529182549060ff821615611836577f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa602084868560ff1916905551338152a1005b8251638dfc202b60e01b8152fd5b823461022a575f36600319011261022a5760209061ffff600c54169051908152f35b823461022a578060031936011261022a5761187f611d1d565b90336001600160a01b0383160361189b575061046f9135612441565b5163334bd91960e11b81529050fd5b823461022a57602036600319011261022a576020906001600160a01b036118cf611d07565b165f5260018252805f20549051908152f35b90503461022a57602090602060031936011261022a5780359067ffffffffffffffff821161022a5761191591369101611d33565b61192093919361204c565b5f5b81811061192b57005b6001906001600160a01b0361194881610ece61032685888c611ec3565b8061195761032684878b611ec3565b165f5286868492838252875f20549385888661197b575b5050505050505001611922565b61032684926119ab927ffb50459538be3884cca2319f1e6c6ff046647e39d9687c512c405b4db5d4b37598611ec3565b165f5282525f8881205560066119c2858254611f2a565b90556119d261032686898d611ec3565b16928751908152a25f86868280858861196e565b823461022a575f36600319011261022a57602090517f0b84ee281e5cf521a9ad54a86fafe78946b157177e231bd8ae785af4d3b3620f8152f35b823461022a578060031936011261022a5761046f9135611a43600161044c611d1d565b6123d1565b823461022a57602036600319011261022a576020906001600160a01b0380611a6e611d07565b165f5260038352815f2054169051908152f35b823461022a57602036600319011261022a57602091355f525f805160206125af83398151915282526001815f2001549051908152f35b823461022a57602036600319011261022a575f6020833593611ada851515611e2a565b825484516323b872dd60e01b81523392810192909252306024830152604482018690529092839160649183916001600160a01b03165af19081156105c95790611b2a915f91611b56575b50611e7f565b519081527fd3c7d1587d8f93bae0fcfc5a820d2187e755a75ec28a06c20a9da047e1b01eb560203392a2005b611b6f915060203d6020116107da576107cc8183611d64565b84611b24565b90503461022a576020908160031936011261022a57803590611b95611ff6565b611ba0821515611e2a565b5f5484516370a0823160e01b815230838201526001600160a01b0390911691908481602481865afa8015611cab5784915f91611c7a575b5010611c4057845163a9059cbb60e01b81523391810191825260208201939093528391839182905f90829060400103925af19182156107e15761046f93505f92611c23575b5050611e7f565b611c399250803d106107da576107cc8183611d64565b5f80611c1c565b845162461bcd60e51b81529081018490526014602482015273496e73756666696369656e742062616c616e636560601b6044820152606490fd5b809250868092503d8311611ca4575b611c938183611d64565b8101031261022a578390515f611bd7565b503d611c89565b86513d5f823e3d90fd5b903461022a57602036600319011261022a57359063ffffffff60e01b821680920361022a57602091637965db0b60e01b8114908115611cf6575b5015158152f35b6301ffc9a760e01b14905083611cef565b600435906001600160a01b038216820361022a57565b602435906001600160a01b038216820361022a57565b9181601f8401121561022a5782359167ffffffffffffffff831161022a576020808501948460051b01011161022a57565b90601f8019910116810190811067ffffffffffffffff821117611d8657604052565b634e487b7160e01b5f52604160045260245ffd5b67ffffffffffffffff8111611d8657601f01601f191660200190565b604060031982011261022a5767ffffffffffffffff9160043583811161022a5782611de391600401611d33565b9390939260243591821161022a57611dfd91600401611d33565b9091565b8054821015611e16575f5260205f2001905f90565b634e487b7160e01b5f52603260045260245ffd5b15611e3157565b60405162461bcd60e51b815260206004820152600e60248201526d125b9d985b1a5908185b5bdd5b9d60921b6044820152606490fd5b9081602091031261022a5751801515810361022a5790565b15611e8657565b60405162461bcd60e51b81526020600482015260156024820152741e541554d1081d1c985b9cd9995c8819985a5b1959605a1b6044820152606490fd5b9190811015611e165760051b0190565b356001600160a01b038116810361022a5790565b15611eee57565b60405162461bcd60e51b8152602060048201526014602482015273496e76616c69642075736572206164647265737360601b6044820152606490fd5b91908203918211611f3757565b634e487b7160e01b5f52601160045260245ffd5b15611f5257565b60405162461bcd60e51b8152602060048201526015602482015274082e4e4c2f240d8cadccee8d040dad2e6dac2e8c6d605b1b6044820152606490fd5b15611f9657565b60405162461bcd60e51b815260206004820152601b60248201527f45786365656473206d61782072657761726420706572207573657200000000006044820152606490fd5b91908201809211611f3757565b5f198114611f375760010190565b335f9081527fb7db2dd08fcb62d0c9e08c51941cae53c267786a0b75803fb7960902fc8ef97d602052604090205460ff161561202e57565b60405163e2517d3f60e01b81523360048201525f6024820152604490fd5b335f9081527f98dafdfbfb28cc7383b6dba209d26804e020844a26b28880cd6eded2cf7b5d3560205260409020547f0f51adb3f49e4a9bbb17b3783f025995eaf8c24be2c8eefff214bdfda05ef94d9060ff16156120a75750565b6044906040519063e2517d3f60e01b82523360048301526024820152fd5b335f9081527f75442b0a96088b5456bc4ed01394c96a4feec0f883c9494257d76b96ab1c9b6b60205260409020547f65d7a28e3265b37a6474929f336521b332c1681b933f6cb9f3376673440d862a9060ff16156120a75750565b805f525f805160206125af83398151915260205260405f20335f5260205260ff60405f205416156120a75750565b6001600160a01b03165f8181527fb7db2dd08fcb62d0c9e08c51941cae53c267786a0b75803fb7960902fc8ef97d60205260409020545f805160206125af8339815191529060ff166121d2575f805260205260405f20815f5260205260405f20600160ff1982541617905533905f5f8051602061258f8339815191528180a4600190565b50505f90565b6001600160a01b03165f8181527f75442b0a96088b5456bc4ed01394c96a4feec0f883c9494257d76b96ab1c9b6b60205260409020547f65d7a28e3265b37a6474929f336521b332c1681b933f6cb9f3376673440d862a91905f805160206125af8339815191529060ff1661227e57825f5260205260405f20815f5260205260405f20600160ff1982541617905533915f8051602061258f8339815191525f80a4600190565b5050505f90565b6001600160a01b03165f8181527f98dafdfbfb28cc7383b6dba209d26804e020844a26b28880cd6eded2cf7b5d3560205260409020547f0f51adb3f49e4a9bbb17b3783f025995eaf8c24be2c8eefff214bdfda05ef94d91905f805160206125af8339815191529060ff1661227e57825f5260205260405f20815f5260205260405f20600160ff1982541617905533915f8051602061258f8339815191525f80a4600190565b6001600160a01b03165f8181527f48ce10f4f7616f06256ebdb7d81e03056199b13e52d95951937b97a5e9b84d6760205260409020547f0b84ee281e5cf521a9ad54a86fafe78946b157177e231bd8ae785af4d3b3620f91905f805160206125af8339815191529060ff1661227e57825f5260205260405f20815f5260205260405f20600160ff1982541617905533915f8051602061258f8339815191525f80a4600190565b90815f525f805160206125af8339815191528060205260405f209160018060a01b031691825f5260205260ff60405f205416155f1461227e57825f5260205260405f20815f5260205260405f20600160ff1982541617905533915f8051602061258f8339815191525f80a4600190565b90815f525f805160206125af8339815191528060205260405f209160018060a01b031691825f5260205260ff60405f2054165f1461227e57825f5260205260405f20815f5260205260405f2060ff19815416905533917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b5f80a4600190565b60ff5f805160206125cf83398151915254166124d857565b60405163d93c066560e01b8152600490fd5b60ff7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005460401c161561251957565b604051631afcd79f60e31b8152600490fd5b90612552575080511561254057805190602001fd5b604051630a12f52160e11b8152600490fd5b81511580612585575b612563575090565b604051639996b31560e01b81526001600160a01b039091166004820152602490fd5b50803b1561255b56fe2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800cd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f03300a264697066735822122012992a754d9cd888c7022f5e6688f95d115a6ecceb568756ac44f77f88e09b2064736f6c63430008160033",
}

// ReferralRewardManagerABI is the input ABI used to generate the binding from.
// Deprecated: Use ReferralRewardManagerMetaData.ABI instead.
var ReferralRewardManagerABI = ReferralRewardManagerMetaData.ABI

// ReferralRewardManagerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ReferralRewardManagerMetaData.Bin instead.
var ReferralRewardManagerBin = ReferralRewardManagerMetaData.Bin

// DeployReferralRewardManager deploys a new Ethereum contract, binding an instance of ReferralRewardManager to it.
func DeployReferralRewardManager(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ReferralRewardManager, error) {
	parsed, err := ReferralRewardManagerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ReferralRewardManagerBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ReferralRewardManager{ReferralRewardManagerCaller: ReferralRewardManagerCaller{contract: contract}, ReferralRewardManagerTransactor: ReferralRewardManagerTransactor{contract: contract}, ReferralRewardManagerFilterer: ReferralRewardManagerFilterer{contract: contract}}, nil
}

// ReferralRewardManager is an auto generated Go binding around an Ethereum contract.
type ReferralRewardManager struct {
	ReferralRewardManagerCaller     // Read-only binding to the contract
	ReferralRewardManagerTransactor // Write-only binding to the contract
	ReferralRewardManagerFilterer   // Log filterer for contract events
}

// ReferralRewardManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type ReferralRewardManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReferralRewardManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ReferralRewardManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReferralRewardManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ReferralRewardManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReferralRewardManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ReferralRewardManagerSession struct {
	Contract     *ReferralRewardManager // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// ReferralRewardManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ReferralRewardManagerCallerSession struct {
	Contract *ReferralRewardManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// ReferralRewardManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ReferralRewardManagerTransactorSession struct {
	Contract     *ReferralRewardManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// ReferralRewardManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type ReferralRewardManagerRaw struct {
	Contract *ReferralRewardManager // Generic contract binding to access the raw methods on
}

// ReferralRewardManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ReferralRewardManagerCallerRaw struct {
	Contract *ReferralRewardManagerCaller // Generic read-only contract binding to access the raw methods on
}

// ReferralRewardManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ReferralRewardManagerTransactorRaw struct {
	Contract *ReferralRewardManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewReferralRewardManager creates a new instance of ReferralRewardManager, bound to a specific deployed contract.
func NewReferralRewardManager(address common.Address, backend bind.ContractBackend) (*ReferralRewardManager, error) {
	contract, err := bindReferralRewardManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ReferralRewardManager{ReferralRewardManagerCaller: ReferralRewardManagerCaller{contract: contract}, ReferralRewardManagerTransactor: ReferralRewardManagerTransactor{contract: contract}, ReferralRewardManagerFilterer: ReferralRewardManagerFilterer{contract: contract}}, nil
}

// NewReferralRewardManagerCaller creates a new read-only instance of ReferralRewardManager, bound to a specific deployed contract.
func NewReferralRewardManagerCaller(address common.Address, caller bind.ContractCaller) (*ReferralRewardManagerCaller, error) {
	contract, err := bindReferralRewardManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ReferralRewardManagerCaller{contract: contract}, nil
}

// NewReferralRewardManagerTransactor creates a new write-only instance of ReferralRewardManager, bound to a specific deployed contract.
func NewReferralRewardManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*ReferralRewardManagerTransactor, error) {
	contract, err := bindReferralRewardManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ReferralRewardManagerTransactor{contract: contract}, nil
}

// NewReferralRewardManagerFilterer creates a new log filterer instance of ReferralRewardManager, bound to a specific deployed contract.
func NewReferralRewardManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*ReferralRewardManagerFilterer, error) {
	contract, err := bindReferralRewardManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ReferralRewardManagerFilterer{contract: contract}, nil
}

// bindReferralRewardManager binds a generic wrapper to an already deployed contract.
func bindReferralRewardManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ReferralRewardManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ReferralRewardManager *ReferralRewardManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ReferralRewardManager.Contract.ReferralRewardManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ReferralRewardManager *ReferralRewardManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ReferralRewardManager.Contract.ReferralRewardManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ReferralRewardManager *ReferralRewardManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ReferralRewardManager.Contract.ReferralRewardManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ReferralRewardManager *ReferralRewardManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ReferralRewardManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ReferralRewardManager *ReferralRewardManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ReferralRewardManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ReferralRewardManager *ReferralRewardManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ReferralRewardManager.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_ReferralRewardManager *ReferralRewardManagerCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ReferralRewardManager.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_ReferralRewardManager *ReferralRewardManagerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _ReferralRewardManager.Contract.DEFAULTADMINROLE(&_ReferralRewardManager.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_ReferralRewardManager *ReferralRewardManagerCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _ReferralRewardManager.Contract.DEFAULTADMINROLE(&_ReferralRewardManager.CallOpts)
}

// FUNDMANAGERROLE is a free data retrieval call binding the contract method 0x2f83a6bc.
//
// Solidity: function FUND_MANAGER_ROLE() view returns(bytes32)
func (_ReferralRewardManager *ReferralRewardManagerCaller) FUNDMANAGERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ReferralRewardManager.contract.Call(opts, &out, "FUND_MANAGER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// FUNDMANAGERROLE is a free data retrieval call binding the contract method 0x2f83a6bc.
//
// Solidity: function FUND_MANAGER_ROLE() view returns(bytes32)
func (_ReferralRewardManager *ReferralRewardManagerSession) FUNDMANAGERROLE() ([32]byte, error) {
	return _ReferralRewardManager.Contract.FUNDMANAGERROLE(&_ReferralRewardManager.CallOpts)
}

// FUNDMANAGERROLE is a free data retrieval call binding the contract method 0x2f83a6bc.
//
// Solidity: function FUND_MANAGER_ROLE() view returns(bytes32)
func (_ReferralRewardManager *ReferralRewardManagerCallerSession) FUNDMANAGERROLE() ([32]byte, error) {
	return _ReferralRewardManager.Contract.FUNDMANAGERROLE(&_ReferralRewardManager.CallOpts)
}

// PAUSERROLE is a free data retrieval call binding the contract method 0xe63ab1e9.
//
// Solidity: function PAUSER_ROLE() view returns(bytes32)
func (_ReferralRewardManager *ReferralRewardManagerCaller) PAUSERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ReferralRewardManager.contract.Call(opts, &out, "PAUSER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// PAUSERROLE is a free data retrieval call binding the contract method 0xe63ab1e9.
//
// Solidity: function PAUSER_ROLE() view returns(bytes32)
func (_ReferralRewardManager *ReferralRewardManagerSession) PAUSERROLE() ([32]byte, error) {
	return _ReferralRewardManager.Contract.PAUSERROLE(&_ReferralRewardManager.CallOpts)
}

// PAUSERROLE is a free data retrieval call binding the contract method 0xe63ab1e9.
//
// Solidity: function PAUSER_ROLE() view returns(bytes32)
func (_ReferralRewardManager *ReferralRewardManagerCallerSession) PAUSERROLE() ([32]byte, error) {
	return _ReferralRewardManager.Contract.PAUSERROLE(&_ReferralRewardManager.CallOpts)
}

// REWARDMANAGERROLE is a free data retrieval call binding the contract method 0x6d750d80.
//
// Solidity: function REWARD_MANAGER_ROLE() view returns(bytes32)
func (_ReferralRewardManager *ReferralRewardManagerCaller) REWARDMANAGERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ReferralRewardManager.contract.Call(opts, &out, "REWARD_MANAGER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// REWARDMANAGERROLE is a free data retrieval call binding the contract method 0x6d750d80.
//
// Solidity: function REWARD_MANAGER_ROLE() view returns(bytes32)
func (_ReferralRewardManager *ReferralRewardManagerSession) REWARDMANAGERROLE() ([32]byte, error) {
	return _ReferralRewardManager.Contract.REWARDMANAGERROLE(&_ReferralRewardManager.CallOpts)
}

// REWARDMANAGERROLE is a free data retrieval call binding the contract method 0x6d750d80.
//
// Solidity: function REWARD_MANAGER_ROLE() view returns(bytes32)
func (_ReferralRewardManager *ReferralRewardManagerCallerSession) REWARDMANAGERROLE() ([32]byte, error) {
	return _ReferralRewardManager.Contract.REWARDMANAGERROLE(&_ReferralRewardManager.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_ReferralRewardManager *ReferralRewardManagerCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ReferralRewardManager.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_ReferralRewardManager *ReferralRewardManagerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _ReferralRewardManager.Contract.UPGRADEINTERFACEVERSION(&_ReferralRewardManager.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_ReferralRewardManager *ReferralRewardManagerCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _ReferralRewardManager.Contract.UPGRADEINTERFACEVERSION(&_ReferralRewardManager.CallOpts)
}

// GetConfig is a free data retrieval call binding the contract method 0xc3f909d4.
//
// Solidity: function getConfig() view returns(uint256 _minClaimAmount, uint256 _maxRewardPerUser, uint256 _maxReferralsPerUser)
func (_ReferralRewardManager *ReferralRewardManagerCaller) GetConfig(opts *bind.CallOpts) (struct {
	MinClaimAmount      *big.Int
	MaxRewardPerUser    *big.Int
	MaxReferralsPerUser *big.Int
}, error) {
	var out []interface{}
	err := _ReferralRewardManager.contract.Call(opts, &out, "getConfig")

	outstruct := new(struct {
		MinClaimAmount      *big.Int
		MaxRewardPerUser    *big.Int
		MaxReferralsPerUser *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.MinClaimAmount = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.MaxRewardPerUser = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.MaxReferralsPerUser = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetConfig is a free data retrieval call binding the contract method 0xc3f909d4.
//
// Solidity: function getConfig() view returns(uint256 _minClaimAmount, uint256 _maxRewardPerUser, uint256 _maxReferralsPerUser)
func (_ReferralRewardManager *ReferralRewardManagerSession) GetConfig() (struct {
	MinClaimAmount      *big.Int
	MaxRewardPerUser    *big.Int
	MaxReferralsPerUser *big.Int
}, error) {
	return _ReferralRewardManager.Contract.GetConfig(&_ReferralRewardManager.CallOpts)
}

// GetConfig is a free data retrieval call binding the contract method 0xc3f909d4.
//
// Solidity: function getConfig() view returns(uint256 _minClaimAmount, uint256 _maxRewardPerUser, uint256 _maxReferralsPerUser)
func (_ReferralRewardManager *ReferralRewardManagerCallerSession) GetConfig() (struct {
	MinClaimAmount      *big.Int
	MaxRewardPerUser    *big.Int
	MaxReferralsPerUser *big.Int
}, error) {
	return _ReferralRewardManager.Contract.GetConfig(&_ReferralRewardManager.CallOpts)
}

// GetReferrals is a free data retrieval call binding the contract method 0x41a0894d.
//
// Solidity: function getReferrals(address _referrer) view returns(address[])
func (_ReferralRewardManager *ReferralRewardManagerCaller) GetReferrals(opts *bind.CallOpts, _referrer common.Address) ([]common.Address, error) {
	var out []interface{}
	err := _ReferralRewardManager.contract.Call(opts, &out, "getReferrals", _referrer)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetReferrals is a free data retrieval call binding the contract method 0x41a0894d.
//
// Solidity: function getReferrals(address _referrer) view returns(address[])
func (_ReferralRewardManager *ReferralRewardManagerSession) GetReferrals(_referrer common.Address) ([]common.Address, error) {
	return _ReferralRewardManager.Contract.GetReferrals(&_ReferralRewardManager.CallOpts, _referrer)
}

// GetReferrals is a free data retrieval call binding the contract method 0x41a0894d.
//
// Solidity: function getReferrals(address _referrer) view returns(address[])
func (_ReferralRewardManager *ReferralRewardManagerCallerSession) GetReferrals(_referrer common.Address) ([]common.Address, error) {
	return _ReferralRewardManager.Contract.GetReferrals(&_ReferralRewardManager.CallOpts, _referrer)
}

// GetRewardPoolStatus is a free data retrieval call binding the contract method 0xcb05347a.
//
// Solidity: function getRewardPoolStatus() view returns(address poolAddress, uint256 balance, uint256 totalPending, uint256 totalClaimed)
func (_ReferralRewardManager *ReferralRewardManagerCaller) GetRewardPoolStatus(opts *bind.CallOpts) (struct {
	PoolAddress  common.Address
	Balance      *big.Int
	TotalPending *big.Int
	TotalClaimed *big.Int
}, error) {
	var out []interface{}
	err := _ReferralRewardManager.contract.Call(opts, &out, "getRewardPoolStatus")

	outstruct := new(struct {
		PoolAddress  common.Address
		Balance      *big.Int
		TotalPending *big.Int
		TotalClaimed *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.PoolAddress = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Balance = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.TotalPending = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.TotalClaimed = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetRewardPoolStatus is a free data retrieval call binding the contract method 0xcb05347a.
//
// Solidity: function getRewardPoolStatus() view returns(address poolAddress, uint256 balance, uint256 totalPending, uint256 totalClaimed)
func (_ReferralRewardManager *ReferralRewardManagerSession) GetRewardPoolStatus() (struct {
	PoolAddress  common.Address
	Balance      *big.Int
	TotalPending *big.Int
	TotalClaimed *big.Int
}, error) {
	return _ReferralRewardManager.Contract.GetRewardPoolStatus(&_ReferralRewardManager.CallOpts)
}

// GetRewardPoolStatus is a free data retrieval call binding the contract method 0xcb05347a.
//
// Solidity: function getRewardPoolStatus() view returns(address poolAddress, uint256 balance, uint256 totalPending, uint256 totalClaimed)
func (_ReferralRewardManager *ReferralRewardManagerCallerSession) GetRewardPoolStatus() (struct {
	PoolAddress  common.Address
	Balance      *big.Int
	TotalPending *big.Int
	TotalClaimed *big.Int
}, error) {
	return _ReferralRewardManager.Contract.GetRewardPoolStatus(&_ReferralRewardManager.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_ReferralRewardManager *ReferralRewardManagerCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _ReferralRewardManager.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_ReferralRewardManager *ReferralRewardManagerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _ReferralRewardManager.Contract.GetRoleAdmin(&_ReferralRewardManager.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_ReferralRewardManager *ReferralRewardManagerCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _ReferralRewardManager.Contract.GetRoleAdmin(&_ReferralRewardManager.CallOpts, role)
}

// GetSystemStats is a free data retrieval call binding the contract method 0xd1ce59a7.
//
// Solidity: function getSystemStats() view returns(uint256 _totalUsers, uint256 _totalReferrers, uint256 _totalPending, uint256 _totalClaimed)
func (_ReferralRewardManager *ReferralRewardManagerCaller) GetSystemStats(opts *bind.CallOpts) (struct {
	TotalUsers     *big.Int
	TotalReferrers *big.Int
	TotalPending   *big.Int
	TotalClaimed   *big.Int
}, error) {
	var out []interface{}
	err := _ReferralRewardManager.contract.Call(opts, &out, "getSystemStats")

	outstruct := new(struct {
		TotalUsers     *big.Int
		TotalReferrers *big.Int
		TotalPending   *big.Int
		TotalClaimed   *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.TotalUsers = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.TotalReferrers = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.TotalPending = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.TotalClaimed = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetSystemStats is a free data retrieval call binding the contract method 0xd1ce59a7.
//
// Solidity: function getSystemStats() view returns(uint256 _totalUsers, uint256 _totalReferrers, uint256 _totalPending, uint256 _totalClaimed)
func (_ReferralRewardManager *ReferralRewardManagerSession) GetSystemStats() (struct {
	TotalUsers     *big.Int
	TotalReferrers *big.Int
	TotalPending   *big.Int
	TotalClaimed   *big.Int
}, error) {
	return _ReferralRewardManager.Contract.GetSystemStats(&_ReferralRewardManager.CallOpts)
}

// GetSystemStats is a free data retrieval call binding the contract method 0xd1ce59a7.
//
// Solidity: function getSystemStats() view returns(uint256 _totalUsers, uint256 _totalReferrers, uint256 _totalPending, uint256 _totalClaimed)
func (_ReferralRewardManager *ReferralRewardManagerCallerSession) GetSystemStats() (struct {
	TotalUsers     *big.Int
	TotalReferrers *big.Int
	TotalPending   *big.Int
	TotalClaimed   *big.Int
}, error) {
	return _ReferralRewardManager.Contract.GetSystemStats(&_ReferralRewardManager.CallOpts)
}

// GetUserReferralInfo is a free data retrieval call binding the contract method 0xd3552712.
//
// Solidity: function getUserReferralInfo(address user) view returns(address _referrer, uint256 _referralCount, uint256 _pendingReward, uint256 _totalClaimed)
func (_ReferralRewardManager *ReferralRewardManagerCaller) GetUserReferralInfo(opts *bind.CallOpts, user common.Address) (struct {
	Referrer      common.Address
	ReferralCount *big.Int
	PendingReward *big.Int
	TotalClaimed  *big.Int
}, error) {
	var out []interface{}
	err := _ReferralRewardManager.contract.Call(opts, &out, "getUserReferralInfo", user)

	outstruct := new(struct {
		Referrer      common.Address
		ReferralCount *big.Int
		PendingReward *big.Int
		TotalClaimed  *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Referrer = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.ReferralCount = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.PendingReward = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.TotalClaimed = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetUserReferralInfo is a free data retrieval call binding the contract method 0xd3552712.
//
// Solidity: function getUserReferralInfo(address user) view returns(address _referrer, uint256 _referralCount, uint256 _pendingReward, uint256 _totalClaimed)
func (_ReferralRewardManager *ReferralRewardManagerSession) GetUserReferralInfo(user common.Address) (struct {
	Referrer      common.Address
	ReferralCount *big.Int
	PendingReward *big.Int
	TotalClaimed  *big.Int
}, error) {
	return _ReferralRewardManager.Contract.GetUserReferralInfo(&_ReferralRewardManager.CallOpts, user)
}

// GetUserReferralInfo is a free data retrieval call binding the contract method 0xd3552712.
//
// Solidity: function getUserReferralInfo(address user) view returns(address _referrer, uint256 _referralCount, uint256 _pendingReward, uint256 _totalClaimed)
func (_ReferralRewardManager *ReferralRewardManagerCallerSession) GetUserReferralInfo(user common.Address) (struct {
	Referrer      common.Address
	ReferralCount *big.Int
	PendingReward *big.Int
	TotalClaimed  *big.Int
}, error) {
	return _ReferralRewardManager.Contract.GetUserReferralInfo(&_ReferralRewardManager.CallOpts, user)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_ReferralRewardManager *ReferralRewardManagerCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _ReferralRewardManager.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_ReferralRewardManager *ReferralRewardManagerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _ReferralRewardManager.Contract.HasRole(&_ReferralRewardManager.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_ReferralRewardManager *ReferralRewardManagerCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _ReferralRewardManager.Contract.HasRole(&_ReferralRewardManager.CallOpts, role, account)
}

// MaxReferralsPerUser is a free data retrieval call binding the contract method 0x38196bc7.
//
// Solidity: function maxReferralsPerUser() view returns(uint16)
func (_ReferralRewardManager *ReferralRewardManagerCaller) MaxReferralsPerUser(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _ReferralRewardManager.contract.Call(opts, &out, "maxReferralsPerUser")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// MaxReferralsPerUser is a free data retrieval call binding the contract method 0x38196bc7.
//
// Solidity: function maxReferralsPerUser() view returns(uint16)
func (_ReferralRewardManager *ReferralRewardManagerSession) MaxReferralsPerUser() (uint16, error) {
	return _ReferralRewardManager.Contract.MaxReferralsPerUser(&_ReferralRewardManager.CallOpts)
}

// MaxReferralsPerUser is a free data retrieval call binding the contract method 0x38196bc7.
//
// Solidity: function maxReferralsPerUser() view returns(uint16)
func (_ReferralRewardManager *ReferralRewardManagerCallerSession) MaxReferralsPerUser() (uint16, error) {
	return _ReferralRewardManager.Contract.MaxReferralsPerUser(&_ReferralRewardManager.CallOpts)
}

// MaxRewardPerUser is a free data retrieval call binding the contract method 0xad8d6f2c.
//
// Solidity: function maxRewardPerUser() view returns(uint256)
func (_ReferralRewardManager *ReferralRewardManagerCaller) MaxRewardPerUser(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ReferralRewardManager.contract.Call(opts, &out, "maxRewardPerUser")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxRewardPerUser is a free data retrieval call binding the contract method 0xad8d6f2c.
//
// Solidity: function maxRewardPerUser() view returns(uint256)
func (_ReferralRewardManager *ReferralRewardManagerSession) MaxRewardPerUser() (*big.Int, error) {
	return _ReferralRewardManager.Contract.MaxRewardPerUser(&_ReferralRewardManager.CallOpts)
}

// MaxRewardPerUser is a free data retrieval call binding the contract method 0xad8d6f2c.
//
// Solidity: function maxRewardPerUser() view returns(uint256)
func (_ReferralRewardManager *ReferralRewardManagerCallerSession) MaxRewardPerUser() (*big.Int, error) {
	return _ReferralRewardManager.Contract.MaxRewardPerUser(&_ReferralRewardManager.CallOpts)
}

// MinClaimAmount is a free data retrieval call binding the contract method 0xadc25bde.
//
// Solidity: function minClaimAmount() view returns(uint256)
func (_ReferralRewardManager *ReferralRewardManagerCaller) MinClaimAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ReferralRewardManager.contract.Call(opts, &out, "minClaimAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinClaimAmount is a free data retrieval call binding the contract method 0xadc25bde.
//
// Solidity: function minClaimAmount() view returns(uint256)
func (_ReferralRewardManager *ReferralRewardManagerSession) MinClaimAmount() (*big.Int, error) {
	return _ReferralRewardManager.Contract.MinClaimAmount(&_ReferralRewardManager.CallOpts)
}

// MinClaimAmount is a free data retrieval call binding the contract method 0xadc25bde.
//
// Solidity: function minClaimAmount() view returns(uint256)
func (_ReferralRewardManager *ReferralRewardManagerCallerSession) MinClaimAmount() (*big.Int, error) {
	return _ReferralRewardManager.Contract.MinClaimAmount(&_ReferralRewardManager.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_ReferralRewardManager *ReferralRewardManagerCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _ReferralRewardManager.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_ReferralRewardManager *ReferralRewardManagerSession) Paused() (bool, error) {
	return _ReferralRewardManager.Contract.Paused(&_ReferralRewardManager.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_ReferralRewardManager *ReferralRewardManagerCallerSession) Paused() (bool, error) {
	return _ReferralRewardManager.Contract.Paused(&_ReferralRewardManager.CallOpts)
}

// PendingRewards is a free data retrieval call binding the contract method 0x31d7a262.
//
// Solidity: function pendingRewards(address ) view returns(uint256)
func (_ReferralRewardManager *ReferralRewardManagerCaller) PendingRewards(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ReferralRewardManager.contract.Call(opts, &out, "pendingRewards", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PendingRewards is a free data retrieval call binding the contract method 0x31d7a262.
//
// Solidity: function pendingRewards(address ) view returns(uint256)
func (_ReferralRewardManager *ReferralRewardManagerSession) PendingRewards(arg0 common.Address) (*big.Int, error) {
	return _ReferralRewardManager.Contract.PendingRewards(&_ReferralRewardManager.CallOpts, arg0)
}

// PendingRewards is a free data retrieval call binding the contract method 0x31d7a262.
//
// Solidity: function pendingRewards(address ) view returns(uint256)
func (_ReferralRewardManager *ReferralRewardManagerCallerSession) PendingRewards(arg0 common.Address) (*big.Int, error) {
	return _ReferralRewardManager.Contract.PendingRewards(&_ReferralRewardManager.CallOpts, arg0)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_ReferralRewardManager *ReferralRewardManagerCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ReferralRewardManager.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_ReferralRewardManager *ReferralRewardManagerSession) ProxiableUUID() ([32]byte, error) {
	return _ReferralRewardManager.Contract.ProxiableUUID(&_ReferralRewardManager.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_ReferralRewardManager *ReferralRewardManagerCallerSession) ProxiableUUID() ([32]byte, error) {
	return _ReferralRewardManager.Contract.ProxiableUUID(&_ReferralRewardManager.CallOpts)
}

// ReferralCount is a free data retrieval call binding the contract method 0xdb74559b.
//
// Solidity: function referralCount(address ) view returns(uint256)
func (_ReferralRewardManager *ReferralRewardManagerCaller) ReferralCount(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ReferralRewardManager.contract.Call(opts, &out, "referralCount", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ReferralCount is a free data retrieval call binding the contract method 0xdb74559b.
//
// Solidity: function referralCount(address ) view returns(uint256)
func (_ReferralRewardManager *ReferralRewardManagerSession) ReferralCount(arg0 common.Address) (*big.Int, error) {
	return _ReferralRewardManager.Contract.ReferralCount(&_ReferralRewardManager.CallOpts, arg0)
}

// ReferralCount is a free data retrieval call binding the contract method 0xdb74559b.
//
// Solidity: function referralCount(address ) view returns(uint256)
func (_ReferralRewardManager *ReferralRewardManagerCallerSession) ReferralCount(arg0 common.Address) (*big.Int, error) {
	return _ReferralRewardManager.Contract.ReferralCount(&_ReferralRewardManager.CallOpts, arg0)
}

// Referrals is a free data retrieval call binding the contract method 0xe697b5d8.
//
// Solidity: function referrals(address , uint256 ) view returns(address)
func (_ReferralRewardManager *ReferralRewardManagerCaller) Referrals(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _ReferralRewardManager.contract.Call(opts, &out, "referrals", arg0, arg1)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Referrals is a free data retrieval call binding the contract method 0xe697b5d8.
//
// Solidity: function referrals(address , uint256 ) view returns(address)
func (_ReferralRewardManager *ReferralRewardManagerSession) Referrals(arg0 common.Address, arg1 *big.Int) (common.Address, error) {
	return _ReferralRewardManager.Contract.Referrals(&_ReferralRewardManager.CallOpts, arg0, arg1)
}

// Referrals is a free data retrieval call binding the contract method 0xe697b5d8.
//
// Solidity: function referrals(address , uint256 ) view returns(address)
func (_ReferralRewardManager *ReferralRewardManagerCallerSession) Referrals(arg0 common.Address, arg1 *big.Int) (common.Address, error) {
	return _ReferralRewardManager.Contract.Referrals(&_ReferralRewardManager.CallOpts, arg0, arg1)
}

// Referrer is a free data retrieval call binding the contract method 0x2cf003c2.
//
// Solidity: function referrer(address ) view returns(address)
func (_ReferralRewardManager *ReferralRewardManagerCaller) Referrer(opts *bind.CallOpts, arg0 common.Address) (common.Address, error) {
	var out []interface{}
	err := _ReferralRewardManager.contract.Call(opts, &out, "referrer", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Referrer is a free data retrieval call binding the contract method 0x2cf003c2.
//
// Solidity: function referrer(address ) view returns(address)
func (_ReferralRewardManager *ReferralRewardManagerSession) Referrer(arg0 common.Address) (common.Address, error) {
	return _ReferralRewardManager.Contract.Referrer(&_ReferralRewardManager.CallOpts, arg0)
}

// Referrer is a free data retrieval call binding the contract method 0x2cf003c2.
//
// Solidity: function referrer(address ) view returns(address)
func (_ReferralRewardManager *ReferralRewardManagerCallerSession) Referrer(arg0 common.Address) (common.Address, error) {
	return _ReferralRewardManager.Contract.Referrer(&_ReferralRewardManager.CallOpts, arg0)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_ReferralRewardManager *ReferralRewardManagerCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _ReferralRewardManager.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_ReferralRewardManager *ReferralRewardManagerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ReferralRewardManager.Contract.SupportsInterface(&_ReferralRewardManager.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_ReferralRewardManager *ReferralRewardManagerCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ReferralRewardManager.Contract.SupportsInterface(&_ReferralRewardManager.CallOpts, interfaceId)
}

// TotalClaimedRewards is a free data retrieval call binding the contract method 0xa110b93f.
//
// Solidity: function totalClaimedRewards(address ) view returns(uint256)
func (_ReferralRewardManager *ReferralRewardManagerCaller) TotalClaimedRewards(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ReferralRewardManager.contract.Call(opts, &out, "totalClaimedRewards", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalClaimedRewards is a free data retrieval call binding the contract method 0xa110b93f.
//
// Solidity: function totalClaimedRewards(address ) view returns(uint256)
func (_ReferralRewardManager *ReferralRewardManagerSession) TotalClaimedRewards(arg0 common.Address) (*big.Int, error) {
	return _ReferralRewardManager.Contract.TotalClaimedRewards(&_ReferralRewardManager.CallOpts, arg0)
}

// TotalClaimedRewards is a free data retrieval call binding the contract method 0xa110b93f.
//
// Solidity: function totalClaimedRewards(address ) view returns(uint256)
func (_ReferralRewardManager *ReferralRewardManagerCallerSession) TotalClaimedRewards(arg0 common.Address) (*big.Int, error) {
	return _ReferralRewardManager.Contract.TotalClaimedRewards(&_ReferralRewardManager.CallOpts, arg0)
}

// TotalClaimedRewardsGlobal is a free data retrieval call binding the contract method 0x9e6c6124.
//
// Solidity: function totalClaimedRewardsGlobal() view returns(uint256)
func (_ReferralRewardManager *ReferralRewardManagerCaller) TotalClaimedRewardsGlobal(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ReferralRewardManager.contract.Call(opts, &out, "totalClaimedRewardsGlobal")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalClaimedRewardsGlobal is a free data retrieval call binding the contract method 0x9e6c6124.
//
// Solidity: function totalClaimedRewardsGlobal() view returns(uint256)
func (_ReferralRewardManager *ReferralRewardManagerSession) TotalClaimedRewardsGlobal() (*big.Int, error) {
	return _ReferralRewardManager.Contract.TotalClaimedRewardsGlobal(&_ReferralRewardManager.CallOpts)
}

// TotalClaimedRewardsGlobal is a free data retrieval call binding the contract method 0x9e6c6124.
//
// Solidity: function totalClaimedRewardsGlobal() view returns(uint256)
func (_ReferralRewardManager *ReferralRewardManagerCallerSession) TotalClaimedRewardsGlobal() (*big.Int, error) {
	return _ReferralRewardManager.Contract.TotalClaimedRewardsGlobal(&_ReferralRewardManager.CallOpts)
}

// TotalPendingRewards is a free data retrieval call binding the contract method 0xec1371f2.
//
// Solidity: function totalPendingRewards() view returns(uint256)
func (_ReferralRewardManager *ReferralRewardManagerCaller) TotalPendingRewards(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ReferralRewardManager.contract.Call(opts, &out, "totalPendingRewards")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalPendingRewards is a free data retrieval call binding the contract method 0xec1371f2.
//
// Solidity: function totalPendingRewards() view returns(uint256)
func (_ReferralRewardManager *ReferralRewardManagerSession) TotalPendingRewards() (*big.Int, error) {
	return _ReferralRewardManager.Contract.TotalPendingRewards(&_ReferralRewardManager.CallOpts)
}

// TotalPendingRewards is a free data retrieval call binding the contract method 0xec1371f2.
//
// Solidity: function totalPendingRewards() view returns(uint256)
func (_ReferralRewardManager *ReferralRewardManagerCallerSession) TotalPendingRewards() (*big.Int, error) {
	return _ReferralRewardManager.Contract.TotalPendingRewards(&_ReferralRewardManager.CallOpts)
}

// TotalReferrers is a free data retrieval call binding the contract method 0x62d03cb7.
//
// Solidity: function totalReferrers() view returns(uint256)
func (_ReferralRewardManager *ReferralRewardManagerCaller) TotalReferrers(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ReferralRewardManager.contract.Call(opts, &out, "totalReferrers")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalReferrers is a free data retrieval call binding the contract method 0x62d03cb7.
//
// Solidity: function totalReferrers() view returns(uint256)
func (_ReferralRewardManager *ReferralRewardManagerSession) TotalReferrers() (*big.Int, error) {
	return _ReferralRewardManager.Contract.TotalReferrers(&_ReferralRewardManager.CallOpts)
}

// TotalReferrers is a free data retrieval call binding the contract method 0x62d03cb7.
//
// Solidity: function totalReferrers() view returns(uint256)
func (_ReferralRewardManager *ReferralRewardManagerCallerSession) TotalReferrers() (*big.Int, error) {
	return _ReferralRewardManager.Contract.TotalReferrers(&_ReferralRewardManager.CallOpts)
}

// TotalUsers is a free data retrieval call binding the contract method 0xbff1f9e1.
//
// Solidity: function totalUsers() view returns(uint256)
func (_ReferralRewardManager *ReferralRewardManagerCaller) TotalUsers(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ReferralRewardManager.contract.Call(opts, &out, "totalUsers")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalUsers is a free data retrieval call binding the contract method 0xbff1f9e1.
//
// Solidity: function totalUsers() view returns(uint256)
func (_ReferralRewardManager *ReferralRewardManagerSession) TotalUsers() (*big.Int, error) {
	return _ReferralRewardManager.Contract.TotalUsers(&_ReferralRewardManager.CallOpts)
}

// TotalUsers is a free data retrieval call binding the contract method 0xbff1f9e1.
//
// Solidity: function totalUsers() view returns(uint256)
func (_ReferralRewardManager *ReferralRewardManagerCallerSession) TotalUsers() (*big.Int, error) {
	return _ReferralRewardManager.Contract.TotalUsers(&_ReferralRewardManager.CallOpts)
}

// YpusdToken is a free data retrieval call binding the contract method 0xf339ae78.
//
// Solidity: function ypusdToken() view returns(address)
func (_ReferralRewardManager *ReferralRewardManagerCaller) YpusdToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ReferralRewardManager.contract.Call(opts, &out, "ypusdToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// YpusdToken is a free data retrieval call binding the contract method 0xf339ae78.
//
// Solidity: function ypusdToken() view returns(address)
func (_ReferralRewardManager *ReferralRewardManagerSession) YpusdToken() (common.Address, error) {
	return _ReferralRewardManager.Contract.YpusdToken(&_ReferralRewardManager.CallOpts)
}

// YpusdToken is a free data retrieval call binding the contract method 0xf339ae78.
//
// Solidity: function ypusdToken() view returns(address)
func (_ReferralRewardManager *ReferralRewardManagerCallerSession) YpusdToken() (common.Address, error) {
	return _ReferralRewardManager.Contract.YpusdToken(&_ReferralRewardManager.CallOpts)
}

// BatchAddRewards is a paid mutator transaction binding the contract method 0xe3e9dfba.
//
// Solidity: function batchAddRewards(address[] users, uint256[] amounts) returns()
func (_ReferralRewardManager *ReferralRewardManagerTransactor) BatchAddRewards(opts *bind.TransactOpts, users []common.Address, amounts []*big.Int) (*types.Transaction, error) {
	return _ReferralRewardManager.contract.Transact(opts, "batchAddRewards", users, amounts)
}

// BatchAddRewards is a paid mutator transaction binding the contract method 0xe3e9dfba.
//
// Solidity: function batchAddRewards(address[] users, uint256[] amounts) returns()
func (_ReferralRewardManager *ReferralRewardManagerSession) BatchAddRewards(users []common.Address, amounts []*big.Int) (*types.Transaction, error) {
	return _ReferralRewardManager.Contract.BatchAddRewards(&_ReferralRewardManager.TransactOpts, users, amounts)
}

// BatchAddRewards is a paid mutator transaction binding the contract method 0xe3e9dfba.
//
// Solidity: function batchAddRewards(address[] users, uint256[] amounts) returns()
func (_ReferralRewardManager *ReferralRewardManagerTransactorSession) BatchAddRewards(users []common.Address, amounts []*big.Int) (*types.Transaction, error) {
	return _ReferralRewardManager.Contract.BatchAddRewards(&_ReferralRewardManager.TransactOpts, users, amounts)
}

// BatchClearRewards is a paid mutator transaction binding the contract method 0x30eae166.
//
// Solidity: function batchClearRewards(address[] users) returns()
func (_ReferralRewardManager *ReferralRewardManagerTransactor) BatchClearRewards(opts *bind.TransactOpts, users []common.Address) (*types.Transaction, error) {
	return _ReferralRewardManager.contract.Transact(opts, "batchClearRewards", users)
}

// BatchClearRewards is a paid mutator transaction binding the contract method 0x30eae166.
//
// Solidity: function batchClearRewards(address[] users) returns()
func (_ReferralRewardManager *ReferralRewardManagerSession) BatchClearRewards(users []common.Address) (*types.Transaction, error) {
	return _ReferralRewardManager.Contract.BatchClearRewards(&_ReferralRewardManager.TransactOpts, users)
}

// BatchClearRewards is a paid mutator transaction binding the contract method 0x30eae166.
//
// Solidity: function batchClearRewards(address[] users) returns()
func (_ReferralRewardManager *ReferralRewardManagerTransactorSession) BatchClearRewards(users []common.Address) (*types.Transaction, error) {
	return _ReferralRewardManager.Contract.BatchClearRewards(&_ReferralRewardManager.TransactOpts, users)
}

// BatchReduceRewards is a paid mutator transaction binding the contract method 0x6fdebdc0.
//
// Solidity: function batchReduceRewards(address[] users, uint256[] amounts) returns()
func (_ReferralRewardManager *ReferralRewardManagerTransactor) BatchReduceRewards(opts *bind.TransactOpts, users []common.Address, amounts []*big.Int) (*types.Transaction, error) {
	return _ReferralRewardManager.contract.Transact(opts, "batchReduceRewards", users, amounts)
}

// BatchReduceRewards is a paid mutator transaction binding the contract method 0x6fdebdc0.
//
// Solidity: function batchReduceRewards(address[] users, uint256[] amounts) returns()
func (_ReferralRewardManager *ReferralRewardManagerSession) BatchReduceRewards(users []common.Address, amounts []*big.Int) (*types.Transaction, error) {
	return _ReferralRewardManager.Contract.BatchReduceRewards(&_ReferralRewardManager.TransactOpts, users, amounts)
}

// BatchReduceRewards is a paid mutator transaction binding the contract method 0x6fdebdc0.
//
// Solidity: function batchReduceRewards(address[] users, uint256[] amounts) returns()
func (_ReferralRewardManager *ReferralRewardManagerTransactorSession) BatchReduceRewards(users []common.Address, amounts []*big.Int) (*types.Transaction, error) {
	return _ReferralRewardManager.Contract.BatchReduceRewards(&_ReferralRewardManager.TransactOpts, users, amounts)
}

// BatchSetRewards is a paid mutator transaction binding the contract method 0x6983e74d.
//
// Solidity: function batchSetRewards(address[] users, uint256[] amounts) returns()
func (_ReferralRewardManager *ReferralRewardManagerTransactor) BatchSetRewards(opts *bind.TransactOpts, users []common.Address, amounts []*big.Int) (*types.Transaction, error) {
	return _ReferralRewardManager.contract.Transact(opts, "batchSetRewards", users, amounts)
}

// BatchSetRewards is a paid mutator transaction binding the contract method 0x6983e74d.
//
// Solidity: function batchSetRewards(address[] users, uint256[] amounts) returns()
func (_ReferralRewardManager *ReferralRewardManagerSession) BatchSetRewards(users []common.Address, amounts []*big.Int) (*types.Transaction, error) {
	return _ReferralRewardManager.Contract.BatchSetRewards(&_ReferralRewardManager.TransactOpts, users, amounts)
}

// BatchSetRewards is a paid mutator transaction binding the contract method 0x6983e74d.
//
// Solidity: function batchSetRewards(address[] users, uint256[] amounts) returns()
func (_ReferralRewardManager *ReferralRewardManagerTransactorSession) BatchSetRewards(users []common.Address, amounts []*big.Int) (*types.Transaction, error) {
	return _ReferralRewardManager.Contract.BatchSetRewards(&_ReferralRewardManager.TransactOpts, users, amounts)
}

// ClaimReward is a paid mutator transaction binding the contract method 0xb88a802f.
//
// Solidity: function claimReward() returns()
func (_ReferralRewardManager *ReferralRewardManagerTransactor) ClaimReward(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ReferralRewardManager.contract.Transact(opts, "claimReward")
}

// ClaimReward is a paid mutator transaction binding the contract method 0xb88a802f.
//
// Solidity: function claimReward() returns()
func (_ReferralRewardManager *ReferralRewardManagerSession) ClaimReward() (*types.Transaction, error) {
	return _ReferralRewardManager.Contract.ClaimReward(&_ReferralRewardManager.TransactOpts)
}

// ClaimReward is a paid mutator transaction binding the contract method 0xb88a802f.
//
// Solidity: function claimReward() returns()
func (_ReferralRewardManager *ReferralRewardManagerTransactorSession) ClaimReward() (*types.Transaction, error) {
	return _ReferralRewardManager.Contract.ClaimReward(&_ReferralRewardManager.TransactOpts)
}

// FundRewardPool is a paid mutator transaction binding the contract method 0x1d583e0d.
//
// Solidity: function fundRewardPool(uint256 amount) returns()
func (_ReferralRewardManager *ReferralRewardManagerTransactor) FundRewardPool(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _ReferralRewardManager.contract.Transact(opts, "fundRewardPool", amount)
}

// FundRewardPool is a paid mutator transaction binding the contract method 0x1d583e0d.
//
// Solidity: function fundRewardPool(uint256 amount) returns()
func (_ReferralRewardManager *ReferralRewardManagerSession) FundRewardPool(amount *big.Int) (*types.Transaction, error) {
	return _ReferralRewardManager.Contract.FundRewardPool(&_ReferralRewardManager.TransactOpts, amount)
}

// FundRewardPool is a paid mutator transaction binding the contract method 0x1d583e0d.
//
// Solidity: function fundRewardPool(uint256 amount) returns()
func (_ReferralRewardManager *ReferralRewardManagerTransactorSession) FundRewardPool(amount *big.Int) (*types.Transaction, error) {
	return _ReferralRewardManager.Contract.FundRewardPool(&_ReferralRewardManager.TransactOpts, amount)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_ReferralRewardManager *ReferralRewardManagerTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ReferralRewardManager.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_ReferralRewardManager *ReferralRewardManagerSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ReferralRewardManager.Contract.GrantRole(&_ReferralRewardManager.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_ReferralRewardManager *ReferralRewardManagerTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ReferralRewardManager.Contract.GrantRole(&_ReferralRewardManager.TransactOpts, role, account)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address admin, address _ypusdToken) returns()
func (_ReferralRewardManager *ReferralRewardManagerTransactor) Initialize(opts *bind.TransactOpts, admin common.Address, _ypusdToken common.Address) (*types.Transaction, error) {
	return _ReferralRewardManager.contract.Transact(opts, "initialize", admin, _ypusdToken)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address admin, address _ypusdToken) returns()
func (_ReferralRewardManager *ReferralRewardManagerSession) Initialize(admin common.Address, _ypusdToken common.Address) (*types.Transaction, error) {
	return _ReferralRewardManager.Contract.Initialize(&_ReferralRewardManager.TransactOpts, admin, _ypusdToken)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address admin, address _ypusdToken) returns()
func (_ReferralRewardManager *ReferralRewardManagerTransactorSession) Initialize(admin common.Address, _ypusdToken common.Address) (*types.Transaction, error) {
	return _ReferralRewardManager.Contract.Initialize(&_ReferralRewardManager.TransactOpts, admin, _ypusdToken)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_ReferralRewardManager *ReferralRewardManagerTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ReferralRewardManager.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_ReferralRewardManager *ReferralRewardManagerSession) Pause() (*types.Transaction, error) {
	return _ReferralRewardManager.Contract.Pause(&_ReferralRewardManager.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_ReferralRewardManager *ReferralRewardManagerTransactorSession) Pause() (*types.Transaction, error) {
	return _ReferralRewardManager.Contract.Pause(&_ReferralRewardManager.TransactOpts)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_ReferralRewardManager *ReferralRewardManagerTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _ReferralRewardManager.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_ReferralRewardManager *ReferralRewardManagerSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _ReferralRewardManager.Contract.RenounceRole(&_ReferralRewardManager.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_ReferralRewardManager *ReferralRewardManagerTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _ReferralRewardManager.Contract.RenounceRole(&_ReferralRewardManager.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_ReferralRewardManager *ReferralRewardManagerTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ReferralRewardManager.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_ReferralRewardManager *ReferralRewardManagerSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ReferralRewardManager.Contract.RevokeRole(&_ReferralRewardManager.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_ReferralRewardManager *ReferralRewardManagerTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ReferralRewardManager.Contract.RevokeRole(&_ReferralRewardManager.TransactOpts, role, account)
}

// SetReferrer is a paid mutator transaction binding the contract method 0xa18a7bfc.
//
// Solidity: function setReferrer(address _referrer) returns()
func (_ReferralRewardManager *ReferralRewardManagerTransactor) SetReferrer(opts *bind.TransactOpts, _referrer common.Address) (*types.Transaction, error) {
	return _ReferralRewardManager.contract.Transact(opts, "setReferrer", _referrer)
}

// SetReferrer is a paid mutator transaction binding the contract method 0xa18a7bfc.
//
// Solidity: function setReferrer(address _referrer) returns()
func (_ReferralRewardManager *ReferralRewardManagerSession) SetReferrer(_referrer common.Address) (*types.Transaction, error) {
	return _ReferralRewardManager.Contract.SetReferrer(&_ReferralRewardManager.TransactOpts, _referrer)
}

// SetReferrer is a paid mutator transaction binding the contract method 0xa18a7bfc.
//
// Solidity: function setReferrer(address _referrer) returns()
func (_ReferralRewardManager *ReferralRewardManagerTransactorSession) SetReferrer(_referrer common.Address) (*types.Transaction, error) {
	return _ReferralRewardManager.Contract.SetReferrer(&_ReferralRewardManager.TransactOpts, _referrer)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_ReferralRewardManager *ReferralRewardManagerTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ReferralRewardManager.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_ReferralRewardManager *ReferralRewardManagerSession) Unpause() (*types.Transaction, error) {
	return _ReferralRewardManager.Contract.Unpause(&_ReferralRewardManager.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_ReferralRewardManager *ReferralRewardManagerTransactorSession) Unpause() (*types.Transaction, error) {
	return _ReferralRewardManager.Contract.Unpause(&_ReferralRewardManager.TransactOpts)
}

// UpdateConfig is a paid mutator transaction binding the contract method 0x73ad4946.
//
// Solidity: function updateConfig(uint256 _minClaimAmount, uint256 _maxRewardPerUser, uint256 _maxReferralsPerUser) returns()
func (_ReferralRewardManager *ReferralRewardManagerTransactor) UpdateConfig(opts *bind.TransactOpts, _minClaimAmount *big.Int, _maxRewardPerUser *big.Int, _maxReferralsPerUser *big.Int) (*types.Transaction, error) {
	return _ReferralRewardManager.contract.Transact(opts, "updateConfig", _minClaimAmount, _maxRewardPerUser, _maxReferralsPerUser)
}

// UpdateConfig is a paid mutator transaction binding the contract method 0x73ad4946.
//
// Solidity: function updateConfig(uint256 _minClaimAmount, uint256 _maxRewardPerUser, uint256 _maxReferralsPerUser) returns()
func (_ReferralRewardManager *ReferralRewardManagerSession) UpdateConfig(_minClaimAmount *big.Int, _maxRewardPerUser *big.Int, _maxReferralsPerUser *big.Int) (*types.Transaction, error) {
	return _ReferralRewardManager.Contract.UpdateConfig(&_ReferralRewardManager.TransactOpts, _minClaimAmount, _maxRewardPerUser, _maxReferralsPerUser)
}

// UpdateConfig is a paid mutator transaction binding the contract method 0x73ad4946.
//
// Solidity: function updateConfig(uint256 _minClaimAmount, uint256 _maxRewardPerUser, uint256 _maxReferralsPerUser) returns()
func (_ReferralRewardManager *ReferralRewardManagerTransactorSession) UpdateConfig(_minClaimAmount *big.Int, _maxRewardPerUser *big.Int, _maxReferralsPerUser *big.Int) (*types.Transaction, error) {
	return _ReferralRewardManager.Contract.UpdateConfig(&_ReferralRewardManager.TransactOpts, _minClaimAmount, _maxRewardPerUser, _maxReferralsPerUser)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_ReferralRewardManager *ReferralRewardManagerTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _ReferralRewardManager.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_ReferralRewardManager *ReferralRewardManagerSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _ReferralRewardManager.Contract.UpgradeToAndCall(&_ReferralRewardManager.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_ReferralRewardManager *ReferralRewardManagerTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _ReferralRewardManager.Contract.UpgradeToAndCall(&_ReferralRewardManager.TransactOpts, newImplementation, data)
}

// WithdrawFunds is a paid mutator transaction binding the contract method 0x155dd5ee.
//
// Solidity: function withdrawFunds(uint256 amount) returns()
func (_ReferralRewardManager *ReferralRewardManagerTransactor) WithdrawFunds(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _ReferralRewardManager.contract.Transact(opts, "withdrawFunds", amount)
}

// WithdrawFunds is a paid mutator transaction binding the contract method 0x155dd5ee.
//
// Solidity: function withdrawFunds(uint256 amount) returns()
func (_ReferralRewardManager *ReferralRewardManagerSession) WithdrawFunds(amount *big.Int) (*types.Transaction, error) {
	return _ReferralRewardManager.Contract.WithdrawFunds(&_ReferralRewardManager.TransactOpts, amount)
}

// WithdrawFunds is a paid mutator transaction binding the contract method 0x155dd5ee.
//
// Solidity: function withdrawFunds(uint256 amount) returns()
func (_ReferralRewardManager *ReferralRewardManagerTransactorSession) WithdrawFunds(amount *big.Int) (*types.Transaction, error) {
	return _ReferralRewardManager.Contract.WithdrawFunds(&_ReferralRewardManager.TransactOpts, amount)
}

// ReferralRewardManagerConfigUpdatedIterator is returned from FilterConfigUpdated and is used to iterate over the raw logs and unpacked data for ConfigUpdated events raised by the ReferralRewardManager contract.
type ReferralRewardManagerConfigUpdatedIterator struct {
	Event *ReferralRewardManagerConfigUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ReferralRewardManagerConfigUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReferralRewardManagerConfigUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ReferralRewardManagerConfigUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ReferralRewardManagerConfigUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReferralRewardManagerConfigUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReferralRewardManagerConfigUpdated represents a ConfigUpdated event raised by the ReferralRewardManager contract.
type ReferralRewardManagerConfigUpdated struct {
	MinClaimAmount      *big.Int
	MaxRewardPerUser    *big.Int
	MaxReferralsPerUser *big.Int
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterConfigUpdated is a free log retrieval operation binding the contract event 0x6d05db271e19f930af71c4765de54ef86294762644c20f4d6fd2609d057d3c7b.
//
// Solidity: event ConfigUpdated(uint256 minClaimAmount, uint256 maxRewardPerUser, uint256 maxReferralsPerUser)
func (_ReferralRewardManager *ReferralRewardManagerFilterer) FilterConfigUpdated(opts *bind.FilterOpts) (*ReferralRewardManagerConfigUpdatedIterator, error) {

	logs, sub, err := _ReferralRewardManager.contract.FilterLogs(opts, "ConfigUpdated")
	if err != nil {
		return nil, err
	}
	return &ReferralRewardManagerConfigUpdatedIterator{contract: _ReferralRewardManager.contract, event: "ConfigUpdated", logs: logs, sub: sub}, nil
}

// WatchConfigUpdated is a free log subscription operation binding the contract event 0x6d05db271e19f930af71c4765de54ef86294762644c20f4d6fd2609d057d3c7b.
//
// Solidity: event ConfigUpdated(uint256 minClaimAmount, uint256 maxRewardPerUser, uint256 maxReferralsPerUser)
func (_ReferralRewardManager *ReferralRewardManagerFilterer) WatchConfigUpdated(opts *bind.WatchOpts, sink chan<- *ReferralRewardManagerConfigUpdated) (event.Subscription, error) {

	logs, sub, err := _ReferralRewardManager.contract.WatchLogs(opts, "ConfigUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReferralRewardManagerConfigUpdated)
				if err := _ReferralRewardManager.contract.UnpackLog(event, "ConfigUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseConfigUpdated is a log parse operation binding the contract event 0x6d05db271e19f930af71c4765de54ef86294762644c20f4d6fd2609d057d3c7b.
//
// Solidity: event ConfigUpdated(uint256 minClaimAmount, uint256 maxRewardPerUser, uint256 maxReferralsPerUser)
func (_ReferralRewardManager *ReferralRewardManagerFilterer) ParseConfigUpdated(log types.Log) (*ReferralRewardManagerConfigUpdated, error) {
	event := new(ReferralRewardManagerConfigUpdated)
	if err := _ReferralRewardManager.contract.UnpackLog(event, "ConfigUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReferralRewardManagerInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the ReferralRewardManager contract.
type ReferralRewardManagerInitializedIterator struct {
	Event *ReferralRewardManagerInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ReferralRewardManagerInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReferralRewardManagerInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ReferralRewardManagerInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ReferralRewardManagerInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReferralRewardManagerInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReferralRewardManagerInitialized represents a Initialized event raised by the ReferralRewardManager contract.
type ReferralRewardManagerInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_ReferralRewardManager *ReferralRewardManagerFilterer) FilterInitialized(opts *bind.FilterOpts) (*ReferralRewardManagerInitializedIterator, error) {

	logs, sub, err := _ReferralRewardManager.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ReferralRewardManagerInitializedIterator{contract: _ReferralRewardManager.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_ReferralRewardManager *ReferralRewardManagerFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ReferralRewardManagerInitialized) (event.Subscription, error) {

	logs, sub, err := _ReferralRewardManager.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReferralRewardManagerInitialized)
				if err := _ReferralRewardManager.contract.UnpackLog(event, "Initialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitialized is a log parse operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_ReferralRewardManager *ReferralRewardManagerFilterer) ParseInitialized(log types.Log) (*ReferralRewardManagerInitialized, error) {
	event := new(ReferralRewardManagerInitialized)
	if err := _ReferralRewardManager.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReferralRewardManagerPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the ReferralRewardManager contract.
type ReferralRewardManagerPausedIterator struct {
	Event *ReferralRewardManagerPaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ReferralRewardManagerPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReferralRewardManagerPaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ReferralRewardManagerPaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ReferralRewardManagerPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReferralRewardManagerPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReferralRewardManagerPaused represents a Paused event raised by the ReferralRewardManager contract.
type ReferralRewardManagerPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_ReferralRewardManager *ReferralRewardManagerFilterer) FilterPaused(opts *bind.FilterOpts) (*ReferralRewardManagerPausedIterator, error) {

	logs, sub, err := _ReferralRewardManager.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &ReferralRewardManagerPausedIterator{contract: _ReferralRewardManager.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_ReferralRewardManager *ReferralRewardManagerFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *ReferralRewardManagerPaused) (event.Subscription, error) {

	logs, sub, err := _ReferralRewardManager.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReferralRewardManagerPaused)
				if err := _ReferralRewardManager.contract.UnpackLog(event, "Paused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_ReferralRewardManager *ReferralRewardManagerFilterer) ParsePaused(log types.Log) (*ReferralRewardManagerPaused, error) {
	event := new(ReferralRewardManagerPaused)
	if err := _ReferralRewardManager.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReferralRewardManagerReferrerSetIterator is returned from FilterReferrerSet and is used to iterate over the raw logs and unpacked data for ReferrerSet events raised by the ReferralRewardManager contract.
type ReferralRewardManagerReferrerSetIterator struct {
	Event *ReferralRewardManagerReferrerSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ReferralRewardManagerReferrerSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReferralRewardManagerReferrerSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ReferralRewardManagerReferrerSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ReferralRewardManagerReferrerSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReferralRewardManagerReferrerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReferralRewardManagerReferrerSet represents a ReferrerSet event raised by the ReferralRewardManager contract.
type ReferralRewardManagerReferrerSet struct {
	User     common.Address
	Referrer common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterReferrerSet is a free log retrieval operation binding the contract event 0x5f7165288eef601591cf549e15ff19ef9060b7f71b9c115be946fa1fe7ebf68a.
//
// Solidity: event ReferrerSet(address indexed user, address indexed referrer)
func (_ReferralRewardManager *ReferralRewardManagerFilterer) FilterReferrerSet(opts *bind.FilterOpts, user []common.Address, referrer []common.Address) (*ReferralRewardManagerReferrerSetIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var referrerRule []interface{}
	for _, referrerItem := range referrer {
		referrerRule = append(referrerRule, referrerItem)
	}

	logs, sub, err := _ReferralRewardManager.contract.FilterLogs(opts, "ReferrerSet", userRule, referrerRule)
	if err != nil {
		return nil, err
	}
	return &ReferralRewardManagerReferrerSetIterator{contract: _ReferralRewardManager.contract, event: "ReferrerSet", logs: logs, sub: sub}, nil
}

// WatchReferrerSet is a free log subscription operation binding the contract event 0x5f7165288eef601591cf549e15ff19ef9060b7f71b9c115be946fa1fe7ebf68a.
//
// Solidity: event ReferrerSet(address indexed user, address indexed referrer)
func (_ReferralRewardManager *ReferralRewardManagerFilterer) WatchReferrerSet(opts *bind.WatchOpts, sink chan<- *ReferralRewardManagerReferrerSet, user []common.Address, referrer []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var referrerRule []interface{}
	for _, referrerItem := range referrer {
		referrerRule = append(referrerRule, referrerItem)
	}

	logs, sub, err := _ReferralRewardManager.contract.WatchLogs(opts, "ReferrerSet", userRule, referrerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReferralRewardManagerReferrerSet)
				if err := _ReferralRewardManager.contract.UnpackLog(event, "ReferrerSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseReferrerSet is a log parse operation binding the contract event 0x5f7165288eef601591cf549e15ff19ef9060b7f71b9c115be946fa1fe7ebf68a.
//
// Solidity: event ReferrerSet(address indexed user, address indexed referrer)
func (_ReferralRewardManager *ReferralRewardManagerFilterer) ParseReferrerSet(log types.Log) (*ReferralRewardManagerReferrerSet, error) {
	event := new(ReferralRewardManagerReferrerSet)
	if err := _ReferralRewardManager.contract.UnpackLog(event, "ReferrerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReferralRewardManagerRewardAddedIterator is returned from FilterRewardAdded and is used to iterate over the raw logs and unpacked data for RewardAdded events raised by the ReferralRewardManager contract.
type ReferralRewardManagerRewardAddedIterator struct {
	Event *ReferralRewardManagerRewardAdded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ReferralRewardManagerRewardAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReferralRewardManagerRewardAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ReferralRewardManagerRewardAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ReferralRewardManagerRewardAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReferralRewardManagerRewardAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReferralRewardManagerRewardAdded represents a RewardAdded event raised by the ReferralRewardManager contract.
type ReferralRewardManagerRewardAdded struct {
	User    common.Address
	Amount  *big.Int
	Manager common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRewardAdded is a free log retrieval operation binding the contract event 0x54dba80f86e498df5a2fcdb5c089d051d97ce8ca9ddb708c70da66557a954de7.
//
// Solidity: event RewardAdded(address indexed user, uint256 amount, address indexed manager)
func (_ReferralRewardManager *ReferralRewardManagerFilterer) FilterRewardAdded(opts *bind.FilterOpts, user []common.Address, manager []common.Address) (*ReferralRewardManagerRewardAddedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	var managerRule []interface{}
	for _, managerItem := range manager {
		managerRule = append(managerRule, managerItem)
	}

	logs, sub, err := _ReferralRewardManager.contract.FilterLogs(opts, "RewardAdded", userRule, managerRule)
	if err != nil {
		return nil, err
	}
	return &ReferralRewardManagerRewardAddedIterator{contract: _ReferralRewardManager.contract, event: "RewardAdded", logs: logs, sub: sub}, nil
}

// WatchRewardAdded is a free log subscription operation binding the contract event 0x54dba80f86e498df5a2fcdb5c089d051d97ce8ca9ddb708c70da66557a954de7.
//
// Solidity: event RewardAdded(address indexed user, uint256 amount, address indexed manager)
func (_ReferralRewardManager *ReferralRewardManagerFilterer) WatchRewardAdded(opts *bind.WatchOpts, sink chan<- *ReferralRewardManagerRewardAdded, user []common.Address, manager []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	var managerRule []interface{}
	for _, managerItem := range manager {
		managerRule = append(managerRule, managerItem)
	}

	logs, sub, err := _ReferralRewardManager.contract.WatchLogs(opts, "RewardAdded", userRule, managerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReferralRewardManagerRewardAdded)
				if err := _ReferralRewardManager.contract.UnpackLog(event, "RewardAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRewardAdded is a log parse operation binding the contract event 0x54dba80f86e498df5a2fcdb5c089d051d97ce8ca9ddb708c70da66557a954de7.
//
// Solidity: event RewardAdded(address indexed user, uint256 amount, address indexed manager)
func (_ReferralRewardManager *ReferralRewardManagerFilterer) ParseRewardAdded(log types.Log) (*ReferralRewardManagerRewardAdded, error) {
	event := new(ReferralRewardManagerRewardAdded)
	if err := _ReferralRewardManager.contract.UnpackLog(event, "RewardAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReferralRewardManagerRewardClaimedIterator is returned from FilterRewardClaimed and is used to iterate over the raw logs and unpacked data for RewardClaimed events raised by the ReferralRewardManager contract.
type ReferralRewardManagerRewardClaimedIterator struct {
	Event *ReferralRewardManagerRewardClaimed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ReferralRewardManagerRewardClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReferralRewardManagerRewardClaimed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ReferralRewardManagerRewardClaimed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ReferralRewardManagerRewardClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReferralRewardManagerRewardClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReferralRewardManagerRewardClaimed represents a RewardClaimed event raised by the ReferralRewardManager contract.
type ReferralRewardManagerRewardClaimed struct {
	User   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRewardClaimed is a free log retrieval operation binding the contract event 0x106f923f993c2149d49b4255ff723acafa1f2d94393f561d3eda32ae348f7241.
//
// Solidity: event RewardClaimed(address indexed user, uint256 amount)
func (_ReferralRewardManager *ReferralRewardManagerFilterer) FilterRewardClaimed(opts *bind.FilterOpts, user []common.Address) (*ReferralRewardManagerRewardClaimedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _ReferralRewardManager.contract.FilterLogs(opts, "RewardClaimed", userRule)
	if err != nil {
		return nil, err
	}
	return &ReferralRewardManagerRewardClaimedIterator{contract: _ReferralRewardManager.contract, event: "RewardClaimed", logs: logs, sub: sub}, nil
}

// WatchRewardClaimed is a free log subscription operation binding the contract event 0x106f923f993c2149d49b4255ff723acafa1f2d94393f561d3eda32ae348f7241.
//
// Solidity: event RewardClaimed(address indexed user, uint256 amount)
func (_ReferralRewardManager *ReferralRewardManagerFilterer) WatchRewardClaimed(opts *bind.WatchOpts, sink chan<- *ReferralRewardManagerRewardClaimed, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _ReferralRewardManager.contract.WatchLogs(opts, "RewardClaimed", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReferralRewardManagerRewardClaimed)
				if err := _ReferralRewardManager.contract.UnpackLog(event, "RewardClaimed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRewardClaimed is a log parse operation binding the contract event 0x106f923f993c2149d49b4255ff723acafa1f2d94393f561d3eda32ae348f7241.
//
// Solidity: event RewardClaimed(address indexed user, uint256 amount)
func (_ReferralRewardManager *ReferralRewardManagerFilterer) ParseRewardClaimed(log types.Log) (*ReferralRewardManagerRewardClaimed, error) {
	event := new(ReferralRewardManagerRewardClaimed)
	if err := _ReferralRewardManager.contract.UnpackLog(event, "RewardClaimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReferralRewardManagerRewardClearedIterator is returned from FilterRewardCleared and is used to iterate over the raw logs and unpacked data for RewardCleared events raised by the ReferralRewardManager contract.
type ReferralRewardManagerRewardClearedIterator struct {
	Event *ReferralRewardManagerRewardCleared // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ReferralRewardManagerRewardClearedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReferralRewardManagerRewardCleared)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ReferralRewardManagerRewardCleared)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ReferralRewardManagerRewardClearedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReferralRewardManagerRewardClearedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReferralRewardManagerRewardCleared represents a RewardCleared event raised by the ReferralRewardManager contract.
type ReferralRewardManagerRewardCleared struct {
	User   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRewardCleared is a free log retrieval operation binding the contract event 0xfb50459538be3884cca2319f1e6c6ff046647e39d9687c512c405b4db5d4b375.
//
// Solidity: event RewardCleared(address indexed user, uint256 amount)
func (_ReferralRewardManager *ReferralRewardManagerFilterer) FilterRewardCleared(opts *bind.FilterOpts, user []common.Address) (*ReferralRewardManagerRewardClearedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _ReferralRewardManager.contract.FilterLogs(opts, "RewardCleared", userRule)
	if err != nil {
		return nil, err
	}
	return &ReferralRewardManagerRewardClearedIterator{contract: _ReferralRewardManager.contract, event: "RewardCleared", logs: logs, sub: sub}, nil
}

// WatchRewardCleared is a free log subscription operation binding the contract event 0xfb50459538be3884cca2319f1e6c6ff046647e39d9687c512c405b4db5d4b375.
//
// Solidity: event RewardCleared(address indexed user, uint256 amount)
func (_ReferralRewardManager *ReferralRewardManagerFilterer) WatchRewardCleared(opts *bind.WatchOpts, sink chan<- *ReferralRewardManagerRewardCleared, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _ReferralRewardManager.contract.WatchLogs(opts, "RewardCleared", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReferralRewardManagerRewardCleared)
				if err := _ReferralRewardManager.contract.UnpackLog(event, "RewardCleared", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRewardCleared is a log parse operation binding the contract event 0xfb50459538be3884cca2319f1e6c6ff046647e39d9687c512c405b4db5d4b375.
//
// Solidity: event RewardCleared(address indexed user, uint256 amount)
func (_ReferralRewardManager *ReferralRewardManagerFilterer) ParseRewardCleared(log types.Log) (*ReferralRewardManagerRewardCleared, error) {
	event := new(ReferralRewardManagerRewardCleared)
	if err := _ReferralRewardManager.contract.UnpackLog(event, "RewardCleared", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReferralRewardManagerRewardPoolFundedIterator is returned from FilterRewardPoolFunded and is used to iterate over the raw logs and unpacked data for RewardPoolFunded events raised by the ReferralRewardManager contract.
type ReferralRewardManagerRewardPoolFundedIterator struct {
	Event *ReferralRewardManagerRewardPoolFunded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ReferralRewardManagerRewardPoolFundedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReferralRewardManagerRewardPoolFunded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ReferralRewardManagerRewardPoolFunded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ReferralRewardManagerRewardPoolFundedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReferralRewardManagerRewardPoolFundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReferralRewardManagerRewardPoolFunded represents a RewardPoolFunded event raised by the ReferralRewardManager contract.
type ReferralRewardManagerRewardPoolFunded struct {
	Funder common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRewardPoolFunded is a free log retrieval operation binding the contract event 0xd3c7d1587d8f93bae0fcfc5a820d2187e755a75ec28a06c20a9da047e1b01eb5.
//
// Solidity: event RewardPoolFunded(address indexed funder, uint256 amount)
func (_ReferralRewardManager *ReferralRewardManagerFilterer) FilterRewardPoolFunded(opts *bind.FilterOpts, funder []common.Address) (*ReferralRewardManagerRewardPoolFundedIterator, error) {

	var funderRule []interface{}
	for _, funderItem := range funder {
		funderRule = append(funderRule, funderItem)
	}

	logs, sub, err := _ReferralRewardManager.contract.FilterLogs(opts, "RewardPoolFunded", funderRule)
	if err != nil {
		return nil, err
	}
	return &ReferralRewardManagerRewardPoolFundedIterator{contract: _ReferralRewardManager.contract, event: "RewardPoolFunded", logs: logs, sub: sub}, nil
}

// WatchRewardPoolFunded is a free log subscription operation binding the contract event 0xd3c7d1587d8f93bae0fcfc5a820d2187e755a75ec28a06c20a9da047e1b01eb5.
//
// Solidity: event RewardPoolFunded(address indexed funder, uint256 amount)
func (_ReferralRewardManager *ReferralRewardManagerFilterer) WatchRewardPoolFunded(opts *bind.WatchOpts, sink chan<- *ReferralRewardManagerRewardPoolFunded, funder []common.Address) (event.Subscription, error) {

	var funderRule []interface{}
	for _, funderItem := range funder {
		funderRule = append(funderRule, funderItem)
	}

	logs, sub, err := _ReferralRewardManager.contract.WatchLogs(opts, "RewardPoolFunded", funderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReferralRewardManagerRewardPoolFunded)
				if err := _ReferralRewardManager.contract.UnpackLog(event, "RewardPoolFunded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRewardPoolFunded is a log parse operation binding the contract event 0xd3c7d1587d8f93bae0fcfc5a820d2187e755a75ec28a06c20a9da047e1b01eb5.
//
// Solidity: event RewardPoolFunded(address indexed funder, uint256 amount)
func (_ReferralRewardManager *ReferralRewardManagerFilterer) ParseRewardPoolFunded(log types.Log) (*ReferralRewardManagerRewardPoolFunded, error) {
	event := new(ReferralRewardManagerRewardPoolFunded)
	if err := _ReferralRewardManager.contract.UnpackLog(event, "RewardPoolFunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReferralRewardManagerRewardReducedIterator is returned from FilterRewardReduced and is used to iterate over the raw logs and unpacked data for RewardReduced events raised by the ReferralRewardManager contract.
type ReferralRewardManagerRewardReducedIterator struct {
	Event *ReferralRewardManagerRewardReduced // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ReferralRewardManagerRewardReducedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReferralRewardManagerRewardReduced)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ReferralRewardManagerRewardReduced)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ReferralRewardManagerRewardReducedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReferralRewardManagerRewardReducedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReferralRewardManagerRewardReduced represents a RewardReduced event raised by the ReferralRewardManager contract.
type ReferralRewardManagerRewardReduced struct {
	User    common.Address
	Amount  *big.Int
	Manager common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRewardReduced is a free log retrieval operation binding the contract event 0x4f98a5b5c72400c0bdafc1fe75f6069e95d1355e66133b9e2e517a00fde2b6bc.
//
// Solidity: event RewardReduced(address indexed user, uint256 amount, address indexed manager)
func (_ReferralRewardManager *ReferralRewardManagerFilterer) FilterRewardReduced(opts *bind.FilterOpts, user []common.Address, manager []common.Address) (*ReferralRewardManagerRewardReducedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	var managerRule []interface{}
	for _, managerItem := range manager {
		managerRule = append(managerRule, managerItem)
	}

	logs, sub, err := _ReferralRewardManager.contract.FilterLogs(opts, "RewardReduced", userRule, managerRule)
	if err != nil {
		return nil, err
	}
	return &ReferralRewardManagerRewardReducedIterator{contract: _ReferralRewardManager.contract, event: "RewardReduced", logs: logs, sub: sub}, nil
}

// WatchRewardReduced is a free log subscription operation binding the contract event 0x4f98a5b5c72400c0bdafc1fe75f6069e95d1355e66133b9e2e517a00fde2b6bc.
//
// Solidity: event RewardReduced(address indexed user, uint256 amount, address indexed manager)
func (_ReferralRewardManager *ReferralRewardManagerFilterer) WatchRewardReduced(opts *bind.WatchOpts, sink chan<- *ReferralRewardManagerRewardReduced, user []common.Address, manager []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	var managerRule []interface{}
	for _, managerItem := range manager {
		managerRule = append(managerRule, managerItem)
	}

	logs, sub, err := _ReferralRewardManager.contract.WatchLogs(opts, "RewardReduced", userRule, managerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReferralRewardManagerRewardReduced)
				if err := _ReferralRewardManager.contract.UnpackLog(event, "RewardReduced", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRewardReduced is a log parse operation binding the contract event 0x4f98a5b5c72400c0bdafc1fe75f6069e95d1355e66133b9e2e517a00fde2b6bc.
//
// Solidity: event RewardReduced(address indexed user, uint256 amount, address indexed manager)
func (_ReferralRewardManager *ReferralRewardManagerFilterer) ParseRewardReduced(log types.Log) (*ReferralRewardManagerRewardReduced, error) {
	event := new(ReferralRewardManagerRewardReduced)
	if err := _ReferralRewardManager.contract.UnpackLog(event, "RewardReduced", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReferralRewardManagerRewardSetIterator is returned from FilterRewardSet and is used to iterate over the raw logs and unpacked data for RewardSet events raised by the ReferralRewardManager contract.
type ReferralRewardManagerRewardSetIterator struct {
	Event *ReferralRewardManagerRewardSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ReferralRewardManagerRewardSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReferralRewardManagerRewardSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ReferralRewardManagerRewardSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ReferralRewardManagerRewardSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReferralRewardManagerRewardSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReferralRewardManagerRewardSet represents a RewardSet event raised by the ReferralRewardManager contract.
type ReferralRewardManagerRewardSet struct {
	User      common.Address
	OldAmount *big.Int
	NewAmount *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRewardSet is a free log retrieval operation binding the contract event 0xbd0682fae90263f394bf341cc4c207fc356246f627963b5fdf3580867cb4a74e.
//
// Solidity: event RewardSet(address indexed user, uint256 oldAmount, uint256 newAmount)
func (_ReferralRewardManager *ReferralRewardManagerFilterer) FilterRewardSet(opts *bind.FilterOpts, user []common.Address) (*ReferralRewardManagerRewardSetIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _ReferralRewardManager.contract.FilterLogs(opts, "RewardSet", userRule)
	if err != nil {
		return nil, err
	}
	return &ReferralRewardManagerRewardSetIterator{contract: _ReferralRewardManager.contract, event: "RewardSet", logs: logs, sub: sub}, nil
}

// WatchRewardSet is a free log subscription operation binding the contract event 0xbd0682fae90263f394bf341cc4c207fc356246f627963b5fdf3580867cb4a74e.
//
// Solidity: event RewardSet(address indexed user, uint256 oldAmount, uint256 newAmount)
func (_ReferralRewardManager *ReferralRewardManagerFilterer) WatchRewardSet(opts *bind.WatchOpts, sink chan<- *ReferralRewardManagerRewardSet, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _ReferralRewardManager.contract.WatchLogs(opts, "RewardSet", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReferralRewardManagerRewardSet)
				if err := _ReferralRewardManager.contract.UnpackLog(event, "RewardSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRewardSet is a log parse operation binding the contract event 0xbd0682fae90263f394bf341cc4c207fc356246f627963b5fdf3580867cb4a74e.
//
// Solidity: event RewardSet(address indexed user, uint256 oldAmount, uint256 newAmount)
func (_ReferralRewardManager *ReferralRewardManagerFilterer) ParseRewardSet(log types.Log) (*ReferralRewardManagerRewardSet, error) {
	event := new(ReferralRewardManagerRewardSet)
	if err := _ReferralRewardManager.contract.UnpackLog(event, "RewardSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReferralRewardManagerRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the ReferralRewardManager contract.
type ReferralRewardManagerRoleAdminChangedIterator struct {
	Event *ReferralRewardManagerRoleAdminChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ReferralRewardManagerRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReferralRewardManagerRoleAdminChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ReferralRewardManagerRoleAdminChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ReferralRewardManagerRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReferralRewardManagerRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReferralRewardManagerRoleAdminChanged represents a RoleAdminChanged event raised by the ReferralRewardManager contract.
type ReferralRewardManagerRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_ReferralRewardManager *ReferralRewardManagerFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*ReferralRewardManagerRoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _ReferralRewardManager.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &ReferralRewardManagerRoleAdminChangedIterator{contract: _ReferralRewardManager.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_ReferralRewardManager *ReferralRewardManagerFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *ReferralRewardManagerRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _ReferralRewardManager.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReferralRewardManagerRoleAdminChanged)
				if err := _ReferralRewardManager.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_ReferralRewardManager *ReferralRewardManagerFilterer) ParseRoleAdminChanged(log types.Log) (*ReferralRewardManagerRoleAdminChanged, error) {
	event := new(ReferralRewardManagerRoleAdminChanged)
	if err := _ReferralRewardManager.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReferralRewardManagerRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the ReferralRewardManager contract.
type ReferralRewardManagerRoleGrantedIterator struct {
	Event *ReferralRewardManagerRoleGranted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ReferralRewardManagerRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReferralRewardManagerRoleGranted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ReferralRewardManagerRoleGranted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ReferralRewardManagerRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReferralRewardManagerRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReferralRewardManagerRoleGranted represents a RoleGranted event raised by the ReferralRewardManager contract.
type ReferralRewardManagerRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_ReferralRewardManager *ReferralRewardManagerFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*ReferralRewardManagerRoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _ReferralRewardManager.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &ReferralRewardManagerRoleGrantedIterator{contract: _ReferralRewardManager.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_ReferralRewardManager *ReferralRewardManagerFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *ReferralRewardManagerRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _ReferralRewardManager.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReferralRewardManagerRoleGranted)
				if err := _ReferralRewardManager.contract.UnpackLog(event, "RoleGranted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_ReferralRewardManager *ReferralRewardManagerFilterer) ParseRoleGranted(log types.Log) (*ReferralRewardManagerRoleGranted, error) {
	event := new(ReferralRewardManagerRoleGranted)
	if err := _ReferralRewardManager.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReferralRewardManagerRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the ReferralRewardManager contract.
type ReferralRewardManagerRoleRevokedIterator struct {
	Event *ReferralRewardManagerRoleRevoked // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ReferralRewardManagerRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReferralRewardManagerRoleRevoked)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ReferralRewardManagerRoleRevoked)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ReferralRewardManagerRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReferralRewardManagerRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReferralRewardManagerRoleRevoked represents a RoleRevoked event raised by the ReferralRewardManager contract.
type ReferralRewardManagerRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_ReferralRewardManager *ReferralRewardManagerFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*ReferralRewardManagerRoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _ReferralRewardManager.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &ReferralRewardManagerRoleRevokedIterator{contract: _ReferralRewardManager.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_ReferralRewardManager *ReferralRewardManagerFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *ReferralRewardManagerRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _ReferralRewardManager.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReferralRewardManagerRoleRevoked)
				if err := _ReferralRewardManager.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_ReferralRewardManager *ReferralRewardManagerFilterer) ParseRoleRevoked(log types.Log) (*ReferralRewardManagerRoleRevoked, error) {
	event := new(ReferralRewardManagerRoleRevoked)
	if err := _ReferralRewardManager.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReferralRewardManagerUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the ReferralRewardManager contract.
type ReferralRewardManagerUnpausedIterator struct {
	Event *ReferralRewardManagerUnpaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ReferralRewardManagerUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReferralRewardManagerUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ReferralRewardManagerUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ReferralRewardManagerUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReferralRewardManagerUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReferralRewardManagerUnpaused represents a Unpaused event raised by the ReferralRewardManager contract.
type ReferralRewardManagerUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_ReferralRewardManager *ReferralRewardManagerFilterer) FilterUnpaused(opts *bind.FilterOpts) (*ReferralRewardManagerUnpausedIterator, error) {

	logs, sub, err := _ReferralRewardManager.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &ReferralRewardManagerUnpausedIterator{contract: _ReferralRewardManager.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_ReferralRewardManager *ReferralRewardManagerFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *ReferralRewardManagerUnpaused) (event.Subscription, error) {

	logs, sub, err := _ReferralRewardManager.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReferralRewardManagerUnpaused)
				if err := _ReferralRewardManager.contract.UnpackLog(event, "Unpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_ReferralRewardManager *ReferralRewardManagerFilterer) ParseUnpaused(log types.Log) (*ReferralRewardManagerUnpaused, error) {
	event := new(ReferralRewardManagerUnpaused)
	if err := _ReferralRewardManager.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReferralRewardManagerUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the ReferralRewardManager contract.
type ReferralRewardManagerUpgradedIterator struct {
	Event *ReferralRewardManagerUpgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ReferralRewardManagerUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReferralRewardManagerUpgraded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ReferralRewardManagerUpgraded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ReferralRewardManagerUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReferralRewardManagerUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReferralRewardManagerUpgraded represents a Upgraded event raised by the ReferralRewardManager contract.
type ReferralRewardManagerUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_ReferralRewardManager *ReferralRewardManagerFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*ReferralRewardManagerUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _ReferralRewardManager.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &ReferralRewardManagerUpgradedIterator{contract: _ReferralRewardManager.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_ReferralRewardManager *ReferralRewardManagerFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *ReferralRewardManagerUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _ReferralRewardManager.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReferralRewardManagerUpgraded)
				if err := _ReferralRewardManager.contract.UnpackLog(event, "Upgraded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_ReferralRewardManager *ReferralRewardManagerFilterer) ParseUpgraded(log types.Log) (*ReferralRewardManagerUpgraded, error) {
	event := new(ReferralRewardManagerUpgraded)
	if err := _ReferralRewardManager.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
