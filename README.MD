# words-parser

This script parses "known words", "unknown words" and saves the diff list to csv file with word frequency.

### Installation
```sh
git clone https://github.com/igor-pgmt/words-parser.git
cd words-parser/
go build
sudo chmod +x words-parser
```

### Usage

##### You can run this program with next flags:
```sh
Usage of ./words-parser:
  -help
        Please, specify to show this message. [bool] (default: false
  -knownWordsDir string
        Please, specify path to resulting file. [string] (default: ./knownWords/ (default "./knownWords/")
  -resultingFile string
        Please, specify path to resulting file. [string] (default: ./result.csv (default "result.csv")
  -unknownWordsDir string
```

##### Examples:

Connect to the fastest server:
```sh
./words-parser 
```

Connect to the fastest server:
```sh
./words-parser -knownWordsDir ./fds -unknownWordsDir ~/dir2 -resultingFile res.csv
```
