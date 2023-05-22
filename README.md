<h1 align="center">
  <img style="height: 100px; width: 100px;" src="https://github.com/daobrussels/cw/blob/main/logos/logo.png" alt="citizen wallet logo"/><br/>
  Citizen Wallet
</h1>

Receive and send citizen coins to pay at participating events.

Move your leftovers coins to your Citizen Wallet on your smartphone.

[Read more.](https://citizenwallet.xyz/)

# cw

## Intro

The repo's purpose is to provide gas station functionality and event listening on a given chain for the Citizen Wallet app. It should be able to integrate with any EVM compatible chain.

Each individual program should be stateless and be able to be scaled up/down to multiple instances.

## Roadmap

⚪️ pending ⏳ in progress ✅ done

✅ Setup

- project structure
- main functions for programs
- router
- health check (middleware)
- signature (middleware)
- services
- modules
- env

✅ Gas Station (ERC-4337) v1

- Re-sign user transactions and take fees from master wallet
- Gateway (EntryPoint)
- Paymaster [**onlyEntryPoint**]
- Account Factory [**onlyEntryPoint**]
- Gratitude Token Factory [**onlyEntryPoint**]
- Profile Factory [**onlyEntryPoint**]
- Account [**onlyOwnerOrEntryPoint**]
- Gratitude [**onlyOwnerOrEntryPoint**]
- Profile [**onlyOwnerOrEntryPoint**]

⚪️ Gas Station (ERC-4337) v1.1

- Notification Subscriber Factory [**onlyEntryPoint**]
- Notification Subscriber [**onlyOwnerOrEntryPoint**]

⚪️ Event Listener v1

- New block with relevant transactions [**event**]
- Notify all tokens of an associated address [**onEvent**]

## Installation

`go get ./...`

## Set up environment

`cp .example.env .env`

Replace values in `.env` for your setup

## Run Gas Station

`go run cmd/station/main.go -url endpoint`

## Run Blockchain Event Handler

`go run cmd/events/main.go -url endpoint`

## Spin up a TestChain

```
docker run --publish 8545:8545 trufflesuite/ganache:latest --account "0x429321276245f7d39855c8040f498af9392cafed95e1e4f50d158b2b39faa9cc,100000000000000000000000" --account "0xe1b5da7d6c2009c09dcb30781ec1dc4e9f73598a26b57e742d706102b69a1716,100000000000000000000000" --account "0x45c532f2bcb9a21f1a25b1d739bd9d3d65209e86836f370897c94e2e571ec18d,100000000000000000000000" --chain.chainId 1682515751360 --chain.networkId 1682515751360 --unlock "0x0b772F674eD6fB67C5647Be0fbBd2FBe95156D60" --unlock "0xBa711ff057dfAC08E4568Bb972EeC2313454f55A" --unlock "0x664ce0F7785E4bA5Ff422C77314eF982F193BeF5"
b6f582e89807d1a6529d7d724bd5d3e1188b46904fa8bffcaf4102edcb27687b
```

## Additional links

[Repo style guide inspiration](https://www.gobeyond.dev/standard-package-layout/)
