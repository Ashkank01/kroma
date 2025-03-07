// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Votes.sol";

import { KromaMintableERC20 } from "../universal/KromaMintableERC20.sol";

/**
 * @custom:proxied
 * @custom:predeploy 0x42000000000000000000000000000000000000ff
 * @title GovernanceToken
 * @notice The KRO token used in governance, supporting voting and delegation. Implements
 *         EIP 2612 allowing signed approvals. It can be bridged to other specified chain.
 *         This token is minted once every block by the MintManager.
 */
contract GovernanceToken is KromaMintableERC20, ERC20Burnable, ERC20Votes {
    /**
     * @notice Address of the MintManager contract that will have minter role.
     */
    address public immutable MINT_MANAGER;

    /**
     * @notice Constructs the GovernanceToken contract.
     *
     * @param _bridge      Address of the StandardBridge on this network.
     * @param _remoteToken Address of the corresponding token on the remote chain.
     * @param _mintManager Address of the MintManager contract.
     */
    constructor(
        address _bridge,
        address _remoteToken,
        address _mintManager
    ) KromaMintableERC20(_bridge, _remoteToken, "Kroma", "KRO") ERC20Permit("Kroma") {
        MINT_MANAGER = _mintManager;
    }

    /**
     * @notice Allows the owner to mint tokens.
     *
     * @param _account The account receiving minted tokens.
     * @param _amount  The amount of tokens to mint.
     */
    function mint(address _account, uint256 _amount) external override {
        require(
            msg.sender == BRIDGE || msg.sender == MINT_MANAGER,
            "GovernanceToken: only minter can mint"
        );

        _mint(_account, _amount);
    }

    /**
     * @notice Allows the StandardBridge on this network to burn tokens.
     *
     * @param _from   Address to burn tokens from.
     * @param _amount Amount of tokens to burn.
     */
    function burn(
        address _from,
        uint256 _amount
    ) external override onlyBridge {
        _burn(_from, _amount);
    }

    /**
     * @notice Callback called after a token transfer.
     *
     * @param _from   The account sending tokens.
     * @param _to     The account receiving tokens.
     * @param _amount The amount of tokens being transferred.
     */
    function _afterTokenTransfer(
        address _from,
        address _to,
        uint256 _amount
    ) internal override(ERC20, ERC20Votes) {
        super._afterTokenTransfer(_from, _to, _amount);
    }

    /**
     * @notice Internal mint function.
     *
     * @param _account The account receiving minted tokens.
     * @param _amount  The amount of tokens to mint.
     */
    function _mint(address _account, uint256 _amount) internal override(ERC20, ERC20Votes) {
        super._mint(_account, _amount);
    }

    /**
     * @notice Internal burn function.
     *
     * @param _account The account that tokens will be burned from.
     * @param _amount  The amount of tokens that will be burned.
     */
    function _burn(address _account, uint256 _amount) internal override(ERC20, ERC20Votes) {
        super._burn(_account, _amount);
    }
}
