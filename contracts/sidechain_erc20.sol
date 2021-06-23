// Abstract contract for the full ERC 20 Token standard
// https://github.com/ethereum/EIPs/issues/20
pragma solidity >=0.4.23 <=0.4.25;

contract Token {
    /* This is a slight change to the ERC20 base standard.*/
    /// total amount of tokens
    uint256 public totalSupply;

    /// @param _owner The address from which the balance will be retrieved
    /// @return The balance
    function balanceOf(address _owner) public view returns (uint256 balance);

    /// @notice send `_value` token to `_to` from `msg.sender`
    /// @param _to The address of the recipient
    /// @param _value The amount of token to be transferred
    /// @return Whether the transfer was successful or not
    function transfer(address _to, uint256 _value) public returns (bool success);

    /// @notice send `_value` token to `_to` from `_from` on the condition it is approved by `_from`
    /// @param _from The address of the sender
    /// @param _to The address of the recipient
    /// @param _value The amount of token to be transferred
    /// @return Whether the transfer was successful or not
    function transferFrom(address _from, address _to, uint256 _value) public returns (bool success);

    /// @notice `msg.sender` approves `_spender` to spend `_value` tokens
    /// @param _spender The address of the account able to transfer the tokens
    /// @param _value The amount of tokens to be approved for transfer
    /// @return Whether the approval was successful or not
    function approve(address _spender, uint256 _value) public returns (bool success);

    /// @param _owner The address of the account owning tokens
    /// @param _spender The address of the account able to transfer the tokens
    /// @return Amount of remaining tokens allowed to spent
    function allowance(address _owner, address _spender) public view returns (uint256 remaining);

    event Transfer(address indexed _from, address indexed _to, uint256 _value);
    event Approval(address indexed _owner, address indexed _spender, uint256 _value);
}

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

contract StandardToken is Token,Owned {

    function transfer(address _to, uint256 _value) public  returns (bool success) {
        //Default assumes totalSupply can't be over max (2^256 - 1).
        //If your token leaves out totalSupply and can issue more tokens as time goes on, you need to check if it doesn't wrap.
        //Replace the if with this one instead.
        if (balances[msg.sender] >= _value && balances[_to] + _value > balances[_to]) {
            balances[msg.sender] -= _value;
            balances[_to] += _value;
            emit Transfer(msg.sender, _to, _value);
            return true;
        } else { return false; }
    }

    function transferFrom(address _from, address _to, uint256 _value) public returns (bool success) {
        //same as above. Replace this line with the following if you want to protect against wrapping uints.
        if (balances[_from] >= _value && allowed[_from][msg.sender] >= _value && balances[_to] + _value > balances[_to]) {
            balances[_to] += _value;
            balances[_from] -= _value;
            allowed[_from][msg.sender] -= _value;
            emit Transfer(_from, _to, _value);
            return true;
        } else { return false; }
    }

    function balanceOf(address _owner) public view returns (uint256 balance) {
        return balances[_owner];
    }

    function approve(address _spender, uint256 _value) public returns (bool success) {
        allowed[msg.sender][_spender] = _value;
        emit Approval(msg.sender, _spender, _value);
        return true;
    }

    function allowance(address _owner, address _spender) public view returns (uint256 remaining) {
        return allowed[_owner][_spender];
    }

    mapping (address => uint256) balances;
    mapping (address => mapping (address => uint256)) allowed;
}

contract AtmosphereToken is StandardToken {

    function () public {
        revert();
    }

    string public name;                   //fancy name
    uint8 public decimals = 18;                //How many decimals to show. ie. There could 1000 base units with 3 decimals. Meaning 0.980 SBX = 980 base units. It's like comparing 1 wei to 1 ether.
    string public version = 'v0.1';
    string public symbol ;

    struct LockinInfo  {
     bytes32 SecretHash; //这是lockin发起人提供的hash
        uint256 Expiration; //锁过期时间
        uint256 value; //转入金额
    }
    mapping(address=>LockinInfo) public lockin_htlc; //lockin过程中的htlc
    event PrepareLockin(address account, uint256 value);
    event LockinSecret(address account,bytes32 secret);
    event PrepareLockout(address account, uint256 _value);
    event Lockout(address account, bytes32 secretHash);
    event CancelLockin(address account, bytes32 secretHash);
    event CancelLockout(address account, bytes32 secretHash);
    struct LockoutInfo {
        bytes32 SecretHash; //转出时指定的密码hash
        uint256 Expiration; //超期以后可以撤销
        uint256 value; //金额是多少
    }
    mapping(address=>LockoutInfo) public lockout_htlc; //lockout 过程中的HTLC

    constructor(string tokenName) public {
        totalSupply=1; //保证total supply大于等于0,符合ERC20规范
        name = tokenName;
        symbol = tokenName;
    }
    //ze:mainchain expiration
    //ce:side chain expiration
    //z:主链确认块数转换到侧链的确认块数(比如spectrum和以太坊都是15秒,那转换比率就是1)
    //c:侧链确认块数
    //用户:交易发起人
    //只能公证人提供HTLC
    // lockin分为两步,
    //第一步:公证人收到用户lockin请求,并且观察到主链上发生了prepareLockin,在侧链进行prepareLockin  要求ze>ce+z
    //第二步:  公证人观察到侧连上真正发生了lockin(由用户发起),就会知道密码,这时公证人可以在有效期内将主链资产转移到指定合约中去
    //如果交易发起人没有在规定时间内在侧连上进行相应的lockin,公证人(任何人)可以在过期以后,在侧链撤销HTLC(cancel lockin)
    function prepareLockin(address account,bytes32 secret_hash,uint256 expiration,uint256 value) onlyOwner public{
        require(account != 0);
        require(lockin_htlc[account].value==0);
        LockinInfo storage li=lockin_htlc[account];
        li.SecretHash=secret_hash;
        li.Expiration=expiration;
        li.value=value;
        emit PrepareLockin(account,value);
    }
    //由用户提供密码,真正的为自己的账户分配token,其他任何知道密码的人也都可以做.
    function lockin(address account,bytes32 secret)   public {
        LockinInfo storage li=lockin_htlc[account];
        //验证密码匹配,并且在有效期内
        require(li.value>0);
        require(li.SecretHash==sha256(abi.encodePacked(secret)));
        require(li.Expiration>block.number);

        //根据HTLC信息,为这个账户分配相应的token
        uint256 value=li.value;
        uint256 oldValue=balances[account];
        require(oldValue+value>=value);
        balances[account]=oldValue+value;
        totalSupply=totalSupply+value;
        require(totalSupply>value); // 部署合约时的那个1一直存在,还不能溢出

        //清空记录,也节省gas
        li.value=0;
        li.SecretHash=bytes32(0);
        li.Expiration=0;
        emit LockinSecret(account, secret);
    }
    //lockin过程出错,expiration过期以后,任何人都可以撤销此次交易
    function cancelLockin(address account)   public {
        LockinInfo storage li=lockin_htlc[account];
        //已经过期了
        require(li.value>0);
        require(block.number>li.Expiration);
        // 下发事件用
        bytes32 secretHash = li.SecretHash;

        //清空记录,也节省gas
        li.value=0;
        li.SecretHash=bytes32(0);
        li.Expiration=0;
        emit CancelLockin(account,secretHash);
    }

    //退出的过程和lockin过程类似,第一步是退出人(用户)在侧链PrepareLockout,公证人在接收到用户请求,并且收到相应的事件以后(需要足够的确认块数),会在主链上PrepareLockout 要求ce>c+ze
    //第二步 用户观察到主链上发生了PrepareLockout,会在过期时间之内,用密码在主链上进行lockout
    //第三部,公证人则根据主链上观察到的密码,销毁相应的token
    //准备退出,需要在合约里记录,公证人需要监控这里的事件,采用相应的操作, 后续应该提供PrepareLockoutProxy函数,可以保证交易方没有侧链主币的情况下,仍然可以发起合约
    function prepareLockout(bytes32 secret_hash,uint256 expiration,uint256 value) public{
        LockoutInfo storage li=lockout_htlc[msg.sender];
        require(value>0);
        require(li.value==0); // 没有正在退出的历史交易
        require(balances[msg.sender]>=value);
        li.value=value;
        li.SecretHash=secret_hash;
        li.Expiration=expiration;

        balances[msg.sender]-=value;

        emit PrepareLockout(msg.sender,value);
    }
    //用户在主链上提交secret,资产已经转移走, 公证人(任何人)观察到密码,销毁相应的token
    function lockout(address from,bytes32 secret)   public {
        LockoutInfo storage li=lockout_htlc[from];
        uint256 value=li.value;
        require(value>0);
        require(li.Expiration>block.number);
        require(li.SecretHash==sha256(abi.encodePacked(secret)));
        // 下发事件用
        bytes32 secretHash = li.SecretHash;

        li.value=0;
        li.SecretHash=bytes32(0);
        li.Expiration=0;
        //侧链发行总量要降低
        totalSupply-=value;
        require(totalSupply>=1);
        emit Lockout(from, secretHash);
    }
    //锁过期以后,由用户取消 其他任何人也都可以做
    function cancelLockOut(address account) public {
        LockoutInfo storage li=lockout_htlc[account];
        uint256 value=li.value;
        require(value>0);
        require(block.number>li.Expiration);
        // 下发事件用
        bytes32 secretHash = li.SecretHash;
        li.value=0;
        li.SecretHash=bytes32(0);
        li.Expiration=0;
        //退回到个人账户上
        balances[account]+=value;
        emit CancelLockout(account, secretHash);
    }
    function queryLockin(address account) view external returns(bytes32,uint256,uint256) {
        LockinInfo storage li=lockin_htlc[account];
        return (li.SecretHash, li.Expiration,li.value);
    }
    function queryLockout(address account)   view external returns(bytes32,uint256,uint256) {
        LockoutInfo storage li=lockout_htlc[account];
        return  (li.SecretHash, li.Expiration,li.value);
    }

}
