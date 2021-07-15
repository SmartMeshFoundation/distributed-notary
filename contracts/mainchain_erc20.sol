// Abstract contract for the full ERC 20 Token standard
// https://github.com/ethereum/EIPs/issues/20
pragma solidity >=0.4.23 <=0.4.25;

contract Owned {

    /// `owner` is the only address that can call a function with this
    /// modifier
    modifier onlyOwner() {
        require(msg.sender == owner);
        _;
    }

    address public owner;

    /// @notice The Constructor assigns the message sender to be `owner`
    constructor() public {
        owner = msg.sender;
    }

    address newOwner=0x0;

    event OwnerUpdate(address _prevOwner, address _newOwner);

    ///change the owner
    function changeOwner(address _newOwner) public onlyOwner {
        require(_newOwner != owner);
        newOwner = _newOwner;
    }

    /// accept the ownership
    function acceptOwnership() public{
        require(msg.sender == newOwner);
        emit OwnerUpdate(owner, newOwner);
        owner = newOwner;
        newOwner = 0x0;
    }
}



contract LockedSpectrum is Owned {

    function () public {
        revert();
    }

    string public name = "LockedSpectrum for atmosphere";                   //fancy name
    string public version = 'v0.1';
    string public constant symbol = "LSMT";

    struct LockinInfo  {
        bytes32 SecretHash; //这是lockin发起人提供的hash
        uint256 Expiration; //锁过期时间
        uint256 value; //转入金额
    }
    mapping(address=>LockinInfo) public lockin_htlc; //lockin过程中的htlc
    event PrepareLockin(address account,uint256 value,bytes32 secret_hash,uint256 expiration);
    event LockoutSecret(address account, bytes32 secret);
    event PrepareLockout(address account, uint256 _value);
    event Lockin(address account, bytes32 secretHash);
    event CancelLockin(address account, bytes32 secretHash);
    event CancelLockout(address account, bytes32 secretHash);
    struct LockoutInfo {
        bytes32 SecretHash; //转出时指定的密码hash
        uint256 Expiration; //超期以后可以撤销
        uint256 value; //金额是多少
    }
    mapping(address=>LockoutInfo) public lockout_htlc; //lockout 过程中的HTLC

    constructor() public {

    }
    //ze:主链expiration
    //ce:侧链expiration
    //z:主链确认块数转换到侧链的确认块数(比如spectrum和以太坊都是15秒,那转换比率就是1)
    //c:侧链确认块数
    //用户:交易发起人
    // 主链lockin分为两步,
    //第一步:用户在主链上主动prepareLockin到合约中, 指定ze
    //第二步: 用户观察到侧链合约中公证人也进行了prepareLockin,并且金额,过期时间合理(金额考虑手续费,过期时间考虑自己是否有足够时间)
    //第三步: 用户依据密码在侧连上发起lockin,在侧连上获取到相应的token
    //第四步:  公证人观察到侧连上真正发生了lockin(由用户发起),就会知道密码,这时公证人可以在有效期内将主链资产转移到指定合约中去
    //如果交易发起人没有在规定时间内在侧连上进行相应的lockin,公证人(任何人)可以在过期以后在主链cance lockout
    function prepareLockin( bytes32 secret_hash,uint256 expiration)  payable public{
        require(lockin_htlc[msg.sender].value==0);
        require(msg.value > 0);
        LockinInfo storage li=lockin_htlc[msg.sender];
        li.SecretHash=secret_hash;
        li.Expiration=expiration;
        li.value=msg.value;
        emit PrepareLockin(msg.sender,msg.value,secret_hash,expiration);
    }
    //公证人观察到侧链上用户提供的密码,凭密码销毁凭据,防止用户在过期以后再次获取到token
    function lockin(address account,bytes32 secret)   public {
        LockinInfo storage li=lockin_htlc[account];
        //验证密码匹配,并且在有效期内
        require(li.value>0);
        require(li.SecretHash==sha256(abi.encodePacked(secret)));
        require(li.Expiration>block.number);
        // 下发事件用
        bytes32 secretHash = li.SecretHash;

        //根据HTLC信息,为这个账户分配相应的token
        li.value=0;
        li.SecretHash=bytes32(0);
        li.Expiration=0;
        emit Lockin(account, secretHash);
    }
    //lockin过程出错,expiration过期以后,任何人都可以撤销此次交易,实际上最可能的就是用户自己
    function cancelLockin(address account)   public {
        LockinInfo storage li=lockin_htlc[account];
        //已经过期了
        uint256 value=li.value;
        require(li.value>0);
        require(block.number>li.Expiration);
        // 下发事件用
        bytes32 secretHash = li.SecretHash;
        //清空记录,也节省gas
        li.value=0;
        li.SecretHash=bytes32(0);
        li.Expiration=0;

        //清空后在进行,防止递归调用.
        account.transfer(value);
        emit CancelLockin(account, secretHash);
    }

    //退出的过程和lockin过程类似,
    //第一步是退出人(用户)在侧链PrepareLockout,公证人在收到相应的事件以后,会在主链上发起PrepareLockout 要求ce>c+ze
    //第二步 用户观察到主链上发生了PrepareLockout,会在过期时间之内,用密码解锁交易,
    //第三部,公证人则根据主链上观察到的密码,销毁相应的token
    function prepareLockoutHTLC(address account,bytes32 secret_hash,uint256 expiration,uint256 value ) onlyOwner public{
        require(account != 0);
        LockoutInfo storage li=lockout_htlc[account];
        require(value>0);
        require(li.value==0);
        require(expiration > block.number + 300);
        //        require(expiration>50000); //不能低于三天,这样一旦公证人做出了错误的lockout,也应该
        li.value=value;
        li.SecretHash=secret_hash;
        li.Expiration=expiration;
        emit PrepareLockout(account,value);
    }
    //用户提交secret,转移资产,  知道密码的任何人都可以做,允许代理
    function lockout(address account,bytes32 secret)   public {
        LockoutInfo storage li=lockout_htlc[account];
        uint256 value=li.value;
        require(value>0);
        require(li.Expiration>block.number);
        require(li.SecretHash==sha256(abi.encodePacked(secret)));

        li.value=0;
        li.SecretHash=bytes32(0);
        li.Expiration=0;

        account.transfer(value);
        emit LockoutSecret(account,secret);
    }
    //锁过期以后,由公证人取消(任何人)
    function cancleLockOut(address account) public {
        LockoutInfo storage li=lockout_htlc[account];
        uint256 value=li.value;
        require(value>0);
        require(block.number>li.Expiration);
        // 下发事件用
        bytes32 secretHash = li.SecretHash;
        li.value=0;
        li.SecretHash=bytes32(0);
        li.Expiration=0;
        emit CancelLockout(account, secretHash);
    }
    function queryLockin(address account)   view external returns(bytes32,uint256,uint256) {
        LockinInfo storage li=lockin_htlc[account];
        return  (li.SecretHash, li.Expiration,li.value);
    }
    function queryLockout(address account) view external returns(bytes32,uint256,uint256) {
        LockoutInfo storage li=lockout_htlc[account];
        return (li.SecretHash, li.Expiration,li.value);
    }
}
