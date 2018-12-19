// Abstract contract for the full ERC 20 Token standard
// https://github.com/ethereum/EIPs/issues/20
pragma solidity ^0.4.24;

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
    function Owned() public {
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

contract EthereumToken is StandardToken {

    function () public {
        revert();
    }

    string public name = "Ethereum Token for atmosphere";                   //fancy name
    uint8 public decimals = 18;                //How many decimals to show. ie. There could 1000 base units with 3 decimals. Meaning 0.980 SBX = 980 base units. It's like comparing 1 wei to 1 ether.
    string public symbol = "EthereumToken";                 //An identifier
    string public version = 'v0.1';       //SMT 0.1 standard. Just an arbitrary versioning scheme.

    // The nonce for avoid transfer replay attacks
    mapping(address => uint256) nonces;
    struct LockinInfo  {
     bytes32 SecretHash; //这是lockin发起人提供的hash
        int64 Expiration; //锁过期时间
        uint256 value; //转入金额
    }
    mapping(address=>LockinInfo) public lockin_htlc; //lockin过程中的htlc
    event LockinHTLC(address account,bytes32 secret_hash,int64 expiration,uint256 value);

    struct LockoutInfo {
        bytes32 SecretHash; //转出时指定的密码hash
        int64 Expiration; //超期以后可以撤销
        uint256 value; //金额是多少
        bytes32 Data;  //附加数据
    }
    mapping(address=>LockoutInfo) public lockout_htlc; //lockout 过程中的HTLC

    function EthereumToken() public {
        totalSupply=1; //保证total supply大于等于0,符合ERC20规范
    }

    //只能公证人提供HTLC
    // lockin分为两步,
    //第一步:公证人为根据收到的请求在侧连上prepareLockin,
    //第二步: 公证人在观察到主链上发生了相应的交易,会执行lockin
    //如果交易发起人没有在规定时间内在主链上进行相应的转入交易,公证人可以在过期以后撤销HTLC
    function prepareLockinHTLC(address account,bytes32 secret_hash,int64 expiration,uint256 value) onlyOwner public{
        require(lockin_htlc[account].value==0);
        LockinInfo storage li=lockin_htlc[account];
        li.SecretHash=secret_hash;
        li.Expiration=expiration;
        li.value=value;
        emit LockinHTLC(account,secret_hash,expiration,value);
    }
    //正常应该有公证人来为账户分配token,如果公证人没有在指定期限内做,其他任何知道密码的人都可以做.
    function lockin(address account,bytes32 secret)   public {
        LockinInfo storage li=lockin_htlc[account];
        //验证密码匹配,并且在有效期内
        require(li.value>0);
        require(li.SecretHash==keccak256(abi.encodePacked(secret)));
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
    }
    function cancelLockin(address account) onlyOwner public {
        LockinInfo storage li=lockin_htlc[account];
        //已经过期了
        require(li.value>0);
        require(block.number>li.Expiration);

        //清空记录,也节省gas
        li.value=0;
        li.SecretHash=bytes32(0);
        li.Expiration=0;
    }

    //退出的过程和lockin过程类似,也分为两步,第一步是退出人提出退出请求,公证人在收到相应的事件以后,会在主链上发起相应交易
    //第二步 退出人观察到主链上发生了相应的交易,会用密码解锁交易,
    //第三部,公证人则根据观察到的密码,销毁相应的token
    event PrePareLockedOut(address indexed _from, uint256 _value);
    //准备退出,需要在合约里记录,公证人需要监控这里的事件,采用相应的操作, 后续应该提供PrePareLockedOutProxy函数,可以保证交易方没有侧链主币的情况下,仍然可以发起合约
    function prePareLockedOutHTLC(bytes32 secret_hash,int64 expiration,uint256 value, bytes32 data) public{
        LockoutInfo storage li=lockout_htlc[msg.sender];
        require(li.value==0); // 没有正在退出的历史交易
        require(balances[msg.sender]>=value);
        li.value=value;
        li.SecretHash=secret_hash;
        li.Expiration=expiration;
        li.Data=data;

        balances[msg.sender]-=value;

        emit PrePareLockedOut(msg.sender,value);
    }
    //公证人在主链上转账到具体地址后,销毁相应记录
    function lockedOut(address from,bytes32 secret)   public {
        LockoutInfo storage li=lockout_htlc[msg.sender];
        uint256 value=li.value;
        require(value>0);
        require(li.Expiration>block.number);
        require(li.SecretHash==keccak256(abi.encodePacked(secret)));

        li.value=0;
        li.SecretHash=bytes32(0);
        li.Expiration=0;
        li.Data=bytes32(0);
    }
    //锁过期以后,由自己取消
    function cancleLockOut( ) public {
        LockoutInfo storage li=lockout_htlc[msg.sender];
        uint256 value=li.value;
        require(value>0);
        require(block.number>li.Expiration);
        li.value=0;
        li.SecretHash=bytes32(0);
        li.Expiration=0;
        li.Data=bytes32(0);
        //退回到个人账户上
        balances[msg.sender]+=value;
    }
}