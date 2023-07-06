# crypto_keygen
Multithreaded generation of BTC or ETH keys, matched to vanity patterns or existing addresses and saved to a sqlite db

## Overview

This program will generate crypto address keys and look for matches of any pattern or address in any text file saved in one of the match directories. The text files can be named anything, as long as they are ".txt" files, with one address or pattern per line. Due to the nature of comparison on a map type object, matching a full address is much quicker than looking for a leading vanity pattern, so each of these types of matching are processed differently. Place lists of addresses in the appropriate "_address" folder, and vanity patterns to try and match in the appropriate "_pattern" folder.

For BTC addresses, the program creates segwit addresses (which have a leading "bc1" in the address, and are case insensitive).

For ETH pattern matches, any characters that should be represented as numbers will be converted to check for a match automatically. For example, you could have "Decaf coffee" in a .txt file in the "match_eth_pattern" folder. This will be automatically converted to "0xdecafc0ffee" when looking for a generated matching address.

Matches found (public and private keys) will be written to the console, and also saved to a local sqlite database.

## Clone the repository

```console
git clone https://github.com/dmisino/crypto_keygen.git
cd crypto_keygen
```

## Installation and setup

To run this project, you'll need Go installed on your machine. You can install it from [go.dev](https://go.dev/doc/install).

Install required go modules, and build the project:

```bash
go get
go build
```

## Usage

For generating keys, run one of the following:

```bash
crypto_keygen btc
crypto_keygen eth
```

## Settings: page size and worker threads

The number of keys to be generated on each iteration of the program loop are configured in the .env file, set by default as 1000000. This can also be overridden as an additional command line argument:

```bash
crypto_keygen eth 50000
```

The slowest part of trying to brute force match existing addresses or in creating new vanity addresses is the act of creating the new addresses. Once created, the comparison can be fairly quick. To speed up the process of generating addresses, multiple threads are spawned. A .env file setting defaults this to 8 worker threads, though the best number for your system will depend on a number of factors, such as total CPU cores available. You can play with this setting to try and achieve the fastest processing possible.

## Creating match lists

The lists you create to try and match will depend on what you are trying to accomplish. Here is some info that might be useful.

### ETH vanity addresses

Ethereum addresses can have any number 0-9, or the letters a-f, since they are hexidecimal numbers.

Using numbers as letters, you could use the following letters when creating a vanity word or phrase: abcdefois (o,i and s represented by 0, 1 and 5 respectively). Find words containing specific letters in English, French, German or Spanish: https://www.dcode.fr/words-containing

### Chance of getting an address collision

The odds of generating an address, and having that address match another specific address are ridiculously small.

For BTC: 1 in 2^160

For ETH: 1 in 2^256

Even if for example, you had a list of 1 million BTC addresses you were trying to find, and you managed to generate 1 trillion addresses per day to check for a match. And say you did that every day for 50 years. What would the odds be of ever finding the private keys for one of your target addresses? The math gets a little messy, but the answer in short is zero. From any practical standpoint, the odds are still effectively zero.

### ETH addresses

An ethereum address is "0x" + exactly 40 characters. Try to discover the private key for these addresses:

0x3141592653589793238462643383279502884197

0x1618033988749894848204586834365638117720

Those addresses would be perfect representations of Pi and The Golden Ratio, and might be the ultimate vanity addresses. The chances of ANYONE ever finding those, without some encryption breaking advance in technology, are still basically zero.

### BTC addresses

Download a text file with all BTC addresses that have any balance (over 47 million addresses currently): http://addresses.loyce.club/ 
