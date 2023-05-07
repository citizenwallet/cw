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

⏳ Gas Station v1

- Re-sign user transactions and take fees from master wallet
- Call the contract mint function (address, amount) [**onlyOwner**]
- Call the contract burn function (address, amount) [**onlyOwner**]
- Associate/Dissociate an address with a push notification token [**only if address and sender match**]
- Incoming webhook to mint/burn [**onlyOwner**]

⚪️ Event Listener v1

- New pending transaction [**event**]
- New confirmed transaction [**event**]
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

## Additional links

[Repo style guide inspiration](https://www.gobeyond.dev/standard-package-layout/)
