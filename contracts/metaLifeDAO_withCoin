// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./metaLifeDAO_Simple.sol";

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/utils/Checkpoints.sol";
import "@openzeppelin/contracts/utils/Timers.sol";


abstract contract _metaLifeDAOwithCoin is ERC20, metaLifeDAOiwthDelegation, metaLifeDAOSimple {
    using Checkpoints for Checkpoints.History;
    using Timers for Timers.BlockNumber;

    //------------
    //Override Interface
    //------------

    string private constant _version="ABSTRACT";

    function quorum(uint256 blockNumber) public view virtual override returns (uint256){
        return getPastTotalSupply(blockNumber) * quorumFactorInBP()/ 10000;
    }

    function quorum() public view virtual returns (uint256){
        return totalSupply() * quorumFactorInBP()/ 10000;
    }

    function getVotes(address account, uint256 blockNumber) public view override returns (uint256){
        return getPastVotes(account, blockNumber);
    }

    function cancelDAO() external override {
        require(getVotes(msg.sender) == totalSupply());
        emit DAOCanceled();
    }

    //------------
    //ERC20 Module
    //------------

    function mint(address to, uint256 amount) public onlyGovernance{
        _mint(to, amount);
    }

    function _afterTokenTransfer(
        address from,
        address to,
        uint256 amount
    ) internal virtual override{
        _transferVotingUnits(from, to, amount);
    }

    //------------
    //Vote Module
    //------------
    mapping(address => address) private _delegation;
    mapping(address => Checkpoints.History) private _delegateCheckpoints;
    Checkpoints.History private _totalCheckpoints;

    function delegates(address account) public view override returns(address){
        if (_delegation[account] == address(0)){
            return account;
        } else {
            return _delegation[account];
        }
    }

    function getVotes(address account) public view virtual returns (uint256) {
        return _delegateCheckpoints[account].latest();
    }

    function getPastVotes(address account, uint256 blockNumber) internal view virtual returns (uint256) {
        if(blockNumber < proposalCounts){ //reuse blockNumber as proposalId
            blockNumber = _proposals[blockNumber].voteStart.getDeadline();
        }
        return _delegateCheckpoints[account].getAtBlock(blockNumber);
    }

    function getPastTotalSupply(uint256 blockNumber) public view virtual returns (uint256) {
        require(blockNumber <= block.number, "Block not mined");
        if (blockNumber == block.number){
            return totalSupply();
        }
        return _totalCheckpoints.getAtBlock(blockNumber);
    }

    function _delegate(address account, address delegatee) internal virtual override {
        address oldDelegate = delegates(account);
        _delegation[account] = delegatee;

        emit DelegateChanged(account, oldDelegate, delegatee);
        _moveDelegateVotes(oldDelegate, delegatee, balanceOf(account));
    }

    function _transferVotingUnits(
        address from,
        address to,
        uint256 amount
    ) internal virtual {
        if (from == address(0)) {
            _totalCheckpoints.push(_totalCheckpoints.latest() + amount);
        }
        if (to == address(0)) {
            _totalCheckpoints.push(_totalCheckpoints.latest() - amount);
        }
        _moveDelegateVotes(delegates(from), delegates(to), amount);
    }

    function _moveDelegateVotes(
        address from,
        address to,
        uint256 amount
    ) private {
        if (from != to && amount > 0) {
            if (from != address(0)) {
                (uint256 oldValue, uint256 newValue) = _delegateCheckpoints[from].push(_delegateCheckpoints[from].latest() - amount);
                emit DelegateVotesChanged(from, oldValue, newValue);
            }
            if (to != address(0)) {
                (uint256 oldValue, uint256 newValue) = _delegateCheckpoints[to].push(_delegateCheckpoints[to].latest() + amount);
                emit DelegateVotesChanged(to, oldValue, newValue);
            }
        }
    }

    constructor (string memory _daoName,
        string memory _daoURI,
        string memory _daoInfo,
        string memory version_,
        uint64 votingPeriod_,
        uint256 quorumFactorInBP_
    ) metaLifeDAOConfig(_daoName, _daoURI, _daoInfo, version_) ERC20(_daoName, _daoName){
        _votingPeriod = votingPeriod_;
        _quorumFactorInBP = quorumFactorInBP_;
    }

}


contract metaLifeDAOwithCoin is _metaLifeDAOwithCoin{
    string private constant _version="MetaLifeDAO:101:withCoin";

    constructor (string memory _daoName,
        string memory _daoURI,
        string memory _daoInfo,
        uint64 votingPeriod_,
        uint256 quorumFactorInBP_,
        address[] memory initialMembers,
        uint256[] memory initialSupplies
    ) _metaLifeDAOwithCoin(_daoName, _daoURI, _daoInfo, _version, votingPeriod_, quorumFactorInBP_) {
        require(initialMembers.length == initialSupplies.length, "input");
        for (uint i; i < initialMembers.length; i++){
            _mint(initialMembers[i], initialSupplies[i]);
        }
    }
}
