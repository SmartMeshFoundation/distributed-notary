// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./metaLifeDAO_Base.sol";


abstract contract metaLifeDAOConfig is metaLifeDAOBase{
    uint256 internal _proposalThreshold;
    uint256 internal _votingDelay;
    uint256 internal _votingPeriod;
    uint256 internal _quorumFactorInBP;

    event UpdateDAOMetaInfo();

    function proposalThreshold() public view virtual override returns(uint256){
        return _proposalThreshold;
    }

    function votingDelay() public view virtual override returns(uint256){
        return 0;
    }

    function votingPeriod() public view virtual override returns(uint256){
        return _votingPeriod;
    }

    function quorumFactorInBP() public view virtual returns(uint256){
        return _quorumFactorInBP;
    }

    function setProposalThreshold (uint256 proposalThreshold_) public onlyGovernance{
        _proposalThreshold = proposalThreshold_;
        emit UpdateDAOMetaInfo();
    }

    function setVotingPeriod (uint256 votingPeriod_) public onlyGovernance{
        _votingPeriod = votingPeriod_;
        emit UpdateDAOMetaInfo();
    }

    function setQuorumFactorInBP(uint256 quorumFactorInBP_) public onlyGovernance{
        _quorumFactorInBP = quorumFactorInBP_;
        emit UpdateDAOMetaInfo();
    }

    function setName (string memory _name) public onlyGovernance{
        daoName = _name;
        emit UpdateDAOMetaInfo();
    }

    function setURI (string memory _uri) public onlyGovernance{
        daoURI = _uri;
        emit UpdateDAOMetaInfo();
    }

    function setInfo (string memory _info) public onlyGovernance{
        daoInfo = _info;
        emit UpdateDAOMetaInfo();
    }

    constructor (
        string memory daoName_,
        string memory daoURI_,
        string memory daoInfo_,
        string memory version_
    ) metaLifeDAOBase(daoName_, daoURI_, daoInfo_) EIP712(daoName_, version_) {
        _proposalThreshold = 1;
    }
}

abstract contract metaLifeDAOiwthDelegation is metaLifeDAOBase{
    event DelegateChanged(address indexed delegator, address indexed fromDelegate, address indexed toDelegate);
    event DelegateVotesChanged(address indexed delegate, uint256 previousBalance, uint256 newBalance);

    mapping(address => address) private _delegation;

    function delegates(address account) public view virtual returns (address) {
        return _delegation[account];
    }

    function delegate(address delegatee) public virtual {
        address account = msg.sender;
        _delegate(account, delegatee);
    }

    function delegateBySig(
        address delegatee,
        uint256 nonce,
        uint256 expiry,
        uint8 v,
        bytes32 r,
        bytes32 s
    ) public virtual {
        require(block.timestamp <= expiry, "Expired");
        address signer = ECDSA.recover(
            _hashTypedDataV4(keccak256(abi.encode(keccak256("Delegation(address delegatee,uint256 nonce,uint256 expiry)"), delegatee, nonce, expiry))),
            v,
            r,
            s
        );
        _delegate(signer, delegatee);
    }

    function _delegate(address account, address delegatee) internal virtual {
        address oldDelegate = delegates(account);
        _delegation[account] = delegatee;

        emit DelegateChanged(account, oldDelegate, delegatee);
    }

    function castVoteByDelegation(address delegator, uint256 proposalId, uint8 support) public virtual returns (uint256) {
        require(delegates(delegator) == msg.sender, "Rejected");

        return _castVote(proposalId, delegator, support);
    }
}

abstract contract metaLifeDAOSimple is metaLifeDAOiwthDelegation, metaLifeDAOConfig{
    enum VoteType {
        Against,
        For,
        Abstain
    }

    struct ProposalVote {
        uint256 againstVotes;
        uint256 forVotes;
        uint256 abstainVotes;
        mapping(address => bool) hasVoted;
    }

    mapping(uint256 => ProposalVote) private _proposalVotes;

    function hasVoted(uint256 proposalId, address account) public view virtual override returns (bool) {
        return _proposalVotes[proposalId].hasVoted[account];
    }

    function proposalVotes(uint256 proposalId)
            public
            view
            virtual
            returns (
                uint256 againstVotes,
                uint256 forVotes,
                uint256 abstainVotes
            )
        {
            ProposalVote storage proposalvote = _proposalVotes[proposalId];
            return (proposalvote.againstVotes, proposalvote.forVotes, proposalvote.abstainVotes);
        }

    function _checkVote(uint256 proposalId, address account) internal virtual {
        ProposalVote storage proposalvote = _proposalVotes[proposalId];

        require(!proposalvote.hasVoted[account], "Already cast");
        proposalvote.hasVoted[account] = true;
    }

    function _countVote(
        uint256 proposalId,
        address account,
        uint8 support,
        uint256 weight
    ) internal virtual override {
        ProposalVote storage proposalvote = _proposalVotes[proposalId];

        _checkVote(proposalId, account);

        if (support == uint8(VoteType.Against)) {
            proposalvote.againstVotes += weight;
        } else if (support == uint8(VoteType.For)) {
            proposalvote.forVotes += weight;
        } else if (support == uint8(VoteType.Abstain)) {
            proposalvote.abstainVotes += weight;
        } else {
            revert("Invalid vote");
        }
    }

    function _quorumReached(uint256 proposalId) internal view override virtual returns (bool){
        ProposalVote storage proposalvote = _proposalVotes[proposalId];

        return quorum(proposalSnapshot(proposalId)) <= proposalvote.againstVotes + proposalvote.forVotes + proposalvote.abstainVotes;
    }

    function _voteSucceeded(uint256 proposalId) internal view override virtual returns (bool){
        ProposalVote storage proposalvote = _proposalVotes[proposalId];

        return proposalvote.forVotes > proposalvote.againstVotes;
    }

    event DAOCanceled();

    function cancelDAO() external virtual;
}
