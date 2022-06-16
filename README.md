# dns-magic

# Install
git clone git@github.com:AndrejHalaj/dns-magic.git\
cd dns-magic\
go build

# Basic Usage
./dns-magic encode [-t A/AAAA/... (default: A)] [-v] hostname\
./dns-magic decode request\
./dns-magic file-encode [-t A/AAAA/... (default: A)] [-o outputcsv (default: output.csv)] inputfile
