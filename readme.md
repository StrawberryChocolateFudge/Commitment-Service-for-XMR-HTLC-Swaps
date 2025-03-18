# Commitment Provider - HTLC support for Monero Trading
WIP
Example UIs:
![current ui](current_ui.png "current ui")


Example:

![Get Secret current UI](GetSecret.png "get secret example ui")

The commitment provider aims to tackle a problem with Monero Atomic Swaps, they are difficult to develop.

The commitment provider is a secret provider infrastructure that provides secrets for Hash Time Lock Contracts, to use a simple commitment-reveal scheme while trading monero.

The HTLC contracts can be deployed on Bitcoin, Litecoin, Ethereum, Solana or any chain that can support it.

Monero doesn't have smart contracts and developing multisig is complex and requires the users to download a specialized wallet and complex messaging.

Instead of that, the Commitment Provider is a (decentralizable) third party that verifies Monero transactions and reveals secrets in-exchange.

The commitment provider:

1. Never has access to any funds
2. Never knows where the commitments and secrets are actually used
3. Can't be used to steal funds by third party, if the database is hacked.
4. It's not an escrow and doesn't provide dispute resolution
5. It's more similar to an Oracle that checks Monero payments and reveals secrets to specific users.
6. Discards the monero payment proofs after the secret is revealed.

The commitment provider reveals the pre-image of a commitment to a user that can provide a valid proof of monero transaction. Payment proofs are checked using `check_tx_key TXID TXKEY ADDRESS`

### Example Flow
```
* Alice wants to exchange her XMR to Sol with Bob.
* Alice sends Bob her Sol address and they agree on the exchange rate
* Bob contacts the Commitment Provider and requests a Commitment, When requesting the commitment Bob configures the commitment provider to reveal the secret of the commitment for a payment proof of X amount made to his address
* Bob interacts with a smart contract on Solana and deposits the Sol to trade with the commitment and Alice's address
* The HTLC on Solana unlocks the Sol if the secret is provided by Alice else it refunds it in 1 day back to Bob
* Alice makes the XMR deposit to Bobs address
* Alice then creates a payment proof and uses it to get the secret from the commitment provider 
* Alice then pulls payments from the Smart contract using the secret
```


The explained flow works on every smart contract chain and Solana was just an example.

### Benefits
1. The commitment provider offers a simple flow for trading
2. It's more easy to develop for than a Multisig Escrow
3. Trading is always P2P, Alice transfers directly to Bob, then Bob transfer to Alice via Smart Contract
4. No local app to run for users. No local cli apps needed to swap
5. It works in the browser and uses existing monero wallet


### Considerations
* The users need to have 2 browser windows open. One to use the commitment provider's interface and one to do a swap with a browser wallet on chain
* It can however be used as a JSON API and embedded into applications/websites or served as an IFRAME
* It must stay online, else there can be loss of funds due to the HTLC timelock expiration, but secret recovery mechanisms can be implemented in many ways
* The Commitments can be used to query for the XMR address  (only while the commitment is valid)  The XMR addresses are never stored on smart contract chains this way.

### Monetizing
The commitment provider is a paid service that works with a `pay per commitment` model. Payments are needed also to stop malicious users from requesting too many commitments.

Example: 
* Bob deposits XMR to the commitment provider's address
* Bob uses the Proof of Deposit to get an API key from the Commitment Provider
* API key has a quota
* Bob uses the API Key to create commitments at the provider and when he uses his quota he needs to get a new key
* Example fees could be: 1000 commitments per $0.1 worth of XMR 


### Decentralization
The commitment provider infrastructure can be decentralized, however instead of building a complex network to create it, a simpler approach should be used, it should work as a DAO.

* N amount of people can become commitment providers
* They create an N of M multisig account via Monero CLI and communicate via chat
* The commitment providers configure their instances to pay into the Multisig account
* At the end of the month, commitment providers withdraw their profits from the multisig in consensus
* If a commitment service goes down and this results in loss of funds, reimbursement could be requested via chat
* A provider that goes down too often could be kicked out if it results in loss of funds

This would allow for a simple decentralization where humans are representing each instance and large networks of nodes don't exist. Having 4-5 providers are more than enough.

It would be also possible to implement decentralized storage for secrets via threshold homomorphic encryption which is up for research (golang latigo library can work for this, but I don't want to over complicate things)

## API

The API serves HTML without Javascript and there is a JSON API

//TODO: This is under development


## How commitments are computed

The commitment is a sha256 hash of a 32 byte secret buffer.
The commitment and the secret can be converted string from to a buffer using using Javascript with this:

```
function decodeString(hexString) {
  return Buffer.from(hexString, 'hex');
}

// Encode a byte array to a hex string
function encodeToString(byteArray) {
  return Buffer.from(byteArray).toString('hex');
}

```

Following this will make sure you stay consistent with the Go implementation which uses `hex.EncodeToString` and `hex.DecodeString` and `sha256.Sum256`

The corresponding hash in javascript is

`const hashBuffer = await crypto.subtle.digest("SHA-256", buff);`

So you can recreate the hash provided by the commitment provider in Javascript easily. This should be compatible with all sha256 hashes provided by on-chain contracts too!


## Contract
Here is an example Hash time lock contract in solidity that would be compatible with this

```
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract HTLC {
    address payable public sender;
    address payable public recipient;
    bytes32 public hashLock;
    uint256 public timeLock;
    bool public isClaimed;

    event Funded(uint256 amount);
    event Claimed(bytes32 secret);
    event Refunded();

    modifier onlySender() {
        require(msg.sender == sender, "Not the sender");
        _;
    }

    modifier onlyRecipient() {
        require(msg.sender == recipient, "Not the recipient");
        _;
    }

    constructor(
        address payable _recipient,
        bytes32 _hashLock,
        uint256 _timeLock
    ) payable {
        sender = payable(msg.sender);
        recipient = _recipient;
        hashLock = _hashLock;
        timeLock = block.timestamp + _timeLock;
        isClaimed = false;

        require(msg.value > 0, "No funds provided");
        emit Funded(msg.value);
    }

    function claim(bytes32 _secret) public onlyRecipient {
        require(sha256(abi.encodePacked(_secret)) == hashLock, "Invalid secret");
        require(block.timestamp < timeLock, "Time lock expired");
        require(!isClaimed, "Already claimed");

        isClaimed = true;
        recipient.transfer(address(this).balance);
        emit Claimed(_secret);
    }

    function refund() public onlySender {
        require(block.timestamp >= timeLock, "Time lock not expired");
        require(!isClaimed, "Already claimed");

        sender.transfer(address(this).balance);
        emit Refunded();
    }
}

```