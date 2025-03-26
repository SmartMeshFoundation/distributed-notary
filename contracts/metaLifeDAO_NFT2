// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./metaLifeDAO_Simple.sol";
import "./metaLifeDAO_NFT.sol";
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "@openzeppelin/contracts/token/ERC721/IERC721Receiver.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/IERC721Enumerable.sol";

interface metaMaster{
    function mint(address _collection, address _recipient, string memory _hash) external returns(uint);

    function mintAndSell(address _collection, string memory _hash, uint market, address _token, uint _price, uint _duration) external;

    function sell(address collection, uint token_id, uint market, address _token, uint _price, uint _duration) external;

    function createCollection(string memory _name, string memory _symbol, string memory _baseURI, uint _max_supply, string memory _metaInfo,
        address royaltiesReceiver, uint royaltiesPercentageInBips) external returns(address);
}

contract metaLifeDAONFT2 is _metaLifeDAONFT, IERC721Receiver {

    string private constant _version="MetaLifeDAO:212:NFT2";

    address public metaMasterContract;

    address public starter;

    constructor (string memory _daoName,
        string memory _daoURI,
        string memory _daoInfo,
        uint64 votingPeriod_,
        uint256 quorumFactorInBP_,
        address metaMaster_,
        string memory _baseURI,
        uint256 _max_supply,
        address starter_
    ) metaLifeDAOConfig(_daoName, _daoURI, _daoInfo, _version) {
        metaMasterContract = metaMaster_;
        _votingPeriod = votingPeriod_;
        _quorumFactorInBP = quorumFactorInBP_;
        bindedNFT = metaMaster(metaMasterContract).createCollection(_daoName, _daoName, _baseURI, _max_supply, _daoInfo, address(this), 250);
        starter = starter_;
    }

    function mintAndSell(string memory _hash, uint market, address _token, uint _price, uint _duration) external{
        require(msg.sender == starter);

        metaMaster(metaMasterContract).mintAndSell(bindedNFT, _hash, market, _token, _price, _duration);
    }

    function onERC721Received(
        address,
        address,
        uint256,
        bytes memory
    ) public virtual override returns (bytes4) {
        return this.onERC721Received.selector;
    }
}
